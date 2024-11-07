package dbml

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
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
	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	for _, tbl := range params.Schema.Tables.List() {
		dbmlFile := &dbml.File{}
		table, refs := e.mapTable(ctx, tbl)
		dbmlFile.Tables = []*dbml.Table{table}
		dbmlFile.Refs = refs

		pages = append(pages, &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s.dbml", tbl.Name.Value),
			Content:  []byte(dbmlFile.Render()),
		})
	}

	return pages, nil
}

func (e *Exporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	dbmlFile := &dbml.File{
		Tables: make([]*dbml.Table, 0, params.Schema.Tables.Len()),
		Refs:   make([]*dbml.Ref, 0),
	}

	for _, tbl := range params.Schema.Tables.List() {
		table, refs := e.mapTable(ctx, tbl)
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

func (e *Exporter) mapTable(ctx context.Context, tbl *schema.Table) (*dbml.Table, []*dbml.Ref) {
	table := &dbml.Table{
		Name:    tbl.Name.Value,
		Columns: make([]*dbml.Column, 0, len(tbl.Columns)),
		Note:    tbl.Comment,
	}

	for _, col := range tbl.Columns {
		column := &dbml.Column{
			Name: col.Name.Value,
			Type: col.Type.Value,
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
			From: tbl.Name.Append(".").Append(fk.ColumnsNames.Join(",").Value).Value,
			Type: ">",
			To:   fk.Table.Append(".").Append(fk.Table.Value).Value,
		})
	}

	return table, refs
}

func (e *Exporter) mapDefault(col *schema.Column) (dbml.ColumnDefault, error) {
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
