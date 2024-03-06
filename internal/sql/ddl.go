package sql

import (
	"fmt"
	"strings"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/sqlquery"
)

type DDLBuilder struct {
}

func NewDDLBuilder() *DDLBuilder {
	return &DDLBuilder{}
}

type isLastLine func() bool

func (b *DDLBuilder) BuildDDL(table *schema.Table) []string {
	var upQueries []string

	createTableQuery := []string{
		fmt.Sprintf("CREATE TABLE %s", table.Name.Value),
		"(",
	}

	lines := len(table.Columns) + len(table.ForeignKeys) + len(table.UniqueKeys)
	if table.PrimaryKey != nil {
		lines++
	}
	lineID := 0

	isLast := func() bool {
		lineID++

		return lineID == lines
	}

	maxColumnLen := 0
	for _, column := range table.Columns {
		if column.Name.Len() > maxColumnLen {
			maxColumnLen = column.Name.Len()
		}
	}

	for _, column := range table.Columns {
		lineID++

		notNull := ""
		if !column.Nullable {
			notNull = " NOT NULL"
		}

		comma := ","
		if lineID == lines {
			comma = ""
		}

		spacesAfterColumnName := maxColumnLen - column.Name.Len() + 1

		line := fmt.Sprintf(
			"    %s%s%s%s%s",
			column.Name.Value,
			strings.Repeat(" ", spacesAfterColumnName),
			column.Type.Value,
			notNull,
			comma,
		)

		createTableQuery = append(createTableQuery, line)

		if column.Comment.IsNotEmpty() {
			upQueries = append(upQueries, sqlquery.BuildCommentOnColumn(
				column.TableName.Value,
				column.Name.Value,
				column.Comment.Value,
			))
		}
	}

	if lineID != lines {
		createTableQuery = append(createTableQuery, "")
	}

	if table.PrimaryKey != nil {
		createTableQuery = append(createTableQuery, b.buildPrimaryKey(table, isLast))
	}

	if len(table.ForeignKeys) > 0 {
		createTableQuery = append(createTableQuery, b.buildForeignKeys(table, isLast)...)
	}

	if len(table.UniqueKeys) > 0 {
		createTableQuery = append(createTableQuery, b.buildUniqueKeys(table, isLast)...)
	}

	createTableQuery = append(createTableQuery, ");")

	upSQL := strings.Join(createTableQuery, "\n")

	upQueries = append([]string{upSQL}, upQueries...)

	return upQueries
}

func (b *DDLBuilder) buildPrimaryKey(table *schema.Table, isLast isLastLine) string {
	comma := ","
	if isLast() {
		comma = ""
	}

	return fmt.Sprintf(
		"    %s%s",
		sqlquery.BuildPK(table.PrimaryKey.Name.Value, table.PrimaryKey.ColumnsNames),
		comma,
	)
}

func (b *DDLBuilder) buildForeignKeys(table *schema.Table, isLast isLastLine) []string {
	queries := make([]string, 0, len(table.ForeignKeys))

	for _, fk := range table.ForeignKeys {
		comma := ","
		if isLast() {
			comma = ""
		}

		q := fmt.Sprintf(
			"    CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)%s",
			fk.Name.Value,
			fk.ColumnsNames.Join(", "),
			fk.ForeignTable.Value,
			fk.ForeignColumn.Value,
			comma,
		)

		queries = append(queries, q)
	}

	return queries
}

func (b *DDLBuilder) buildUniqueKeys(table *schema.Table, isLast isLastLine) []string {
	queries := make([]string, 0, len(table.UniqueKeys))

	for _, uk := range table.UniqueKeys {
		comma := ","
		if isLast() {
			comma = ""
		}

		q := fmt.Sprintf("%s%s", sqlquery.BuildUK(uk.Name.Value, uk.ColumnsNames), comma)

		queries = append(queries, q)
	}

	return queries
}
