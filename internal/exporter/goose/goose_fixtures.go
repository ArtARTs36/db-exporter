package goose

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/exporter/migrations"
	"github.com/artarts36/db-exporter/internal/infrastructure/data"
	"log/slog"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/shared/goose"
	"github.com/artarts36/db-exporter/internal/sql"
)

type FixturesExporter struct {
	page         *common.Page
	dataLoader   *data.Loader
	queryBuilder *sql.QueryBuilder
}

func NewFixturesExporter(
	pager *common.Pager,
	dataLoader *data.Loader,
	insertBuilder *sql.QueryBuilder,
) *FixturesExporter {
	return &FixturesExporter{
		page:         pager.Of("goose/migration.sql"),
		dataLoader:   dataLoader,
		queryBuilder: insertBuilder,
	}
}

func (e *FixturesExporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	slog.DebugContext(ctx, "[goose-fixtures-exporter] building queries and rendering migration files")

	for i, table := range params.Schema.Tables.List() {
		data, err := e.dataLoader.Load(ctx, params.Conn, table.Name.Value)
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

		migration := &migrations.Migration{
			UpQueries:   []string{upQuery},
			DownQueries: e.queryBuilder.BuildDeleteQueries(table, data),
		}

		p, err := e.page.Export(
			goose.CreateMigrationFilename(fmt.Sprintf(
				"inserts_into_%s_table",
				table.Name.Value,
			), i),
			map[string]stick.Value{
				"migration": migration,
			},
		)
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (e *FixturesExporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	migration := &migrations.Migration{
		UpQueries:   make([]string, 0, params.Schema.Tables.Len()),
		DownQueries: make([]string, 0, params.Schema.Tables.Len()),
	}

	slog.DebugContext(ctx, "[goose-fixtures-exporter] building queries")

	for _, table := range params.Schema.Tables.List() {
		data, err := e.dataLoader.Load(ctx, params.Conn, table.Name.Value)
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

		migration.UpQueries = append(migration.UpQueries, upQuery)
		migration.DownQueries = append(migration.DownQueries, e.queryBuilder.BuildDeleteQueries(table, data)...)
	}

	slog.DebugContext(ctx, "[goose-fixtures-exporter] rendering migration file")

	p, err := e.page.Export(
		goose.CreateMigrationFilename("inserts", 1),
		map[string]stick.Value{
			"migration": migration,
		},
	)
	if err != nil {
		return nil, err
	}

	return []*exporter.ExportedPage{
		p,
	}, nil
}
