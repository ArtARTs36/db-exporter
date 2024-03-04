package exporter

import (
	"context"

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

func (e *DiagramExporter) Export(_ context.Context, schema *schema.Schema, params *ExportParams) ([]*ExportedPage, error) {
	c, err := buildGraphviz(e.renderer, schema)
	if err != nil {
		return nil, err
	}

	return []*ExportedPage{
		{
			FileName: "diagram.svg",
			Content:  c,
		},
	}, nil
}
