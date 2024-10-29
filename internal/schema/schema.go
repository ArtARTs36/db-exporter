package schema

import (
	"github.com/artarts36/db-exporter/internal/shared/ds"
)

type Schema struct {
	Tables    *TableMap
	Sequences map[string]*Sequence
	Enums     map[string]*Enum
}

type ForeignKey struct {
	Name          ds.String
	Table         ds.String
	ColumnsNames  *ds.Strings
	ForeignTable  ds.String
	ForeignColumn ds.String

	IsDeferrable        bool
	IsInitiallyDeferred bool
}

type PrimaryKey struct {
	Name         ds.String
	ColumnsNames *ds.Strings
}

type UniqueKey struct {
	Name         ds.String
	ColumnsNames *ds.Strings
}

func (s *Schema) Clone() *Schema {
	return &Schema{
		Tables:    s.Tables.Clone(),
		Sequences: s.Sequences,
		Enums:     s.Enums,
	}
}

func (s *Schema) TablesNames() []string {
	names := make([]string, 0, s.Tables.Len())

	s.Tables.Each(func(table *Table) {
		names = append(names, table.Name.Value)
	})

	return names
}

func (s *Schema) SortByRelations() {
	tableList := s.Tables.List()
	tableMap := NewTableMap()

	for _, table := range tableList {
		s.appendTableMapByRelations(tableMap, table)
	}

	s.Tables = tableMap
}

func (s *Schema) appendTableMapByRelations(tableMap *TableMap, table *Table) {
	for _, key := range table.ForeignKeys {
		foreignTable, exists := s.Tables.Get(key.ForeignTable)
		if !exists {
			continue
		}

		if tableMap.Has(foreignTable.Name) {
			continue
		}

		if foreignTable.HasForeignKeyTo(table.Name.Value) {
			continue
		}

		s.appendTableMapByRelations(tableMap, foreignTable)
	}

	tableMap.Add(table)
}
