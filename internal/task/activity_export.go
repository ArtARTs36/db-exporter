package task

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter"
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

	generatedFiles, err := r.pageStorage.Save(pages, &savePageParams{
		Dir:        expParams.Activity.Out.Dir,
		FilePrefix: expParams.Activity.Out.FilePrefix,
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
	var pages []*exporter.ExportedPage
	var err error

	exp, exists := r.exporters[params.Activity.Export]
	if !exists {
		return nil, fmt.Errorf("exporter for format %q not found", params.Activity.Export)
	}

	sc := params.Schema
	if len(params.Activity.Tables) > 0 {
		sc = sc.Clone()
		sc.Tables = sc.Tables.Only(params.Activity.Tables)
	}

	exporterParams := &exporter.ExportParams{
		Schema: sc,
		Spec:   params.Activity.Spec,
	}

	if params.Activity.TablePerFile {
		pages, err = exp.ExportPerFile(ctx, exporterParams)
	} else {
		pages, err = exp.Export(ctx, exporterParams)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to doImport: %w", err)
	}

	return pages, nil
}
