package exporter

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/laravel"
	"github.com/artarts36/db-exporter/internal/shared/sqlquery"
	"github.com/artarts36/db-exporter/internal/sql"
	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

type LaravelMigrationsExporter struct {
	renderer *template.Renderer
}

type laravelMigration struct {
	Name    string
	Queries *laravelMigrationQueries
}

type laravelMigrationQueries struct {
	Up   []string
	Down []string
}

func NewLaravelMigrationsExporter(renderer *template.Renderer) *LaravelMigrationsExporter {
	return &LaravelMigrationsExporter{
		renderer: renderer,
	}
}

func (e *LaravelMigrationsExporter) ExportPerFile(_ context.Context, sch *schema.Schema, _ *ExportParams) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, len(sch.Tables))
	i := 0

	for _, table := range sch.Tables {
		queries := e.makeMigrationQueries(table)

		migration := &laravelMigration{
			Name: fmt.Sprintf(
				"Create%sTable",
				table.Name.Pascal().Value,
			),
			Queries: &laravelMigrationQueries{
				Up:   queries.Up,
				Down: queries.Down,
			},
		}

		page, err := render(
			e.renderer,
			"laravel/migration-raw.php",
			laravel.CreateMigrationFilename(fmt.Sprintf("create_%s_table", table.Name.Value), i),
			map[string]stick.Value{
				"migration": migration,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to render: %w", err)
		}

		pages = append(pages, page)
		i++
	}

	return pages, nil
}

func (e *LaravelMigrationsExporter) Export(_ context.Context, schema *schema.Schema, _ *ExportParams) ([]*ExportedPage, error) {
	migration := &laravelMigration{
		Name: "InitMigration",
		Queries: &laravelMigrationQueries{
			Up:   []string{},
			Down: []string{},
		},
	}

	for _, table := range schema.Tables {
		queries := e.makeMigrationQueries(table)

		migration.Queries.Up = append(migration.Queries.Up, queries.Up...)
		migration.Queries.Down = append(migration.Queries.Down, queries.Down...)
	}

	page, err := render(
		e.renderer,
		"laravel/migration-raw.php",
		laravel.CreateMigrationFilename("init", 0),
		map[string]stick.Value{
			"migration": migration,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to render: %w", err)
	}

	return []*ExportedPage{
		page,
	}, nil
}

func (e *LaravelMigrationsExporter) makeMigrationQueries(table *schema.Table) *laravelMigrationQueries {
	return &laravelMigrationQueries{
		Up: sql.BuildDDL(table),
		Down: []string{
			sqlquery.BuildDropTable(table.Name.Value),
		},
	}
}
