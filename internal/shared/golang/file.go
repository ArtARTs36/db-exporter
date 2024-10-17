package golang

type File struct {
	Name    string
	Package *Package
	Imports *ImportGroups
}

func NewFile(name string, pkg *Package) File {
	return File{Name: name, Package: pkg, Imports: NewImportGroups()}
}

func (f *File) ImportStd(pkg *Package) {
	if f.Package.FullName == pkg.FullName {
		return
	}

	f.Imports.AddStd(pkg.FullName)
}

func (f *File) ImportShared(pkg *Package) {
	if f.Package.FullName == pkg.FullName {
		return
	}

	f.Imports.AddShared(pkg.FullName)
}

func (f *File) ImportLocal(pkg *Package) {
	if f.Package.FullName == pkg.FullName {
		return
	}

	f.Imports.AddLocal(pkg.FullName)
}
