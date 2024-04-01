package sql

import "fmt"

func (b *QueryBuilder) mapValue(val interface{}) string {
	colValStr := "null"

	switch tColVal := val.(type) {
	case string:
		colValStr = fmt.Sprintf("'%s'", tColVal)
	case bool:
		if tColVal {
			colValStr = "true"
		} else {
			colValStr = "false"
		}
	case int, int8, int16, int32, int64, uint, uint8, uint32, uint64:
		colValStr = fmt.Sprintf("%d", tColVal)
	}

	return colValStr
}
