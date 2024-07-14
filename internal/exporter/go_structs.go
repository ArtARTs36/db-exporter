package exporter

import (
	"context"
	"fmt"
	"strings"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/shared/golang"
	"github.com/artarts36/db-exporter/internal/template"
)

const GoStructsExporterName = "go-structs"

var goAbbreviationsSet = map[string]bool{
	"id":   true,
	"uuid": true,
	"json": true,
	"db":   true,
}

type GoStructsExporter struct {
	unimplementedImporter
	renderer *template.Renderer
}

type goSchema struct {
	Tables  []*goStruct
	Imports *ds.Set
}

type goStruct struct {
	Name       ds.String
	Properties []*goProperty
}

type goProperty struct {
	Name       string
	NameOffset int
	ColumnName string
	Type       string
	TypeOffset int
	Pointer    bool

	Column *schema.Column
}

func NewGoStructsExporter(renderer *template.Renderer) Exporter {
	return &GoStructsExporter{
		renderer: renderer,
	}
}

func (e *GoStructsExporter) ExportPerFile(
	_ context.Context,
	sch *schema.Schema,
	params *ExportParams,
) ([]*ExportedPage, error) {
	pages := make([]*ExportedPage, 0, sch.Tables.Len())
	pkg := e.selectPackage(params)

	for _, table := range sch.Tables.List() {
		goSch := e.makeGoSchema([]*schema.Table{
			table,
		})

		page, err := render(
			e.renderer,
			"go-structs/model.go.tpl",
			fmt.Sprintf("%s.go", table.Name.Singular().Lower()),
			map[string]stick.Value{
				"schema":  goSch,
				"package": pkg,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to render: %w", err)
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func (e *GoStructsExporter) Export(
	_ context.Context,
	schema *schema.Schema,
	params *ExportParams,
) ([]*ExportedPage, error) {
	goSch := e.makeGoSchema(schema.Tables.List())
	pkg := e.selectPackage(params)

	page, err := render(e.renderer, "go-structs/model.go.tpl", "models.go", map[string]stick.Value{
		"schema":  goSch,
		"package": pkg,
	})
	if err != nil {
		return nil, err
	}

	return []*ExportedPage{
		page,
	}, nil
}

func (e *GoStructsExporter) selectPackage(params *ExportParams) string {
	if params.Package != "" {
		return strings.ToLower(params.Package)
	}

	return "models"
}

func (e *GoStructsExporter) mapGoType(col *schema.Column, imports *ds.Set) string {
	switch col.PreparedType {
	case schema.ColumnTypeInteger64, schema.ColumnTypeInteger:
		if col.Nullable {
			imports.Add("database/sql")

			return "sql.NullInt64"
		}

		return golang.TypeInt64
	case schema.ColumnTypeInteger16:
		if col.Nullable {
			imports.Add("database/sql")

			return "sql.NullInt16"
		}

		return golang.TypeInt16
	case schema.ColumnTypeString:
		if col.Nullable {
			imports.Add("database/sql")

			return "sql.NullString"
		}

		return golang.TypeString
	case schema.ColumnTypeTimestamp:
		if col.Nullable {
			imports.Add("database/sql")

			return "sql.NullTime"
		}

		imports.Add("time")

		return "time.Time"
	case schema.ColumnTypeBoolean:
		if col.Nullable {
			imports.Add("database/sql")

			return "sql.NullBool"
		}

		return golang.TypeBool
	case schema.ColumnTypeFloat64:
		if col.Nullable {
			imports.Add("database/sql")

			return "sql.NullFloat64"
		}

		return golang.TypeFloat64
	case schema.ColumnTypeFloat32:
		if col.Nullable {
			imports.Add("database/sql")

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

func (e *GoStructsExporter) makeGoSchema(tables []*schema.Table) *goSchema {
	goSch := &goSchema{
		Tables:  make([]*goStruct, 0, len(tables)),
		Imports: ds.NewSet(),
	}

	for _, t := range tables {
		str := &goStruct{
			Name:       *t.Name.Singular().Pascal().FixAbbreviations(goAbbreviationsSet),
			Properties: make([]*goProperty, 0, len(t.Columns)),
		}

		goSch.Tables = append(goSch.Tables, str)

		propNameOffset := 0
		propTypeOffset := 0
		for _, c := range t.Columns {
			prop := &goProperty{
				Name:       c.Name.Pascal().FixAbbreviations(goAbbreviationsSet).Val,
				Type:       e.mapGoType(c, goSch.Imports),
				ColumnName: c.Name.Val,
			}

			str.Properties = append(str.Properties, prop)

			if len(prop.Name) > propNameOffset {
				propNameOffset = c.Name.Pascal().Len()
			}

			if len(prop.Type) > propTypeOffset {
				propTypeOffset = len(prop.Type)
			}
		}

		for _, prop := range str.Properties {
			prop.NameOffset = propNameOffset - len(prop.Name)
			prop.TypeOffset = propTypeOffset - len(prop.Type)
		}
	}

	return goSch
}
