package goentity

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/artarts36/gds"
)

type addImportCallback func(pkg string)

type GoPropertyMapper struct {
}

type GoProperty struct {
	Name       *gds.String
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

func (m *GoPropertyMapper) mapColumns(
	columns []*schema.Column,
	enums map[string]*golang.StringEnum,
	addImportCallback addImportCallback,
) *goProperties {
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
			Type:       m.mapGoType(column, enums, addImportCallback),
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

func (m *GoPropertyMapper) mapGoType(
	col *schema.Column,
	enums map[string]*golang.StringEnum,
	addImport func(pkg string),
) string {
	if e, ok := enums[col.Name.Value]; ok {
		return e.Name.Value
	}

	switch col.PreparedType {
	case schema.DataTypeInteger64, schema.DataTypeInteger:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullInt64
		}

		return golang.TypeInt64
	case schema.DataTypeInteger16:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullInt16
		}

		return golang.TypeInt16
	case schema.DataTypeString:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullString
		}

		return golang.TypeString
	case schema.DataTypeTimestamp:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullTime
		}

		addImport("time")

		return golang.TypeTimeTime
	case schema.DataTypeBoolean:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullBool
		}

		return golang.TypeBool
	case schema.DataTypeFloat64:
		if col.Nullable {
			addImport("database/sql")

			return golang.TypeSQLNullFloat64
		}

		return golang.TypeFloat64
	case schema.DataTypeFloat32:
		if col.Nullable {
			addImport("database/sql")

			return golang.Ptr(golang.TypeFloat32)
		}

		return golang.TypeFloat32
	case schema.DataTypeBytes:
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
