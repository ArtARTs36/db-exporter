package exporter

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/goose"
	"github.com/artarts36/db-exporter/internal/shared/sqlquery"
	"github.com/artarts36/db-exporter/internal/sql"
	"github.com/artarts36/db-exporter/internal/template"
)

const GooseExporterName = "goose"

type GooseExporter struct {
	renderer   *template.Renderer
	ddlBuilder *sql.DDLBuilder
}

type gooseMigration struct {
	upQueries   []string
	downQueries []string
}

func NewGooseExporter(renderer *template.Renderer, ddlBuilder *sql.DDLBuilder) *GooseExporter {
	return &GooseExporter{
		renderer:   renderer,
		ddlBuilder: ddlBuilder,
	}
}

func (e *GooseExporter) ExportPerFile(
	ctx context.Context,
	sch *schema.Schema,
	_ *ExportParams,
) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, sch.Tables.Len())

	slog.DebugContext(ctx, "[goose-exporter] building queries and rendering migration files")

	for i, table := range sch.Tables.List() {
		migration := e.makeMigration(table)

		p, err := render(
			e.renderer,
			"goose/migration.sql",
			goose.CreateMigrationFilename(fmt.Sprintf(
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

func (e *GooseExporter) Export(ctx context.Context, sch *schema.Schema, _ *ExportParams) ([]*ExportedPage, error) {
	upQueries := make([]string, 0, sch.Tables.Len())
	downQueries := make([]string, 0, sch.Tables.Len())

	slog.DebugContext(ctx, "[gooseexporter] building queries")

	for _, table := range sch.Tables.List() {
		migration := e.makeMigration(table)

		upQueries = append(upQueries, migration.upQueries...)
		downQueries = append(downQueries, migration.downQueries...)
	}

	slog.DebugContext(ctx, "[goose-exporter] rendering migration file")

	p, err := render(
		e.renderer,
		"goose/migration.sql",
		goose.CreateMigrationFilename("init", 1),
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
		upQueries: e.ddlBuilder.BuildDDL(table),
		downQueries: []string{
			sqlquery.BuildDropTable(table.Name.Value),
		},
	}
}
