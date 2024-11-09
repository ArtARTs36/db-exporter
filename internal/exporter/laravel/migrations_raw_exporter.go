package laravel

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/migrations"
	"github.com/artarts36/db-exporter/internal/infrastructure/sql"
	"github.com/artarts36/db-exporter/internal/shared/laravel"
	"github.com/artarts36/gds"
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
			func(i int, tableName gds.String) *migrations.MigrationMeta {
				return &migrations.MigrationMeta{
					Filename: laravel.CreateMigrationFilename(fmt.Sprintf("create_%s_table", tableName.Value), i),
					Attrs: map[string]interface{}{
						"name": fmt.Sprintf("Create%sTable", tableName.Pascal().Value),
					},
				}
			},
			func() *migrations.MigrationMeta {
				return &migrations.MigrationMeta{
					Filename: laravel.CreateMigrationFilename("init", 1),
					Attrs: map[string]interface{}{
						"name": "InitMigration",
					},
				}
			},
		),
	)
}
