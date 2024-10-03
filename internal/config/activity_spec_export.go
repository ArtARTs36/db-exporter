package config

type ExporterName string

const (
	ExporterNameMd                   ExporterName = "md"
	ExporterNameDiagram              ExporterName = "diagram"
	ExporterNameGoStructs            ExporterName = "go-structs"
	ExporterNameGoose                ExporterName = "goose"
	ExporterNameGooseFixtures        ExporterName = "goose-fixtures"
	ExporterNameGoSQLMigrate         ExporterName = "go-sql-migrate"
	ExporterNameLaravelMigrationsRaw ExporterName = "laravel-migrations-raw"
	ExporterNameGrpcCrud             ExporterName = "grpc-crud"
	ExporterNameYamlFixtures         ExporterName = "yaml-fixtures"
)

type GoStructsExportSpec struct {
	Package string `yaml:"package"`
}

type GRPCCrudExportSpec struct {
	GoPackage      string `yaml:"go_package"`
	ProtoGoPackage string `yaml:"proto_go_package"`
}

type MarkdownExportSpec struct {
	WithDiagram bool `yaml:"with_diagram"`
}
