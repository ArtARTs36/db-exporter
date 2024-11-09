package jsonschema

import (
	"context"
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/exporter"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/jsonschema"
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
		content, err := e.buildJSONSchema(params, []*schema.Table{table})
		if err != nil {
			return nil, fmt.Errorf("failed to build json schema for table %q: %w", table.Name, err)
		}

		pages = append(pages, &exporter.ExportedPage{
			FileName: fmt.Sprintf("%s.json", table.Name.Lower()),
			Content:  content,
		})
	}

	return pages, nil
}

func (e *Exporter) Export(
	_ context.Context,
	params *exporter.ExportParams,
) ([]*exporter.ExportedPage, error) {
	content, err := e.buildJSONSchema(params, params.Schema.Tables.List())
	if err != nil {
		return nil, fmt.Errorf("failed to build json schema: %w", err)
	}

	return []*exporter.ExportedPage{
		{
			FileName: "schema.json",
			Content:  content,
		},
	}, nil
}

func (e *Exporter) buildJSONSchema(params *exporter.ExportParams, tables []*schema.Table) ([]byte, error) {
	spec, ok := params.Spec.(*config.JSONSchemaExportSpec)
	if !ok {
		return nil, errors.New("got invalid spec")
	}

	sch := jsonschema.Draft04()
	sch.Title = spec.Schema.Title
	sch.Description = spec.Schema.Description

	for _, table := range tables {
		tableObject := jsonschema.NewProperty(jsonschema.TypeObject)

		required := make([]string, 0, len(table.Columns))

		for _, column := range table.Columns {
			jsonType := mapJSONType(column)

			colProp := jsonschema.NewProperty(jsonType)
			colProp.Description = column.Comment.Value
			colProp.Format = e.mapFormat(column)
			colProp.Default = e.prepareDefaultValue(column)

			if !column.Nullable {
				required = append(required, column.Name.Value)
			}

			if column.Enum != nil {
				colProp.Enum = column.Enum.Values
			}

			tableObject.Properties[column.Name.Value] = colProp
		}

		tableObject.Required = required
		tableObject.DisableAdditionalProperties()

		sch.Properties[table.Name.Value] = jsonschema.Property{
			Ref: sch.AddDefinition(table.Name.Value, tableObject),
		}
	}

	marshaller := sch.Marshal
	if spec.Pretty {
		marshaller = sch.MarshallPretty
	}

	content, err := marshaller()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	return content, nil
}

func (e *Exporter) prepareDefaultValue(col *schema.Column) interface{} {
	if col.Default == nil || col.Default.Type != schema.ColumnDefaultTypeValue {
		return nil
	}

	return col.Default.Value
}

func (e *Exporter) mapFormat(column *schema.Column) jsonschema.Format {
	if column.PreparedType == schema.DataTypeTimestamp {
		return jsonschema.FormatDateTime
	} else if column.PreparedType == schema.DataTypeString {
		if column.Type.IsUUID { //nolint:gocritic // not need
			return jsonschema.FormatUUID
		} else if column.Name.Equal("email") || column.Name.Ends("_email") {
			return jsonschema.FormatEmail
		} else if column.Name.Equal("uri") || column.Name.Ends("_uri") {
			return jsonschema.FormatURI
		}
	}

	return jsonschema.FormatUnknown
}
