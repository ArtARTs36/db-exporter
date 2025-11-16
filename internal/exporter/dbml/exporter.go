package dbml

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/shared/dbml"
)

type Exporter struct {
	mapper *mapper
}

func NewExporter() *Exporter {
	return &Exporter{
		mapper: &mapper{},
	}
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
		if tbl.IsPartition() {
			continue
		}

		dbmlFile := &dbml.File{}
		table, refs, err := e.mapper.mapTable(ctx, tbl, params.Schema.Driver)
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

	enumFile := &dbml.File{Enums: e.mapper.mapEnums(params.Schema)}
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
		Enums:  e.mapper.mapEnums(params.Schema),
	}

	for _, tbl := range params.Schema.Tables.List() {
		if tbl.IsPartition() {
			continue
		}

		table, refs, err := e.mapper.mapTable(ctx, tbl, params.Schema.Driver)
		if err != nil {
			return nil, fmt.Errorf("failed to map table %q: %w", tbl.Name, err)
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
