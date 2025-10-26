package factory

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/csv"
	"github.com/artarts36/db-exporter/internal/exporter/custom"
	"github.com/artarts36/db-exporter/internal/exporter/dbml"
	"github.com/artarts36/db-exporter/internal/exporter/ddl"
	"github.com/artarts36/db-exporter/internal/exporter/diagram"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	goentity "github.com/artarts36/db-exporter/internal/exporter/go-entity"
	gosqlmigrate "github.com/artarts36/db-exporter/internal/exporter/go-sql-migrate"
	"github.com/artarts36/db-exporter/internal/exporter/goose"
	"github.com/artarts36/db-exporter/internal/exporter/graphql"
	grpccrud "github.com/artarts36/db-exporter/internal/exporter/grpc-crud"
	"github.com/artarts36/db-exporter/internal/exporter/jsonschema"
	"github.com/artarts36/db-exporter/internal/exporter/laravel"
	"github.com/artarts36/db-exporter/internal/exporter/markdown"
	"github.com/artarts36/db-exporter/internal/exporter/mermaid"
	"github.com/artarts36/db-exporter/internal/infrastructure/data"
	"github.com/artarts36/db-exporter/internal/infrastructure/sql"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/artarts36/db-exporter/internal/template"
)

func CreateExporters(renderer *template.Renderer) map[config.ExporterName]exporter.Exporter {
	pager := common.NewPager(renderer)
	dataLoader := data.NewLoader()
	dataTransformers := []csv.DataTransformer{
		csv.OnlyColumnsDataTransformer(),
		csv.SkipColumnsDataTransformer(),
		csv.RenameColumnsDataTransformer(),
	}
	goEntityMapper := goentity.NewEntityMapper(goentity.NewGoPropertyMapper())
	goEntityGenerator := goentity.NewEntityGenerator(pager)
	goModFinder := golang.NewModFinder()
	goPropertyMapper := goentity.NewGoPropertyMapper()
	ddlBuilderManager := sql.NewDDLBuilderManager()
	diagramCreator := diagram.NewCreator(diagram.NewGraphBuilder(renderer))

	return map[config.ExporterName]exporter.Exporter{
		config.ExporterNameMd:      markdown.NewMarkdownExporter(pager, diagramCreator),
		config.ExporterNameDiagram: diagram.NewDiagramExporter(diagramCreator),
		config.ExporterNameGoEntities: goentity.NewEntitiesExporter(
			pager,
			goEntityMapper,
			goEntityGenerator,
			goModFinder,
		),
		config.ExporterNameGoose: goose.NewMigrationsExporter(pager, ddlBuilderManager),
		config.ExporterNameGoSQLMigrate: gosqlmigrate.NewSQLMigrateExporter(
			pager,
			ddlBuilderManager,
		),
		config.ExporterNameLaravelMigrationsRaw: laravel.NewLaravelMigrationsRawExporter(pager, ddlBuilderManager),
		config.ExporterNameLaravelModels:        laravel.NewLaravelModelsExporter(pager),
		config.ExporterNameGrpcCrud:             grpccrud.NewExporter(),
		config.ExporterNameGooseFixtures:        goose.NewFixturesExporter(pager, dataLoader, sql.NewInsertBuilder()),
		config.ExporterNameCSV:                  csv.NewExporter(dataLoader, dataTransformers),
		config.ExporterNameGoEntityRepository: goentity.NewRepositoryExporter(
			pager,
			goModFinder,
			goEntityMapper,
			goEntityGenerator,
			goPropertyMapper,
		),
		config.ExporterNameJSONSchema: jsonschema.NewExporter(),
		config.ExporterNameGraphql:    graphql.NewExporter(),
		config.ExporterNameDBML:       dbml.NewExporter(),
		config.ExporterNameCustom:     custom.NewExporter(renderer),
		config.ExporterNameDDL:        ddl.NewExporter(pager, ddlBuilderManager),
		config.ExporterNameMermaid:    mermaid.NewExporter(),
	}
}
