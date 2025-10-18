package schema

import "github.com/artarts36/gds"

type Enum struct {
	Name   *gds.String
	Values []string
	Used   int

	// List of names of tables, which using this enum.
	UsingInTables []string
}

func (e *Enum) UsingInSingleTable() bool {
	return len(e.UsingInTables) == 1
}
