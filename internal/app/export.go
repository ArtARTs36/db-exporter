package app

import (
	"context"
	"fmt"
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

	TablePerFile bool
	WithDiagram  bool
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

	pages, err := exp.Export(ctx, sc, &exporter.ExportParams{
		TablePerFile: params.TablePerFile,
		WithDiagram:  params.WithDiagram,
	})
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
