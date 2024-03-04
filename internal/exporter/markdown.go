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
	var diagram *ExportedPage

	if params.WithDiagram {
		diagramContent, err := buildGraphviz(e.renderer, schema)
		if err != nil {
			return nil, fmt.Errorf("failed to build diagram: %w", err)
		}

		diagram = &ExportedPage{
			FileName: "diagram.svg",
			Content:  diagramContent,
		}
	}

	if params.TablePerFile {
		return e.exportPerFile(schema, diagram)
	}

	return e.exportToSingleFile(schema, diagram)
}

func (e *MarkdownExporter) exportPerFile(sc *schema.Schema, diagram *ExportedPage) ([]*ExportedPage, error) {
	pagesCap := len(sc.Tables) + 1
	if diagram != nil {
		pagesCap += 1
	}

	pages := make([]*ExportedPage, 0, pagesCap)
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
		"tables":  preparedTables,
		"diagram": diagram,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to export index file: %w", err)
	}

	pages = append(pages, indexPage)

	if diagram != nil {
		pages = append(pages, diagram)
	}

	return pages, nil
}

func (e *MarkdownExporter) exportToSingleFile(sc *schema.Schema, diagram *ExportedPage) ([]*ExportedPage, error) {
	page, err := render(e.renderer, "markdown/single-tables.md", "tables.md", map[string]stick.Value{
		"schema":        sc,
		"diagram":       diagram,
		"diagramExists": diagram != nil,
	})
	if err != nil {
		return nil, err
	}

	pages := []*ExportedPage{
		page,
	}

	if diagram != nil {
		pages = append(pages, diagram)
	}

	return pages, nil
}
