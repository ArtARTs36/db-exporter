package dbml

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/dbml"
	"log/slog"
)

type Exporter struct {
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	pagesLen := params.Schema.Tables.Len()
	if len(params.Schema.Enums) > 0 {
		pagesLen++
	}

	pages := make([]*exporter.ExportedPage, 0, pagesLen)

	for _, tbl := range params.Schema.Tables.List() {
		dbmlFile := &dbml.File{}
		table, refs, err := e.mapTable(ctx, tbl, params.Schema.Driver)
		if err != nil {
			return nil, fmt.Errorf("failed to map table %q: w", err)
		}
		dbmlFile.Tables = []*dbml.Table{table}
		dbmlFile.Refs = refs

		pages = append(pages, &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s.dbml", tbl.Name.Value),
			Content:  []byte(dbmlFile.Render()),
		})
	}

	enumFile := &dbml.File{Enums: e.mapEnums(params.Schema)}
	pages = append(pages, &exporter.ExportedPage{
		FileName: "enums.dbml",
		Content:  []byte(enumFile.Render()),
	})

	return pages, nil
}

func (e *Exporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	dbmlFile := &dbml.File{
		Tables: make([]*dbml.Table, 0, params.Schema.Tables.Len()),
		Refs:   make([]*dbml.Ref, 0),
		Enums:  e.mapEnums(params.Schema),
	}

	for _, tbl := range params.Schema.Tables.List() {
		table, refs, err := e.mapTable(ctx, tbl, params.Schema.Driver)
		if err != nil {
			return nil, fmt.Errorf("failed to map table %q: %w", table.Name, err)
		}
		dbmlFile.Tables = append(dbmlFile.Tables, table)
		dbmlFile.Refs = append(dbmlFile.Refs, refs...)
	}

	return []*exporter.ExportedPage{
		{
			FileName: "schema.dbml",
			Content:  []byte(dbmlFile.Render()),
		},
	}, nil
}

func (e *Exporter) mapTable(
	ctx context.Context,
	tbl *schema.Table,
	source config.DatabaseDriver,
) (*dbml.Table, []*dbml.Ref, error) {
	table := &dbml.Table{
		Name:    tbl.Name.Value,
		Columns: make([]*dbml.Column, 0, len(tbl.Columns)),
		Note:    tbl.Comment,
	}

	for _, col := range tbl.Columns {
		typ, err := sqltype.TransitSQLType(source, config.DatabaseDriverDBML, col.Type)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to map column %q type: %w", col.Name, err)
		}

		column := &dbml.Column{
			Name: col.Name.Value,
			Type: typ.Name,
			Settings: dbml.ColumnSettings{
				PrimaryKey: col.IsPrimaryKey(),
				Increment:  col.IsAutoincrement,
				Note:       col.Comment.Value,
				Unique:     col.IsUniqueKey(),
			},
		}

		if col.Nullable {
			column.AsNullable()
		}

		def, err := e.mapDefault(col)
		if err != nil {
			slog.
				With(slog.String("table_name", table.Name)).
				With(slog.String("column_name", column.Name)).
				WarnContext(ctx, "[dbml-exporter] failed to map default value of column")

			column.Settings.Default.Value = col.DefaultRaw.String
		} else {
			column.Settings.Default = def
		}

		table.Columns = append(table.Columns, column)
	}

	refs := make([]*dbml.Ref, 0, len(tbl.ForeignKeys))

	for _, fk := range tbl.ForeignKeys {
		refs = append(refs, &dbml.Ref{
			From: fk.Table.Append(".").Append(fk.ColumnsNames.Join(",").Value).Value,
			Type: ">",
			To:   fk.ForeignTable.Append(".").Append(fk.ForeignColumn.Value).Value,
		})
	}

	return table, refs, nil
}

func (e *Exporter) mapDefault(col *schema.Column) (dbml.ColumnDefault, error) {
	if col.Default == nil {
		return dbml.ColumnDefault{}, nil
	}

	if col.Default.Type == schema.ColumnDefaultTypeValue {
		switch v := col.Default.Value.(type) {
		case bool:
			boolVal := "false"
			if v {
				boolVal = "true"
			}

			return dbml.ColumnDefault{
				Value: boolVal,
				Type:  dbml.ColumnDefaultTypeBoolean,
			}, nil
		case string:
			return dbml.ColumnDefault{
				Value: v,
				Type:  dbml.ColumnDefaultTypeString,
			}, nil
		case int:
			return dbml.ColumnDefault{
				Value: fmt.Sprintf("%d", v),
				Type:  dbml.ColumnDefaultTypeNumber,
			}, nil
		default:
			return dbml.ColumnDefault{}, fmt.Errorf("value of %T unsupported", col.Default.Value)
		}
	}

	if col.Default.Type == schema.ColumnDefaultTypeFunc {
		return dbml.ColumnDefault{
			Type:  dbml.ColumnDefaultTypeExpression,
			Value: col.Default.Value.(string),
		}, nil
	}

	return dbml.ColumnDefault{}, nil
}

func (e *Exporter) mapEnums(sch *schema.Schema) []*dbml.Enum {
	enums := make([]*dbml.Enum, 0, len(sch.Enums))

	for _, enum := range sch.Enums {
		enums = append(enums, e.mapEnum(enum))
	}

	return enums
}

func (e *Exporter) mapEnum(en *schema.Enum) *dbml.Enum {
	enum := &dbml.Enum{
		Name:   en.Name.Value,
		Values: make([]dbml.EnumValue, 0, len(en.Values)),
	}

	for _, val := range en.Values {
		enum.Values = append(enum.Values, dbml.EnumValue{
			Name: val,
		})
	}

	return enum
}
