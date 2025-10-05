package task

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"github.com/artarts36/db-exporter/internal/shared/migrations"
	"log/slog"
	"strings"
)

type ImportActivityRunner struct {
	migrationsTblDetector *migrations.TableDetector
	fs                    fs.Driver
	importers             map[config.ImporterName]exporter.Importer
}

func NewImportActivityRunner(fs fs.Driver, importers map[config.ImporterName]exporter.Importer) *ImportActivityRunner {
	return &ImportActivityRunner{
		migrationsTblDetector: migrations.NewTableDetector(),
		fs:                    fs,
		importers:             importers,
	}
}

func (a *ImportActivityRunner) Run(ctx context.Context, params *ActivityRunParams) (*ActivityResult, error) {
	importer, exists := a.importers[params.Activity.Import.Format]
	if !exists {
		return nil, fmt.Errorf("importer for format %q not found", params.Activity.Import.Format)
	}

	files, err := a.doImport(ctx, importer, params)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		slog.InfoContext(ctx, "[importcmd] no files to import")

		return &ActivityResult{
			Import: &ImportActivityResult{
				files:            make([]exporter.ImportedFile, 0),
				tableRowCountMap: map[string]int64{},
			},
		}, nil
	}

	filesPaths := strings.Builder{}
	countsMap := map[string]int64{}

	for _, file := range files {
		if filesPaths.Len() > 0 {
			filesPaths.WriteRune(',')
		}

		filesPaths.WriteString(file.Name)

		for table, ar := range file.AffectedRows {
			countsMap[table] += ar
		}
	}

	slog.InfoContext(
		ctx,
		fmt.Sprintf("[importcmd] successfully imported from %d files: %s", len(files), filesPaths.String()),
	)

	return &ActivityResult{
		Import: &ImportActivityResult{
			files:            files,
			tableRowCountMap: countsMap,
		},
	}, nil
}

func (a *ImportActivityRunner) doImport(
	ctx context.Context,
	exp exporter.Importer,
	params *ActivityRunParams,
) ([]exporter.ImportedFile, error) {
	var pages []exporter.ImportedFile
	var err error
	importerParams := &exporter.ImportParams{
		Directory: fs.NewDirectory(a.fs, params.Activity.Import.From),
		TableFilter: func(tableName string) bool {
			if params.Activity.Tables.List.Len() > 0 && !params.Activity.Tables.List.Has(tableName) {
				return false
			}

			if params.Activity.Tables.Prefix != "" && !strings.HasPrefix(tableName, params.Activity.Tables.Prefix) {
				return false
			}

			return true
		},
		Conn:   params.Conn,
		Schema: params.Schema,
	}

	pages, err = exp.Import(ctx, importerParams)
	if err != nil {
		return nil, fmt.Errorf("failed to import: %w", err)
	}

	return pages, nil
}
