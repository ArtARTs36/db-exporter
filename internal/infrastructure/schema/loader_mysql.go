package schema

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/gds"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log/slog"

	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/schema"
)

type MySQLLoader struct {
}

func NewMySQLLoader() *MySQLLoader {
	return &MySQLLoader{}
}

func (l *MySQLLoader) Load(ctx context.Context, cn *conn.Connection) (*schema.Schema, error) {
	db, err := cn.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("connect to db: %w", err)
	}

	columns, err := l.selectColumns(ctx, cn, db)
	if err != nil {
		return nil, fmt.Errorf("select columns: %w", err)
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

		if col.IsPrimaryKey {
			schemaColumn.PrimaryKey = schema.CreatePrimaryKeyForColumn(schemaColumn) // @todo get real key name
			table.PrimaryKey = schemaColumn.PrimaryKey
		}
		if col.IsUniqueKey {
			schemaColumn.UniqueKey = schema.CreateUniqueKeyForColumn(schemaColumn)
			table.UniqueKeys[schemaColumn.UniqueKey.Name.Value] = schemaColumn.UniqueKey
		}

		schemaColumn.Type = sqltype.MapMySQLType(col.DataType.Value)
		if col.DataTypeLength.Valid {
			schemaColumn.Type = schemaColumn.Type.WithLength(fmt.Sprintf("%d", col.DataTypeLength.Int16))
		}
		if col.AutoIncrement {
			schemaColumn.IsAutoincrement = col.AutoIncrement
		}
		if col.DefaultValue.Valid {
			schemaColumn.DefaultRaw = col.DefaultValue
		}

		table.AddColumn(schemaColumn)
	}

	return sch, nil
}

func (l *MySQLLoader) selectColumns(
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
       IF(COLUMN_KEY = 'PRI', true, false) as is_primary_key,
       IF(COLUMN_KEY = 'UNI', true, false) as is_unique_key,
       COLUMN_DEFAULT as default_value
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

	Comment gds.String `db:"comment"`

	Nullable      bool           `db:"nullable"`
	AutoIncrement bool           `db:"auto_increment"`
	IsPrimaryKey  bool           `db:"is_primary_key"`
	IsUniqueKey   bool           `db:"is_unique_key"`
	DefaultValue  sql.NullString `db:"default_value"`
}
