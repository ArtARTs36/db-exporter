package task

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/schema"
	"log/slog"

	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"github.com/artarts36/db-exporter/internal/shared/migrations"
	"github.com/artarts36/db-exporter/internal/template"
)

type ExportActivityRunner struct {
	migrationsTblDetector *migrations.TableDetector
	pageStorage           *pageStorage
	fs                    fs.Driver
	renderer              *template.Renderer
	exporters             map[config.ExporterName]exporter.Exporter
}

func NewExportActivityRunner(
	fs fs.Driver,
	renderer *template.Renderer,
	exporters map[config.ExporterName]exporter.Exporter,
) *ExportActivityRunner {
	return &ExportActivityRunner{
		migrationsTblDetector: migrations.NewTableDetector(),
		pageStorage:           &pageStorage{fs},
		fs:                    fs,
		renderer:              renderer,
		exporters:             exporters,
	}
}

func (r *ExportActivityRunner) Run(ctx context.Context, expParams *ActivityRunParams) (*ActivityResult, error) {
	pages, err := r.export(ctx, expParams)
	if err != nil {
		return nil, err
	}

	generatedFiles, err := r.pageStorage.Save(ctx, pages, &savePageParams{
		Dir:        expParams.Activity.Export.Out.Dir,
		FilePrefix: expParams.Activity.Export.Out.FilePrefix,
		SkipExists: expParams.Activity.Export.SkipExists,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save generated pages: %w", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("[exportcmd] successful generated %d files", len(pages)))

	return &ActivityResult{
		Export: &ExportActivityResult{
			files: generatedFiles,
		},
	}, nil
}

func (r *ExportActivityRunner) export(
	ctx context.Context,
	params *ActivityRunParams,
) ([]*exporter.ExportedPage, error) {
	exp, exists := r.exporters[params.Activity.Export.Format]
	if !exists {
		return nil, fmt.Errorf("exporter for format %q not found", params.Activity.Export.Format)
	}

	sc := r.filterTables(params.Schema, params)

	exporterParams := &exporter.ExportParams{
		Schema:    sc,
		Spec:      params.Activity.Export.Spec,
		Conn:      params.Conn,
		Directory: fs.NewDirectory(r.fs, params.Activity.Export.Out.Dir),
	}

	export := func() ([]*exporter.ExportedPage, error) {
		if params.Activity.Export.TablePerFile {
			return exp.ExportPerFile(ctx, exporterParams)
		}

		return exp.Export(ctx, exporterParams)
	}

	pages, err := export()
	if err != nil {
		return nil, fmt.Errorf("exporter %q unable to export: %w", params.Activity.Export.Format, err)
	}

	return pages, nil
}

func (r *ExportActivityRunner) filterTables(sc *schema.Schema, params *ActivityRunParams) *schema.Schema {
	if params.Activity.Tables.List.Len() > 0 {
		sc = sc.OnlyTables(params.Activity.Tables.List.List())
	} else if params.Activity.Tables.Prefix != "" {
		sc = sc.WithoutTable(func(table *schema.Table) bool {
			return !table.Name.Starts(params.Activity.Tables.Prefix)
		})
	}

	return sc
}
