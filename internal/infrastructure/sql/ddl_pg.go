package sql

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/gds"
	"slices"
	"strings"

	"github.com/artarts36/db-exporter/internal/schema"
)

const expIfNotExists = "IF NOT EXISTS "

type PostgresDDLBuilder struct {
}

type BuildDDLOpts struct {
	UseIfNotExists bool
	UseIfExists    bool

	Source config.DatabaseDriver
}

func NewPostgresDDLBuilder() *PostgresDDLBuilder {
	return &PostgresDDLBuilder{}
}

type isLastLine func() bool

func (b *PostgresDDLBuilder) buildCreateEmptyTable(table *schema.Table, useIfNotExists bool) string {
	ifne := ""
	if useIfNotExists {
		ifne = expIfNotExists
	}

	return fmt.Sprintf("CREATE TABLE %s%s()", ifne, table.Name.Value)
}

func (b *PostgresDDLBuilder) Build(schema *schema.Schema, params BuildDDLOpts) (*DDL, error) {
	ddl := &DDL{
		Name:        "init",
		UpQueries:   make([]string, 0, len(schema.Enums)+schema.Tables.Len()),
		DownQueries: []string{},
	}

	for _, enum := range schema.Enums {
		ddl.UpQueries = append(ddl.UpQueries, b.CreateEnum(enum))
	}

	for _, sequence := range schema.Sequences {
		seqSQL, err := b.CreateSequence(sequence, CreateSequenceParams{
			UseIfNotExists: params.UseIfNotExists,
		})
		if err != nil {
			return nil, err
		}

		ddl.UpQueries = append(ddl.UpQueries, seqSQL)
	}

	ddls, err := b.BuildPerTable(schema, params)
	if err != nil {
		return nil, err
	}

	for _, tableDDL := range ddls {
		ddl.UpQueries = append(ddl.UpQueries, tableDDL.UpQueries...)
		ddl.DownQueries = append(ddl.DownQueries, tableDDL.DownQueries...)
	}

	for _, enum := range schema.Enums {
		ddl.DownQueries = append(ddl.DownQueries, b.DropType(enum.Name.Value, params.UseIfExists))
	}

	for _, seq := range schema.Sequences {
		ddl.DownQueries = append(ddl.DownQueries, b.DropSequence(seq, params.UseIfExists))
	}

	return ddl, nil
}

func (b *PostgresDDLBuilder) BuildPerTable(sch *schema.Schema, params BuildDDLOpts) ([]*DDL, error) { //nolint:funlen,lll // not need
	build := func(table *schema.Table) (*DDL, error) {
		if len(table.Columns) == 0 {
			return &DDL{
				Name:        table.Name.Value,
				UpQueries:   []string{b.buildCreateEmptyTable(table, params.UseIfNotExists)},
				DownQueries: []string{b.buildDropTable(table, params.UseIfNotExists)},
			}, nil
		}

		ddl := &DDL{
			Name:        table.Name.Value,
			UpQueries:   []string{},
			DownQueries: []string{},
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

			colType, err := sqltype.TransitSQLType(sch.Driver, config.DatabaseDriverPostgres, column.Type)
			if err != nil {
				return nil, fmt.Errorf("failed to map column type: %w", err)
			}

			line := fmt.Sprintf(
				"    %s%s%s%s%s%s",
				column.Name.Value,
				strings.Repeat(" ", spacesAfterColumnName),
				colType.Name,
				notNull,
				defaultValue,
				comma,
			)

			createTableQuery = append(createTableQuery, line)

			if column.Comment.IsNotEmpty() {
				ddl.UpQueries = append(ddl.UpQueries, b.CommentOnColumn(column))
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

		ddl.UpQueries = append([]string{upSQL}, ddl.UpQueries...)
		ddl.DownQueries = append(ddl.DownQueries, b.buildDropTable(table, params.UseIfExists))

		return ddl, nil
	}

	ddls := make([]*DDL, 0)

	for _, table := range sch.Tables.List() {
		ddl, err := build(table)
		if err != nil {
			return nil, err
		}

		ddls = append(ddls, ddl)
	}

	return ddls, nil
}

type CreateSequenceParams struct {
	UseIfNotExists bool

	Source config.DatabaseDriver
	Target config.DatabaseDriver
}

func (b *PostgresDDLBuilder) ifne(use bool) string {
	if use {
		return expIfNotExists
	}
	return ""
}

func (b *PostgresDDLBuilder) ife(use bool) string {
	if use {
		return "IF EXISTS"
	}
	return ""
}

func (b *PostgresDDLBuilder) CreateSequence(seq *schema.Sequence, params CreateSequenceParams) (string, error) {
	dType, err := sqltype.TransitSQLType(params.Source, params.Target, seq.DataType)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("CREATE SEQUENCE %s%s as %s;", b.ifne(params.UseIfNotExists), seq.Name, dType.Name), nil
}

func (b *PostgresDDLBuilder) buildDropTable(table *schema.Table, useIfExists bool) string {
	return fmt.Sprintf("DROP TABLE %s%s;", b.ife(useIfExists), table.Name.Value)
}

func (b *PostgresDDLBuilder) CreateEnum(enum *schema.Enum) string {
	valuesString := ""
	for i, value := range enum.Values {
		valuesString += fmt.Sprintf("'%s'", value)

		if i < len(enum.Values)-1 {
			valuesString += ", "
		}
	}

	return fmt.Sprintf(`CREATE TYPE %s AS ENUM (%s);`, enum.Name.Value, valuesString)
}

func (b *PostgresDDLBuilder) DropType(name string, ifExists bool) string {
	return fmt.Sprintf("DROP TYPE %s%s;", b.ife(ifExists), name)
}

func (b *PostgresDDLBuilder) DropSequence(seq *schema.Sequence, ifExists bool) string {
	return fmt.Sprintf("DROP SEQUENCE %s%s;", b.ife(ifExists), seq.Name)
}

func (b *PostgresDDLBuilder) CommentOnColumn(col *schema.Column) string {
	return fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';", col.TableName.Value, col.Name.Value, col.Comment.Value)
}

func (b *PostgresDDLBuilder) buildPrimaryKey(table *schema.Table, isLast isLastLine) string {
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

func (b *PostgresDDLBuilder) CreatePrimaryKey(name string, columns *gds.Strings) string {
	return fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", name, columns.Join(", ").Value)
}

func (b *PostgresDDLBuilder) buildForeignKeys(table *schema.Table, isLast isLastLine) []string {
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

func (b *PostgresDDLBuilder) buildUniqueKeys(table *schema.Table, isLast isLastLine) []string {
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

func (b *PostgresDDLBuilder) CreateUniqueKey(name string, columns *gds.Strings) string {
	return fmt.Sprintf("    CONSTRAINT %s UNIQUE (%s)", name, columns.Join(", ").Value)
}
