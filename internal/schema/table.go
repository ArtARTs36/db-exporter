package schema

import (
	"github.com/artarts36/gds"
)

type Table struct {
	Name    gds.String `db:"name"`
	Columns []*Column  `db:"-"`

	PrimaryKey  *PrimaryKey            `db:"-"`
	ForeignKeys map[string]*ForeignKey `db:"-"`
	UniqueKeys  map[string]*UniqueKey  `db:"-"`

	UsingSequences map[string]*Sequence `db:"-"`
	UsingEnums     map[string]*Enum     `db:"-"`

	columnsNames []string `db:"-"`

	Comment string
}

func (t *Table) ColumnsNames() []string {
	if t.columnsNames == nil {
		t.columnsNames = make([]string, 0, len(t.Columns))

		for _, column := range t.Columns {
			t.columnsNames = append(t.columnsNames, column.Name.Value)
		}
	}

	return t.columnsNames
}

func (t *Table) GetColumn(name string) *Column {
	for _, column := range t.Columns {
		if column.Name.Equal(name) {
			return column
		}
	}

	return nil
}

func (t *Table) HasUniqueKeys() bool {
	return len(t.UniqueKeys) > 0
}

func (t *Table) HasSingleUniqueKey() bool {
	return len(t.UniqueKeys) == 1
}

func (t *Table) GetFirstUniqueKey() *UniqueKey {
	for _, key := range t.UniqueKeys {
		return key
	}

	return nil
}

func (t *Table) HasForeignKeyTo(tableName string) bool {
	for _, key := range t.ForeignKeys {
		if key.ForeignTable.Equal(tableName) {
			return true
		}
	}

	return false
}

func (t *Table) GetPKColumns() []*Column {
	if t.PrimaryKey == nil {
		return []*Column{}
	}

	cols := make([]*Column, 0, t.PrimaryKey.ColumnsNames.Len())

	for _, col := range t.Columns {
		if t.PrimaryKey.ColumnsNames.Contains(col.Name.Value) {
			cols = append(cols, col)
		}
	}

	return cols
}
