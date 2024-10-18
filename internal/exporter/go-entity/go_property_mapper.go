package goentity

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

type addImportCallback func(pkg string)

type GoPropertyMapper struct {
}

type GoProperty struct {
	Name       *ds.String
	PluralName string
	Type       string

	Column *schema.Column
}

type goProperties struct {
	List                    []*GoProperty
	MaxPropNameLength       int
	MaxPropPluralNameLength int
	MaxTypeNameLength       int
}

func NewGoPropertyMapper() *GoPropertyMapper {
	return &GoPropertyMapper{}
}

func (m *GoPropertyMapper) mapColumns(columns []*schema.Column, addImportCallback addImportCallback) *goProperties {
	props := &goProperties{
		List: make([]*GoProperty, len(columns)),
	}

	if addImportCallback == nil {
		addImportCallback = func(_ string) {}
	}

	maxNameLength := 0
	maxPluralNameLength := 0
	maxTypeLength := 0
	for i, column := range columns {
		prop := &GoProperty{
			Name:       column.Name.Pascal().FixAbbreviations(goAbbreviationsSet),
			PluralName: column.Name.Pascal().PluralFixAbbreviations(goAbbreviationsPluralsSet).Value,
			Type:       m.mapGoType(column, addImportCallback),
			Column:     column,
		}

		props.List[i] = prop

		if prop.Name.Len() > maxNameLength {
			maxNameLength = column.Name.Pascal().Len()
		}

		if len(prop.PluralName) > maxPluralNameLength {
			maxPluralNameLength = len(prop.PluralName)
		}

		if len(prop.Type) > maxTypeLength {
			maxTypeLength = len(prop.Type)
		}
	}

	props.MaxPropNameLength = maxNameLength
	props.MaxPropPluralNameLength = maxPluralNameLength
	props.MaxTypeNameLength = maxTypeLength

	return props
}

func (m *GoPropertyMapper) mapGoType(col *schema.Column, addImport func(pkg string)) string {
	switch col.PreparedType {
	case schema.ColumnTypeInteger64, schema.ColumnTypeInteger:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullInt64
		}

		return golang.TypeInt64
	case schema.ColumnTypeInteger16:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullInt16
		}

		return golang.TypeInt16
	case schema.ColumnTypeString:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullString
		}

		return golang.TypeString
	case schema.ColumnTypeTimestamp:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullTime
		}

		addImport("time")

		return golang.TypeTimeTime
	case schema.ColumnTypeBoolean:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullBool
		}

		return golang.TypeBool
	case schema.ColumnTypeFloat64:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullFloat64
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

func (p *GoProperty) IsString() bool {
	return p.Type == golang.TypeString
}
