package exporter

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/template"
)

type MarkdownExporter struct {
	unimplementedImporter
	renderer     *template.Renderer
	graphBuilder *graphBuilder
}

type markdownPreparedTable struct {
	*schema.Table
	FileName string
}

func NewMarkdownExporter(renderer *template.Renderer) Exporter {
	return &MarkdownExporter{
		renderer:     renderer,
		graphBuilder: &graphBuilder{renderer: renderer},
	}
}

func (e *MarkdownExporter) ExportPerFile(
	_ context.Context,
	params *ExportParams,
) ([]*ExportedPage, error) {
	spec := params.Spec.(*config.MarkdownExportSpec)

	var diagram *ExportedPage
	pagesCap := params.Schema.Tables.Len() + 1
	if spec.WithDiagram {
		pagesCap++
		var err error
		diagram, err = buildDiagramPage(e.graphBuilder, params.Schema.Tables, "diagram.svg")
		if err != nil {
			return nil, fmt.Errorf("failed to build diagram: %w", err)
		}
	}

	pages := make([]*ExportedPage, 0, pagesCap)
	preparedTables := make([]*markdownPreparedTable, 0, params.Schema.Tables.Len())

	for _, table := range params.Schema.Tables.List() {
		fileName := fmt.Sprintf("%s.md", table.Name.Value)

		page, err := render(e.renderer, "md/per-table.md", fileName, map[string]stick.Value{
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

	indexPage, err := render(e.renderer, "md/per-index.md", "index.md", map[string]stick.Value{
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

func (e *MarkdownExporter) Export(
	_ context.Context,
	params *ExportParams,
) ([]*ExportedPage, error) {
	var diagram *ExportedPage

	spec := params.Spec.(*config.MarkdownExportSpec)

	if spec.WithDiagram {
		var err error

		diagram, err = buildDiagramPage(e.graphBuilder, params.Schema.Tables, "diagram.svg")
		if err != nil {
			return nil, fmt.Errorf("failed to build diagram: %w", err)
		}
	}

	page, err := render(
		e.renderer,
		"md/single-tables.md",
		e.createIndexPageName(params.Schema),
		map[string]stick.Value{
			"schema":        params.Schema,
			"diagram":       diagram,
			"diagramExists": diagram != nil,
		},
	)
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

func (e *MarkdownExporter) createIndexPageName(sch *schema.Schema) string {
	if sch.Tables.Has(ds.String{Value: "INDEX"}) {
		return "index.md"
	}

	return "INDEX.md"
}
