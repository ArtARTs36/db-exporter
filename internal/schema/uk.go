package schema

import (
	"github.com/artarts36/gds"
)

type UniqueKey struct {
	Name         gds.String
	ColumnsNames *gds.Strings
}
