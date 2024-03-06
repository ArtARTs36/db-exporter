package schema

import "github.com/artarts36/db-exporter/internal/shared/ds"

type Schema struct {
	Tables map[String]*Table
}

type Table struct {
	Name    String    `db:"name"`
	Columns []*Column `db:"-"`

	PrimaryKey  *PrimaryKey            `db:"-"`
	ForeignKeys map[string]*ForeignKey `db:"-"`
	UniqueKeys  map[string]*UniqueKey  `db:"-"`

	columnsNames []string `db:"-"`
}

type ForeignKey struct {
	Name          String
	Table         String
	ColumnsNames  *ds.Strings
	ForeignTable  String
	ForeignColumn String
}

type PrimaryKey struct {
	Name         String
	ColumnsNames *ds.Strings
}

type UniqueKey struct {
	Name         String
	ColumnsNames *ds.Strings
}

func (s *Schema) TablesNames() []string {
	names := make([]string, 0, len(s.Tables))
	for _, table := range s.Tables {
		names = append(names, table.Name.Value)
	}

	return names
}

func (s *Schema) SortByRelations() {
	tables := map[String]*Table{}
	queue := map[String]*Table{}

	for _, table := range s.Tables {
		if len(table.ForeignKeys) == 0 {
			tables[table.Name] = table

			continue
		}

		queue[table.Name] = table
	}

	for len(queue) > 0 {
		for _, t := range queue {
			for _, fk := range t.ForeignKeys {
				if _, exists := tables[fk.Table]; !exists {
					break
				}
			}

			tables[t.Name] = t
			delete(queue, t.Name)
		}
	}

	s.Tables = tables
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
