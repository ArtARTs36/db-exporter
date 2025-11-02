package custom

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/infrastructure/workspace"
	"github.com/artarts36/db-exporter/internal/shared/iox"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

type Exporter struct {
	renderer *template.Renderer
}

func NewExporter(
	renderer *template.Renderer,
) *Exporter {
	return &Exporter{renderer: renderer}
}

func (e *Exporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*Specification)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected Specification, got %T", params.Spec)
	}

	err := params.Workspace.Write(ctx, &workspace.WritingFile{
		Filename: e.filenameCreator(spec)("result"),
		Writer: func(buffer iox.Writer) error {
			return e.renderer.RenderTo(spec.Template, map[string]stick.Value{
				"schema": &exportingSchema{
					Tables: params.Schema.Tables.List(),
				},
			}, buffer)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("write table to workspace: %w", err)
	}

	return nil, nil
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

	createFilename := e.filenameCreator(spec)

	err := params.Schema.Tables.EachWithErr(func(table *schema.Table) error {
		err := params.Workspace.Write(ctx, &workspace.WritingFile{
			Filename: createFilename(table.Name.Value),
			Writer: func(buffer iox.Writer) error {
				return e.renderer.RenderTo(spec.Template, map[string]stick.Value{
					"table": table,
				}, buffer)
			},
		})
		if err != nil {
			return fmt.Errorf("write table to workspace: %w", err)
		}

		return nil
	})

	return pages, err
}

func (e *Exporter) filenameCreator(spec *Specification) func(tableName string) string {
	if spec.Output.Extension == "" {
		return func(tableName string) string {
			return tableName
		}
	}

	return func(tableName string) string {
		return fmt.Sprintf("%s.%s", tableName, spec.Output.Extension)
	}
}
