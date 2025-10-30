package csv

import (
	"context"
	"errors"
	"fmt"

	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/infrastructure/data"
)

type Exporter struct {
	dataLoader       *data.Loader
	dataTransformers []DataTransformer
	generator        *generator
}

func NewExporter(
	dataLoader *data.Loader,
	dataTransformers []DataTransformer,
) *Exporter {
	return &Exporter{
		dataLoader:       dataLoader,
		dataTransformers: dataTransformers,
	}
}

func (c *Exporter) ExportPerFile(ctx context.Context, params *exporter.ExportParams) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*Specification)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	delimiter := spec.Delimiter
	if delimiter == "" {
		delimiter = ","
	}

	for _, table := range params.Schema.Tables.List() {
		tableData, err := c.dataLoader.Load(ctx, params.Conn, table)
		if err != nil {
			return nil, fmt.Errorf("load data from table %q: %w", table.Name.Value, err)
		}
		if len(tableData) == 0 {
			continue
		}

		trData := &transformingData{
			cols: table.ColumnsNames(),
			rows: tableData,
		}

		if len(spec.Transform) > 0 {
			for _, transformer := range c.dataTransformers {
				for _, transformSpec := range spec.Transform[table.Name.Value] {
					trData, err = transformer(trData, transformSpec)
					if err != nil {
						return nil, fmt.Errorf("failed to transform tableData: %w", err)
					}
				}
			}
		}

		pageContent := c.generator.generate(trData, delimiter)

		p := &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s.csv", table.Name.String()),
			Content:  []byte(pageContent),
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (c *Exporter) Export(ctx context.Context, params *exporter.ExportParams) ([]*exporter.ExportedPage, error) {
	return c.ExportPerFile(ctx, params)
}
