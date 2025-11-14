package goentity

import (
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
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
	sourceDriver schema.DatabaseDriver,
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
			Type:       m.mapGoType(sourceDriver, column, enums, addImportCallback),
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
	sourceDriver schema.DatabaseDriver,
	col *schema.Column,
	enums map[string]*golang.StringEnum,
	addImport func(pkg string),
) string {
	if e, ok := enums[col.Name.Value]; ok {
		return e.Name.Value
	}

	goType := sqltype.MapGoType(sourceDriver, col.DataType)

	if !col.Nullable {
		if goType.PackagePath != "" {
			addImport(goType.PackagePath)
		}

		return goType.Call()
	}

	if goType.Null != nil {
		if goType.Null.PackagePath != "" {
			addImport(goType.Null.PackagePath)
		}

		return goType.Null.Call()
	}

	return golang.Ptr(goType.Call())
}

func (p *GoProperty) IsString() bool {
	return p.Type == golang.TypeString.Name
}
