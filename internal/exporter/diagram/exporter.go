package diagram

import (
	"context"
	"fmt"

	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

type Exporter struct {
	graphBuilder *GraphBuilder
}

func NewDiagramExporter(renderer *template.Renderer) exporter.Exporter {
	return &Exporter{
		graphBuilder: &GraphBuilder{renderer: renderer},
	}
}

func (e *Exporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

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

func (e *Exporter) Export(_ context.Context, params *exporter.ExportParams) ([]*exporter.ExportedPage, error) {
	diagram, err := buildDiagramPage(e.graphBuilder, params.Schema.Tables, "diagram.svg")
	if err != nil {
		return nil, err
	}

	return []*exporter.ExportedPage{diagram}, nil
}

func buildDiagramPage(
	builder *GraphBuilder,
	tables *schema.TableMap,
	filename string,
) (*exporter.ExportedPage, error) {
	c, err := builder.BuildSVG(tables)
	if err != nil {
		return nil, err
	}

	return &exporter.ExportedPage{
		FileName: filename,
		Content:  c,
	}, nil
}
