package dbml

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/dbml"
)

type Exporter struct {
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	for _, tbl := range params.Schema.Tables.List() {
		dbmlFile := &dbml.File{}
		table, refs := e.mapTable(tbl)
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
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	dbmlFile := &dbml.File{
		Tables: make([]*dbml.Table, 0, params.Schema.Tables.Len()),
	}

	for _, tbl := range params.Schema.Tables.List() {
		table, refs := e.mapTable(tbl)
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

func (e *Exporter) mapTable(tbl *schema.Table) (*dbml.Table, []*dbml.Ref) {
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
