package schema

import "github.com/artarts36/db-exporter/internal/shared/ds"

type Schema struct {
	Tables map[ds.String]*Table
}

type ForeignKey struct {
	Name          ds.String
	Table         ds.String
	ColumnsNames  *ds.Strings
	ForeignTable  ds.String
	ForeignColumn ds.String
}

type PrimaryKey struct {
	Name         ds.String
	ColumnsNames *ds.Strings
}

type UniqueKey struct {
	Name         ds.String
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
	tables := map[ds.String]*Table{}
	queue := map[ds.String]*Table{}

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
