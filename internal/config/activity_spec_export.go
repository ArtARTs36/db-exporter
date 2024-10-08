package config

import orderedmap "github.com/wk8/go-ordered-map/v2"

type ExporterName string

const (
	ExporterNameMd                   ExporterName = "md"
	ExporterNameDiagram              ExporterName = "diagram"
	ExporterNameGoStructs            ExporterName = "go-structs"
	ExporterNameGoose                ExporterName = "goose"
	ExporterNameGooseFixtures        ExporterName = "goose-fixtures"
	ExporterNameGoSQLMigrate         ExporterName = "go-sql-migrate"
	ExporterNameLaravelMigrationsRaw ExporterName = "laravel-migrations-raw"
	ExporterNameLaravelModels        ExporterName = "laravel-models"
	ExporterNameGrpcCrud             ExporterName = "grpc-crud"
	ExporterNameYamlFixtures         ExporterName = "yaml-fixtures"
	ExporterNameCSV                  ExporterName = "csv"
)

type GoStructsExportSpec struct {
	Package string `yaml:"package"`
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
