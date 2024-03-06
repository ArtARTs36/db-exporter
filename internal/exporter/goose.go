package exporter

import (
	"context"
	"fmt"
	"log"

	"github.com/artarts36/db-exporter/internal/shared/sqlquery"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/goose"
	"github.com/artarts36/db-exporter/internal/sql"
	"github.com/artarts36/db-exporter/internal/template"
)

type GooseExporter struct {
	renderer *template.Renderer
}

type gooseMigration struct {
	upQueries   []string
	downQueries []string
}

func NewGooseExporter(renderer *template.Renderer) *GooseExporter {
	return &GooseExporter{
		renderer: renderer,
	}
}

func (e *GooseExporter) ExportPerFile(_ context.Context, sch *schema.Schema, _ *ExportParams) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, len(sch.Tables))

	log.Printf("[gooseexporter] building queries and rendering migration files")

	for _, table := range sch.Tables {
		migration := e.makeMigration(table)

		p, err := render(
			e.renderer,
			"goose/migration.sql",
			goose.CreateMigrationFilename(fmt.Sprintf(
				"create_%s_table",
				table.Name.Value,
			)),
			map[string]stick.Value{
				"up_queries":   migration.upQueries,
				"down_queries": migration.downQueries,
			},
		)
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (e *GooseExporter) Export(_ context.Context, sch *schema.Schema, params *ExportParams) ([]*ExportedPage, error) {
	upQueries := make([]string, 0, len(sch.Tables))
	downQueries := make([]string, 0, len(sch.Tables))

	log.Printf("[gooseexporter] sorting tables")

	sch.SortByRelations()

	log.Printf("[gooseexporter] building queries")

	for _, table := range sch.Tables {
		if params.WithoutMigrationsTable && goose.IsMigrationsTable(table.Name.Value) {
			continue
		}

		migration := e.makeMigration(table)

		upQueries = append(upQueries, migration.upQueries...)
		downQueries = append(downQueries, migration.downQueries...)
	}

	log.Printf("[gooseexporter] rendering migration file")

	p, err := render(
		e.renderer,
		"goose/migration.sql",
		goose.CreateMigrationFilename("init"),
		map[string]stick.Value{
			"up_queries":   upQueries,
			"down_queries": downQueries,
		},
	)
	if err != nil {
		return nil, err
	}

	return []*ExportedPage{
		p,
	}, nil
}

func (e *GooseExporter) makeMigration(table *schema.Table) *gooseMigration {
	return &gooseMigration{
		upQueries: sql.BuildDDL(table),
		downQueries: []string{
			sqlquery.BuildDropTable(table.Name.Value),
		},
	}
}
