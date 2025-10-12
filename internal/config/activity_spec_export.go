package config

import (
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/imageencoder"
	"github.com/artarts36/db-exporter/internal/shared/webcolor"
	"github.com/artarts36/specw"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"golang.org/x/image/colornames"
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

type GRPCCrudExportSpec struct {
	Package string                                     `yaml:"package" json:"package"`
	Options orderedmap.OrderedMap[string, interface{}] `yaml:"options" json:"options"`
}

type MarkdownExportSpec struct {
	WithDiagram bool              `yaml:"with_diagram" json:"with_diagram"`
	Diagram     DiagramExportSpec `yaml:"diagram" json:"diagram"`
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
	Target DatabaseDriver `yaml:"target" json:"target"`
}

type DiagramExportSpec struct {
	Image struct {
		Format      imageencoder.Format           `yaml:"format" json:"format"`
		Compression imageencoder.CompressionLevel `yaml:"compression" json:"compression"`
	} `yaml:"image" json:"image"`
	Style struct {
		Background struct {
			Grid *struct {
				LineColor *specw.Color `yaml:"line_color" json:"line_color"`
				CellSize  int          `yaml:"cell_size" json:"cell_size"`
			} `yaml:"grid" json:"grid"`
			Color *specw.Color `yaml:"color" json:"color"`
		} `yaml:"background" json:"background"`
		Table struct {
			Name struct {
				BackgroundColor string `yaml:"background_color" json:"background_color"` // #hex
				TextColor       string `yaml:"text_color" json:"text_color"`             // #hex
			} `yaml:"name" json:"name"`
		} `yaml:"table" json:"table"`
		Font struct {
			Family string  `yaml:"family" json:"family"`
			Size   float64 `yaml:"size" json:"size"`
		} `yaml:"font" json:"font"`
	} `yaml:"style" json:"style"`
}

type CustomExportSpec struct {
	Template string `yaml:"template" json:"template"`
	Output   struct {
		Extension string `yaml:"extension" json:"extension"`
	} `yaml:"output" json:"output"`
}

func (s *DiagramExportSpec) Validate() error {
	const (
		defaultGridCellSize = 30
		defaultFontSize     = 32
	)

	if !s.Image.Format.Valid() {
		return fmt.Errorf("unknown image format: %s", s.Image.Format)
	}

	if s.Image.Format == imageencoder.FormatUnspecified {
		s.Image.Format = imageencoder.FormatPNG
	}

	if s.Image.Compression == imageencoder.CompressionLevelUnspecified {
		s.Image.Compression = imageencoder.CompressionLevelLow
	}

	if s.Style.Table.Name.BackgroundColor == "" {
		s.Style.Table.Name.BackgroundColor = "#3498db"
	}

	if s.Style.Table.Name.TextColor == "" {
		s.Style.Table.Name.TextColor = "white"
	}

	if s.Style.Background.Color == nil {
		s.Style.Background.Color = &specw.Color{
			Color: colornames.White,
		}
	}

	if s.Style.Background.Grid != nil {
		if s.Style.Background.Grid.LineColor == nil {
			s.Style.Background.Grid.LineColor = &specw.Color{
				Color: webcolor.ColorEEE,
			}
		}
		if s.Style.Background.Grid.CellSize == 0 {
			s.Style.Background.Grid.CellSize = defaultGridCellSize
		}
	}

	if s.Style.Font.Size == 0 {
		s.Style.Font.Size = defaultFontSize
	}

	if s.Style.Font.Family == "" {
		s.Style.Font.Family = "Arial"
	}

	return nil
}

func (s *CustomExportSpec) Validate() error {
	if s.Template == "" {
		return fmt.Errorf("custom export template is required")
	}

	return nil
}

func (m *MigrationsSpec) InjectDatabaseDriver(driver DatabaseDriver) {
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
			writeableDatabaseDrivers,
		)
	}

	if !m.Target.CanMigrate() {
		return fmt.Errorf("target have driver %q, which unsupported migrate queries", m.Target)
	}

	return nil
}

func (s *MarkdownExportSpec) Validate() error {
	if err := s.Diagram.Validate(); err != nil {
		return fmt.Errorf("diagram: %w", err)
	}
	return nil
}
