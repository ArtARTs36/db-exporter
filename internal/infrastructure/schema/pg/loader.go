package pg

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/db-exporter/internal/shared/regex"
	"github.com/artarts36/gds"
	"log/slog"
	"regexp"
	"strconv"

	_ "github.com/lib/pq" // for pg driver

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/pg"
)

type Loader struct{}

var (
	pgColumnDefaultValueStringRegexp   = regexp.MustCompile(`^'(.*)'::character varying$`)
	pgColumnDefaultValueFuncRegexp     = regexp.MustCompile(`^(.*)\(\)$`)
	pgColumnDefaultValueSequenceRegexp = regexp.MustCompile(`^nextval\('(.*)'::regclass\)$`)
)

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
	ColumnsNames *gds.Strings
	Type         string

	ForeignTableName  string
	ForeignColumnName string

	IsDeferrable        bool
	IsInitiallyDeferred bool
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(ctx context.Context, cn *conn.Connection) (*schema.Schema, error) { //nolint:funlen // not need
	db, err := cn.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}

	query := `
select c.column_name as name,
       c.table_name,
       c.domain_name as domain_name,
       case
			when (c.data_type = 'USER-DEFINED') then c.udt_name
			else c.data_type
	   END as type_raw,
       pg_catalog.col_description(format('%s.%s',c.table_schema,c.table_name)::regclass::oid,c.ordinal_position)
           as "comment",
       case
			when is_nullable = 'NO' THEN false
			else true
	   END as nullable,
       c.column_default as default_value,
       case 
           when c.character_maximum_length is null THEN 0
           else c.character_maximum_length
       END as character_length
from information_schema.columns c
where c.table_schema = $1
order by c.ordinal_position`

	sch := schema.NewSchema(schema.DatabaseDriverPostgres)

	var cols []*schema.Column

	slog.DebugContext(ctx, "[pgloader] loading columns")

	err = db.SelectContext(ctx, &cols, query, cn.Database().Schema)
	if err != nil {
		return nil, err
	}

	slog.DebugContext(ctx, fmt.Sprintf("[pgloader] loaded %d columns", len(cols)))

	slog.DebugContext(ctx, "[pgloader] loading constraints")

	constraints, constraintsCount, err := l.loadConstraints(ctx, cn)
	if err != nil {
		return nil, fmt.Errorf("failed to load constraints: %w", err)
	}

	slog.DebugContext(ctx, fmt.Sprintf("[pgloader] loaded %d constraints", constraintsCount))

	slog.DebugContext(ctx, "[pgloader] loading sequences")

	sch.Sequences, err = l.loadSequences(ctx, cn)
	if err != nil {
		return nil, fmt.Errorf("failed to load sequences: %w", err)
	}

	sch.Enums, err = l.loadEnums(ctx, cn)
	if err != nil {
		return nil, fmt.Errorf("failed to load enums: %w", err)
	}

	slog.DebugContext(ctx, fmt.Sprintf("[pgloader] loaded %d sequences", len(sch.Sequences)))

	// loading domains
	slog.DebugContext(ctx, "[pgloader] loading domains")

	sch.Domains, err = l.loadDomains(ctx, cn)
	if err != nil {
		return nil, fmt.Errorf("load domains: %w", err)
	}

	slog.DebugContext(ctx, "[pgloader] domains loaded", slog.Int("domains_count", sch.Domains.Len()))

	// mapping columns
	for _, col := range cols {
		table, tableExists := sch.Tables.Get(col.TableName)
		if !tableExists {
			table = schema.NewTable(col.TableName)

			sch.Tables.Add(table)
		}

		col.DataType = sqltype.MapPGType(col.TypeRaw.Value)
		col.Default = l.parseColumnDefault(col)

		enum, enumExists := sch.Enums[col.TypeRaw.Value]
		if enumExists {
			col.Enum = enum
			enum.Used++
			enum.UsingInTables = append(enum.UsingInTables, table.Name.Value)
			table.UsingEnums[enum.Name.Value] = enum
		}

		domain, domainExists := sch.Domains.Get(col.TypeRaw.Value)
		if domainExists {
			col.Domain = domain
			domain.UsingInTables = append(domain.UsingInTables, table.Name.Value)
			table.UsingDomains[domain.Name] = domain
		}

		if col.Default != nil && col.Default.Type == schema.ColumnDefaultTypeAutoincrement {
			col.IsAutoincrement = true

			seqName, ok := col.Default.Value.(string)
			if !ok {
				return nil, fmt.Errorf("failed to get sequence name for %s.%s", table.Name, col.Name)
			}

			seq, seqExists := sch.Sequences[seqName]
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

		table.AddColumn(col)
	}

	return sch, nil
}

func (l *Loader) parseColumnDefault(col *schema.Column) *schema.ColumnDefault {
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

	if col.DefaultRaw.String == "CURRENT_TIMESTAMP" {
		return &schema.ColumnDefault{
			Type:  schema.ColumnDefaultTypeFunc,
			Value: col.DefaultRaw.String,
		}
	}

	if col.DataType.IsInteger {
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

func (l *Loader) applyConstraints(table *schema.Table, col *schema.Column, constraints []*squashedConstraint) {
	for _, constr := range constraints {
		switch constr.Type {
		case pg.ConstraintPKName:
			pk := table.PrimaryKey
			if pk == nil {
				pk = &schema.PrimaryKey{
					Name: gds.String{
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
					Name: gds.String{
						Value: constr.Name,
					},
					Table:        table.Name,
					ColumnsNames: constr.ColumnsNames,
					ForeignTable: gds.String{
						Value: constr.ForeignTableName,
					},
					ForeignColumn: gds.String{
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
					Name: gds.String{
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

func (l *Loader) loadEnums(ctx context.Context, conn *conn.Connection) (map[string]*schema.Enum, error) {
	query := `select
       t.typname as enum_name,
       e.enumlabel as enum_value
from pg_type t
         join pg_enum e on t.oid = e.enumtypid
         join pg_catalog.pg_namespace n ON n.oid = t.typnamespace
where n.nspname = $1`

	type enumValue struct {
		EnumName  string `db:"enum_name"`
		EnumValue string `db:"enum_value"`
	}

	var enumValues []*enumValue

	db, err := conn.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = db.SelectContext(ctx, &enumValues, query, conn.Database().Schema)
	if err != nil {
		return nil, err
	}

	enums := map[string]*schema.Enum{}

	for _, value := range enumValues {
		enum, ok := enums[value.EnumName]
		if !ok {
			enum = &schema.Enum{
				Name:          gds.NewString(value.EnumName),
				Values:        make([]string, 0),
				UsingInTables: make([]string, 0),
			}
		}

		enum.Values = append(enum.Values, value.EnumValue)

		enums[value.EnumName] = enum
	}

	return enums, nil
}

func (l *Loader) loadSequences(ctx context.Context, conn *conn.Connection) (
	map[string]*schema.Sequence,
	error,
) {
	query := `select
    s.sequence_name as name,
    s.data_type as data_type_raw from information_schema.sequences s
where s.sequence_schema = $1`

	var sequences []*schema.Sequence

	db, err := conn.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = db.SelectContext(ctx, &sequences, query, conn.Database().Schema)
	if err != nil {
		return nil, err
	}

	sequenceMap := map[string]*schema.Sequence{}

	for _, sequence := range sequences {
		sequence.DataType = sqltype.MapPGType(sequence.DataTypeRaw)

		sequenceMap[sequence.Name] = sequence
	}

	return sequenceMap, nil
}

func (l *Loader) loadConstraints(
	ctx context.Context,
	cn *conn.Connection,
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

	db, err := cn.Connect(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = db.SelectContext(ctx, &constraints, query, cn.Database().Schema)
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
				ColumnsNames:        gds.NewStrings(constr.ColumnName),
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
