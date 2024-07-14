package sql

import (
	"fmt"
	"strings"

	"github.com/artarts36/db-exporter/internal/schema"
)

func (b *QueryBuilder) BuildDeleteQueries(table *schema.Table, rows []map[string]interface{}) []string {
	if table.PrimaryKey == nil {
		queries := make([]string, 0, len(rows))

		for _, row := range rows {
			queries = append(queries, b.BuildDeleteQuery(table.Name.Val, row))
		}

		return queries
	}

	if table.PrimaryKey.ColumnsNames.Len() == 1 && len(rows) > 1 {
		col := table.PrimaryKey.ColumnsNames.First()
		values := make([]interface{}, 0, len(rows))

		for _, row := range rows {
			values = append(values, row[col])
		}

		return []string{b.BuildDeleteInQuery(table.Name.Val, col, values)}
	}

	queries := make([]string, 0, len(rows))

	for _, row := range rows {
		pk := map[string]interface{}{}

		for _, col := range table.PrimaryKey.ColumnsNames.List() {
			pk[col] = row[col]
		}

		queries = append(queries, b.BuildDeleteQuery(table.Name.Val, pk))
	}

	return queries
}

func (b *QueryBuilder) BuildDeleteQuery(table string, fields map[string]interface{}) string {
	q := make([]string, 0, len(fields))

	for field, val := range fields {
		q = append(q, fmt.Sprintf("%s = %s", field, b.mapValue(val)))
	}

	return fmt.Sprintf("DELETE FROM %s WHERE %s;", table, strings.Join(q, " AND "))
}

func (b *QueryBuilder) BuildDeleteInQuery(table string, field string, values []interface{}) string {
	valuesStr := make([]string, 0, len(values))

	for _, val := range values {
		valuesStr = append(valuesStr, b.mapValue(val))
	}

	return fmt.Sprintf("DELETE FROM %s WHERE %s IN (%s);", table, field, strings.Join(valuesStr, ", "))
}
