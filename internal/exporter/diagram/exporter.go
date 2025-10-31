package diagram

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
)

type Exporter struct {
	creator *Creator
}

func NewDiagramExporter(creator *Creator) exporter.Exporter {
	return &Exporter{
		creator: creator,
	}
}

func (e *Exporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*Specification)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected Specification, got %T", params.Spec)
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	err := params.Schema.Tables.EachWithErr(func(table *schema.Table) error {
		p, err := e.buildDiagramPage(
			ctx,
			schema.NewTableMap(table),
			fmt.Sprintf("diagram_%s.%s", table.Name.Value, string(spec.Image.Format)),
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

func (e *Exporter) Export(ctx context.Context, params *exporter.ExportParams) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*Specification)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected Specification, got %T", params.Spec)
	}

	diagram, err := e.buildDiagramPage(
		ctx,
		params.Schema.Tables,
		fmt.Sprintf("diagram.%s", spec.Image.Format),
		spec,
	)
	if err != nil {
		return nil, err
	}

	return []*exporter.ExportedPage{diagram}, nil
}

func (e *Exporter) buildDiagramPage(
	ctx context.Context,
	tables *schema.TableMap,
	filename string,
	spec *Specification,
) (*exporter.ExportedPage, error) {
	c, err := e.creator.Create(ctx, tables, spec)
	if err != nil {
		return nil, err
	}

	return &exporter.ExportedPage{
		FileName: filename,
		Content:  c,
	}, nil
}
