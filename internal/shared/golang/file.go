package golang

import "github.com/artarts36/db-exporter/internal/shared/ds"

type File struct {
	Name    string
	Package *Package
	Imports *ds.Set
}
