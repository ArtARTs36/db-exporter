package diagram

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"

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
	spec, ok := params.Spec.(*config.DiagramExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected DiagramExportSpec, got %T", params.Spec)
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	err := params.Schema.Tables.EachWithErr(func(table *schema.Table) error {
		p, err := e.buildDiagramPage(
			schema.NewTableMap(table),
			fmt.Sprintf("diagram_%s.svg", table.Name.Value),
			spec,
		)
		if err != nil {
			return err
		}

		pages = append(pages, p)

		return nil
	})

	return pages, err
}

func (e *Exporter) Export(_ context.Context, params *exporter.ExportParams) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.DiagramExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected DiagramExportSpec, got %T", params.Spec)
	}

	diagram, err := e.buildDiagramPage(params.Schema.Tables, "diagram.png", spec)
	if err != nil {
		return nil, err
	}

	return []*exporter.ExportedPage{diagram}, nil
}

func (e *Exporter) buildDiagramPage(
	tables *schema.TableMap,
	filename string,
	spec *config.DiagramExportSpec,
) (*exporter.ExportedPage, error) {
	c, err := e.graphBuilder.BuildSVG(tables, spec)
	if err != nil {
		return nil, err
	}

	return &exporter.ExportedPage{
		FileName: filename,
		Content:  c,
	}, nil
}
