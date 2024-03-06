package exporter

import (
	"fmt"

	"github.com/artarts36/db-exporter/internal/template"
)

var Names = []string{
	"md",
	"diagram",
	"go-structs",
	"goose",
	"laravel-migrations-raw",
}

func CreateExporter(name string, renderer *template.Renderer) (Exporter, error) {
	if name == "md" {
		return NewMarkdownExporter(renderer), nil
	}

	if name == "diagram" {
		return NewDiagramExporter(renderer), nil
	}

	if name == "go-structs" {
		return NewGoStructsExporter(renderer), nil
	}

	if name == "goose" {
		return NewGooseExporter(renderer), nil
	}

	if name == "laravel-migrations-raw" {
		return NewLaravelMigrationsExporter(renderer), nil
	}

	return nil, fmt.Errorf("format %q unsupported", name)
}
