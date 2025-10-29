package csv

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/infrastructure/data"
	"github.com/artarts36/db-exporter/internal/infrastructure/workspace"
	"github.com/artarts36/db-exporter/internal/schema"
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
		tableData, err := c.dataLoader.Load(ctx, params.Conn, table)
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

				err = c.transform(table, trData, spec)
				if err != nil {
					return err
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

func (c *Exporter) transform(table *schema.Table, data *transformingData, spec *config.CSVExportSpec) error {
	if len(spec.Transform) == 0 {
		return nil
	}

	for _, transformer := range c.dataTransformers {
		for _, transformSpec := range spec.Transform[table.Name.Value] {
			var err error
			data, err = transformer(data, transformSpec)
			if err != nil {
				return fmt.Errorf("transform table data: %w", err)
			}
		}
	}
	return nil
}

func (c *Exporter) Export(ctx context.Context, params *exporter.ExportParams) ([]*exporter.ExportedPage, error) {
	return c.ExportPerFile(ctx, params)
}
