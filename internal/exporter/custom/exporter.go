package custom

import (
	"context"
	"fmt"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

type Exporter struct {
	renderer *template.Renderer
	pager    *common.Pager
}

func NewExporter(
	renderer *template.Renderer,
	pager *common.Pager,
) *Exporter {
	return &Exporter{renderer: renderer, pager: pager}
}

func (e *Exporter) Export(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.CustomExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected CustomExportSpec, got %T", params.Spec)
	}

	pager := e.pager.Of(spec.Template)

	p, err := pager.Export(e.filenameCreator(spec)("result"), map[string]stick.Value{
		"schema": &exportingSchema{
			Tables: params.Schema.Tables.List(),
		},
	})
	if err != nil {
		return nil, err
	}

	return []*exporter.ExportedPage{p}, nil
}

func (e *Exporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.CustomExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected CustomExportSpec, got %T", params.Spec)
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	pager := e.pager.Of(spec.Template)

	createFilename := e.filenameCreator(spec)

	err := params.Schema.Tables.EachWithErr(func(table *schema.Table) error {
		p, err := pager.Export(createFilename(table.Name.Value), map[string]stick.Value{
			"table": table,
		})
		if err != nil {
			return err
		}

		pages = append(pages, p)

		return nil
	})

	return pages, err
}

func (e *Exporter) filenameCreator(spec *config.CustomExportSpec) func(tableName string) string {
	if spec.Output.Extension == "" {
		return func(tableName string) string {
			return tableName
		}
	}

	return func(tableName string) string {
		return fmt.Sprintf("%s.%s", tableName, spec.Output.Extension)
	}
}
