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

type GoEntitiesExportSpec struct {
	GoModule string `yaml:"go_module" json:"go_module"`
	Package  string `yaml:"package" json:"package"` // default: entities
}

type GoEntityRepositorySpecRepoInterfacesPlace string

const (
	GoEntityRepositorySpecRepoInterfacesPlaceUnspecified    = ""
	GoEntityRepositorySpecRepoInterfacesPlaceWithEntity     = "with_entity"
	GoEntityRepositorySpecRepoInterfacesPlaceWithRepository = "with_repository"
	GoEntityRepositorySpecRepoInterfacesPlaceEntity         = "entity"
)

type GoEntityRepositorySpec struct {
	GoModule string `yaml:"go_module" json:"go_module"`
	Entities struct {
		Package string `yaml:"package" json:"package"`
	} `yaml:"entities" json:"entities"`
	Repositories struct {
		Package   string `yaml:"package" json:"package"`
		Container struct {
			StructName string `yaml:"struct_name" json:"struct_name"`
		} `yaml:"container" json:"container"`
		Interfaces struct {
			Place     GoEntityRepositorySpecRepoInterfacesPlace `yaml:"place" json:"place"`
			WithMocks bool                                      `yaml:"with_mocks" json:"with_mocks"`
		} `yaml:"interfaces" json:"interfaces"`
	} `yaml:"repositories" json:"repositories"`
}
