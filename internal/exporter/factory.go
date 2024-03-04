package exporter

import (
	"fmt"

	"github.com/artarts36/db-exporter/internal/template"
)

func CreateExporter(name string, renderer *template.Renderer) (Exporter, error) {
	if name == "md" || name == "markdown" {
		return NewMarkdownExporter(renderer), nil
	}

	return nil, fmt.Errorf("format %q unsupported", name)
}
