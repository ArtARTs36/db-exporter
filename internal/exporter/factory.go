package exporter

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/sql"

	"github.com/artarts36/db-exporter/internal/template"
)

var Names = []string{
	MarkdownExporterName,
	DiagramExporterName,
	GoStructsExporterName,
	GooseExporterName,
	GoSQLMigrateExporterName,
	LaravelMigrationsRawExporterName,
	GrpcCrudExporterName,
	GooseFixturesExporterName,
	YamlFixturesExporterName,
}

func CreateExporter(name string, renderer *template.Renderer, connection *db.Connection) (Exporter, error) {
	if name == MarkdownExporterName {
		return NewMarkdownExporter(renderer), nil
	}

	if name == DiagramExporterName {
		return NewDiagramExporter(renderer), nil
	}

	if name == GoStructsExporterName {
		return NewGoStructsExporter(renderer), nil
	}

	if name == GooseExporterName {
		return NewGooseExporter(renderer, sql.NewDDLBuilder()), nil
	}

	if name == GoSQLMigrateExporterName {
		return NewSQLMigrateExporter(renderer, sql.NewDDLBuilder()), nil
	}

	if name == LaravelMigrationsRawExporterName {
		return NewLaravelMigrationsRawExporter(renderer, sql.NewDDLBuilder()), nil
	}

	if name == GrpcCrudExporterName {
		return NewGrpcCrudExporter(renderer), nil
	}

	if name == GooseFixturesExporterName {
		return NewGooseFixturesExporter(
			db.NewDataLoader(connection),
			renderer,
			sql.NewInsertBuilder(),
		), nil
	}

	if name == YamlFixturesExporterName {
		return NewYamlFixturesExporter(
			db.NewDataLoader(connection),
			renderer,
			db.NewInserter(connection),
			connection,
		), nil
	}

	return nil, fmt.Errorf("format %q unsupported", name)
}
