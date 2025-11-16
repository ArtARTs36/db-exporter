package graphql

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/graphql"
	"github.com/artarts36/db-exporter/internal/shared/iox"
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
		if table.IsPartition() {
			continue
		}

		entity := e.buildEntity(table)

		w := iox.NewWriter()
		entity.Type.Build(w)

		pages = append(pages, &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s.graphql", table.Name.Value),
			Content:  w.Bytes(),
		})
	}

	return pages, nil
}

func (e *Exporter) Export(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	file := graphql.File{
		Types: make([]*graphql.Object, 0, params.Schema.Tables.Len()),
	}

	for _, table := range params.Schema.Tables.List() {
		if table.IsPartition() {
			continue
		}

		entity := e.buildEntity(table)
		file.Types = append(file.Types, entity.Type)
	}

	w := iox.NewWriter()
	file.Build(w)

	return []*exporter.ExportedPage{
		{
			FileName: "schema.graphql",
			Content:  w.Bytes(),
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
	switch {
	case col.DataType.IsInteger:
		return graphql.TypeInt
	case col.DataType.IsFloat:
		return graphql.TypeFloat
	case col.DataType.IsBoolean:
		return graphql.TypeBoolean
	}

	return graphql.TypeString
}
