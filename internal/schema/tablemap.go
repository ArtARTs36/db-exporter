package schema

import (
	"github.com/artarts36/gds"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type TableMap struct {
	oMap *orderedmap.OrderedMap[gds.String, *Table]

	list []*Table
}

func NewTableMap(table ...*Table) *TableMap {
	oMap := orderedmap.New[gds.String, *Table]()

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

func (m *TableMap) Get(name gds.String) (*Table, bool) {
	return m.oMap.Get(name)
}

func (m *TableMap) Has(name gds.String) bool {
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

func (m *TableMap) Only(tableNames []string) *TableMap {
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

func (m *TableMap) Clone() *TableMap {
	return NewTableMap(m.list...)
}
