package factory

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/csv"
	"github.com/artarts36/db-exporter/internal/exporter/diagram"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	goentity "github.com/artarts36/db-exporter/internal/exporter/go-entity"
	gosqlmigrate "github.com/artarts36/db-exporter/internal/exporter/go-sql-migrate"
	"github.com/artarts36/db-exporter/internal/exporter/goose"
	grpccrud "github.com/artarts36/db-exporter/internal/exporter/grpc-crud"
	"github.com/artarts36/db-exporter/internal/exporter/laravel"
	"github.com/artarts36/db-exporter/internal/exporter/markdown"
	"github.com/artarts36/db-exporter/internal/exporter/yaml"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/artarts36/db-exporter/internal/sql"

	"github.com/artarts36/db-exporter/internal/template"
)

func CreateExporters(renderer *template.Renderer) map[config.ExporterName]exporter.Exporter {
	pager := common.NewPager(renderer)
	dataLoader := db.NewDataLoader()
	dataTransformers := []csv.DataTransformer{
		csv.OnlyColumnsDataTransformer(),
		csv.SkipColumnsDataTransformer(),
		csv.RenameColumnsDataTransformer(),
	}
	goEntityMapper := goentity.NewEntityMapper(goentity.NewGoPropertyMapper())
	goEntityGenerator := goentity.NewEntityGenerator(pager)
	goModFinder := golang.NewModFinder()
	graphBuilder := diagram.NewGraphBuilder(renderer)

	return map[config.ExporterName]exporter.Exporter{
		config.ExporterNameMd:      markdown.NewMarkdownExporter(pager, graphBuilder),
		config.ExporterNameDiagram: diagram.NewDiagramExporter(renderer),
		config.ExporterNameGoEntities: goentity.NewEntitiesExporter(
			pager,
			goEntityMapper,
			goEntityGenerator,
			goModFinder,
		),
		config.ExporterNameGoose: goose.NewMigrationsExporter(pager, sql.NewDDLBuilder()),
		config.ExporterNameGoSQLMigrate: gosqlmigrate.NewSQLMigrateExporter(
			pager,
			renderer,
			sql.NewDDLBuilder(),
		),
		config.ExporterNameLaravelMigrationsRaw: laravel.NewLaravelMigrationsRawExporter(pager, sql.NewDDLBuilder()),
		config.ExporterNameLaravelModels:        laravel.NewLaravelModelsExporter(pager),
		config.ExporterNameGrpcCrud:             grpccrud.NewCrudExporter(pager),
		config.ExporterNameGooseFixtures:        goose.NewFixturesExporter(pager, dataLoader, sql.NewInsertBuilder()),
		config.ExporterNameYamlFixtures:         yaml.NewFixturesExporter(dataLoader, db.NewInserter()),
		config.ExporterNameCSV:                  csv.NewExporter(dataLoader, pager, dataTransformers),
		config.ExporterNameGoEntityRepository: goentity.NewRepositoryExporter(
			pager,
			goModFinder,
			goEntityMapper,
			goEntityGenerator,
		),
	}
}

func CreateImporters() map[config.ImporterName]exporter.Importer {
	return map[config.ImporterName]exporter.Importer{
		config.ImporterNameYamlFixtures: yaml.NewFixturesExporter(db.NewDataLoader(), db.NewInserter()),
	}
}
