package csv

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type generator struct{}

func (c *generator) generate(data *transformingData, columnDelimiter string) string {
	sb := strings.Builder{}

	for i, col := range data.cols {
		sb.WriteString(col)
		if i < len(data.cols)-1 {
			sb.WriteString(columnDelimiter)
		}
	}

	for _, row := range data.rows {
		sb.WriteString("\n")

		for colID, col := range data.cols {
			val := "\"\""
			if v, ok := row[col]; ok {
				val = c.mapValue(v)
			}

			sb.WriteString(val)
			if colID < len(data.cols)-1 {
				sb.WriteString(columnDelimiter)
			}
		}
	}

	return sb.String()
}

func (*generator) mapValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case time.Time:
		return v.Format(time.RFC3339Nano)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", value)
	}
}
