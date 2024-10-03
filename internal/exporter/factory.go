package exporter

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/sql"

	"github.com/artarts36/db-exporter/internal/template"
)

func CreateExporters(renderer *template.Renderer) map[config.ExporterName]Exporter {
	return map[config.ExporterName]Exporter{
		config.ExporterNameMd:                   NewMarkdownExporter(renderer),
		config.ExporterNameDiagram:              NewDiagramExporter(renderer),
		config.ExporterNameGoStructs:            NewGoStructsExporter(renderer),
		config.ExporterNameGoose:                NewGooseExporter(renderer, sql.NewDDLBuilder()),
		config.ExporterNameGoSqlMigrate:         NewSQLMigrateExporter(renderer, sql.NewDDLBuilder()),
		config.ExporterNameLaravelMigrationsRaw: NewLaravelMigrationsRawExporter(renderer, sql.NewDDLBuilder()),
		config.ExporterNameGrpcCrud:             NewGrpcCrudExporter(renderer),
		config.ExporterNameGooseFixtures:        NewGooseFixturesExporter(db.NewDataLoader(), renderer, sql.NewInsertBuilder()),
		config.ExporterNameYamlFixtures:         NewYamlFixturesExporter(db.NewDataLoader(), renderer, db.NewInserter()),
	}
}
