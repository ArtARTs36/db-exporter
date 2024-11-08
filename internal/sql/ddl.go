package sql

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/infrastructure/typemap"
	"slices"
	"strings"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/ds"
)

const expIfNotExists = "IF NOT EXISTS "

type DDLBuilder struct {
}

type BuildDDLParams struct {
	UseIfNotExists bool

	Source config.DatabaseDriver
	Target config.DatabaseDriver
}

func NewDDLBuilder() *DDLBuilder {
	return &DDLBuilder{}
}

type isLastLine func() bool

func (b *DDLBuilder) BuildDDL(table *schema.Table, params BuildDDLParams) ([]string, error) { //nolint:funlen,lll // not need
	var upQueries []string

	if len(table.Columns) == 0 {
		ifne := ""
		if params.UseIfNotExists {
			ifne = expIfNotExists
		}

		return []string{
			fmt.Sprintf("CREATE TABLE %s%s()", ifne, table.Name.Value),
		}, nil
	}

	ifne := ""
	if params.UseIfNotExists {
		ifne = expIfNotExists
	}

	createTableQuery := []string{
		fmt.Sprintf("CREATE TABLE %s%s", ifne, table.Name.Value),
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

		defaultValue := ""
		if column.DefaultRaw.Valid {
			defaultValue = fmt.Sprintf(" DEFAULT %s", column.DefaultRaw.String)
		}

		colType, err := typemap.TransitSQLType(params.Source, params.Target, column.Type.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to map column type: %w", err)
		}

		line := fmt.Sprintf(
			"    %s%s%s%s%s%s",
			column.Name.Value,
			strings.Repeat(" ", spacesAfterColumnName),
			colType,
			notNull,
			defaultValue,
			comma,
		)

		createTableQuery = append(createTableQuery, line)

		if column.Comment.IsNotEmpty() {
			upQueries = append(upQueries, b.CommentOnColumn(column))
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

	return upQueries, nil
}

type CreateSequenceParams struct {
	UseIfNotExists bool

	Source config.DatabaseDriver
	Target config.DatabaseDriver
}

func (b *DDLBuilder) CreateSequence(seq *schema.Sequence, params CreateSequenceParams) (string, error) {
	ifne := ""
	if params.UseIfNotExists {
		ifne = expIfNotExists
	}

	dType, err := typemap.TransitSQLType(params.Source, params.Target, seq.DataType)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("CREATE SEQUENCE %s%s as %s;", ifne, seq.Name, dType), nil
}

func (b *DDLBuilder) DropTable(table *schema.Table, useIfExists bool) string {
	ife := ""
	if useIfExists {
		ife = "IF EXISTS "
	}

	return fmt.Sprintf("DROP TABLE %s%s;", ife, table.Name.Value)
}

func (b *DDLBuilder) CreateEnum(enum *schema.Enum) string {
	valuesString := ""
	for i, value := range enum.Values {
		valuesString += fmt.Sprintf("'%s'", value)

		if i < len(enum.Values)-1 {
			valuesString += ", "
		}
	}

	return fmt.Sprintf(`CREATE TYPE %s AS ENUM (%s);`, enum.Name.Value, valuesString)
}

func (b *DDLBuilder) DropType(name string, ifExists bool) string {
	ife := ""
	if ifExists {
		ife = "IF EXISTS"
	}

	return fmt.Sprintf("DROP TYPE %s%s;", ife, name)
}

func (b *DDLBuilder) DropSequence(seq *schema.Sequence, ifExists bool) string {
	ife := ""
	if ifExists {
		ife = "IF EXISTS"
	}

	return fmt.Sprintf("DROP SEQUENCE %s%s;", ife, seq.Name)
}

func (b *DDLBuilder) CommentOnColumn(col *schema.Column) string {
	return fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';", col.TableName.Value, col.Name.Value, col.Comment.Value)
}

func (b *DDLBuilder) buildPrimaryKey(table *schema.Table, isLast isLastLine) string {
	comma := ","
	if isLast() {
		comma = ""
	}

	return fmt.Sprintf(
		"    %s%s",
		b.CreatePrimaryKey(table.PrimaryKey.Name.Value, table.PrimaryKey.ColumnsNames),
		comma,
	)
}

func (b *DDLBuilder) CreatePrimaryKey(name string, columns *ds.Strings) string {
	return fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", name, columns.Join(", ").Value)
}

func (b *DDLBuilder) buildForeignKeys(table *schema.Table, isLast isLastLine) []string {
	queries := make([]string, 0, len(table.ForeignKeys))

	fks := make([]*schema.ForeignKey, 0, len(table.ForeignKeys))
	for _, fk := range table.ForeignKeys {
		fks = append(fks, fk)
	}

	slices.SortFunc(fks, func(a, b *schema.ForeignKey) int {
		return strings.Compare(a.Name.Value, b.Name.Value)
	})

	for _, fk := range fks {
		comma := ","
		if isLast() {
			comma = ""
		}

		deferrableString := ""
		if fk.IsDeferrable {
			deferrableString = " DEFERRABLE"

			if fk.IsInitiallyDeferred {
				deferrableString += " INITIALLY DEFERRED"
			}
		}

		q := fmt.Sprintf(
			"    CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)%s%s",
			fk.Name.Value,
			fk.ColumnsNames.Join(", ").Value,
			fk.ForeignTable.Value,
			fk.ForeignColumn.Value,
			deferrableString,
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

		q := fmt.Sprintf("%s%s", b.CreateUniqueKey(uk.Name.Value, uk.ColumnsNames), comma)

		queries = append(queries, q)
	}

	return queries
}

func (b *DDLBuilder) CreateUniqueKey(name string, columns *ds.Strings) string {
	return fmt.Sprintf("    CONSTRAINT %s UNIQUE (%s)", name, columns.Join(", ").Value)
}
