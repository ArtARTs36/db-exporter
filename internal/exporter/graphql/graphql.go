package graphql

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/graphql"
	"strings"
)

type Exporter struct {
}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) ExportPerFile(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	pages := make([]*exporter.ExportedPage, 0, params.Schema.Tables.Len())

	for _, table := range params.Schema.Tables.List() {
		entity := e.buildEntity(table)

		pages = append(pages, &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s.graphql", table.Name.Value),
			Content:  []byte(entity.Type.Build()),
		})
	}

	return pages, nil
}

func (e *Exporter) Export(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	page := make([]string, 0, params.Schema.Tables.Len())

	for _, table := range params.Schema.Tables.List() {
		entity := e.buildEntity(table)

		page = append(page, entity.Type.Build())
	}

	return []*exporter.ExportedPage{
		{
			FileName: "schema.graphql",
			Content:  []byte(strings.Join(page, "\n\n")),
		},
	}, nil
}

type gEntity struct {
	Type   *graphql.Object
	Inputs struct {
		Create *graphql.Object
	}
}

func (e *Exporter) buildEntity(table *schema.Table) *gEntity {
	entity := &gEntity{
		Type: graphql.NewType(table.Name.Pascal().Singular().String()),
	}

	for _, col := range table.Columns {
		prop := entity.Type.
			AddField(col.Name.Camel().Value).
			Of(e.mapGraphqlPropertyType(col))

		if !col.Nullable {
			prop.Require()
		}

		if col.Comment.IsNotEmpty() {
			prop.Comment(col.Comment.Value)
		}
	}

	return entity
}

func (e *Exporter) mapGraphqlPropertyType(col *schema.Column) graphql.Type {
	if col.PreparedType.IsInteger() {
		return graphql.TypeInt
	}

	if col.PreparedType.IsFloat() {
		return graphql.TypeFloat
	}

	if col.PreparedType.IsBool() {
		return graphql.TypeBoolean
	}

	return graphql.TypeString
}
