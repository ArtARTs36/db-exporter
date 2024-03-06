package exporter

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/ds"

	"github.com/tyler-sommer/stick"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/template"
)

var goAbbreviationsSet = map[string]bool{
	"id":   true,
	"uuid": true,
	"json": true,
}

type GoStructsExporter struct {
	renderer *template.Renderer
}

type goSchema struct {
	Tables []*goStruct
}

type goStruct struct {
	Name       schema.String
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

func (e *GoStructsExporter) ExportPerFile(_ context.Context, sc *schema.Schema, params *ExportParams) ([]*ExportedPage, error) {
	return nil, fmt.Errorf("export per file unsupported")
}

func (e *GoStructsExporter) Export(_ context.Context, schema *schema.Schema, _ *ExportParams) ([]*ExportedPage, error) {
	goSch := &goSchema{
		Tables: make([]*goStruct, 0, len(schema.Tables)),
	}

	imports := ds.NewSet()

	for _, t := range schema.Tables {
		str := &goStruct{
			Name:       *t.Name.Singular(),
			Properties: make([]*goProperty, 0, len(t.Columns)),
		}

		goSch.Tables = append(goSch.Tables, str)

		propNameOffset := 0
		propTypeOffset := 0
		for _, c := range t.Columns {
			prop := &goProperty{
				Name:       c.Name.Pascal().FixAbbreviations(goAbbreviationsSet).Value,
				Type:       e.mapGoType(c, imports),
				ColumnName: c.Name.Value,
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

	content, err := e.renderer.Render("gostructs/single-models.tpl", map[string]stick.Value{
		"schema":      goSch,
		"imports":     imports.List(),
		"has_imports": imports.Valid(),
	})
	if err != nil {
		return nil, err
	}

	return []*ExportedPage{
		{
			FileName: "models.go",
			Content:  content,
		},
	}, nil
}

func (e *GoStructsExporter) mapGoType(col *schema.Column, imports *ds.Set) string {
	switch col.PreparedType {
	case schema.ColumnTypeInteger:
		if col.Nullable {
			imports.Add("database/sql")

			return "sql.NullInt64"
		}

		return "int64"
	case schema.ColumnTypeString:
		if col.Nullable {
			imports.Add("database/sql")

			return "sql.NullString"
		}

		return "string"
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

			return "*sql.NullBool"
		}

		return "bool"
	case schema.ColumnTypeFloat:
		if col.Nullable {
			imports.Add("database/sql")

			return "sql.NullFloat64"
		}

		return "float64"
	default:
		return "string"
	}
}
