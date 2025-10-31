package csv

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/infrastructure/data"
)

type DataTransformer func(data *transformingData, spec SpecificationTransform) (*transformingData, error)

type transformingData struct {
	cols []string
	rows data.TableData
}

func SkipColumnsDataTransformer() DataTransformer {
	return func(trData *transformingData, spec SpecificationTransform) (*transformingData, error) {
		if len(spec.SkipColumns) == 0 {
			return trData, nil
		}

		skipMap := map[string]bool{}
		for _, col := range spec.SkipColumns {
			skipMap[col] = true
		}

		return filterColumnDataTransformer(trData, func(col string) bool {
			return !skipMap[col]
		}), nil
	}
}

func OnlyColumnsDataTransformer() DataTransformer {
	return func(data *transformingData, spec SpecificationTransform) (*transformingData, error) {
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
	return func(trData *transformingData, spec SpecificationTransform) (*transformingData, error) {
		if len(spec.RenameColumns) == 0 {
			return trData, nil
		}

		if len(trData.rows) == 0 {
			return trData, nil
		}

		for col := range spec.RenameColumns {
			if _, exists := trData.rows[0][col]; !exists {
				return nil, fmt.Errorf("column %q not found", col)
			}
		}

		cols := make([]string, 0)
		for _, col := range trData.cols {
			if newCol, exists := spec.RenameColumns[col]; exists {
				cols = append(cols, newCol)
			} else {
				cols = append(cols, col)
			}
		}

		rows := make(data.TableData, len(trData.rows))
		for i, row := range trData.rows {
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

func filterColumnDataTransformer(trData *transformingData, filter func(col string) bool) *transformingData {
	cols := make([]string, 0)
	rows := make(data.TableData, len(trData.rows))

	for _, col := range trData.cols {
		if filter(col) {
			cols = append(cols, col)
		}
	}

	for i, row := range trData.rows {
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
