package config

import (
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

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
	ExporterNameGraphql              ExporterName = "graphql"
	ExporterNameDBML                 ExporterName = "dbml"
	ExporterNameCustom               ExporterName = "custom"
)

type GoEntitiesExportSpec struct {
	GoModule string `yaml:"go_module" json:"go_module"`
	Package  string `yaml:"package" json:"package"` // default: entities
}

type GRPCCrudExportSpec struct {
	Package string                                     `yaml:"package" json:"package"`
	Options orderedmap.OrderedMap[string, interface{}] `yaml:"options" json:"options"`
}

type MarkdownExportSpec struct {
	WithDiagram bool `yaml:"with_diagram" json:"with_diagram"`
}

type CSVExportSpec struct {
	Delimiter string                           `yaml:"delimiter" json:"delimiter"`
	Transform map[string][]ExportSpecTransform `yaml:"transform" json:"transform"`
}

type ExportSpecTransform struct {
	OnlyColumns   []string          `yaml:"only_columns" json:"only_columns"`
	SkipColumns   []string          `yaml:"skip_columns" json:"skip_columns"`
	RenameColumns map[string]string `yaml:"rename_columns" json:"rename_columns"`
}

type LaravelModelsExportSpec struct {
	Namespace string `yaml:"namespace" json:"namespace"`
	TimeAs    string `yaml:"time_as" json:"time_as"` // datetime, carbon
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

type JSONSchemaExportSpec struct {
	Pretty bool `yaml:"pretty" json:"pretty"`
	Schema struct {
		Title       string `yaml:"title" json:"title"`
		Description string `yaml:"description" json:"description"`
	} `yaml:"schema" json:"schema"`
}

type MigrationsSpec struct {
	Use struct {
		IfNotExists bool `yaml:"if_not_exists" json:"if_not_exists"`
		IfExists    bool `yaml:"if_exists" json:"if_exists"`
	} `yaml:"use"`
	Target DatabaseDriver
}

type CustomExportSpec struct {
	Template string `yaml:"template" json:"template"`
	Output   struct {
		Extension string `yaml:"extension" json:"extension"`
	} `yaml:"output" json:"output"`
}

func (s *CustomExportSpec) Validate() error {
	if s.Template == "" {
		return fmt.Errorf("custom export template is required")
	}

	return nil
}

func (m *MigrationsSpec) Validate() error {
	if m.Target == "" {
		m.Target = DatabaseDriverPostgres
		return nil
	}

	if !m.Target.Valid() {
		return fmt.Errorf(
			"target have unsupported driver %q. Available: %v",
			m.Target,
			writeableDatabaseDrivers,
		)
	}

	if !m.Target.CanMigrate() {
		return fmt.Errorf("target have driver %q, which unsupported migrate queries", m.Target)
	}

	return nil
}
