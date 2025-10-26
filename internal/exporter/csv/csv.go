package csv

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/infrastructure/data"
	"github.com/artarts36/db-exporter/internal/infrastructure/workspace"
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
	spec, ok := params.Spec.(*config.CSVExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	for _, table := range params.Schema.Tables.List() {
		tableData, err := c.dataLoader.Load(ctx, params.Conn, table.Name.Value)
		if err != nil {
			return nil, fmt.Errorf("load data from table %q: %w", table.Name.Value, err)
		}
		if len(tableData) == 0 {
			continue
		}

		err = params.Workspace.Write(
			ctx,
			fmt.Sprintf("%s.csv", table.Name.String()),
			func(buffer workspace.Buffer) error {
				trData := &transformingData{
					cols: table.ColumnsNames(),
					rows: tableData,
				}

				if len(spec.Transform) > 0 {
					for _, transformer := range c.dataTransformers {
						for _, transformSpec := range spec.Transform[table.Name.Value] {
							trData, err = transformer(trData, transformSpec)
							if err != nil {
								return fmt.Errorf("failed to transform tableData: %w", err)
							}
						}
					}
				}

				c.generator.generate(trData, spec.Delimiter, buffer)

				return nil
			})
		if err != nil {
			return nil, fmt.Errorf("write table data to workspace: %w", err)
		}
	}

	return []*exporter.ExportedPage{}, nil
}

func (c *Exporter) Export(ctx context.Context, params *exporter.ExportParams) ([]*exporter.ExportedPage, error) {
	return c.ExportPerFile(ctx, params)
}
