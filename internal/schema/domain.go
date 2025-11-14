package schema

type Domain struct {
	Name     string
	DataType DataType

	ConstraintName string
	CheckClause    string

	// List of names of tables, which using this enum.
	UsingInTables []string
}

func (e *Domain) UsingInSingleTable() bool {
	return len(e.UsingInTables) == 1
}
