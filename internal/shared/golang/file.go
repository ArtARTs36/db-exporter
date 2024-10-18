package golang

import (
	"fmt"
	"log/slog"
	"path/filepath"
)

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

func (f *File) CallRelativePath(that File, namePrefix string) string {
	newFileName := fmt.Sprintf("%s%s", namePrefix, f.Name)

	if f.Package.FullName == that.Package.FullName {
		return newFileName
	}

	path, err := filepath.Rel(that.Package.ProjectRelativePath, f.Package.ProjectRelativePath)
	if err != nil {
		slog.Error(
			"[go-file] failed to create relative path",
			slog.String("from", that.Package.ProjectRelativePath),
			slog.String("to", f.Package.ProjectRelativePath),
			slog.Any("err", err),
		)

		return ""
	}

	return fmt.Sprintf("%s/%s", path, newFileName)
}
