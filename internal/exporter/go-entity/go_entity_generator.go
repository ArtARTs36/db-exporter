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

func (g *EntityGenerator) Generate(entity *Entity, pkg golang.Package) (*exporter.ExportedPage, error) {
	return g.pager.Of("go-entities/entity.go.tpl").Export(
		fmt.Sprintf("%s/%s.go", pkg.ProjectRelativePath, entity.Table.Name.Singular().Lower()),
		map[string]stick.Value{
			"schema": &Entities{
				Entities: []*Entity{entity},
				Imports:  entity.Imports,
			},
			"package": pkg,
		},
	)
}
