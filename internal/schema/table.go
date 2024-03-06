package schema

import "github.com/artarts36/db-exporter/internal/shared/ds"

type Table struct {
	Name    ds.String `db:"name"`
	Columns []*Column `db:"-"`

	PrimaryKey  *PrimaryKey            `db:"-"`
	ForeignKeys map[string]*ForeignKey `db:"-"`
	UniqueKeys  map[string]*UniqueKey  `db:"-"`

	columnsNames []string `db:"-"`
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
