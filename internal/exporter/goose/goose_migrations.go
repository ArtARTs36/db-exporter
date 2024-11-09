package goose

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/migrations"
	"github.com/artarts36/db-exporter/internal/shared/goose"
	"github.com/artarts36/db-exporter/internal/sql"
	"github.com/artarts36/gds"
)

func NewMigrationsExporter(
	pager *common.Pager,
	ddlBuilder *sql.DDLBuilder,
) *migrations.Exporter {
	return migrations.NewExporter(
		"goose-migrations-exporter",
		pager,
		"goose/migration.sql",
		ddlBuilder,
		migrations.NewFuncMigrationMaker(
			func(i int, tableName gds.String) *migrations.MigrationMeta {
				return &migrations.MigrationMeta{
					Filename: goose.CreateMigrationFilename(fmt.Sprintf("create_%s_table", tableName), i),
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
