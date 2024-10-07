package exporter

import (
	"context"
	"fmt"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

type DiagramExporter struct {
	graphBuilder *graphBuilder
}

func NewDiagramExporter(renderer *template.Renderer) Exporter {
	return &DiagramExporter{
		graphBuilder: &graphBuilder{renderer: renderer},
	}
}

func (e *DiagramExporter) ExportPerFile(
	_ context.Context,
	params *ExportParams,
) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, params.Schema.Tables.Len())

	err := params.Schema.Tables.EachWithErr(func(table *schema.Table) error {
		p, err := buildDiagramPage(e.graphBuilder, schema.NewTableMap(table), fmt.Sprintf("diagram_%s.svg", table.Name.Value))
		if err != nil {
			return err
		}

		pages = append(pages, p)

		return nil
	})

	return pages, err
}

func (e *DiagramExporter) Export(_ context.Context, params *ExportParams) ([]*ExportedPage, error) {
	diagram, err := buildDiagramPage(e.graphBuilder, params.Schema.Tables, "diagram.svg")
	if err != nil {
		return nil, err
	}

	return []*ExportedPage{diagram}, nil
}

func buildDiagramPage(
	builder *graphBuilder,
	tables *schema.TableMap,
	filename string,
) (*ExportedPage, error) {
	c, err := builder.BuildSVG(tables)
	if err != nil {
		return nil, err
	}

	return &ExportedPage{
		FileName: filename,
		Content:  c,
	}, nil
}
