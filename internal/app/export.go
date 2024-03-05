package app

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/schema"
	"log"
	"os"
	"strings"

	"github.com/artarts36/db-exporter/internal/exporter"
	"github.com/artarts36/db-exporter/internal/schemaloader"
	"github.com/artarts36/db-exporter/internal/template"
)

type ExportCmd struct {
}

type ExportParams struct {
	DriverName string
	DSN        string
	Format     string
	OutDir     string

	TablePerFile           bool
	WithDiagram            bool
	WithoutMigrationsTable bool
}

func (a *ExportCmd) Export(ctx context.Context, params *ExportParams) error {
	loader, err := schemaloader.CreateLoader(params.DriverName)
	if err != nil {
		return fmt.Errorf("unable to create schema loader: %w", err)
	}

	renderer := template.InitRenderer("./templates")

	exp, err := exporter.CreateExporter(params.Format, renderer)
	if err != nil {
		return fmt.Errorf("failed to create exporter: %w", err)
	}

	// processing

	log.Printf("[exportcmd] loading db schema from %s", params.DriverName)

	sc, err := loader.Load(ctx, params.DSN)
	if err != nil {
		return fmt.Errorf("unable to load schema: %w", err)
	}

	log.Printf("[exportcmd] loaded %d tables: [%s]", len(sc.Tables), strings.Join(sc.TablesNames(), ","))

	pages, err := a.export(exp, ctx, sc, params)
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
	exp exporter.Exporter,
	ctx context.Context,
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
	if _, err := os.Stat(params.OutDir); err != nil {
		log.Printf("creating directory %q", params.OutDir)

		err := os.Mkdir(params.OutDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	for _, page := range pages {
		path := fmt.Sprintf("%s/%s", params.OutDir, page.FileName)

		log.Printf("[exportcmd] saving %q", path)

		err := os.WriteFile(path, page.Content, 0755)
		if err != nil {
			return fmt.Errorf("unable to write file: %w", err)
		}
	}

	return nil
}
