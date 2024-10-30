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

type MapEntitiesParams struct {
	Tables  []*schema.Table
	Package *golang.Package
	Enums   map[string]*golang.StringEnum
}

func (m *EntityMapper) MapEntities(params *MapEntitiesParams) *Entities {
	ents := &Entities{
		Entities: make([]*Entity, len(params.Tables)),
		Imports:  golang.NewImportGroups(),
	}
	addImportCallback := func(pkg string) {
		ents.Imports.AddStd(pkg)
	}

	for i, table := range params.Tables {
		ents.Entities[i] = m.mapEntity(&MapEntityParams{
			Table:   table,
			Package: params.Package,
			Enums:   params.Enums,
		}, addImportCallback)
	}

	return ents
}

type MapEntityParams struct {
	Table   *schema.Table
	Package *golang.Package
	Enums   map[string]*golang.StringEnum
}

func (m *EntityMapper) MapEntity(
	params *MapEntityParams,
) *Entity {
	return m.mapEntity(params, func(_ string) {})
}

func (m *EntityMapper) mapEntity(params *MapEntityParams, addImportCallback func(pkg string)) *Entity {
	entity := &Entity{
		Name:      params.Table.Name.Singular().Pascal().FixAbbreviations(goAbbreviationsSet),
		Table:     params.Table,
		Imports:   golang.NewImportGroups(),
		AsVarName: params.Table.Name.Singular().Camel().Value,
		Package:   params.Package,
	}

	addImport := func(pkg string) {
		entity.Imports.AddStd(pkg)
		addImportCallback(pkg)
	}

	entity.Properties = m.propertyMapper.mapColumns(params.Table.Columns, params.Enums, addImport)

	return entity
}
