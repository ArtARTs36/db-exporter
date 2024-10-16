package golang

import "github.com/artarts36/db-exporter/internal/shared/ds"

type File struct {
	Name    string
	Package *Package
	Imports *ds.Set[string]
}

func NewFile(name string, pkg *Package) File {
	return File{Name: name, Package: pkg, Imports: ds.NewSet[string]()}
}

func (f *File) Import(pkg *Package) {
	if f.Package.FullName == pkg.FullName {
		return
	}

	f.Imports.Add(pkg.FullName)
}
