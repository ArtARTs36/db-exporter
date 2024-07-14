package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/app/actions"
	"github.com/artarts36/db-exporter/internal/app/params"
	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"github.com/artarts36/db-exporter/internal/shared/migrations"
	"github.com/artarts36/db-exporter/internal/template"
	"github.com/artarts36/db-exporter/templates"
)

const localTemplatesFolder = "./db-exporter-templates"

type ExportCmd struct {
	migrationsTblDetector *migrations.TableDetector
	pageStorage           *pageStorage
	fs                    fs.Driver
	actions               map[string]actions.Action
}

func NewExportCmd(fs fs.Driver, actions map[string]actions.Action) *ExportCmd {
	return &ExportCmd{
		migrationsTblDetector: migrations.NewTableDetector(),
		pageStorage:           &pageStorage{fs},
		fs:                    fs,
		actions:               actions,
	}
}

func (a *ExportCmd) Export(ctx context.Context, expParams *params.ExportParams) error {
	startedAt := time.Now()

	driverName, err := db.CreateDriverName(expParams.DriverName)
	if err != nil {
		return err
	}

	connection := db.NewConnection(driverName, expParams.DSN)
	defer func() {
		closeErr := connection.Close()

		if closeErr == nil {
			slog.InfoContext(ctx, "[exportcmd] db connection closed")

			return
		}

		slog.ErrorContext(ctx, fmt.Sprintf("failed to close db connection: %s", closeErr))
	}()

	loader, err := db.CreateSchemaLoader(connection)
	if err != nil {
		return fmt.Errorf("unable to create schema loader: %w", err)
	}

	renderer := a.createRenderer()

	exp, err := exporter.CreateExporter(expParams.Format, renderer, connection)
	if err != nil {
		return fmt.Errorf("failed to create exporter: %w", err)
	}

	// processing

	slog.DebugContext(ctx, fmt.Sprintf("[exportcmd] loading db schema from %s", expParams.DriverName))

	sc, err := a.loadSchema(ctx, loader, expParams)
	if err != nil {
		return fmt.Errorf("unable to load schema: %w", err)
	}

	slog.DebugContext(
		ctx,
		fmt.Sprintf(
			"[exportcmd] loaded %d tables: [%s]", sc.Tables.Len(), strings.Join(sc.TablesNames(), ","),
		),
	)

	pages, err := a.export(ctx, exp, sc, expParams)
	if err != nil {
		return err
	}

	generatedFiles, err := a.pageStorage.Save(pages, &savePageParams{
		Dir:        expParams.OutDir,
		FilePrefix: expParams.FilePrefix,
	})
	if err != nil {
		return fmt.Errorf("failed to save generated pages: %w", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("[exportcmd] successful generated %d files", len(pages)))

	err = a.runActions(ctx, &params.ActionParams{
		StartedAt:      startedAt,
		ExportParams:   expParams,
		GeneratedFiles: generatedFiles,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *ExportCmd) export(
	ctx context.Context,
	exp exporter.Exporter,
	sc *schema.Schema,
	params *params.ExportParams,
) ([]*exporter.ExportedPage, error) {
	var pages []*exporter.ExportedPage
	var err error
	exporterParams := &exporter.ExportParams{
		WithDiagram:            params.WithDiagram,
		WithoutMigrationsTable: params.WithoutMigrationsTable,
		Package:                params.Package,
		Directory:              fs.NewDirectory(a.fs, params.OutDir),
	}

	if params.Import {
		if params.TablePerFile {
			err = exp.ImportPerFile(ctx, sc, exporterParams)
		} else {
			err = exp.Import(ctx, sc, exporterParams)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to import: %w", err)
		}

		return pages, nil
	} else {
		if params.TablePerFile {
			pages, err = exp.ExportPerFile(ctx, sc, exporterParams)
		} else {
			pages, err = exp.Export(ctx, sc, exporterParams)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export: %w", err)
	}

	return pages, nil
}

func (a *ExportCmd) loadSchema(
	ctx context.Context,
	loader db.SchemaLoader,
	params *params.ExportParams,
) (*schema.Schema, error) {
	sc, err := loader.Load(ctx)
	if err != nil {
		return nil, err
	}

	if len(params.Tables) > 0 {
		slog.DebugContext(ctx, "[exportcmd] filtering tables")

		sc.Tables = sc.Tables.Without(params.Tables)
	}

	slog.DebugContext(ctx, "[exportcmd] sorting tables by relations")

	sc.SortByRelations()

	if !params.WithoutMigrationsTable {
		return sc, nil
	}

	sc.Tables = sc.Tables.Reject(func(table *schema.Table) bool {
		return a.migrationsTblDetector.IsMigrationsTable(table.Name.Val, table.ColumnsNames())
	})

	return sc, nil
}

func (a *ExportCmd) createRenderer() *template.Renderer {
	var templateLoader stick.Loader

	if a.fs.Exists(localTemplatesFolder) {
		slog.Debug(fmt.Sprintf("[exportcmd] loading templates from folder %q", localTemplatesFolder))

		templateLoader = stick.NewFilesystemLoader(localTemplatesFolder)
	} else {
		slog.Debug("[exportcmd] loading templates from embedded files")

		templateLoader = template.NewEmbedLoader(templates.FS)
	}

	return template.NewRenderer(templateLoader)
}

func (a *ExportCmd) runActions(ctx context.Context, p *params.ActionParams) error {
	for actionName, action := range a.actions {
		if action.Supports(p) {
			slog.DebugContext(ctx, fmt.Sprintf("[exportcmd] running action %q", actionName))

			err := action.Run(ctx, p)
			if err != nil {
				return fmt.Errorf("failed to run action %q: %w", actionName, err)
			}
		}
	}

	return nil
}
