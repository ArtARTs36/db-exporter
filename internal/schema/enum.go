package schema

import "github.com/artarts36/db-exporter/internal/shared/ds"

type Enum struct {
	Name   *ds.String
	Values []string
	Used   int
}
