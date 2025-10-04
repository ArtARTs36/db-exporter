package csv

import (
	"fmt"
	"strconv"
	"strings"
)

type generator struct{}

func (c *generator) generate(data *transformingData, columnDelimiter string) (string, error) {
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
				_, err := sb.WriteString(columnDelimiter)
				if err != nil {
					return "", fmt.Errorf("write string: %w", err)
				}
			}
		}
	}

	return sb.String(), nil
}

func (*generator) mapValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", value)
	}
}
