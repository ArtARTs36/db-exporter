package schema

import (
	"github.com/artarts36/gds"
)

type Table struct {
	Name      gds.String         `db:"name"` // Name of table.
	Columns   []*Column          `db:"-"`    // List of table columns.
	columnMap map[string]*Column `db:"-"`

	PrimaryKey  *PrimaryKey            `db:"-"` // Primary key of table.
	ForeignKeys map[string]*ForeignKey `db:"-"` // Map of foreign keys.
	UniqueKeys  map[string]*UniqueKey  `db:"-"` // Map of unique keys.

	UsingSequences map[string]*Sequence `db:"-"` // Map of using sequences.
	UsingEnums     map[string]*Enum     `db:"-"` // Map of using enums.
	UsingDomains   map[string]*Domain   `db:"-"` // Map of using domains.

	columnsNames []string `db:"-"`

	Comment string // Comment for table.

	PartitionOf *Table // Parent table.
	Partitions  []*Table
}

func NewTable(name gds.String) *Table {
	return &Table{
		Name:           name,
		ForeignKeys:    map[string]*ForeignKey{},
		UniqueKeys:     map[string]*UniqueKey{},
		UsingSequences: map[string]*Sequence{},
		UsingEnums:     map[string]*Enum{},
		UsingDomains:   map[string]*Domain{},
		columnMap:      map[string]*Column{},
		Partitions:     []*Table{},
	}
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
	if c, ok := t.columnMap[name]; ok {
		return c
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

func (t *Table) AddColumn(col *Column) {
	t.Columns = append(t.Columns, col)
	t.columnMap[col.Name.Value] = col
}

func (t *Table) AddEnum(enum *Enum) {
	t.UsingEnums[enum.Name.Value] = enum
}

var softDeletedColumnNames = []string{
	"deleted_at",
	"delete_time",
}

func (t *Table) SupportsSoftDelete() bool {
	for _, name := range softDeletedColumnNames {
		col, ok := t.columnMap[name]
		if !ok {
			continue
		}

		if col.DataType.IsDatetime || col.DataType.IsDate {
			return true
		}
	}

	return false
}

func (t *Table) IsPartition() bool {
	return t.PartitionOf != nil
}
