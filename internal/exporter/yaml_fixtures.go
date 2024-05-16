package exporter

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
	"gopkg.in/yaml.v3"
)

const YamlFixturesExporterName = "yaml-fixtures"

type YamlFixturesExporter struct {
	dataLoader *db.DataLoader
	renderer   *template.Renderer
}

func NewYamlFixturesExporter(
	dataLoader *db.DataLoader,
	renderer *template.Renderer,
) *YamlFixturesExporter {
	return &YamlFixturesExporter{
		dataLoader: dataLoader,
		renderer:   renderer,
	}
}

func (e *YamlFixturesExporter) ExportPerFile(
	ctx context.Context,
	sch *schema.Schema,
	_ *ExportParams,
) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, sch.Tables.Len())

	for _, table := range sch.Tables.List() {
		data, err := e.dataLoader.Load(ctx, table.Name.Value)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			continue
		}

		content, err := yaml.Marshal(map[string]interface{}{
			"rows": data,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to build yaml content: %w", err)
		}

		p := &ExportedPage{
			FileName: fmt.Sprintf("%s.yaml", table.Name.String()),
			Content:  content,
		}

		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (e *YamlFixturesExporter) Export(
	ctx context.Context,
	sch *schema.Schema,
	_ *ExportParams,
) ([]*ExportedPage, error) {
	tablesData := map[string]map[string]interface{}{}

	for _, table := range sch.Tables.List() {
		data, err := e.dataLoader.Load(ctx, table.Name.Value)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			continue
		}

		_, tableDataExists := tablesData[table.Name.Value]
		if !tableDataExists {
			tablesData[table.Name.Value] = map[string]interface{}{}
		}

		tablesData[table.Name.Value]["rows"] = data
	}

	content, err := yaml.Marshal(map[string]interface{}{
		"tables": tablesData,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to build yaml content: %w", err)
	}

	p := &ExportedPage{
		FileName: "fixtures.yaml",
		Content:  content,
	}

	return []*ExportedPage{
		p,
	}, nil
}
