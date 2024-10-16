package goentity

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/ds"
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
	Imports  *ds.Set[string]
}

type Entity struct {
	Name       *ds.String
	Table      *schema.Table
	Properties *goProperties
	Imports    *ds.Set[string]

	AsVarName string
}

func NewEntityMapper(propertyMapper *GoPropertyMapper) *EntityMapper {
	return &EntityMapper{propertyMapper: propertyMapper}
}

func (m *EntityMapper) MapEntities(tables []*schema.Table) *Entities {
	ents := &Entities{
		Entities: make([]*Entity, len(tables)),
		Imports:  ds.NewSet[string](),
	}
	addImportCallback := func(pkg string) {
		ents.Imports.Add(pkg)
	}

	for i, table := range tables {
		ents.Entities[i] = m.mapEntity(table, addImportCallback)
	}

	return ents
}

func (m *EntityMapper) MapEntity(table *schema.Table) *Entity {
	return m.mapEntity(table, func(_ string) {})
}

func (m *EntityMapper) mapEntity(table *schema.Table, addImportCallback func(pkg string)) *Entity {
	entity := &Entity{
		Name:      table.Name.Singular().Pascal().FixAbbreviations(goAbbreviationsSet),
		Table:     table,
		Imports:   ds.NewSet[string](),
		AsVarName: table.Name.Singular().Camel().Value,
	}

	addImport := func(pkg string) {
		entity.Imports.Add(pkg)
		addImportCallback(pkg)
	}

	entity.Properties = m.propertyMapper.mapColumns(table.Columns, addImport)

	return entity
}
