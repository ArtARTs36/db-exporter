package laravel

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/migrations"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/shared/laravel"
	"github.com/artarts36/db-exporter/internal/sql"
)

func NewLaravelMigrationsRawExporter(
	pager *common.Pager,
	ddlBuilder *sql.DDLBuilder,
) *migrations.Exporter {
	return migrations.NewExporter(
		"laravel-raw-migrations-exporter",
		pager,
		"laravel/migration-raw.php",
		ddlBuilder,
		migrations.NewFuncMigrationMaker(
			func(i int, tableName ds.String) *migrations.MigrationMeta {
				return &migrations.MigrationMeta{
					Filename: laravel.CreateMigrationFilename(fmt.Sprintf("create_%s_table", tableName), i),
					Attrs: map[string]interface{}{
						"Name": fmt.Sprintf("Create%sTable", tableName.Pascal().Value),
					},
				}
			},
			func() *migrations.MigrationMeta {
				return &migrations.MigrationMeta{
					Filename: laravel.CreateMigrationFilename("init", 1),
					Attrs: map[string]interface{}{
						"Name": "InitMigration",
					},
				}
			},
		),
	)
}
