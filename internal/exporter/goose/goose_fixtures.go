package goose

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"log/slog"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/shared/goose"
	"github.com/artarts36/db-exporter/internal/sql"
)

type FixturesExporter struct {
	pager        *common.Pager
	dataLoader   *db.DataLoader
	queryBuilder *sql.QueryBuilder
}

func NewFixturesExporter(
	pager *common.Pager,
	dataLoader *db.DataLoader,
	insertBuilder *sql.QueryBuilder,
) *FixturesExporter {
	return &FixturesExporter{
		pager:        pager,
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

		migration := e.makeMigration([]string{upQuery}, e.queryBuilder.BuildDeleteQueries(table, data))

		p, err := e.pager.Of("goose/migration.sql").Export(
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

func (e *FixturesExporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	upQueries := make([]string, 0, params.Schema.Tables.Len())
	downQueries := make([]string, 0, params.Schema.Tables.Len())

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

		upQueries = append(upQueries, upQuery)
		downQueries = append(downQueries, e.queryBuilder.BuildDeleteQueries(table, data)...)
	}

	slog.DebugContext(ctx, "[goose-fixtures-exporter] rendering migration file")

	p, err := e.pager.Of("goose/migration.sql").Export(
		goose.CreateMigrationFilename("inserts", 1),
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

func (e *FixturesExporter) makeMigration(upQueries []string, downQueries []string) *gooseMigration {
	return &gooseMigration{
		upQueries:   upQueries,
		downQueries: downQueries,
	}
}
