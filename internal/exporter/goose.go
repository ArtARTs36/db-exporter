package exporter

import (
	"context"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

type GooseExporter struct {
	renderer *template.Renderer
}

func NewGooseExporter(renderer *template.Renderer) *GooseExporter {
	return &GooseExporter{
		renderer: renderer,
	}
}

func (e *GooseExporter) ExportPerFile(_ context.Context, sc *schema.Schema, params *ExportParams) ([]*ExportedPage, error) {
	return nil, nil
}

func (e *GooseExporter) Export(_ context.Context, schema *schema.Schema, _ *ExportParams) ([]*ExportedPage, error) {
	upQueries := make([]string, 0, len(schema.Tables))
	downQueries := make([]string, 0, len(schema.Tables))
	for _, table := range schema.Tables {
		upQuery := sqlbuilder.CreateTable(table.Name.Value)

		for _, column := range table.Columns {
			defs := []string{
				column.Name.Value,
				column.Type.Value,
			}

			upQuery.Define(defs...)
		}

		upQueries = append(upQueries, upQuery.String())
		downQueries = append(downQueries, fmt.Sprintf("DROP TABLE %s", table.Name.Value))
	}

	p, err := render(e.renderer, "goose/single.sql", "init.sql", map[string]stick.Value{
		"up_queries":   upQueries,
		"down_queries": downQueries,
	})
	if err != nil {
		return nil, err
	}

	return []*ExportedPage{
		p,
	}, nil
}
