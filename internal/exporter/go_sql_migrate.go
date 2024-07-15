package exporter

import (
	"context"
	"fmt"
	"log"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/gosqlmigrate"
	"github.com/artarts36/db-exporter/internal/shared/sqlquery"
	"github.com/artarts36/db-exporter/internal/sql"
	"github.com/artarts36/db-exporter/internal/template"
)

const GoSQLMigrateExporterName = "go-sql-migrate"

type GoSQLMigrateExporter struct {
	unimplementedImporter
	renderer   *template.Renderer
	ddlBuilder *sql.DDLBuilder
}

type goSQLMigrateMigration struct {
	upQueries   []string
	downQueries []string
}

func NewSQLMigrateExporter(renderer *template.Renderer, ddlBuilder *sql.DDLBuilder) *GoSQLMigrateExporter {
	return &GoSQLMigrateExporter{
		renderer:   renderer,
		ddlBuilder: ddlBuilder,
	}
}

func (e *GoSQLMigrateExporter) ExportPerFile(
	_ context.Context,
	sch *schema.Schema,
	_ *ExportParams,
) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, sch.Tables.Len())

	log.Printf("[go-sql-migrate-exporter] building queries and rendering migration files")

	for i, table := range sch.Tables.List() {
		migration := e.makeMigration(table)

		p, err := render(
			e.renderer,
			"go-sql-migrate/migration.sql",
			gosqlmigrate.CreateMigrationFilename(fmt.Sprintf(
				"create_%s_table",
				table.Name.Value,
			), i),
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

func (e *GoSQLMigrateExporter) Export(
	_ context.Context,
	sch *schema.Schema,
	_ *ExportParams,
) ([]*ExportedPage, error) {
	upQueries := make([]string, 0, sch.Tables.Len())
	downQueries := make([]string, 0, sch.Tables.Len())

	log.Printf("[go-sql-migrate-exporter] building queries")

	for _, table := range sch.Tables.List() {
		migration := e.makeMigration(table)

		upQueries = append(upQueries, migration.upQueries...)
		downQueries = append(downQueries, migration.downQueries...)
	}

	log.Printf("[go-sql-migrate-exporter] rendering migration file")

	p, err := render(
		e.renderer,
		"go-sql-migrate/migration.sql",
		gosqlmigrate.CreateMigrationFilename("init", 1),
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

func (e *GoSQLMigrateExporter) makeMigration(table *schema.Table) *goSQLMigrateMigration {
	return &goSQLMigrateMigration{
		upQueries: e.ddlBuilder.BuildDDL(table),
		downQueries: []string{
			sqlquery.BuildDropTable(table.Name.Value),
		},
	}
}
