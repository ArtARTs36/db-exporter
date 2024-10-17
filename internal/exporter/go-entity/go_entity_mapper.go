package goentity

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

var goAbbreviationsSet = map[string]bool{
	"id":   true,
	"uuid": true,
	"json": true,
	"db":   true,
}

var goAbbreviationsPluralsSet = map[string]string{
	"id":   "IDs",
	"uuid": "UUIDs",
	"json": "JSONs",
	"db":   "DBs",
}

type EntityMapper struct {
	propertyMapper *GoPropertyMapper
}

type Entities struct {
	Entities []*Entity
	Imports  *golang.ImportGroups
}

type Entity struct {
	Name       *ds.String
	Table      *schema.Table
	Properties *goProperties
	Imports    *golang.ImportGroups
	Package    *golang.Package

	AsVarName string
}

func NewEntityMapper(propertyMapper *GoPropertyMapper) *EntityMapper {
	return &EntityMapper{propertyMapper: propertyMapper}
}

func (e *Entity) Call(pkg *golang.Package) string {
	return e.Package.CallToStruct(pkg, e.Name.Value)
}

func (m *EntityMapper) MapEntities(tables []*schema.Table, pkg *golang.Package) *Entities {
	ents := &Entities{
		Entities: make([]*Entity, len(tables)),
		Imports:  golang.NewImportGroups(),
	}
	addImportCallback := func(pkg string) {
		ents.Imports.AddStd(pkg)
	}

	for i, table := range tables {
		ents.Entities[i] = m.mapEntity(table, pkg, addImportCallback)
	}

	return ents
}

func (m *EntityMapper) MapEntity(table *schema.Table, pkg *golang.Package) *Entity {
	return m.mapEntity(table, pkg, func(_ string) {})
}

func (m *EntityMapper) mapEntity(table *schema.Table, pkg *golang.Package, addImportCallback func(pkg string)) *Entity {
	entity := &Entity{
		Name:      table.Name.Singular().Pascal().FixAbbreviations(goAbbreviationsSet),
		Table:     table,
		Imports:   golang.NewImportGroups(),
		AsVarName: table.Name.Singular().Camel().Value,
		Package:   pkg,
	}

	addImport := func(pkg string) {
		entity.Imports.AddStd(pkg)
		addImportCallback(pkg)
	}

	entity.Properties = m.propertyMapper.mapColumns(table.Columns, addImport)

	return entity
}
