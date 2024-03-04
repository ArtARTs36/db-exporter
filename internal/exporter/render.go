package exporter

import (
	"fmt"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/template"
)

func render(renderer *template.Renderer, tplName, fileName string, params map[string]stick.Value) (*ExportedPage, error) {
	indexContent, err := renderer.Render(tplName, params)
	if err != nil {
		return nil, fmt.Errorf("unable to render template %q: %w", tplName, err)
	}

	return &ExportedPage{
		FileName: fileName,
		Content:  indexContent,
	}, nil
}
