package gosqlmigrate

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"log"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/gosqlmigrate"
	"github.com/artarts36/db-exporter/internal/shared/sqlquery"
	"github.com/artarts36/db-exporter/internal/sql"
	"github.com/artarts36/db-exporter/internal/template"
)

type Exporter struct {
	pager      *common.Pager
	renderer   *template.Renderer
	ddlBuilder *sql.DDLBuilder
}

type goSQLMigrateMigration struct {
	upQueries   []string
	downQueries []string
}

func NewSQLMigrateExporter(pager *common.Pager, renderer *template.Renderer, ddlBuilder *sql.DDLBuilder) *Exporter {
	return &Exporter{
		pager:      pager,
		renderer:   renderer,
		ddlBuilder: ddlBuilder,
	}
}

func (e *Exporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	log.Printf("[go-sql-migrate-exporter] building queries and rendering migration files")

	migrationPage := e.pager.Of("go-sql-migrate/migration.sql")

	for i, table := range params.Schema.Tables.List() {
		migration := e.makeMigration(table)

		p, err := migrationPage.Export(
			gosqlmigrate.CreateMigrationFilename(fmt.Sprintf("create_%s_table", table.Name.Value), i),
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

func (e *Exporter) Export(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	upQueries := make([]string, 0, params.Schema.Tables.Len())
	downQueries := make([]string, 0, params.Schema.Tables.Len())

	log.Printf("[go-sql-migrate-exporter] building queries")

	for _, table := range params.Schema.Tables.List() {
		migration := e.makeMigration(table)

		upQueries = append(upQueries, migration.upQueries...)
		downQueries = append(downQueries, migration.downQueries...)
	}

	log.Printf("[go-sql-migrate-exporter] rendering migration file")

	p, err := e.pager.Of("go-sql-migrate/migration.sql").Export(
		gosqlmigrate.CreateMigrationFilename("init", 1),
		map[string]stick.Value{
			"up_queries":   upQueries,
			"down_queries": downQueries,
		},
	)
	if err != nil {
		return nil, err
	}

	return []*exporter.ExportedPage{
		p,
	}, nil
}

func (e *Exporter) makeMigration(table *schema.Table) *goSQLMigrateMigration {
	return &goSQLMigrateMigration{
		upQueries: e.ddlBuilder.BuildDDL(table),
		downQueries: []string{
			sqlquery.BuildDropTable(table.Name.Value),
		},
	}
}
