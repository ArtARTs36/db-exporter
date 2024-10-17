package config

import orderedmap "github.com/wk8/go-ordered-map/v2"

type ExporterName string

const (
	ExporterNameMd                   ExporterName = "md"
	ExporterNameDiagram              ExporterName = "diagram"
	ExporterNameGoEntities           ExporterName = "go-entities"
	ExporterNameGoEntityRepository   ExporterName = "go-entity-repository"
	ExporterNameGoose                ExporterName = "goose"
	ExporterNameGooseFixtures        ExporterName = "goose-fixtures"
	ExporterNameGoSQLMigrate         ExporterName = "go-sql-migrate"
	ExporterNameLaravelMigrationsRaw ExporterName = "laravel-migrations-raw"
	ExporterNameLaravelModels        ExporterName = "laravel-models"
	ExporterNameGrpcCrud             ExporterName = "grpc-crud"
	ExporterNameYamlFixtures         ExporterName = "yaml-fixtures"
	ExporterNameCSV                  ExporterName = "csv"
	ExporterNameJSONSchema           ExporterName = "json-schema"
)

type GoEntitiesExportSpec struct {
	GoModule string `yaml:"go_module"`
	Package  string `yaml:"package"` // default: entities
}

type GRPCCrudExportSpec struct {
	Package string                                     `yaml:"package"`
	Options orderedmap.OrderedMap[string, interface{}] `yaml:"options"`
}

type MarkdownExportSpec struct {
	WithDiagram bool `yaml:"with_diagram"`
}

type CSVExportSpec struct {
	Delimiter string                           `yaml:"delimiter"`
	Transform map[string][]ExportSpecTransform `yaml:"transform"`
}

type ExportSpecTransform struct {
	OnlyColumns   []string          `yaml:"only_columns"`
	SkipColumns   []string          `yaml:"skip_columns"`
	RenameColumns map[string]string `yaml:"rename_columns"`
}

type LaravelModelsExportSpec struct {
	Namespace string `yaml:"namespace"`
	TimeAs    string `yaml:"time_as"` // datetime, carbon
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
			Place GoEntityRepositorySpecRepoInterfacesPlace `yaml:"place" json:"place"`
		} `yaml:"interfaces" json:"interfaces"`
	} `yaml:"repositories" json:"repositories"`
}

type JSONSchemaExportSpec struct {
	Pretty bool `yaml:"pretty" json:"pretty"`
	Schema struct {
		Title       string `yaml:"title" json:"title"`
		Description string `yaml:"description" json:"description"`
	} `yaml:"schema" json:"schema"`
}
