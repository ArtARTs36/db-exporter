package mermaid

import (
	"context"

	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/mermaid"
)

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	for _, table := range params.Schema.Tables.List() {
		diagram := mermaid.NewErDiagram()

		e.addTableToDiagram(diagram, table)

		pages = append(pages, &exporter.ExportedPage{
			FileName: table.Name.Append(".mermaid").Value,
			Content:  []byte(diagram.Build()),
		})
	}

	return pages, nil
}

func (e *Exporter) Export(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	diagram := mermaid.NewErDiagram()

	for _, table := range params.Schema.Tables.List() {
		e.addTableToDiagram(diagram, table)
	}

	return []*exporter.ExportedPage{
		{
			FileName: "er_diagram.mermaid",
			Content:  []byte(diagram.Build()),
		},
	}, nil
}

func (e *Exporter) addTableToDiagram(diagram *mermaid.ErDiagram, table *schema.Table) {
	fields := make([]*mermaid.EntityField, len(table.Columns))

	entityName := table.Name.Pascal().Singular().Value

	entity := &mermaid.Entity{
		Name: entityName,
	}

	for i, column := range table.Columns {
		fields[i] = &mermaid.EntityField{
			Name:     column.Name.Pascal().Singular().Value,
			DataType: e.mapFieldType(column.DataType),
			KeyType:  e.mapKeyType(column),
		}
	}

	entity.Fields = fields

	diagram.AddEntity(entity)

	for _, fk := range table.ForeignKeys {
		diagram.AddRelation(&mermaid.Relation{
			Owner:   entityName,
			Related: fk.ForeignTable.Pascal().Singular().Value,
			Action:  "has",
		})
	}
}

func (e *Exporter) mapFieldType(typ schema.DataType) string {
	switch {
	case typ.IsDate:
		return "date"
	case typ.IsDatetime:
		return "datetime"
	case typ.IsBoolean:
		return "bool"
	case typ.IsJSON:
		return "json"
	case typ.IsStringable:
		return "string"
	}
	return typ.Name
}

func (*Exporter) mapKeyType(col *schema.Column) mermaid.KeyType {
	switch {
	case col.IsPrimaryKey():
		return mermaid.KeyTypePK
	case col.HasForeignKey():
		return mermaid.KeyTypeFK
	case col.IsUniqueKey():
		return mermaid.KeyTypeUK
	default:
		return mermaid.KeyTypeUnspecified
	}
}
