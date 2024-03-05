package schema

type Schema struct {
	Tables map[String]*Table
}

type Table struct {
	Name    String    `db:"name"`
	Columns []*Column `db:"-"`

	PrimaryKey *PrimaryKey `db:"-"`
}

type ForeignKey struct {
	Name   String
	Table  String
	Column String
}

type PrimaryKey struct {
	Name         String
	ColumnsNames []string
}

func (s *Schema) TablesNames() []string {
	names := make([]string, 0, len(s.Tables))
	for _, table := range s.Tables {
		names = append(names, table.Name.Value)
	}

	return names
}
