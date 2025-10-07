package mysql

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/mysql"
	"github.com/artarts36/gds"
	"strings"

	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
)

type mysqlConstraint struct {
	Name      string `db:"name"`
	TableName string `db:"table_name"`

	Type        string      `db:"type"`
	ColumnNames sliceString `db:"column_names"`

	ForeignTableName  string `db:"foreign_table_name"`
	ForeignColumnName string `db:"foreign_column_name"`
}

type sliceString struct {
	Strings []string
}

func (s *sliceString) Scan(val any) error {
	switch v := val.(type) {
	case string:
		s.Strings = strings.Split(v, ",")

		return nil
	case []byte:
		s.Strings = strings.Split(string(v), ",")
		return nil
	default:
		return fmt.Errorf("unexpected type %T", val)
	}
}

func (l *Loader) loadConstraints(
	ctx context.Context,
	cn *conn.Connection,
) (map[string]map[string][]*mysqlConstraint, error) {
	const query = `
select
    kcu.CONSTRAINT_NAME as name,
    tco.CONSTRAINT_TYPE as type,
    kcu.TABLE_NAME as table_name,
    GROUP_CONCAT(kcu.COLUMN_NAME) as column_names,
    IF(kcu.REFERENCED_TABLE_NAME != '', kcu.REFERENCED_TABLE_NAME, '') as foreign_table_name,
    IF(kcu.REFERENCED_COLUMN_NAME != '', kcu.REFERENCED_COLUMN_NAME, '') as foreign_column_name
    from information_schema.KEY_COLUMN_USAGE kcu
inner join information_schema.TABLE_CONSTRAINTS tco
    on tco.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME and
       tco.TABLE_NAME = kcu.TABLE_NAME
where tco.TABLE_SCHEMA = ?     
group by
    kcu.TABLE_NAME,
    tco.CONSTRAINT_TYPE,
    kcu.CONSTRAINT_NAME,
    kcu.REFERENCED_TABLE_NAME,
    kcu.REFERENCED_COLUMN_NAME
order by
    kcu.TABLE_NAME,
    kcu.CONSTRAINT_NAME`

	var cs []*mysqlConstraint

	db, err := cn.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = db.SelectContext(ctx, &cs, query, cn.Database().Schema)
	if err != nil {
		return nil, err
	}

	constraintsMap := map[string]map[string][]*mysqlConstraint{}

	for _, constr := range cs {
		if _, ok := constraintsMap[constr.TableName]; !ok {
			constraintsMap[constr.TableName] = map[string][]*mysqlConstraint{}
		}

		for _, columnName := range constr.ColumnNames.Strings {
			if _, ok := constraintsMap[constr.TableName][columnName]; !ok {
				constraintsMap[constr.TableName][columnName] = []*mysqlConstraint{}
			}

			constraintsMap[constr.TableName][columnName] = append(constraintsMap[constr.TableName][columnName], constr)
		}
	}

	return constraintsMap, nil
}

func (l *Loader) applyConstraints(table *schema.Table, col *schema.Column, constraints []*mysqlConstraint) {
	for _, constr := range constraints {
		switch constr.Type {
		case mysql.ConstraintPKName:
			pk := table.PrimaryKey
			if pk == nil {
				pk = &schema.PrimaryKey{
					Name: gds.String{
						Value: constr.Name,
					},
					ColumnsNames: gds.NewStrings(constr.ColumnNames.Strings...),
				}

				table.PrimaryKey = pk
			}

			col.PrimaryKey = pk
		case mysql.ConstraintFKName:
			fk := table.ForeignKeys[constr.Name]

			if fk == nil {
				fk = &schema.ForeignKey{
					Name: gds.String{
						Value: constr.Name,
					},
					Table:        table.Name,
					ColumnsNames: gds.NewStrings(constr.ColumnNames.Strings...),
					ForeignTable: gds.String{
						Value: constr.ForeignTableName,
					},
					ForeignColumn: gds.String{
						Value: constr.ForeignColumnName,
					},
				}

				table.ForeignKeys[constr.Name] = fk
			}

			col.ForeignKey = fk
		case mysql.ConstraintUniqueName:
			uk := table.UniqueKeys[constr.Name]

			if uk == nil {
				uk = &schema.UniqueKey{
					Name: gds.String{
						Value: constr.Name,
					},
					ColumnsNames: gds.NewStrings(constr.ColumnNames.Strings...),
				}

				table.UniqueKeys[constr.Name] = uk
			}

			col.UniqueKey = uk
		}
	}
}
