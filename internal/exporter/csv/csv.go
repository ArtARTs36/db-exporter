package csv

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/infrastructure/data"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/config"
)

type Exporter struct {
	dataLoader       *data.Loader
	pager            *common.Pager
	dataTransformers []DataTransformer
}

func NewExporter(
	dataLoader *data.Loader,
	pager *common.Pager,
	dataTransformers []DataTransformer,
) *Exporter {
	return &Exporter{
		dataLoader:       dataLoader,
		pager:            pager,
		dataTransformers: dataTransformers,
	}
}

func (c *Exporter) ExportPerFile(ctx context.Context, params *exporter.ExportParams) ([]*exporter.ExportedPage, error) {
	spec, ok := params.Spec.(*config.CSVExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	delimiter := spec.Delimiter
	if delimiter == "" {
		delimiter = ","
	}

	csvPage := c.pager.Of("csv/export_single.csv")

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

		p, err := csvPage.Export(fmt.Sprintf("%s.csv", table.Name.String()), map[string]stick.Value{
			"rows":          trData.rows,
			"columns":       trData.cols,
			"col_delimiter": delimiter,
		})
		if err != nil {
			return nil, err
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (c *Exporter) Export(ctx context.Context, params *exporter.ExportParams) ([]*exporter.ExportedPage, error) {
	return c.ExportPerFile(ctx, params)
}
