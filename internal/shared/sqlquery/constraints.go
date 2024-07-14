package sqlquery

import (
	"fmt"

	"github.com/artarts36/db-exporter/internal/shared/ds"
)

func BuildPK(name string, columns *ds.Strings) string {
	return fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", name, columns.Join(", ").Val)
}

func BuildUK(name string, columns *ds.Strings) string {
	return fmt.Sprintf("    CONSTRAINT %s UNIQUE (%s)", name, columns.Join(", ").Val)
}
