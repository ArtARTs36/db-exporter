package app

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/schemaloader"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"github.com/artarts36/db-exporter/internal/shared/migrations"
	"github.com/artarts36/db-exporter/internal/template"
	"github.com/artarts36/db-exporter/templates"
)

const localTemplatesFolder = "./db-exporter-templates"

type ExportCmd struct {
	migrationsTblDetector *migrations.TableDetector
}

type ExportParams struct {
	DriverName string
	DSN        string
	Format     string
	OutDir     string

	TablePerFile           bool
	WithDiagram            bool
	WithoutMigrationsTable bool
	Tables                 []string
}

func NewExportCmd() *ExportCmd {
	return &ExportCmd{
		migrationsTblDetector: migrations.NewTableDetector(),
	}
}

func (a *ExportCmd) Export(ctx context.Context, params *ExportParams) error {
	loader, err := schemaloader.CreateLoader(params.DriverName)
	if err != nil {
		return fmt.Errorf("unable to create schema loader: %w", err)
	}

	renderer := a.createRenderer()

	exp, err := exporter.CreateExporter(params.Format, renderer)
	if err != nil {
		return fmt.Errorf("failed to create exporter: %w", err)
	}

	// processing

	log.Printf("[exportcmd] loading db schema from %s", params.DriverName)

	sc, err := a.loadSchema(ctx, loader, params)
	if err != nil {
		return fmt.Errorf("unable to load schema: %w", err)
	}

	log.Printf("[exportcmd] loaded %d tables: [%s]", sc.Tables.Len(), strings.Join(sc.TablesNames(), ","))

	pages, err := a.export(ctx, exp, sc, params)
	if err != nil {
		return fmt.Errorf("failed to export: %w", err)
	}

	err = a.savePages(pages, params)
	if err != nil {
		return err
	}

	log.Printf("[exportcmd] successful generated %d files", len(pages))

	return nil
}

func (a *ExportCmd) export(
	ctx context.Context,
	exp exporter.Exporter,
	sc *schema.Schema,
	params *ExportParams,
) ([]*exporter.ExportedPage, error) {
	var pages []*exporter.ExportedPage
	var err error
	exporterParams := &exporter.ExportParams{
		WithDiagram:            params.WithDiagram,
		WithoutMigrationsTable: params.WithoutMigrationsTable,
	}

	if params.TablePerFile {
		pages, err = exp.ExportPerFile(ctx, sc, exporterParams)
	} else {
		pages, err = exp.Export(ctx, sc, exporterParams)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export: %w", err)
	}

	return pages, nil
}

func (a *ExportCmd) savePages(pages []*exporter.ExportedPage, params *ExportParams) error {
	if !fs.Exists(params.OutDir) {
		log.Printf("creating directory %q", params.OutDir)

		err := fs.Mkdir(params.OutDir)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	for _, page := range pages {
		path := fmt.Sprintf("%s/%s", params.OutDir, page.FileName)

		log.Printf("[exportcmd] saving %q", path)

		err := fs.CreateFile(path, page.Content)
		if err != nil {
			return fmt.Errorf("unable to write file: %w", err)
		}
	}

	return nil
}

func (a *ExportCmd) loadSchema(
	ctx context.Context,
	loader schemaloader.Loader,
	params *ExportParams,
) (*schema.Schema, error) {
	sc, err := loader.Load(ctx, params.DSN)
	if err != nil {
		return nil, err
	}

	if len(params.Tables) > 0 {
		tableFilter := map[string]bool{}
		for _, table := range params.Tables {
			tableFilter[table] = true
		}

		log.Println("[exportcmd] filtering tables")

		sc.Tables = sc.Tables.Reject(func(table *schema.Table) bool {
			return !tableFilter[table.Name.Value]
		})
	}

	log.Println("[exportcmd] sorting tables by relations")

	sc.SortByRelations()

	if !params.WithoutMigrationsTable {
		return sc, nil
	}

	sc.Tables = sc.Tables.Reject(func(table *schema.Table) bool {
		return a.migrationsTblDetector.IsMigrationsTable(table.Name.Value, table.ColumnsNames())
	})

	return sc, nil
}

func (a *ExportCmd) createRenderer() *template.Renderer {
	var templateLoader stick.Loader

	if fs.Exists(localTemplatesFolder) {
		log.Printf("[exportcmd] loading templates from folder %q", localTemplatesFolder)

		templateLoader = stick.NewFilesystemLoader("./db-exporter-templates")
	} else {
		log.Print("[exportcmd] loading templates from embedded files")

		templateLoader = template.NewEmbedLoader(templates.FS)
	}

	return template.NewRenderer(templateLoader)
}
