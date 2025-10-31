package sql

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/gds"
	"slices"
	"strings"

	"github.com/artarts36/db-exporter/internal/schema"
)

const mySQLColumnNameWrapper = "`"

type MySQLDDLBuilder struct {
}

func NewMySQLDDLBuilder() *MySQLDDLBuilder {
	return &MySQLDDLBuilder{}
}

func (b *MySQLDDLBuilder) Build(schema *schema.Schema, params BuildDDLOpts) (*DDL, error) {
	ddl := &DDL{
		Name:        *gds.NewString("init"),
		UpQueries:   []string{},
		DownQueries: []string{},
	}

	for _, table := range schema.Tables.List() {
		tableDDL, err := b.buildCreateTable(table, schema.Driver, params)
		if err != nil {
			return nil, fmt.Errorf("failed to build DDL for table %q: %w", table.Name.Value, err)
		}

		ddl.UpQueries = append(ddl.UpQueries, tableDDL.UpQueries...)
		ddl.DownQueries = append(ddl.DownQueries, tableDDL.DownQueries...)
	}

	return ddl, nil
}

func (b *MySQLDDLBuilder) buildCreateTable( //nolint:funlen // not need
	table *schema.Table,
	sourceDriver schema.DatabaseDriver,
	params BuildDDLOpts,
) (*DDL, error) {
	if len(table.Columns) == 0 {
		return &DDL{
			Name:        table.Name,
			UpQueries:   []string{buildCreateEmptyTable(table, params.UseIfNotExists)},
			DownQueries: []string{b.buildDropTable(table, params.UseIfNotExists)},
		}, nil
	}

	ddl := &DDL{
		Name:        table.Name,
		UpQueries:   make([]string, 0),
		DownQueries: make([]string, 0),
	}

	ifne := ifne(params.UseIfNotExists)

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

	columnNames, maxColumnLen := b.prepareColumnNames(table.Columns)

	for i, column := range table.Columns {
		notNull := ""
		if !column.Nullable {
			notNull = " NOT NULL"
		}

		comma := ","
		if isLast() {
			comma = ""
		}

		colName := columnNames[i]

		spacesAfterColumnName := maxColumnLen - len(colName) + 1

		defaultValue := ""
		if column.DefaultRaw.Valid && sourceDriver == schema.DatabaseDriverMySQL {
			defaultValue = fmt.Sprintf(" DEFAULT %s", column.DefaultRaw.String)
		}

		colType, err := sqltype.TransitSQLType(sourceDriver, schema.DatabaseDriverMySQL, column.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to map column type: %w", err)
		}

		autoIncrement := ""
		if column.IsAutoincrement {
			autoIncrement = " AUTO_INCREMENT"
		}

		comment := ""
		if column.Comment.IsNotEmpty() {
			comment = fmt.Sprintf(" COMMENT '%s'", column.Comment.Value)
		}

		colTypeDef := colType.String()
		if column.Enum != nil {
			colTypeDef = b.buildEnumColumnType(column.Enum)
		}

		line := fmt.Sprintf(
			"    %s%s%s%s%s%s%s%s",
			colName,
			strings.Repeat(" ", spacesAfterColumnName),
			colTypeDef,
			notNull,
			defaultValue,
			autoIncrement,
			comment,
			comma,
		)

		createTableQuery = append(createTableQuery, line)
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

	ddl.UpQueries = append([]string{
		strings.Join(createTableQuery, "\n"),
	}, ddl.UpQueries...)
	ddl.DownQueries = append(ddl.DownQueries, b.buildDropTable(table, params.UseIfExists))

	return ddl, nil
}

func (b *MySQLDDLBuilder) prepareColumnNames(cols []*schema.Column) ([]string, int) {
	maxColumnLen := 0
	columnNames := make([]string, len(cols))
	for i, column := range cols {
		colName := column.Name.Wrap(mySQLColumnNameWrapper).Value
		if len(colName) > maxColumnLen {
			maxColumnLen = len(colName)
		}
		columnNames[i] = colName
	}
	return columnNames, maxColumnLen
}

func (b *MySQLDDLBuilder) BuildPerTable(sch *schema.Schema, params BuildDDLOpts) ([]*DDL, error) { //nolint:lll // not need
	ddls := make([]*DDL, 0)

	for _, table := range sch.Tables.List() {
		ddl, err := b.buildCreateTable(table, sch.Driver, params)
		if err != nil {
			return nil, err
		}

		ddls = append(ddls, ddl)
	}

	return ddls, nil
}

func (b *MySQLDDLBuilder) buildDropTable(table *schema.Table, useIfExists bool) string {
	return fmt.Sprintf("DROP TABLE %s%s;", ife(useIfExists), table.Name.Value)
}

func (b *MySQLDDLBuilder) buildPrimaryKey(table *schema.Table, isLast isLastLine) string {
	comma := ","
	if isLast() {
		comma = ""
	}

	return fmt.Sprintf(
		"    %s%s",
		b.createPrimaryKey(table.PrimaryKey.Name.Value, table.PrimaryKey.ColumnsNames),
		comma,
	)
}

func (b *MySQLDDLBuilder) createPrimaryKey(name string, columns *gds.Strings) string {
	return fmt.Sprintf("CONSTRAINT `%s` PRIMARY KEY (%s)", name, columns.Wrap(mySQLColumnNameWrapper).Join(", ").Value)
}

func (b *MySQLDDLBuilder) buildForeignKeys(table *schema.Table, isLast isLastLine) []string {
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
			fk.ColumnsNames.Wrap(mySQLColumnNameWrapper).Join(", ").Value,
			fk.ForeignTable.Value,
			fk.ForeignColumn.Wrap(mySQLColumnNameWrapper).Value,
			deferrableString,
			comma,
		)

		queries = append(queries, q)
	}

	return queries
}

func (b *MySQLDDLBuilder) buildUniqueKeys(table *schema.Table, isLast isLastLine) []string {
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

func (b *MySQLDDLBuilder) CreateUniqueKey(name string, columns *gds.Strings) string {
	return fmt.Sprintf(
		"    CONSTRAINT `%s` UNIQUE (%s)",
		name,
		columns.Wrap(mySQLColumnNameWrapper).Join(", ").Value,
	)
}

func (b *MySQLDDLBuilder) buildEnumColumnType(en *schema.Enum) string {
	return fmt.Sprintf("ENUM(%s)", gds.NewStrings(en.Values...).Wrap("'").Join(", "))
}
