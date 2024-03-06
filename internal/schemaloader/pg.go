package schemaloader

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/shared/pg"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for pg driver

	"github.com/artarts36/db-exporter/internal/schema"
)

type PGLoader struct {
}

type constraint struct {
	Name              string `db:"name"`
	TableName         string `db:"table_name"`
	ColumnName        string `db:"column_name"`
	Type              string `db:"type"`
	ForeignTableName  string `db:"foreign_table_name"`
	ForeignColumnName string `db:"foreign_column_name"`
}

type squashedConstraint struct {
	Name              string
	TableName         string
	ColumnsNames      *ds.Strings
	Type              string
	ForeignTableName  string
	ForeignColumnName string
}

func (l *PGLoader) Load(ctx context.Context, dsn string) (*schema.Schema, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed connect to db: %w", err)
	}

	query := `
select c.column_name as name,
       c.table_name,
       c.data_type as type,
       pg_catalog.col_description(format('%s.%s',c.table_schema,c.table_name)::regclass::oid,c.ordinal_position) as "comment",
       case
			when is_nullable = 'NO' THEN false
			else true
	   END as nullable
from information_schema.columns c
where c.table_schema = $1
order by c.ordinal_position`

	var cols []*schema.Column

	err = db.SelectContext(ctx, &cols, query, "public")
	if err != nil {
		return nil, err
	}

	tables := map[schema.String]*schema.Table{}

	constraints, err := l.loadConstraints(ctx, db, "public")
	if err != nil {
		return nil, fmt.Errorf("failed to load constraints: %w", err)
	}

	for _, col := range cols {
		table, tableExists := tables[col.TableName]
		if !tableExists {
			table = &schema.Table{
				Name:        col.TableName,
				ForeignKeys: map[string]*schema.ForeignKey{},
				UniqueKeys:  map[string]*schema.UniqueKey{},
			}
			tables[col.TableName] = table
		}

		col.PreparedType = l.prepareColumnType(col)

		l.applyConstraints(table, col, constraints[col.TableName.Value][col.Name.Value])

		table.Columns = append(table.Columns, col)
	}

	return &schema.Schema{
		Tables: tables,
	}, nil
}

func (l *PGLoader) applyConstraints(table *schema.Table, col *schema.Column, constraints []*squashedConstraint) {
	for _, constr := range constraints {
		if constr.Type == pg.ConstraintPKName {
			pk := table.PrimaryKey
			if pk == nil {
				pk = &schema.PrimaryKey{
					Name: schema.String{
						Value: constr.Name,
					},
					ColumnsNames: constr.ColumnsNames,
				}

				table.PrimaryKey = pk
			}

			col.PrimaryKey = pk
		} else if constr.Type == pg.ConstraintFKName {
			fk := table.ForeignKeys[constr.Name]

			if fk == nil {
				fk = &schema.ForeignKey{
					Name: schema.String{
						Value: constr.Name,
					},
					Table:        table.Name,
					ColumnsNames: constr.ColumnsNames,
					ForeignTable: schema.String{
						Value: constr.ForeignTableName,
					},
					ForeignColumn: schema.String{
						Value: constr.ForeignColumnName,
					},
				}

				table.ForeignKeys[constr.Name] = fk
			}

			col.ForeignKey = fk
		} else if constr.Type == pg.ConstraintUniqueName {
			uk := table.UniqueKeys[constr.Name]

			if uk == nil {
				uk = &schema.UniqueKey{
					Name: schema.String{
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

func (l *PGLoader) prepareColumnType(col *schema.Column) schema.ColumnType {
	switch col.Type.Value {
	case pg.TypeTimestampWithoutTZ:
		return schema.ColumnTypeTimestamp
	case pg.TypeInteger:
		return schema.ColumnTypeInteger
	case pg.TypeBoolean:
		return schema.ColumnTypeBoolean
	case pg.TypeReal:
		return schema.ColumnTypeFloat
	default:
		return schema.ColumnTypeString
	}
}

func (l *PGLoader) loadConstraints(
	ctx context.Context,
	db *sqlx.DB,
	schemaName string,
) (map[string]map[string][]*squashedConstraint, error) {
	query := `select
       tco.constraint_name as "name",
       kcu.table_name,
       kcu.column_name,
       tco.constraint_type as "type",
       ccu.table_name AS foreign_table_name,
       ccu.column_name AS foreign_column_name
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
		return nil, err
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
				Name:              constr.Name,
				TableName:         constr.TableName,
				ColumnsNames:      ds.NewStrings(constr.ColumnName),
				Type:              constr.Type,
				ForeignTableName:  constr.ForeignTableName,
				ForeignColumnName: constr.ForeignColumnName,
			}

			squashed[constr.Name] = sc
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

	return constraintMap, nil
}
