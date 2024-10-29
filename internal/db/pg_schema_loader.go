package db

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/regex"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for pg driver

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/shared/pg"
)

type PGLoader struct{}

var (
	pgColumnDefaultValueStringRegexp   = regexp.MustCompile(`^'(.*)'::character varying$`)
	pgColumnDefaultValueFuncRegexp     = regexp.MustCompile(`^(.*)\(\)$`)
	pgColumnDefaultValueSequenceRegexp = regexp.MustCompile(`^nextval\('(.*)'::regclass\)$`)
)

var pgTypeMap = map[string]schema.DataType{
	pg.TypeText:             schema.DataTypeString,
	pg.TypeUUID:             schema.DataTypeString,
	pg.TypeCharacter:        schema.DataTypeString,
	pg.TypeCharacterVarying: schema.DataTypeString,

	pg.TypeTimestampWithoutTZ: schema.DataTypeTimestamp,
	pg.TypeTimestampWithTZ:    schema.DataTypeTimestamp,

	pg.TypeInteger: schema.DataTypeInteger,
	pg.TypeInt4:    schema.DataTypeInteger,
	pg.TypeInt8:    schema.DataTypeInteger,
	pg.TypeSerial:  schema.DataTypeInteger,

	pg.TypeSmallInt:    schema.DataTypeInteger16,
	pg.TypeSmallSerial: schema.DataTypeInteger16,

	pg.TypeBigint: schema.DataTypeInteger64,

	pg.TypeBoolean: schema.DataTypeBoolean,
	pg.TypeBit:     schema.DataTypeBoolean,

	pg.TypeDoublePrecision: schema.DataTypeFloat32,
	pg.TypeFloat8:          schema.DataTypeFloat32,
	pg.TypeDecimal:         schema.DataTypeFloat32,

	pg.TypeMoney:   schema.DataTypeFloat64,
	pg.TypeReal:    schema.DataTypeFloat64,
	pg.TypeNumeric: schema.DataTypeFloat64,

	pg.TypeBytea: schema.DataTypeBytes,
}

type constraint struct {
	Name       string `db:"name"`
	TableName  string `db:"table_name"`
	ColumnName string `db:"column_name"`
	Type       string `db:"type"`

	ForeignTableName  string `db:"foreign_table_name"`
	ForeignColumnName string `db:"foreign_column_name"`

	IsDeferrable        bool `db:"is_deferrable"`
	IsInitiallyDeferred bool `db:"initially_deferred"`
}

type squashedConstraint struct {
	Name         string
	TableName    string
	ColumnsNames *ds.Strings
	Type         string

	ForeignTableName  string
	ForeignColumnName string

	IsDeferrable        bool
	IsInitiallyDeferred bool
}

func NewPGLoader() *PGLoader {
	return &PGLoader{}
}

func (l *PGLoader) Load(ctx context.Context, conn *Connection) (*schema.Schema, error) {
	db, err := conn.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed connect to db: %w", err)
	}

	query := `
select c.column_name as name,
       c.table_name,
       c.data_type as type,
       pg_catalog.col_description(format('%s.%s',c.table_schema,c.table_name)::regclass::oid,c.ordinal_position)
           as "comment",
       case
			when is_nullable = 'NO' THEN false
			else true
	   END as nullable,
       c.column_default as default_value
from information_schema.columns c
where c.table_schema = $1
order by c.ordinal_position`

	var cols []*schema.Column

	slog.DebugContext(ctx, "[pgloader] loading columns")

	err = db.SelectContext(ctx, &cols, query, conn.cfg.Schema) //nolint: musttag // false-positive
	if err != nil {
		return nil, err
	}

	slog.DebugContext(ctx, fmt.Sprintf("[pgloader] loaded %d columns", len(cols)))

	tables := schema.NewTableMap()

	slog.DebugContext(ctx, "[pgloader] loading constraints")

	constraints, constraintsCount, err := l.loadConstraints(ctx, db, "public")
	if err != nil {
		return nil, fmt.Errorf("failed to load constraints: %w", err)
	}

	slog.DebugContext(ctx, fmt.Sprintf("[pgloader] loaded %d constraints", constraintsCount))

	slog.DebugContext(ctx, "[pgloader] loading sequences")

	sequences, err := l.loadSequences(ctx, db, "public")
	if err != nil {
		return nil, fmt.Errorf("failed to load sequences")
	}

	slog.DebugContext(ctx, fmt.Sprintf("[pgloader] loaded %d sequences", len(sequences)))

	for _, col := range cols {
		table, tableExists := tables.Get(col.TableName)
		if !tableExists {
			table = &schema.Table{
				Name:           col.TableName,
				ForeignKeys:    map[string]*schema.ForeignKey{},
				UniqueKeys:     map[string]*schema.UniqueKey{},
				UsingSequences: map[string]*schema.Sequence{},
			}

			tables.Add(table)
		}

		col.PreparedType = l.prepareDataType(col.Type.Value)
		col.Default = l.parseColumnDefault(col)

		if col.Default != nil && col.Default.Type == schema.ColumnDefaultTypeAutoincrement {
			seqName, ok := col.Default.Value.(string)
			if !ok {
				return nil, fmt.Errorf("failed to get sequence name for %s.%s", table.Name, col.Name)
			}

			seq, seqExists := sequences[seqName]
			if !seqExists {
				return nil, fmt.Errorf("failed to get sequence %q for %s.%s", seqName, table.Name, col.Name)
			}

			seq.Used++

			col.UsingSequences = map[string]*schema.Sequence{
				seqName: seq,
			}

			table.UsingSequences[seqName] = seq
		}

		l.applyConstraints(table, col, constraints[col.TableName.Value][col.Name.Value])

		table.Columns = append(table.Columns, col)
	}

	return &schema.Schema{
		Tables:    tables,
		Sequences: sequences,
	}, nil
}

func (l *PGLoader) parseColumnDefault(col *schema.Column) *schema.ColumnDefault {
	if !col.DefaultRaw.Valid {
		return nil
	}

	if col.DefaultRaw.String == "false" {
		return &schema.ColumnDefault{
			Type:  schema.ColumnDefaultTypeValue,
			Value: false,
		}
	}

	if col.DefaultRaw.String == "true" {
		return &schema.ColumnDefault{
			Type:  schema.ColumnDefaultTypeValue,
			Value: true,
		}
	}

	if col.PreparedType.IsInteger() {
		if parsedInt, intErr := strconv.Atoi(col.DefaultRaw.String); intErr == nil {
			return &schema.ColumnDefault{
				Type:  schema.ColumnDefaultTypeValue,
				Value: parsedInt,
			}
		}
	}

	if val := regex.ParseSingleValue(pgColumnDefaultValueStringRegexp, col.DefaultRaw.String); val != "" {
		return &schema.ColumnDefault{
			Type:  schema.ColumnDefaultTypeValue,
			Value: val,
		}
	}

	if val := regex.ParseSingleValue(pgColumnDefaultValueFuncRegexp, col.DefaultRaw.String); val != "" {
		return &schema.ColumnDefault{
			Type:  schema.ColumnDefaultTypeFunc,
			Value: val,
		}
	}

	if val := regex.ParseSingleValue(pgColumnDefaultValueSequenceRegexp, col.DefaultRaw.String); val != "" {
		return &schema.ColumnDefault{
			Type:  schema.ColumnDefaultTypeAutoincrement,
			Value: val,
		}
	}

	return nil
}

func (l *PGLoader) applyConstraints(table *schema.Table, col *schema.Column, constraints []*squashedConstraint) {
	for _, constr := range constraints {
		switch constr.Type {
		case pg.ConstraintPKName:
			pk := table.PrimaryKey
			if pk == nil {
				pk = &schema.PrimaryKey{
					Name: ds.String{
						Value: constr.Name,
					},
					ColumnsNames: constr.ColumnsNames,
				}

				table.PrimaryKey = pk
			}

			col.PrimaryKey = pk
		case pg.ConstraintFKName:
			fk := table.ForeignKeys[constr.Name]

			if fk == nil {
				fk = &schema.ForeignKey{
					Name: ds.String{
						Value: constr.Name,
					},
					Table:        table.Name,
					ColumnsNames: constr.ColumnsNames,
					ForeignTable: ds.String{
						Value: constr.ForeignTableName,
					},
					ForeignColumn: ds.String{
						Value: constr.ForeignColumnName,
					},
					IsDeferrable:        constr.IsDeferrable,
					IsInitiallyDeferred: constr.IsInitiallyDeferred,
				}

				table.ForeignKeys[constr.Name] = fk
			}

			col.ForeignKey = fk
		case pg.ConstraintUniqueName:
			uk := table.UniqueKeys[constr.Name]

			if uk == nil {
				uk = &schema.UniqueKey{
					Name: ds.String{
						Value: constr.Name,
					},
					ColumnsNames: constr.ColumnsNames,
				}

				table.UniqueKeys[constr.Name] = uk
			}

			col.UniqueKey = uk
		}
	}
}

func (l *PGLoader) prepareDataType(rawType string) schema.DataType {
	t, exists := pgTypeMap[rawType]
	if exists {
		return t
	}

	return schema.DataTypeString
}

func (l *PGLoader) loadSequences(ctx context.Context, db *sqlx.DB, schemaName string) (
	map[string]*schema.Sequence,
	error,
) {
	query := `select
    s.sequence_name as name,
    s.data_type as data_type from information_schema.sequences s
where s.sequence_schema = $1`

	var sequences []*schema.Sequence

	err := db.SelectContext(ctx, &sequences, query, schemaName)
	if err != nil {
		return nil, err
	}

	sequenceMap := map[string]*schema.Sequence{}

	for _, sequence := range sequences {
		sequence.PreparedDataType = l.prepareDataType(sequence.DataType)

		sequenceMap[sequence.Name] = sequence
	}

	return sequenceMap, nil
}

func (l *PGLoader) loadConstraints(
	ctx context.Context,
	db *sqlx.DB,
	schemaName string,
) (map[string]map[string][]*squashedConstraint, int, error) {
	count := 0

	query := `select
       tco.constraint_name as "name",
       kcu.table_name,
       kcu.column_name,
       tco.constraint_type as "type",
       ccu.table_name AS foreign_table_name,
       ccu.column_name AS foreign_column_name,
       case
			when is_deferrable = 'NO' THEN false
			else true
	   END as is_deferrable,
       case
			when initially_deferred = 'NO' THEN false
			else true
	   END as initially_deferred
from information_schema.table_constraints tco
         join information_schema.key_column_usage kcu
              on kcu.constraint_name = tco.constraint_name
                  and kcu.constraint_schema = tco.constraint_schema
                  and kcu.constraint_name = tco.constraint_name
         join information_schema.constraint_column_usage AS ccu
              on ccu.constraint_name = tco.constraint_name
where kcu.table_schema = $1
order by kcu.table_schema,
         kcu.table_name,
         kcu.ordinal_position;`

	var constraints []*constraint

	err := db.SelectContext(ctx, &constraints, query, schemaName)
	if err != nil {
		return nil, count, err
	}

	squashed := map[string]*squashedConstraint{}
	constraintMap := map[string]map[string][]*squashedConstraint{}

	for _, constr := range constraints {
		sc, scExists := squashed[constr.Name]
		if scExists {
			if !sc.ColumnsNames.Contains(constr.ColumnName) {
				sc.ColumnsNames.Add(constr.ColumnName)
			}
		} else {
			sc = &squashedConstraint{
				Name:                constr.Name,
				TableName:           constr.TableName,
				ColumnsNames:        ds.NewStrings(constr.ColumnName),
				Type:                constr.Type,
				ForeignTableName:    constr.ForeignTableName,
				ForeignColumnName:   constr.ForeignColumnName,
				IsDeferrable:        constr.IsDeferrable,
				IsInitiallyDeferred: constr.IsInitiallyDeferred,
			}

			squashed[constr.Name] = sc

			count++
		}

		_, exists := constraintMap[constr.TableName]
		if !exists {
			constraintMap[constr.TableName] = map[string][]*squashedConstraint{}
		}
		constraintMap[constr.TableName][constr.ColumnName] = append(
			constraintMap[constr.TableName][constr.ColumnName],
			sc,
		)
	}

	return constraintMap, count, nil
}
