package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Activity struct {
	Export ExportActivity // fill export or import
	Import ImportActivity

	Database string   `yaml:"database"`
	Tables   []string `yaml:"tables"`
}

type ExportActivity struct {
	Format       ExporterName
	TablePerFile bool `yaml:"table_per_file"`
	SkipExists   bool `yaml:"skip_exists"` // Skip generate already exists files
	Out          struct {
		Dir        string `yaml:"dir"`
		FilePrefix string `yaml:"file_prefix"`
	} `yaml:"out"`
	Spec interface{} `yaml:"-"`
}

type ImportActivity struct {
	Format ImporterName
	Spec   interface{} `yaml:"-"`
	From   string      `yaml:"from"` // path to file or dir
}

func (s *Activity) IsExport() bool {
	return s.Export.Format != ""
}

func (s *Activity) UnmarshalYAML(n *yaml.Node) error {
	type exportOrImport struct {
		Export ExporterName `yaml:"export"`
		Import ImporterName `yaml:"import"`
		Spec   yaml.Node    `yaml:"spec"`

		Database string   `yaml:"database"`
		Tables   []string `yaml:"tables"`
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
		case ExporterNameDiagram, ExporterNameGoose, ExporterNameGooseFixtures, ExporterNameGoSQLMigrate,
			ExporterNameLaravelMigrationsRaw, ExporterNameYamlFixtures:
		case ExporterNameGoStructs:
			exportActivity.Spec = new(GoStructsExportSpec)
		case ExporterNameMd:
			exportActivity.Spec = new(MarkdownExportSpec)
		case ExporterNameGrpcCrud:
			exportActivity.Spec = new(GRPCCrudExportSpec)
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
