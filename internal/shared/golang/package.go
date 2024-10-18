package golang

import (
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/ds"
)

type Package struct {
	Name                string
	ProjectRelativePath string
	FullName            string
}

func BuildPackage(pkgName string, module string) (*Package, error) {
	if len(pkgName) == 0 {
		return nil, errors.New("package is empty")
	}

	pkg := &Package{
		Name:                pkgName,
		ProjectRelativePath: pkgName,
		FullName:            module + "/" + pkgName,
	}

	pkgParts := ds.NewString(pkgName).SplitWords()
	if len(pkgParts) > 0 {
		pkg.Name = pkgParts[len(pkgParts)-1].Word
	}

	return pkg, nil
}

func (p *Package) IsCurrent(currentPackage *Package) bool {
	if currentPackage == nil {
		return false
	}

	return p.FullName == currentPackage.FullName
}

func (p *Package) CallToStruct(currentPackage *Package, structName string) string {
	if p.IsCurrent(currentPackage) {
		return structName
	}

	return fmt.Sprintf("%s.%s", p.Name, structName)
}
