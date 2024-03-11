package schema

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"

	"github.com/artarts36/db-exporter/internal/shared/ds"
)

type TableMap struct {
	oMap *orderedmap.OrderedMap[ds.String, *Table]

	list []*Table
}

func NewTableMap(table ...*Table) *TableMap {
	oMap := orderedmap.New[ds.String, *Table]()

	for _, t := range table {
		oMap.Set(t.Name, t)
	}

	return &TableMap{
		oMap: oMap,
		list: table,
	}
}

func (m *TableMap) Add(table *Table) {
	if _, exists := m.oMap.Get(table.Name); exists {
		return
	}

	m.oMap.Set(table.Name, table)
	m.list = append(m.list, table)
}

func (m *TableMap) Each(callback func(table *Table)) {
	for pair := m.oMap.Oldest(); pair != nil; pair = pair.Next() {
		callback(pair.Value)
	}
}

func (m *TableMap) EachWithErr(callback func(table *Table) error) error {
	for pair := m.oMap.Oldest(); pair != nil; pair = pair.Next() {
		if err := callback(pair.Value); err != nil {
			return err
		}
	}

	return nil
}

func (m *TableMap) Len() int {
	return m.oMap.Len()
}

func (m *TableMap) Get(name ds.String) (*Table, bool) {
	return m.oMap.Get(name)
}

func (m *TableMap) Has(name ds.String) bool {
	_, exists := m.oMap.Get(name)

	return exists
}

func (m *TableMap) Reject(callback func(table *Table) bool) *TableMap {
	tm := NewTableMap()

	for pair := m.oMap.Oldest(); pair != nil; pair = pair.Next() {
		if !callback(pair.Value) {
			tm.Add(pair.Value)
		}
	}

	return tm
}

func (m *TableMap) Without(tableNames []string) *TableMap {
	tableFilter := map[string]bool{}
	for _, table := range tableNames {
		tableFilter[table] = true
	}

	return m.Reject(func(table *Table) bool {
		return !tableFilter[table.Name.Value]
	})
}

func (m *TableMap) List() []*Table {
	return m.list
}
