package config

type ExporterName string

const (
	ExporterNameMd                   ExporterName = "md"
	ExporterNameDiagram              ExporterName = "diagram"
	ExporterNameGoEntities           ExporterName = "go-entities"
	ExporterNameGoEntityRepository   ExporterName = "go-entity-repository"
	ExporterNameGoose                ExporterName = "goose"
	ExporterNameGooseFixtures        ExporterName = "goose-fixtures"
	ExporterNameGoSQLMigrate         ExporterName = "go-sql-migrate"
	ExporterNameDDL                  ExporterName = "ddl"
	ExporterNameLaravelMigrationsRaw ExporterName = "laravel-migrations-raw"
	ExporterNameLaravelModels        ExporterName = "laravel-models"
	ExporterNameGrpcCrud             ExporterName = "grpc-crud"
	ExporterNameCSV                  ExporterName = "csv"
	ExporterNameJSONSchema           ExporterName = "json-schema"
	ExporterNameGraphql              ExporterName = "graphql"
	ExporterNameDBML                 ExporterName = "dbml"
	ExporterNameCustom               ExporterName = "custom"
	ExporterNameMermaid              ExporterName = "mermaid"
)
