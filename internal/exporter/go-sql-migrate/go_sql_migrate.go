package gosqlmigrate

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/migrations"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/shared/gosqlmigrate"
	"github.com/artarts36/db-exporter/internal/sql"
)

func NewSQLMigrateExporter(
	pager *common.Pager,
	ddlBuilder *sql.DDLBuilder,
) *migrations.Exporter {
	return migrations.NewExporter(
		"go-sql-migrate-exporter",
		pager,
		"go-sql-migrate/migration.sql",
		ddlBuilder,
		migrations.NewFuncMigrationMaker(
			func(i int, tableName ds.String) *migrations.MigrationMeta {
				return &migrations.MigrationMeta{
					Filename: gosqlmigrate.CreateMigrationFilename(fmt.Sprintf("create_%s_table", tableName), i),
				}
			},
			func() *migrations.MigrationMeta {
				return &migrations.MigrationMeta{
					Filename: gosqlmigrate.CreateMigrationFilename("init", 1),
				}
			},
		),
	)
}
