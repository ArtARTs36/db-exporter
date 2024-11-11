package schema

import "github.com/artarts36/gds"

type Enum struct {
	Name   *gds.String
	Values []string
	Used   int
}

func (e *Enum) UsedOnce() bool {
	return e.Used == 1
}
