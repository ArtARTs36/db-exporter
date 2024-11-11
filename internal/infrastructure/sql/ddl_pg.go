package sql

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/gds"
	"log/slog"
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
		Name:        *gds.NewString("init"),
		UpQueries:   make([]string, 0, (len(schema.Enums)*2)+(len(schema.Sequences)*2)+(schema.Tables.Len())*2),
		DownQueries: []string{},
	}

	steps := map[string]func() error{
		"build enums create queries": func() error {
			for _, enum := range schema.Enums {
				ddl.UpQueries = append(ddl.UpQueries, b.CreateEnum(enum))
			}
			return nil
		},
		"build sequences create queries": func() error {
			for _, sequence := range schema.Sequences {
				seqSQL, err := b.CreateSequence(sequence, CreateSequenceParams{
					Source:         params.Source,
					UseIfNotExists: params.UseIfNotExists,
				})
				if err != nil {
					return err
				}
				ddl.UpQueries = append(ddl.UpQueries, seqSQL)
			}
			return nil
		},
		"build tables create/drop queries": func() error {
			for _, table := range schema.Tables.List() {
				tableDDL, err := b.buildCreateTable(table, schema.Driver, params)
				if err != nil {
					return err
				}
				ddl.UpQueries = append(ddl.UpQueries, tableDDL.UpQueries...)
				ddl.DownQueries = append(ddl.DownQueries, tableDDL.DownQueries...)
			}
			return nil
		},
		"build enums drop queries": func() error {
			for _, enum := range schema.Enums {
				ddl.DownQueries = append(ddl.DownQueries, b.dropType(enum.Name.Value, params.UseIfExists))
			}
			return nil
		},
		"build sequences drop queries": func() error {
			for _, seq := range schema.Sequences {
				ddl.DownQueries = append(ddl.DownQueries, b.dropSequence(seq, params.UseIfExists))
			}
			return nil
		},
	}

	for name, step := range steps {
		slog.Debug(fmt.Sprintf("[pg-ddl-builder] %s", name))

		if err := step(); err != nil {
			return nil, fmt.Errorf("failed to %s: %w", name, err)
		}
	}

	return ddl, nil
}

func (b *PostgresDDLBuilder) buildCreateTable(
	table *schema.Table,
	sourceDriver config.DatabaseDriver,
	params BuildDDLOpts,
) (*DDL, error) {
	if len(table.Columns) == 0 {
		return &DDL{
			Name:        table.Name,
			UpQueries:   []string{b.buildCreateEmptyTable(table, params.UseIfNotExists)},
			DownQueries: []string{b.buildDropTable(table, params.UseIfNotExists)},
		}, nil
	}

	ddl := &DDL{
		Name:        table.Name,
		UpQueries:   []string{},
		DownQueries: []string{},
	}

	ifne := b.ifne(params.UseIfNotExists)

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

		colType, err := sqltype.TransitSQLType(sourceDriver, config.DatabaseDriverPostgres, column.Type)
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

func (b *PostgresDDLBuilder) BuildPerTable(sch *schema.Schema, params BuildDDLOpts) ([]*DDL, error) { //nolint:funlen,lll // not need
	ddlsLen := sch.Tables.Len()

	enumsDDL := &DDL{
		Name:        *gds.NewString("enums"),
		UpQueries:   make([]string, 0, len(sch.Enums)),
		DownQueries: make([]string, 0, len(sch.Enums)),
	}
	seqDDL := &DDL{
		Name:        *gds.NewString("sequences"),
		UpQueries:   make([]string, 0, len(sch.Sequences)),
		DownQueries: make([]string, 0, len(sch.Sequences)),
	}
	for _, enum := range sch.Enums {
		if enum.Used > 1 {
			enumsDDL.UpQueries = append(enumsDDL.UpQueries, b.CreateEnum(enum))
			enumsDDL.DownQueries = append(enumsDDL.DownQueries, b.dropType(enum.Name.Value, params.UseIfExists))
		}
	}
	for _, seq := range sch.Sequences {
		if seq.Used > 1 {
			seqSQL, err := b.CreateSequence(seq, CreateSequenceParams{
				UseIfNotExists: params.UseIfNotExists,
				Source:         params.Source,
			})
			if err != nil {
				return nil, err
			}
			seqDDL.UpQueries = append(seqDDL.UpQueries, seqSQL)
			seqDDL.DownQueries = append(seqDDL.DownQueries, b.dropSequence(seq, params.UseIfExists))
		}
	}

	if len(enumsDDL.UpQueries) > 0 {
		ddlsLen++
	}
	if len(seqDDL.UpQueries) > 0 {
		ddlsLen++
	}

	ddls := make([]*DDL, 0, ddlsLen)

	if len(enumsDDL.UpQueries) > 0 {
		ddls = append(ddls, enumsDDL)
	}
	if len(seqDDL.UpQueries) > 0 {
		ddls = append(ddls, seqDDL)
	}

	for _, table := range sch.Tables.List() {
		ddl := &DDL{}

		for _, enum := range table.UsingEnums {
			if enum.UsedOnce() {
				ddl.UpQueries = append(ddl.UpQueries, b.CreateEnum(enum))
				ddl.DownQueries = append(ddl.DownQueries, b.dropType(enum.Name.Value, params.UseIfExists))
			}
		}

		for _, sequence := range table.UsingSequences {
			if sequence.UsedOnce() {
				seqSQL, err := b.CreateSequence(sequence, CreateSequenceParams{
					UseIfNotExists: params.UseIfNotExists,
					Source:         params.Source,
				})
				if err != nil {
					return nil, err
				}
				seqDDL.UpQueries = append(seqDDL.UpQueries, seqSQL)
				seqDDL.DownQueries = append(seqDDL.DownQueries, b.dropSequence(sequence, params.UseIfExists))
			}
		}

		tableDDL, err := b.buildCreateTable(table, sch.Driver, params)
		if err != nil {
			return nil, err
		}
		ddl.Name = tableDDL.Name
		ddl.UpQueries = append(ddl.UpQueries, tableDDL.UpQueries...)
		ddl.DownQueries = append(ddl.DownQueries, tableDDL.DownQueries...)

		ddls = append(ddls, ddl)
	}

	return ddls, nil
}

type CreateSequenceParams struct {
	UseIfNotExists bool

	Source config.DatabaseDriver
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
	dType, err := sqltype.TransitSQLType(params.Source, config.DatabaseDriverPostgres, seq.DataType)
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

func (b *PostgresDDLBuilder) dropType(name string, ifExists bool) string {
	return fmt.Sprintf("DROP TYPE %s%s;", b.ife(ifExists), name)
}

func (b *PostgresDDLBuilder) dropSequence(seq *schema.Sequence, ifExists bool) string {
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
