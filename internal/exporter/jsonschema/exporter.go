package jsonschema

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
)

type Exporter struct {
	builder *builder
}

func NewExporter() *Exporter {
	return &Exporter{
		builder: &builder{},
	}
}

func (e *Exporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.JSONSchemaExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	for _, table := range params.Schema.Tables.List() {
		content, err := e.buildJSONSchema(spec, []*schema.Table{table})
		if err != nil {
			return nil, fmt.Errorf("failed to build json schema for table %q: %w", table.Name, err)
		}

		pages = append(pages, &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s.json", table.Name.Lower()),
			Content:  content,
		})
	}

	return pages, nil
}

func (e *Exporter) Export(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.JSONSchemaExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	content, err := e.buildJSONSchema(spec, params.Schema.Tables.List())
	if err != nil {
		return nil, fmt.Errorf("failed to build json schema: %w", err)
	}

	return []*exporter.ExportedPage{
		{
			FileName: "schema.json",
			Content:  content,
		},
	}, nil
}

func (e *Exporter) buildJSONSchema(spec *config.JSONSchemaExportSpec, tables []*schema.Table) ([]byte, error) {
	sch := e.builder.buildJSONSchema(spec, tables)

	marshaller := sch.Marshal
	if spec.Pretty {
		marshaller = sch.MarshallPretty
	}

	content, err := marshaller()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	return content, nil
}
