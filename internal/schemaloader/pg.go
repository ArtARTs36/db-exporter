package schemaloader

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/artarts36/db-exporter/internal/schema"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const pgConstraintPKName = "PRIMARY KEY"
const pgConstraintFKName = "FOREIGN KEY"

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

	constraints, err := l.loadConstraints(db, ctx, "public")
	if err != nil {
		return nil, fmt.Errorf("failed to load constraints: %w", err)
	}

	for _, col := range cols {
		table, tableExists := tables[col.TableName]
		if !tableExists {
			table = &schema.Table{
				Name: col.TableName,
			}
			tables[col.TableName] = table
		}

		for _, constr := range constraints[col.TableName.Value][col.Name.Value] {
			if constr.Type == pgConstraintPKName {
				col.PrimaryKey = sql.NullString{
					String: constr.Name,
					Valid:  true,
				}
			} else if constr.Type == pgConstraintFKName {
				col.ForeignKey = &schema.ForeignKey{
					Name: schema.String{
						Value: constr.Name,
					},
					Table: schema.String{
						Value: constr.ForeignTableName,
					},
					Column: schema.String{
						Value: constr.ForeignColumnName,
					},
				}
			}
		}

		table.Columns = append(table.Columns, col)
	}

	return &schema.Schema{
		Tables: tables,
	}, nil
}

func (l *PGLoader) loadConstraints(db *sqlx.DB, ctx context.Context, schemaName string) (map[string]map[string][]*constraint, error) {
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

	constraintMap := map[string]map[string][]*constraint{}

	for _, constr := range constraints {
		_, exists := constraintMap[constr.TableName]
		if !exists {
			constraintMap[constr.TableName] = map[string][]*constraint{}
		}
		constraintMap[constr.TableName][constr.ColumnName] = append(
			constraintMap[constr.TableName][constr.ColumnName],
			constr,
		)
	}

	return constraintMap, nil
}
