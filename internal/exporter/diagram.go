package exporter

import (
	"context"
	"fmt"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

type DiagramExporter struct {
	renderer *template.Renderer
}

func NewDiagramExporter(renderer *template.Renderer) Exporter {
	return &DiagramExporter{
		renderer: renderer,
	}
}

func (e *DiagramExporter) ExportPerFile(
	_ context.Context,
	sch *schema.Schema,
	_ *ExportParams,
) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, len(sch.Tables))
	for _, table := range sch.Tables {
		p, err := buildDiagramPage(e.renderer, map[schema.String]*schema.Table{
			table.Name: table,
		}, fmt.Sprintf("diagram_%s.svg", table.Name.Value))
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (e *DiagramExporter) Export(_ context.Context, sch *schema.Schema, _ *ExportParams) ([]*ExportedPage, error) {
	diagram, err := buildDiagramPage(e.renderer, sch.Tables, "diagram.svg")
	if err != nil {
		return nil, err
	}

	return []*ExportedPage{diagram}, nil
}

func buildDiagramPage(
	renderer *template.Renderer,
	tables map[schema.String]*schema.Table,
	filename string,
) (*ExportedPage, error) {
	c, err := buildGraphviz(renderer, tables)
	if err != nil {
		return nil, err
	}

	return &ExportedPage{
		FileName: filename,
		Content:  c,
	}, nil
}
