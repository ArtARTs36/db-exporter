package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Activity struct {
	Export ExportActivity // fill export or import
	Import ImportActivity

	Database string         `yaml:"database" json:"database"`
	Tables   ActivityTables `yaml:"tables" json:"tables"`
}

type ActivityTables struct {
	List        stringOrStringSlice `yaml:"list" json:"list"`
	ListFromEnv string              `yaml:"from_env" json:"from_env"`
	Prefix      string              `yaml:"prefix" json:"prefix"`
}

type ExportActivity struct {
	Format       ExporterName
	TablePerFile bool `yaml:"table_per_file" json:"table_per_file"`
	SkipExists   bool `yaml:"skip_exists" json:"skip_exists"`
	Out          struct {
		Dir        string `yaml:"dir" json:"dir"`
		FilePrefix string `yaml:"file_prefix" json:"file_prefix"`
	} `yaml:"out" json:"out"`
	Spec interface{} `yaml:"-" json:"spec"`
}

type ImportActivity struct {
	Format ImporterName
	Spec   interface{} `yaml:"-"`
	From   string      `yaml:"from" json:"from"` // path to file or dir
}

func (s *Activity) IsExport() bool {
	return s.Export.Format != ""
}

func (s *Activity) UnmarshalYAML(n *yaml.Node) error {
	type exportOrImport struct {
		Export ExporterName `yaml:"export" json:"export"`
		Import ImporterName `yaml:"import" json:"import"`
		Spec   yaml.Node    `yaml:"spec" json:"spec"`

		Database    string         `yaml:"database" json:"database"`
		Tables      ActivityTables `yaml:"tables" json:"tables"`
		TablePrefix string         `yaml:"table_prefix" json:"table_prefix"`
	}

	exportOrImportObj := &exportOrImport{}
	if err := n.Decode(exportOrImportObj); err != nil {
		return err
	}

	s.Database = exportOrImportObj.Database
	s.Tables = exportOrImportObj.Tables

	var exportActivity ExportActivity
	var importActivity ImportActivity

	var decodingSpec interface{}

	if exportOrImportObj.Export != "" { //nolint:gocritic // no need
		if err := n.Decode(&exportActivity); err != nil {
			return err
		}

		exportActivity.Format = exportOrImportObj.Export

		switch exportActivity.Format {
		case ExporterNameDiagram, ExporterNameGooseFixtures, ExporterNameYamlFixtures, ExporterNameGraphql, ExporterNameDBML:
		case ExporterNameGoose, ExporterNameGoSQLMigrate, ExporterNameLaravelMigrationsRaw, ExporterNameDDL:
			exportActivity.Spec = new(MigrationsSpec)
		case ExporterNameGoEntities:
			exportActivity.Spec = new(GoEntitiesExportSpec)
		case ExporterNameMd:
			exportActivity.Spec = new(MarkdownExportSpec)
		case ExporterNameGrpcCrud:
			exportActivity.Spec = new(GRPCCrudExportSpec)
		case ExporterNameCSV:
			exportActivity.Spec = new(CSVExportSpec)
		case ExporterNameLaravelModels:
			exportActivity.Spec = new(LaravelModelsExportSpec)
		case ExporterNameGoEntityRepository:
			exportActivity.Spec = new(GoEntityRepositorySpec)
		case ExporterNameJSONSchema:
			exportActivity.Spec = new(JSONSchemaExportSpec)
		case ExporterNameCustom:
			exportActivity.Spec = new(CustomExportSpec)
		default:
			return fmt.Errorf("format %q unsupported", exportActivity.Format)
		}

		decodingSpec = exportActivity.Spec
	} else if exportOrImportObj.Import != "" {
		if err := n.Decode(&importActivity); err != nil {
			return err
		}

		importActivity.Format = exportOrImportObj.Import

		switch exportOrImportObj.Import {
		case ImporterNameYamlFixtures:
		default:
			return fmt.Errorf("format %q unsupported", importActivity.Format)
		}

		decodingSpec = importActivity.Spec
	} else {
		return fmt.Errorf("invalid activity: must be one of export or import")
	}

	s.Export = exportActivity
	s.Import = importActivity

	if decodingSpec != nil {
		err := exportOrImportObj.Spec.Decode(decodingSpec)
		if err != nil {
			return fmt.Errorf("failed to decode activity spec: %w", err)
		}
	}

	return nil
}
