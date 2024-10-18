package yaml

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"

	"github.com/artarts36/db-exporter/internal/shared/ds"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	"gopkg.in/yaml.v3"

	"github.com/artarts36/db-exporter/internal/db"
)

const yamlFixturesFilename = "fixtures.yaml"

type FixturesExporter struct {
	dataLoader *db.DataLoader
	inserter   *db.Inserter
}

type yamlFixture struct {
	Options struct {
		Transaction bool `yaml:"transaction"`
	} `yaml:"options"`
	Tables *orderedmap.OrderedMap[string, *yamlFixtureTable] `yaml:"tables"`
}

type yamlFixtureTable struct {
	Options struct {
		Upsert bool `yaml:"upsert"`
	} `yaml:"options"`
	Rows []map[string]interface{} `yaml:"rows"`
}

func NewFixturesExporter(
	dataLoader *db.DataLoader,
	inserter *db.Inserter,
) *FixturesExporter {
	return &FixturesExporter{
		dataLoader: dataLoader,
		inserter:   inserter,
	}
}

func (e *FixturesExporter) ExportPerFile(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	for _, table := range params.Schema.Tables.List() {
		data, err := e.dataLoader.Load(ctx, params.Conn, table.Name.Value)
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

		p := &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s.yaml", table.Name.String()),
			Content:  content,
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (e *FixturesExporter) Export(
	ctx context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	fixture := &yamlFixture{
		Tables: orderedmap.New[string, *yamlFixtureTable](),
	}

	for _, table := range params.Schema.Tables.List() {
		data, err := e.dataLoader.Load(ctx, params.Conn, table.Name.Value)
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

	p := &exporter.ExportedPage{
		FileName: yamlFixturesFilename,
		Content:  content,
	}

	return []*exporter.ExportedPage{
		p,
	}, nil
}

func (e *FixturesExporter) Import(ctx context.Context, params *exporter.ImportParams) (
	[]exporter.ImportedFile,
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

	doImport := e.doImport
	if fixture.Options.Transaction {
		doImport = func(ctx context.Context, fixture *yamlFixture, params *exporter.ImportParams) (
			exporter.ImportedFile,
			error,
		) {
			var importedFile exporter.ImportedFile

			trErr := params.Conn.Transact(ctx, func(ctx context.Context) error {
				importedFile, err = e.doImport(ctx, fixture, params)
				if err != nil {
					return fmt.Errorf("transaction canceled: %w", err)
				}

				return nil
			})

			return importedFile, trErr
		}
	}

	importedFile, err := doImport(ctx, &fixture, params)
	if err != nil {
		return nil, err
	}

	return []exporter.ImportedFile{importedFile}, nil
}

func (e *FixturesExporter) doImport(
	ctx context.Context,
	fixture *yamlFixture,
	params *exporter.ImportParams,
) (exporter.ImportedFile, error) {
	affectedRows := map[string]int64{}

	for table := fixture.Tables.Oldest(); table != nil; table = table.Next() {
		if !params.TableFilter(table.Key) {
			continue
		}

		var ar int64
		var err error

		if table.Value.Options.Upsert && params.Schema.Tables.Has(*ds.NewString(table.Key)) {
			tbl, _ := params.Schema.Tables.Get(*ds.NewString(table.Key))
			ar, err = e.inserter.Upsert(ctx, params.Conn, tbl, table.Value.Rows)
			if err != nil {
				return exporter.ImportedFile{}, fmt.Errorf("failed to insert: %w", err)
			}
		} else {
			ar, err = e.inserter.Insert(ctx, params.Conn, table.Key, table.Value.Rows)
			if err != nil {
				return exporter.ImportedFile{}, fmt.Errorf("failed to insert: %w", err)
			}
		}

		affectedRows[table.Key] = ar
	}

	return exporter.ImportedFile{
		AffectedRows: affectedRows,
		Name:         yamlFixturesFilename,
	}, nil
}
