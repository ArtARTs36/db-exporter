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

type LaravelMigrationsRawExporter struct {
	renderer   *template.Renderer
	ddlBuilder *sql.DDLBuilder
}

type laravelMigration struct {
	Name    string
	Queries *laravelMigrationQueries
}

type laravelMigrationQueries struct {
	Up   []string
	Down []string
}

func NewLaravelMigrationsRawExporter(
	renderer *template.Renderer,
	ddlBuilder *sql.DDLBuilder,
) *LaravelMigrationsRawExporter {
	return &LaravelMigrationsRawExporter{
		renderer:   renderer,
		ddlBuilder: ddlBuilder,
	}
}

func (e *LaravelMigrationsRawExporter) ExportPerFile(
	_ context.Context,
	params *ExportParams,
) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, params.Schema.Tables.Len())
	i := 0

	for _, table := range params.Schema.Tables.List() {
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

func (e *LaravelMigrationsRawExporter) Export(
	_ context.Context,
	params *ExportParams,
) ([]*ExportedPage, error) {
	migration := &laravelMigration{
		Name: "InitMigration",
		Queries: &laravelMigrationQueries{
			Up:   []string{},
			Down: []string{},
		},
	}

	for _, table := range params.Schema.Tables.List() {
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

func (e *LaravelMigrationsRawExporter) makeMigrationQueries(table *schema.Table) *laravelMigrationQueries {
	return &laravelMigrationQueries{
		Up: e.ddlBuilder.BuildDDL(table),
		Down: []string{
			sqlquery.BuildDropTable(table.Name.Value),
		},
	}
}
