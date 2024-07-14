package exporter

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/goose"
	"github.com/artarts36/db-exporter/internal/sql"
	"github.com/artarts36/db-exporter/internal/template"
)

const GooseFixturesExporterName = "goose-fixtures"

type GooseFixturesExporter struct {
	dataLoader   *db.DataLoader
	renderer     *template.Renderer
	queryBuilder *sql.QueryBuilder
}

func NewGooseFixturesExporter(
	dataLoader *db.DataLoader,
	renderer *template.Renderer,
	insertBuilder *sql.QueryBuilder,
) *GooseFixturesExporter {
	return &GooseFixturesExporter{
		dataLoader:   dataLoader,
		renderer:     renderer,
		queryBuilder: insertBuilder,
	}
}

func (e *GooseFixturesExporter) ExportPerFile(
	ctx context.Context,
	sch *schema.Schema,
	_ *ExportParams,
) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, sch.Tables.Len())

	slog.DebugContext(ctx, "[goose-fixtures-exporter] building queries and rendering migration files")

	for i, table := range sch.Tables.List() {
		data, err := e.dataLoader.Load(ctx, table.Name.Value)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			continue
		}

		upQuery, err := e.queryBuilder.BuildInsertQuery(table, data)
		if err != nil {
			return nil, err
		}

		migration := e.makeMigration([]string{upQuery}, e.queryBuilder.BuildDeleteQueries(table, data))

		p, err := render(
			e.renderer,
			"goose/migration.sql",
			goose.CreateMigrationFilename(fmt.Sprintf(
				"inserts_into_%s_table",
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

func (e *GooseFixturesExporter) Export(
	ctx context.Context,
	sch *schema.Schema,
	_ *ExportParams,
) ([]*ExportedPage, error) {
	upQueries := make([]string, 0, sch.Tables.Len())
	downQueries := make([]string, 0, sch.Tables.Len())

	slog.DebugContext(ctx, "[goose-fixtures-exporter] building queries")

	for _, table := range sch.Tables.List() {
		data, err := e.dataLoader.Load(ctx, table.Name.Value)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			continue
		}

		upQuery, err := e.queryBuilder.BuildInsertQuery(table, data)
		if err != nil {
			return nil, err
		}

		upQueries = append(upQueries, upQuery)
		downQueries = append(downQueries, e.queryBuilder.BuildDeleteQueries(table, data)...)
	}

	slog.DebugContext(ctx, "[goose-fixtures-exporter] rendering migration file")

	p, err := render(
		e.renderer,
		"goose/migration.sql",
		goose.CreateMigrationFilename("inserts", 1),
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

func (e *GooseFixturesExporter) makeMigration(upQueries []string, downQueries []string) *gooseMigration {
	return &gooseMigration{
		upQueries:   upQueries,
		downQueries: downQueries,
	}
}
