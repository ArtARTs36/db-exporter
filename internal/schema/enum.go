package schema

import "github.com/artarts36/gds"

type Enum struct {
	Name   *gds.String
	Values []string
	Used   int
}
