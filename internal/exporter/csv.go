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
	dataLoader *db.DataLoader
	renderer   *template.Renderer
}

func NewCSVExporter(
	dataLoader *db.DataLoader,
	renderer *template.Renderer,
) *CSVExporter {
	return &CSVExporter{
		dataLoader: dataLoader,
		renderer:   renderer,
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

		cols := table.ColumnsNames()

		colFilter := c.columnFilter(spec.TableColumn[table.Name.String()])
		if colFilter != nil {
			newCols := make([]string, 0)
			for _, col := range table.ColumnsNames() {
				if colFilter(col) {
					newCols = append(newCols, col)
				}
			}

			data = data.FilterColumns(colFilter)
		}

		p, err := render(
			c.renderer,
			"csv/export_single.csv",
			fmt.Sprintf("%s.csv", table.Name.String()),
			map[string]stick.Value{
				"rows":          data,
				"columns":       cols,
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

func (c *CSVExporter) columnFilter(cfg config.CSVExportSpecColumnFilter) func(col string) bool {
	if len(cfg.Only) > 0 {
		onlyMap := map[string]bool{}
		for _, col := range cfg.Only {
			onlyMap[col] = true
		}

		return func(col string) bool {
			return onlyMap[col]
		}
	}

	if len(cfg.Skip) > 0 {
		skipMap := map[string]bool{}
		for _, col := range cfg.Skip {
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
