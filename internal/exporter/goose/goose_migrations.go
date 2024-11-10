package goose

import (
	"fmt"

	"github.com/artarts36/gds"

	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/migrations"
	"github.com/artarts36/db-exporter/internal/infrastructure/sql"
	"github.com/artarts36/db-exporter/internal/shared/goose"
)

func NewMigrationsExporter(
	pager *common.Pager,
	ddlBuilder *sql.DDLBuilderManager,
) *migrations.Exporter {
	return migrations.NewExporter(
		"goose-migrations-exporter",
		pager,
		"goose/migration.sql",
		ddlBuilder,
		migrations.NewFuncMigrationMaker(
			func(i int, tableName gds.String) *migrations.MigrationMeta {
				return &migrations.MigrationMeta{
					Filename: goose.CreateMigrationFilename(fmt.Sprintf("create_%s_table", tableName.Value), i),
				}
			},
			func() *migrations.MigrationMeta {
				return &migrations.MigrationMeta{
					Filename: goose.CreateMigrationFilename("init", 1),
				}
			},
		),
	)
}
