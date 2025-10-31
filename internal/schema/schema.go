package schema

import (
	"github.com/artarts36/gds"
)

type Schema struct {
	Tables    *TableMap
	Sequences map[string]*Sequence
	Enums     map[string]*Enum

	Driver DatabaseDriver
}

type ForeignKey struct {
	Name          gds.String
	Table         gds.String
	ColumnsNames  *gds.Strings
	ForeignTable  gds.String
	ForeignColumn gds.String

	IsDeferrable        bool
	IsInitiallyDeferred bool
}

func NewSchema(driver DatabaseDriver) *Schema {
	return &Schema{
		Tables:    NewTableMap(),
		Sequences: map[string]*Sequence{},
		Enums:     map[string]*Enum{},
		Driver:    driver,
	}
}

func (s *Schema) Clone() *Schema {
	return &Schema{
		Tables:    s.Tables.Clone(),
		Sequences: s.Sequences,
		Enums:     s.Enums,
		Driver:    s.Driver,
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

func (s *Schema) OnlyTables(tableNames []string) *Schema {
	tableFilter := map[string]bool{}
	for _, table := range tableNames {
		tableFilter[table] = true
	}

	return s.WithoutTable(func(table *Table) bool {
		return !tableFilter[table.Name.Value]
	})
}

func (s *Schema) WithoutTable(callback func(table *Table) bool) *Schema {
	newSchema := s.Clone()

	wrappedCallback := func(table *Table) bool {
		reject := callback(table)
		if !reject {
			return false
		}

		for enumName := range table.UsingEnums {
			en, exists := newSchema.Enums[enumName]
			if !exists {
				continue
			}

			en.Used--

			if en.Used == 0 {
				delete(newSchema.Enums, enumName)
			}
		}

		for seqName := range table.UsingSequences {
			seq, exists := newSchema.Sequences[seqName]
			if !exists {
				continue
			}

			seq.Used--

			if seq.Used == 0 {
				delete(newSchema.Sequences, seqName)
			}
		}

		return true
	}

	newSchema.Tables = newSchema.Tables.reject(wrappedCallback)

	return newSchema
}
