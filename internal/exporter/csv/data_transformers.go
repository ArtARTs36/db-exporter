package csv

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/db"
)

type DataTransformer func(data *transformingData, spec config.ExportSpecTransform) (*transformingData, error)

type transformingData struct {
	cols []string
	rows db.TableData
}

func SkipColumnsDataTransformer() DataTransformer {
	return func(data *transformingData, spec config.ExportSpecTransform) (*transformingData, error) {
		if len(spec.SkipColumns) == 0 {
			return data, nil
		}

		skipMap := map[string]bool{}
		for _, col := range spec.SkipColumns {
			skipMap[col] = true
		}

		return filterColumnDataTransformer(data, func(col string) bool {
			return !skipMap[col]
		}), nil
	}
}

func OnlyColumnsDataTransformer() DataTransformer {
	return func(data *transformingData, spec config.ExportSpecTransform) (*transformingData, error) {
		if len(spec.OnlyColumns) == 0 {
			return data, nil
		}

		onlyMap := map[string]bool{}
		for _, col := range spec.OnlyColumns {
			onlyMap[col] = true
		}

		return filterColumnDataTransformer(data, func(col string) bool {
			return onlyMap[col]
		}), nil
	}
}

func RenameColumnsDataTransformer() DataTransformer {
	return func(data *transformingData, spec config.ExportSpecTransform) (*transformingData, error) {
		if len(spec.RenameColumns) == 0 {
			return data, nil
		}

		if len(data.rows) == 0 {
			return data, nil
		}

		for col := range spec.RenameColumns {
			if _, exists := data.rows[0][col]; !exists {
				return nil, fmt.Errorf("column %q not found", col)
			}
		}

		cols := make([]string, 0)
		for _, col := range data.cols {
			if newCol, exists := spec.RenameColumns[col]; exists {
				cols = append(cols, newCol)
			} else {
				cols = append(cols, col)
			}
		}

		rows := make(db.TableData, len(data.rows))
		for i, row := range data.rows {
			newRow := make(map[string]interface{}, 0)

			for key, val := range row {
				newKey, newKeyExists := spec.RenameColumns[key]
				if !newKeyExists {
					newKey = key
				}

				newRow[newKey] = val
			}

			rows[i] = newRow
		}

		return &transformingData{
			cols: cols,
			rows: rows,
		}, nil
	}
}

func filterColumnDataTransformer(data *transformingData, filter func(col string) bool) *transformingData {
	cols := make([]string, 0)
	rows := make(db.TableData, len(data.rows))

	for _, col := range data.cols {
		if filter(col) {
			cols = append(cols, col)
		}
	}

	for i, row := range data.rows {
		newRow := make(map[string]interface{}, 0)
		for key, val := range row {
			if filter(key) {
				newRow[key] = val
			}
		}
		rows[i] = newRow
	}

	return &transformingData{
		cols: cols,
		rows: rows,
	}
}
