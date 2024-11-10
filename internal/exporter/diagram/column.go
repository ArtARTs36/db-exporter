package diagram

import "github.com/artarts36/db-exporter/internal/schema"

type diagramTable struct {
	Name    string
	Columns []*diagramColumn
}

type diagramColumn struct {
	Name string
	Type string
}

func mapTable(tbl *schema.Table) *diagramTable {
	table := &diagramTable{
		Name:    tbl.Name.Value,
		Columns: make([]*diagramColumn, 0, len(tbl.Columns)),
	}

	for _, col := range tbl.Columns {
		table.Columns = append(table.Columns, mapColumn(col))
	}

	return table
}

func mapColumn(col *schema.Column) *diagramColumn {
	column := &diagramColumn{
		Name: col.Name.Value,
	}

	switch true {
	case col.Type.IsUUID:
		column.Type = "uuid"
	case col.Type.IsInteger:
		column.Type = "integer"
	case col.Type.IsFloat:
		column.Type = "float"
	case col.Type.IsBoolean:
		column.Type = "boolean"
	case col.Type.IsNumeric:
		column.Type = "number"
	case col.Type.IsStringable:
		column.Type = "string"
	case col.Type.IsDatetime:
		column.Type = "datetime"
	}

	return column
}
