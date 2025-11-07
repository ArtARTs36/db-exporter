package schema

import (
	"github.com/artarts36/gds"
)

type PrimaryKey struct {
	Name         gds.String
	ColumnsNames *gds.Strings
}
