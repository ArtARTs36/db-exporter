package exporter

import (
	"context"
	"errors"
	"fmt"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/db"
	"github.com/artarts36/db-exporter/internal/template"
)

type CSVExporter struct {
	dataLoader       *db.DataLoader
	renderer         *template.Renderer
	dataTransformers []DataTransformer
}

func NewCSVExporter(
	dataLoader *db.DataLoader,
	renderer *template.Renderer,
	dataTransformers []DataTransformer,
) *CSVExporter {
	return &CSVExporter{
		dataLoader:       dataLoader,
		renderer:         renderer,
		dataTransformers: dataTransformers,
	}
}

func (c *CSVExporter) ExportPerFile(ctx context.Context, params *ExportParams) ([]*ExportedPage, error) {
	spec, ok := params.Spec.(*config.CSVExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	pages := make([]*ExportedPage, 0, params.Schema.Tables.Len())

	delimiter := spec.Delimiter
	if delimiter == "" {
		delimiter = ","
	}

	for _, table := range params.Schema.Tables.List() {
		data, err := c.dataLoader.Load(ctx, params.Conn, table.Name.Value)
		if err != nil {
			return nil, err
		}
		if len(data) == 0 {
			continue
		}

		trData := &transformingData{
			cols: table.ColumnsNames(),
			rows: data,
		}

		if len(spec.Transform) > 0 {
			for _, transformer := range c.dataTransformers {
				for _, transformSpec := range spec.Transform[table.Name.Value] {
					trData, err = transformer(trData, transformSpec)
					if err != nil {
						return nil, fmt.Errorf("failed to transform data: %w", err)
					}
				}
			}
		}

		p, err := render(
			c.renderer,
			"csv/export_single.csv",
			fmt.Sprintf("%s.csv", table.Name.String()),
			map[string]stick.Value{
				"rows":          trData.rows,
				"columns":       trData.cols,
				"col_delimiter": delimiter,
			},
		)
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (c *CSVExporter) columnFilter(cfg config.ExportSpecTransform) func(col string) bool {
	if len(cfg.OnlyColumns) > 0 {
		onlyMap := map[string]bool{}
		for _, col := range cfg.OnlyColumns {
			onlyMap[col] = true
		}

		return func(col string) bool {
			return onlyMap[col]
		}
	}

	if len(cfg.OnlyColumns) > 0 {
		skipMap := map[string]bool{}
		for _, col := range cfg.OnlyColumns {
			skipMap[col] = true
		}

		return func(col string) bool {
			return !skipMap[col]
		}
	}

	return nil
}

func (c *CSVExporter) Export(ctx context.Context, params *ExportParams) ([]*ExportedPage, error) {
	return c.ExportPerFile(ctx, params)
}
