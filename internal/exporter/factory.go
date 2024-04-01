package exporter

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/schemaloader"
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
}

func CreateExporter(name string, renderer *template.Renderer, connection *schemaloader.Connection) (Exporter, error) {
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
			schemaloader.NewDataLoader(connection),
			renderer,
			sql.NewInsertBuilder(),
		), nil
	}

	return nil, fmt.Errorf("format %q unsupported", name)
}
