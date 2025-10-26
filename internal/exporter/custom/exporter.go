package custom

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/infrastructure/workspace"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/config"
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
	spec, ok := params.Spec.(*config.CustomExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected CustomExportSpec, got %T", params.Spec)
	}

	err := params.Workspace.Write(ctx, e.filenameCreator(spec)("result"), func(buffer workspace.Buffer) error {
		return e.renderer.RenderTo(spec.Template, map[string]stick.Value{
			"schema": &exportingSchema{
				Tables: params.Schema.Tables.List(),
			},
		}, buffer)
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
	spec, ok := params.Spec.(*config.CustomExportSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec, expected CustomExportSpec, got %T", params.Spec)
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	createFilename := e.filenameCreator(spec)

	err := params.Schema.Tables.EachWithErr(func(table *schema.Table) error {
		err := params.Workspace.Write(ctx, createFilename(table.Name.Value), func(buffer workspace.Buffer) error {
			return e.renderer.RenderTo(spec.Template, map[string]stick.Value{
				"table": table,
			}, buffer)
		})
		if err != nil {
			return fmt.Errorf("write table to workspace: %w", err)
		}

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
