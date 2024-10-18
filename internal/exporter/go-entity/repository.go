package goentity

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

type Repository struct {
	Name       string
	Interface  RepositoryInterface
	Entity     *Entity
	EntityCall string

	Filters struct {
		List   repositoryEntityFilter
		Get    repositoryEntityFilter
		Delete repositoryEntityFilter
	}
	Package *golang.Package

	File golang.File
}

type RepositoryInterface struct {
	Name    string
	Package *golang.Package
}

func buildRepository(
	entity *Entity,
	pkg *golang.Package,
	interfacePkg *golang.Package,
) *Repository {
	repository := &Repository{
		Name:    fmt.Sprintf("PG%sRepository", entity.Name),
		Entity:  entity,
		Package: pkg,
		File:    golang.NewFile(fmt.Sprintf("%s.go", entity.Table.Name.Singular().Lower().Value), pkg),
	}

	repository.Interface.Name = fmt.Sprintf("%sRepository", entity.Name)
	repository.Interface.Package = interfacePkg

	return repository
}

func (i *RepositoryInterface) Call(pkg *golang.Package) string {
	return i.Package.CallToStruct(pkg, i.Name)
}
