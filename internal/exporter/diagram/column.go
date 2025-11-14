package diagram

import "github.com/artarts36/db-exporter/internal/schema"

type diagramTable struct {
	Name                   string
	Columns                []*diagramColumn
	PrimaryKeyColumnsCount int
}

type diagramColumn struct {
	Name string
	Type string

	IsPrimaryKey  bool
	HasForeignKey bool
	IsUniqueKey   bool
}

func mapTable(tbl *schema.Table) *diagramTable {
	table := &diagramTable{
		Name:    tbl.Name.Value,
		Columns: make([]*diagramColumn, 0, len(tbl.Columns)),
	}
	if tbl.PrimaryKey != nil {
		table.PrimaryKeyColumnsCount = tbl.PrimaryKey.ColumnsNames.Len()
	}

	for _, col := range tbl.Columns {
		table.Columns = append(table.Columns, mapColumn(col))
	}

	return table
}

func mapColumn(col *schema.Column) *diagramColumn {
	column := &diagramColumn{
		Name:          col.Name.Value,
		IsPrimaryKey:  col.IsPrimaryKey(),
		HasForeignKey: col.HasForeignKey(),
		IsUniqueKey:   col.IsUniqueKey(),
	}

	switch {
	case col.DataType.IsUUID:
		column.Type = "uuid"
	case col.DataType.IsInteger:
		column.Type = "integer"
	case col.DataType.IsFloat:
		column.Type = "float"
	case col.DataType.IsBoolean:
		column.Type = "boolean"
	case col.DataType.IsNumeric:
		column.Type = "number"
	case col.DataType.IsStringable:
		column.Type = "string"
	case col.DataType.IsDatetime:
		column.Type = "datetime"
	default:
		column.Type = col.DataType.Name
	}

	return column
}
