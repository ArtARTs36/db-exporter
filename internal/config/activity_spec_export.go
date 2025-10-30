package config

import (
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/schema"
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

type MigrationsSpec struct {
	Use struct {
		IfNotExists bool `yaml:"if_not_exists" json:"if_not_exists"`
		IfExists    bool `yaml:"if_exists" json:"if_exists"`
	} `yaml:"use"`
	Target schema.DatabaseDriver `yaml:"target" json:"target"`
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

func (m *MigrationsSpec) InjectDatabaseDriver(driver schema.DatabaseDriver) {
	if m.Target != "" {
		return
	}

	m.Target = driver
}

func (m *MigrationsSpec) Validate() error {
	if m.Target == "" {
		return errors.New("target is required")
	}

	if !m.Target.Valid() {
		return fmt.Errorf(
			"target have unsupported driver %q. Available: %v",
			m.Target,
			schema.GetWriteableDatabaseDrivers(),
		)
	}

	if !m.Target.CanMigrate() {
		return fmt.Errorf("target have driver %q, which unsupported migrate queries", m.Target)
	}

	return nil
}
