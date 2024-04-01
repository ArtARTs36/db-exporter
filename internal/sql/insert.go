package sql

import (
	"fmt"
	"strings"

	"github.com/artarts36/db-exporter/internal/schema"
)

func (b *QueryBuilder) BuildInsertQuery(table *schema.Table, rows []map[string]interface{}) (string, error) {
	const qPlusCapacity = 3

	if len(rows) == 0 {
		return "", fmt.Errorf("rows is empty")
	}

	values := b.buildValues(table, rows)

	q := make([]string, 0, len(values)+qPlusCapacity)
	q = append(q, b.buildInsertInto(table))
	q = append(q, "VALUES")
	q = append(q, values...)

	q[len(q)-1] = fmt.Sprintf("%s;", q[len(q)-1])

	return strings.Join(q, "\n"), nil
}

func (b *QueryBuilder) buildValues(table *schema.Table, rows []map[string]interface{}) []string {
	values := make([]string, 0, len(rows))
	cols := table.ColumnsNames()

	for i, row := range rows {
		value := make([]string, 0, len(cols))

		for _, col := range cols {
			value = append(value, b.mapValue(row[col]))
		}

		comma := ","
		if i == len(rows)-1 {
			comma = ""
		}

		values = append(values, fmt.Sprintf("    (%s)%s", strings.Join(value, ", "), comma))
	}

	return values
}

func (b *QueryBuilder) buildInsertInto(table *schema.Table) string {
	return fmt.Sprintf("INSERT INTO %s (%s)", table.Name.Value, strings.Join(table.ColumnsNames(), ", "))
}

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
