package exporter

import (
	"context"
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/yaml.v3"

	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

const YamlFixturesExporterName = "yaml-fixtures"

type YamlFixturesExporter struct {
	unimplementedImporter
	dataLoader *db.DataLoader
	renderer   *template.Renderer
	inserter   *db.Inserter
}

type yamlFixture struct {
	Tables orderedmap.OrderedMap[string, yamlFixtureTable] `yaml:"tables"`
}

type yamlFixtureTable struct {
	Rows []map[string]interface{} `yaml:"rows"`
}

func NewYamlFixturesExporter(
	dataLoader *db.DataLoader,
	renderer *template.Renderer,
	inserter *db.Inserter,
) *YamlFixturesExporter {
	return &YamlFixturesExporter{
		dataLoader: dataLoader,
		renderer:   renderer,
		inserter:   inserter,
	}
}

func (e *YamlFixturesExporter) ExportPerFile(
	ctx context.Context,
	sch *schema.Schema,
	_ *ExportParams,
) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, sch.Tables.Len())

	for _, table := range sch.Tables.List() {
		data, err := e.dataLoader.Load(ctx, table.Name.Val)
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
		data, err := e.dataLoader.Load(ctx, table.Name.Val)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			continue
		}

		_, tableDataExists := tablesData[table.Name.Val]
		if !tableDataExists {
			tablesData[table.Name.Val] = map[string]interface{}{}
		}

		tablesData[table.Name.Val]["rows"] = data
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

func (e *YamlFixturesExporter) Import(ctx context.Context, _ *schema.Schema, params *ExportParams) error {
	file, err := params.Directory.ReadFile("fixtures.yaml")
	if err != nil {
		return fmt.Errorf("failed to read fixtures.yaml: %w", err)
	}

	var fixture yamlFixture

	if err = yaml.Unmarshal(file, &fixture); err != nil {
		return fmt.Errorf("failed to unmarshal fixtures.yaml: %w", err)
	}

	for table := fixture.Tables.Oldest(); table != nil; table = table.Next() {
		err = e.inserter.Insert(ctx, table.Key, table.Value.Rows)
		if err != nil {
			return fmt.Errorf("failed to insert: %w", err)
		}
	}

	return nil
}
