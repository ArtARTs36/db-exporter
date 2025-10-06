package ddl

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/migrations"
	"github.com/artarts36/db-exporter/internal/infrastructure/sql"
	"github.com/artarts36/gds"
)

func NewExporter(
	pager *common.Pager,
	ddlBuilder *sql.DDLBuilderManager,
) *migrations.Exporter {
	return migrations.NewExporter(
		"ddl-exporter",
		pager,
		"@embed/ddl/queries.ddl",
		ddlBuilder,
		migrations.NewFuncMigrationMaker(
			func(_ int, tableName gds.String) *migrations.MigrationMeta {
				return &migrations.MigrationMeta{
					Filename: fmt.Sprintf("%s.ddl", tableName.Value),
				}
			},
			func() *migrations.MigrationMeta {
				return &migrations.MigrationMeta{
					Filename: "result.ddl",
				}
			},
		),
	)
}
