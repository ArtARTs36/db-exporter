package goentity

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/artarts36/gds"
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
	Name       *gds.String
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
	SourceDriver schema.DatabaseDriver
	Tables       []*schema.Table
	Package      *golang.Package
	Enums        map[string]*golang.StringEnum
}

func (m *EntityMapper) MapEntities(params *MapEntitiesParams) *Entities {
	ents := &Entities{
		Entities: make([]*Entity, 0, len(params.Tables)),
		Imports:  golang.NewImportGroups(),
	}
	addImportCallback := func(pkg string) {
		ents.Imports.AddStd(pkg)
	}

	for _, table := range params.Tables {
		if table.IsPartition() {
			continue
		}

		ents.Entities = append(ents.Entities, m.mapEntity(&MapEntityParams{
			SourceDriver: params.SourceDriver,
			Table:        table,
			Package:      params.Package,
			Enums:        params.Enums,
		}, addImportCallback))
	}

	return ents
}

type MapEntityParams struct {
	SourceDriver schema.DatabaseDriver
	Table        *schema.Table
	Package      *golang.Package
	Enums        map[string]*golang.StringEnum
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
		AsVarName: genVarNameForTable(params.Table),
		Package:   params.Package,
	}

	addImport := func(pkg string) {
		entity.Imports.AddStd(pkg)
		addImportCallback(pkg)
	}

	entity.Properties = m.propertyMapper.mapColumns(params.SourceDriver, params.Table.Columns, params.Enums, addImport)

	return entity
}

func genVarNameForTable(table *schema.Table) string {
	varName := table.Name.Singular().Camel().Value
	if replace, ok := golang.ReservedNameReplaceMap[varName]; ok {
		return replace
	}
	return varName
}
