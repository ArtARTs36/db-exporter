package laravel

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/shared/laravel"
	"github.com/artarts36/db-exporter/internal/shared/sqlquery"
	"github.com/artarts36/db-exporter/internal/sql"
	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
)

type MigrationsRawExporter struct {
	pager      *common.Pager
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
	pager *common.Pager,
	ddlBuilder *sql.DDLBuilder,
) *MigrationsRawExporter {
	return &MigrationsRawExporter{
		pager:      pager,
		ddlBuilder: ddlBuilder,
	}
}

func (e *MigrationsRawExporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())
	i := 0

	migrationPage := e.pager.Of("laravel/migration-raw.php")

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

		page, err := migrationPage.Export(
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

func (e *MigrationsRawExporter) Export(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
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

	page, err := e.pager.Of("laravel/migration-raw.php").Export(
		laravel.CreateMigrationFilename("init", 0),
		map[string]stick.Value{
			"migration": migration,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to render: %w", err)
	}

	return []*exporter.ExportedPage{
		page,
	}, nil
}

func (e *MigrationsRawExporter) makeMigrationQueries(table *schema.Table) *laravelMigrationQueries {
	return &laravelMigrationQueries{
		Up: e.ddlBuilder.BuildDDL(table),
		Down: []string{
			sqlquery.BuildDropTable(table.Name.Value),
		},
	}
}
