package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Activity struct {
	Export       ExporterName `yaml:"export"` // fill export or import
	Import       ImporterName `yaml:"import"`
	Database     string       `yaml:"database"`
	TablePerFile bool         `yaml:"table_per_file"`
	Tables       []string     `yaml:"tables"`
	Out          struct {
		Dir        string `yaml:"dir"`
		FilePrefix string `yaml:"file_prefix"`
	} `yaml:"out"`
	Spec interface{} `yaml:"-"`
}

func (s *Activity) IsExport() bool {
	return s.Export != ""
}

func (s *Activity) UnmarshalYAML(n *yaml.Node) error {
	type act Activity
	type T struct {
		*act `yaml:",inline"`
		Spec yaml.Node `yaml:"spec"`
	}

	obj := &T{act: (*act)(s)}
	if err := n.Decode(obj); err != nil {
		return err
	}

	if s.Export != "" {
		switch s.Export {
		case ExporterNameDiagram, ExporterNameGoose, ExporterNameGooseFixtures, ExporterNameGoSQLMigrate,
			ExporterNameLaravelMigrationsRaw, ExporterNameYamlFixtures:
			return nil
		case ExporterNameGoStructs:
			s.Spec = new(GoStructsExportSpec)
		case ExporterNameMd:
			s.Spec = new(MarkdownExportSpec)
		case ExporterNameGrpcCrud:
			s.Spec = new(GoStructsExportSpec)
		default:
			return fmt.Errorf("format %q unsupported", s.Export)
		}
	} else if s.Import != "" {
		switch s.Import {
		case ImporterNameYamlFixtures:
			return nil
		default:
			return fmt.Errorf("format %q unsupported", s.Export)
		}
	}

	return obj.Spec.Decode(s.Spec)
}
