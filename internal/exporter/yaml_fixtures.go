package exporter

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/ds"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/yaml.v3"

	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

const (
	YamlFixturesExporterName = "yaml-fixtures"
	yamlFixturesFilename     = "fixtures.yaml"
)

type YamlFixturesExporter struct {
	unimplementedImporter
	dataLoader *db.DataLoader
	renderer   *template.Renderer
	inserter   *db.Inserter
}

type yamlFixture struct {
	Tables *orderedmap.OrderedMap[string, *yamlFixtureTable] `yaml:"tables"`
}

type yamlFixtureTable struct {
	Options struct {
		Upsert bool `yaml:"upsert"`
	} `yaml:"options"`
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
		data, err := e.dataLoader.Load(ctx, table.Name.Value)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			continue
		}

		fixture := yamlFixture{
			Tables: orderedmap.New[string, *yamlFixtureTable](),
		}

		fixture.Tables.Set(table.Name.Value, &yamlFixtureTable{Rows: data})

		content, err := yaml.Marshal(fixture)

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
	fixture := &yamlFixture{
		Tables: orderedmap.New[string, *yamlFixtureTable](),
	}

	for _, table := range sch.Tables.List() {
		data, err := e.dataLoader.Load(ctx, table.Name.Value)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			continue
		}

		fixture.Tables.Set(table.Name.Value, &yamlFixtureTable{
			Rows: data,
		})
	}

	content, err := yaml.Marshal(fixture)
	if err != nil {
		return nil, fmt.Errorf("failed to build yaml content: %w", err)
	}

	p := &ExportedPage{
		FileName: yamlFixturesFilename,
		Content:  content,
	}

	return []*ExportedPage{
		p,
	}, nil
}

func (e *YamlFixturesExporter) Import(ctx context.Context, sch *schema.Schema, params *ImportParams) (
	[]ImportedFile,
	error,
) {
	file, err := params.Directory.ReadFile(yamlFixturesFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", yamlFixturesFilename, err)
	}

	var fixture yamlFixture

	if err = yaml.Unmarshal(file, &fixture); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", yamlFixturesFilename, err)
	}

	for table := fixture.Tables.Oldest(); table != nil; table = table.Next() {
		if !params.TableFilter(table.Key) {
			continue
		}

		if table.Value.Options.Upsert && sch.Tables.Has(*ds.NewString(table.Key)) {
			tbl, _ := sch.Tables.Get(*ds.NewString(table.Key))
			err = e.inserter.Upsert(ctx, tbl, table.Value.Rows)
			if err != nil {
				return nil, fmt.Errorf("failed to insert: %w", err)
			}
		} else {
			err = e.inserter.Insert(ctx, table.Key, table.Value.Rows)
			if err != nil {
				return nil, fmt.Errorf("failed to insert: %w", err)
			}
		}
	}

	return []ImportedFile{
		{
			Name: yamlFixturesFilename,
		},
	}, nil
}
