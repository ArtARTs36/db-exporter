package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Activity struct {
	Database     string         `yaml:"database" json:"database"`
	Tables       ActivityTables `yaml:"tables" json:"tables"`
	Format       ExporterName   `yaml:"format" json:"format"`
	TablePerFile bool           `yaml:"table_per_file" json:"table_per_file"`
	SkipExists   bool           `yaml:"skip_exists" json:"skip_exists"`
	Out          struct {
		Dir        string `yaml:"dir" json:"dir"`
		FilePrefix string `yaml:"file_prefix" json:"file_prefix"`
	} `yaml:"out" json:"out"`
	Spec interface{} `yaml:"-" json:"spec"`
}

type ActivityTables struct {
	List        stringOrStringSlice `yaml:"list" json:"list"`
	ListFromEnv string              `yaml:"from_env" json:"from_env"`
	Prefix      string              `yaml:"prefix" json:"prefix"`
}

func (s *Activity) UnmarshalYAML(n *yaml.Node) error {
	var embeddedActivity struct {
		Database     string         `yaml:"database" json:"database"`
		Tables       ActivityTables `yaml:"tables" json:"tables"`
		Format       ExporterName   `yaml:"format" json:"format"`
		TablePerFile bool           `yaml:"table_per_file" json:"table_per_file"`
		SkipExists   bool           `yaml:"skip_exists" json:"skip_exists"`
		Out          struct {
			Dir        string `yaml:"dir" json:"dir"`
			FilePrefix string `yaml:"file_prefix" json:"file_prefix"`
		} `yaml:"out" json:"out"`
		Spec yaml.Node `yaml:"spec" json:"spec"`
	}

	if err := n.Decode(&embeddedActivity); err != nil {
		return err
	}

	s.Database = embeddedActivity.Database
	s.Tables = embeddedActivity.Tables
	s.Format = embeddedActivity.Format
	s.TablePerFile = embeddedActivity.TablePerFile
	s.SkipExists = embeddedActivity.SkipExists
	s.Out = embeddedActivity.Out

	var serr error
	s.Spec, serr = s.newSpec(embeddedActivity.Format)
	if serr != nil {
		return fmt.Errorf("select spec: %w", serr)
	}

	if s.Spec != nil {
		err := embeddedActivity.Spec.Decode(s.Spec)
		if err != nil {
			return fmt.Errorf("failed to decode activity spec: %w", err)
		}
	}

	return nil
}

func (s *Activity) newSpec(format ExporterName) (interface{}, error) {
	var spec interface{}

	switch format {
	case ExporterNameGooseFixtures, ExporterNameGraphql, ExporterNameDBML, ExporterNameMermaid:
	case ExporterNameGoose, ExporterNameGoSQLMigrate, ExporterNameLaravelMigrationsRaw, ExporterNameDDL:
		spec = new(MigrationsSpec)
	case ExporterNameGoEntities:
		spec = new(GoEntitiesExportSpec)
	case ExporterNameMd:
		spec = new(MarkdownExportSpec)
	case ExporterNameGrpcCrud:
		spec = new(GRPCCrudExportSpec)
	case ExporterNameCSV:
		spec = new(CSVExportSpec)
	case ExporterNameLaravelModels:
		spec = new(LaravelModelsExportSpec)
	case ExporterNameGoEntityRepository:
		spec = new(GoEntityRepositorySpec)
	case ExporterNameJSONSchema:
		spec = new(JSONSchemaExportSpec)
	case ExporterNameCustom:
		spec = new(CustomExportSpec)
	case ExporterNameDiagram:
		spec = new(DiagramExportSpec)
	default:
		return nil, fmt.Errorf("format %q unsupported", format)
	}

	return spec, nil
}
