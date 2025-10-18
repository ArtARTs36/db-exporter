package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/artarts36/gds"
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/jmoiron/sqlx"

	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/mysql"
)

type Loader struct {
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(ctx context.Context, cn *conn.Connection) (*schema.Schema, error) {
	db, err := cn.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}

	columns, err := l.selectColumns(ctx, cn, db)
	if err != nil {
		return nil, fmt.Errorf("select columns: %w", err)
	}

	constraintsMap, err := l.loadConstraints(ctx, cn)
	if err != nil {
		return nil, fmt.Errorf("load constraints: %w", err)
	}

	sch := schema.NewSchema(config.DatabaseDriverMySQL)

	for _, col := range columns {
		table, tableExists := sch.Tables.Get(col.TableName)
		if !tableExists {
			table = schema.NewTable(col.TableName)

			sch.Tables.Add(table)
		}

		schemaColumn := &schema.Column{
			Name:           col.Name,
			TableName:      col.TableName,
			TypeRaw:        col.DataType,
			Nullable:       col.Nullable,
			Comment:        col.Comment,
			UsingSequences: make(map[string]*schema.Sequence),
		}

		schemaColumn.Type = sqltype.MapMySQLType(col.DataType.Value)
		if col.DataType.Value == "enum" {
			schemaColumn.Enum = &schema.Enum{
				Name:          col.Name.Append("_enum"),
				Values:        mysql.ParseEnumType(col.ColumnType.Value),
				UsingInTables: []string{table.Name.Value},
				Used:          1,
			}
			table.AddEnum(schemaColumn.Enum)
		} else if col.DataTypeLength.Valid {
			schemaColumn.Type = schemaColumn.Type.WithLength(fmt.Sprintf("%d", col.DataTypeLength.Int16))
		}

		if col.AutoIncrement {
			schemaColumn.IsAutoincrement = col.AutoIncrement
		}
		if col.DefaultValue.Valid {
			schemaColumn.DefaultRaw = col.DefaultValue
		}

		table.AddColumn(schemaColumn)

		if tcs, ok := constraintsMap[table.Name.Value]; ok {
			if cs, csok := tcs[schemaColumn.Name.Value]; csok {
				l.applyConstraints(table, schemaColumn, cs)
			}
		}
	}

	return sch, nil
}

func (l *Loader) selectColumns(
	ctx context.Context,
	cn *conn.Connection,
	db *sqlx.DB,
) ([]*mysqlColumn, error) {
	query := `SELECT COLUMN_NAME as name,
       DATA_TYPE as data_type,
       TABLE_NAME as table_name,
       CHARACTER_MAXIMUM_LENGTH as data_type_length,
       IF(IS_NULLABLE = 'NO', false, true) as nullable,
       IF(EXTRA = 'auto_increment', true, false) as auto_increment,
       COLUMN_COMMENT as comment,
       COLUMN_DEFAULT as default_value,
       COLUMN_TYPE as column_type
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = ?
order by ORDINAL_POSITION
;`

	var cols []*mysqlColumn

	slog.DebugContext(ctx, "[mysql-loader] loading columns")

	err := db.SelectContext(ctx, &cols, query, cn.Database().Schema)
	if err != nil {
		return nil, err
	}

	slog.DebugContext(ctx, fmt.Sprintf("[pgloader] loaded %d columns", len(cols)))

	return cols, nil
}

type mysqlColumn struct {
	Name      gds.String `db:"name"`
	TableName gds.String `db:"table_name"`

	DataType       gds.String    `db:"data_type"`
	DataTypeLength sql.NullInt16 `db:"data_type_length"`
	ColumnType     gds.String    `db:"column_type"`

	Comment gds.String `db:"comment"`

	Nullable      bool           `db:"nullable"`
	AutoIncrement bool           `db:"auto_increment"`
	DefaultValue  sql.NullString `db:"default_value"`
}
