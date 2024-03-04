package app

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/artarts36/db-exporter/internal/exporter"
	"github.com/artarts36/db-exporter/internal/schemaloader"
	"github.com/artarts36/db-exporter/internal/template"
)

type ExportCmd struct {
}

type ExportParams struct {
	DSN    string
	Format string
	OutDir string

	TablePerFile bool
}

func (a *ExportCmd) Export(ctx context.Context, params *ExportParams) error {
	loader, err := schemaloader.CreateLoader("postgres")
	if err != nil {
		return fmt.Errorf("unable to create schema loader: %w", err)
	}

	renderer := template.InitRenderer("./templates")

	exp, err := exporter.CreateExporter(params.Format, renderer)
	if err != nil {
		return fmt.Errorf("failed to create exporter: %w", err)
	}

	// processing

	sc, err := loader.Load(ctx, params.DSN)
	if err != nil {
		return fmt.Errorf("unable to load schema: %w", err)
	}

	log.Printf("loaded %d tables", len(sc.Tables))

	pages, err := exp.Export(ctx, sc, &exporter.ExportParams{
		TablePerFile: params.TablePerFile,
	})
	if err != nil {
		return fmt.Errorf("failed to export: %w", err)
	}

	if _, err = os.Stat(params.OutDir); err != nil {
		log.Printf("creating directory %q", params.OutDir)

		err := os.Mkdir(params.OutDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	for _, page := range pages {
		path := fmt.Sprintf("%s/%s", params.OutDir, page.FileName)

		log.Printf("saving %q", path)

		err := os.WriteFile(path, page.Content, 0755)
		if err != nil {
			return fmt.Errorf("unable to write file: %w", err)
		}
	}

	return nil
}
