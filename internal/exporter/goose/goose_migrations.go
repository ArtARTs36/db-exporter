package goose

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"log/slog"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/goose"
	"github.com/artarts36/db-exporter/internal/shared/sqlquery"
	"github.com/artarts36/db-exporter/internal/sql"
)

type MigrationsExporter struct {
	pager      *common.Pager
	ddlBuilder *sql.DDLBuilder
}

type gooseMigration struct {
	upQueries   []string
	downQueries []string
}

func NewMigrationsExporter(
	pager *common.Pager,
	ddlBuilder *sql.DDLBuilder,
) *MigrationsExporter {
	return &MigrationsExporter{
		pager:      pager,
		ddlBuilder: ddlBuilder,
	}
}

func (e *MigrationsExporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	slog.DebugContext(ctx, "[goose-migrations-exporter] building queries and rendering migration files")

	migrationPage := e.pager.Of("goose/migration.sql")

	for i, table := range params.Schema.Tables.List() {
		migration := e.makeMigration(table)

		p, err := migrationPage.Export(
			goose.CreateMigrationFilename(fmt.Sprintf("create_%s_table", table.Name.Value), i),
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

func (e *MigrationsExporter) Export(ctx context.Context, params *exporter.ExportParams) ([]*exporter.ExportedPage, error) {
	upQueries := make([]string, 0, params.Schema.Tables.Len())
	downQueries := make([]string, 0, params.Schema.Tables.Len())

	slog.DebugContext(ctx, "[goose-migrations-exporter] building queries")

	for _, table := range params.Schema.Tables.List() {
		migration := e.makeMigration(table)

		upQueries = append(upQueries, migration.upQueries...)
		downQueries = append(downQueries, migration.downQueries...)
	}

	slog.DebugContext(ctx, "[goose-migrations-exporter] rendering migration file")

	p, err := e.pager.Of("goose/migration.sql").Export(
		goose.CreateMigrationFilename("init", 1),
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

func (e *MigrationsExporter) makeMigration(table *schema.Table) *gooseMigration {
	return &gooseMigration{
		upQueries: e.ddlBuilder.BuildDDL(table),
		downQueries: []string{
			sqlquery.BuildDropTable(table.Name.Value),
		},
	}
}
