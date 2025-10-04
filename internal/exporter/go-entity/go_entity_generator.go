package goentity

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/common"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/tyler-sommer/stick"
)

type EntityGenerator struct {
	pager *common.Pager
}

func NewEntityGenerator(pager *common.Pager) *EntityGenerator {
	return &EntityGenerator{pager: pager}
}

type GenerateEntityParams struct {
	Entity       *Entity
	Package      *golang.Package
	Repositories []*Repository
	WithMocks    bool
}

func (g *EntityGenerator) GenerateEntity(params *GenerateEntityParams) (*exporter.ExportedPage, error) {
	goFile := golang.File{
		Name:    fmt.Sprintf("%s.go", params.Entity.Table.Name.Singular().Lower()),
		Package: params.Package,
		Imports: params.Entity.Imports,
	}

	return g.pager.Of("@embed/go-entities/entity.go.tpl").Export(
		fmt.Sprintf("%s/%s", params.Package.ProjectRelativePath, goFile.Name),
		map[string]stick.Value{
			"schema": map[string]stick.Value{
				"Entities":     []*Entity{params.Entity},
				"Repositories": params.Repositories,
				"Imports":      params.Entity.Imports,
				"WithMocks":    params.WithMocks,
			},
			"_file": goFile,
		},
	)
}

func (g *EntityGenerator) GenerateEntities(
	entities *Entities,
	pkg *golang.Package,
	enums map[string]*golang.StringEnum,
) (*exporter.ExportedPage, error) {
	return g.pager.Of("@embed/go-entities/entity.go.tpl").Export(
		fmt.Sprintf("%s/entities.go", pkg.ProjectRelativePath),
		map[string]stick.Value{
			"schema": map[string]stick.Value{
				"Entities": entities.Entities,
				"Enums":    enums,
			},
			"_file": golang.File{
				Name:    "entities.go",
				Package: pkg,
				Imports: entities.Imports,
			},
		},
	)
}
