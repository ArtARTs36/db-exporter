package goentity

import (
	"fmt"

	"github.com/artarts36/db-exporter/internal/shared/golang"
)

type repositoryEntityFilter struct {
	Name       string
	Properties *goProperties
	Package    *golang.Package
}

func createRepositoryEntityFilter(
	entity *Entity,
	action string,
	pkg *golang.Package,
	properties *goProperties,
) repositoryEntityFilter {
	name := fmt.Sprintf("%s%sFilter", action, entity.Name)

	return repositoryEntityFilter{
		Name:       name,
		Properties: properties,
		Package:    pkg,
	}
}

func (e *repositoryEntityFilter) Call(pkg *golang.Package) string {
	return e.Package.CallToStruct(pkg, e.Name)
}
