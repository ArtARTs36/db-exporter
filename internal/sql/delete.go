package sql

import (
	"fmt"
	"strings"

	"github.com/artarts36/db-exporter/internal/schema"
)

func (b *QueryBuilder) BuildDeleteQueries(table *schema.Table, rows []map[string]interface{}) []string {
	queries := make([]string, 0, len(rows))

	for _, row := range rows {
		pk := map[string]interface{}{}

		for _, col := range table.PrimaryKey.ColumnsNames.List() {
			pk[col] = row[col]
		}

		queries = append(queries, b.BuildDeleteQuery(table.Name.Value, pk))
	}

	return queries
}

func (b *QueryBuilder) BuildDeleteQuery(table string, fields map[string]interface{}) string {
	q := []string{
		fmt.Sprintf("DELETE FROM %s WHERE", table),
	}

	for field, val := range fields {
		q = append(q, fmt.Sprintf("%s = %s", field, b.mapValue(val)))
	}

	return fmt.Sprintf("%s;", strings.Join(q, " "))
}
