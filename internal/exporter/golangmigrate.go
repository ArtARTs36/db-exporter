package exporter

import (
	"context"
	"fmt"
	"log"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/goose"
	"github.com/artarts36/db-exporter/internal/shared/sqlquery"
	"github.com/artarts36/db-exporter/internal/sql"
	"github.com/artarts36/db-exporter/internal/template"
)

const GolangMigrateExporterName = "golang-migrate"

type GolangMigrateExporter struct {
	renderer   *template.Renderer
	ddlBuilder *sql.DDLBuilder
}

type golangMigrateMigration struct {
	upQueries   []string
	downQueries []string
}

func NewGolangMigrateExporter(renderer *template.Renderer, ddlBuilder *sql.DDLBuilder) *GolangMigrateExporter {
	return &GolangMigrateExporter{
		renderer:   renderer,
		ddlBuilder: ddlBuilder,
	}
}

func (e *GolangMigrateExporter) ExportPerFile(
	_ context.Context,
	sch *schema.Schema,
	_ *ExportParams,
) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, sch.Tables.Len())

	log.Printf("[golang-migrate-exporter] building queries and rendering migration files")

	for _, table := range sch.Tables.List() {
		migration := e.makeMigration(table)

		p, err := render(
			e.renderer,
			"golang-migrate/migration.sql",
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

func (e *GolangMigrateExporter) Export(
	_ context.Context,
	sch *schema.Schema,
	_ *ExportParams,
) ([]*ExportedPage, error) {
	upQueries := make([]string, 0, sch.Tables.Len())
	downQueries := make([]string, 0, sch.Tables.Len())

	log.Printf("[golang-migrate-exporter] building queries")

	for _, table := range sch.Tables.List() {
		migration := e.makeMigration(table)

		upQueries = append(upQueries, migration.upQueries...)
		downQueries = append(downQueries, migration.downQueries...)
	}

	log.Printf("[golang-migrate-exporter] rendering migration file")

	p, err := render(
		e.renderer,
		"golang-migrate/migration.sql",
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

func (e *GolangMigrateExporter) makeMigration(table *schema.Table) *golangMigrateMigration {
	return &golangMigrateMigration{
		upQueries: e.ddlBuilder.BuildDDL(table),
		downQueries: []string{
			sqlquery.BuildDropTable(table.Name.Value),
		},
	}
}
