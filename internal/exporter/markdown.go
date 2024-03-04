package exporter

import (
	"context"
	"fmt"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

type MarkdownExporter struct {
	renderer *template.Renderer
}

type markdownPreparedTable struct {
	*schema.Table
	FileName string
}

func NewMarkdownExporter(renderer *template.Renderer) Exporter {
	return &MarkdownExporter{
		renderer: renderer,
	}
}

func (e *MarkdownExporter) Export(_ context.Context, schema *schema.Schema, params *ExportParams) ([]*ExportedPage, error) {
	if params.TablePerFile {
		return e.exportPerFile(schema)
	}

	return e.exportToSingleFile(schema)
}

func (e *MarkdownExporter) exportPerFile(sc *schema.Schema) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, len(sc.Tables)+1)
	preparedTables := make([]*markdownPreparedTable, 0, len(sc.Tables))

	for _, table := range sc.Tables {
		fileName := fmt.Sprintf("table_%s.md", table.Name.Value)

		page, err := render(e.renderer, "markdown/per-table.md", fileName, map[string]stick.Value{
			"table": table,
		})
		if err != nil {
			return nil, err
		}

		pages = append(pages, page)

		preparedTables = append(preparedTables, &markdownPreparedTable{
			Table:    table,
			FileName: fileName,
		})
	}

	indexPage, err := render(e.renderer, "markdown/per-index.md", "index.md", map[string]stick.Value{
		"tables": preparedTables,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to export index file: %w", err)
	}

	pages = append(pages, indexPage)

	return pages, nil
}

func (e *MarkdownExporter) exportToSingleFile(sc *schema.Schema) ([]*ExportedPage, error) {
	page, err := render(e.renderer, "markdown/single-tables.md", "tables.md", map[string]stick.Value{
		"schema": sc,
	})
	if err != nil {
		return nil, err
	}

	return []*ExportedPage{
		page,
	}, nil
}
