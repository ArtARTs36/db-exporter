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

type EntityMapper struct {
}

type Entities struct {
	Entities []*Entity
	Imports  *ds.Set
}

type Entity struct {
	Name       *ds.String
	Table      *schema.Table
	Properties []*GoProperty
	Imports    *ds.Set
}

type GoProperty struct {
	Name       string
	NameOffset int
	ColumnName string
	Type       string
	TypeOffset int

	Column *schema.Column
}

func NewEntityMapper() *EntityMapper {
	return &EntityMapper{}
}

func (m *EntityMapper) MapEntities(tables []*schema.Table) *Entities {
	ents := &Entities{
		Entities: make([]*Entity, len(tables)),
		Imports:  ds.NewSet(),
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
	return m.mapEntity(table, func(pkg string) {})
}

func (m *EntityMapper) mapEntity(table *schema.Table, addImportCallback func(pkg string)) *Entity {
	entity := &Entity{
		Name:       table.Name.Singular().Pascal().FixAbbreviations(goAbbreviationsSet),
		Table:      table,
		Properties: make([]*GoProperty, 0, len(table.Columns)),
		Imports:    ds.NewSet(),
	}
	addImport := func(pkg string) {
		entity.Imports.Add(pkg)
		addImportCallback(pkg)
	}

	propNameOffset := 0
	propTypeOffset := 0
	for _, c := range table.Columns {
		prop := &GoProperty{
			Name:       c.Name.Pascal().FixAbbreviations(goAbbreviationsSet).Value,
			Type:       m.mapGoType(c, addImport),
			ColumnName: c.Name.Value,
		}

		entity.Properties = append(entity.Properties, prop)

		if len(prop.Name) > propNameOffset {
			propNameOffset = c.Name.Pascal().Len()
		}

		if len(prop.Type) > propTypeOffset {
			propTypeOffset = len(prop.Type)
		}
	}

	for _, prop := range entity.Properties {
		prop.NameOffset = propNameOffset - len(prop.Name)
		prop.TypeOffset = propTypeOffset - len(prop.Type)
	}

	return entity
}

func (m *EntityMapper) mapGoType(col *schema.Column, addImport func(pkg string)) string {
	switch col.PreparedType {
	case schema.ColumnTypeInteger64, schema.ColumnTypeInteger:
		if col.Nullable {
			addImport("database/sql")

			return "sql.NullInt64"
		}

		return golang.TypeInt64
	case schema.ColumnTypeInteger16:
		if col.Nullable {
			addImport("database/sql")

			return "sql.NullInt16"
		}

		return golang.TypeInt16
	case schema.ColumnTypeString:
		if col.Nullable {
			addImport("database/sql")

			return "sql.NullString"
		}

		return golang.TypeString
	case schema.ColumnTypeTimestamp:
		if col.Nullable {
			addImport("database/sql")

			return "sql.NullTime"
		}

		addImport("time")

		return "time.Time"
	case schema.ColumnTypeBoolean:
		if col.Nullable {
			addImport("database/sql")

			return "sql.NullBool"
		}

		return golang.TypeBool
	case schema.ColumnTypeFloat64:
		if col.Nullable {
			addImport("database/sql")

			return "sql.NullFloat64"
		}

		return golang.TypeFloat64
	case schema.ColumnTypeFloat32:
		if col.Nullable {
			addImport("database/sql")

			return golang.Ptr(golang.TypeFloat32)
		}

		return golang.TypeFloat32
	case schema.ColumnTypeBytes:
		if col.Nullable {
			return golang.Ptr(golang.TypeByteSlice)
		}

		return golang.TypeByteSlice
	}

	return golang.TypeString
}
