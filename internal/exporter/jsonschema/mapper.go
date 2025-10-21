package jsonschema

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/jsonschema"
)

type mapper struct{}

func (e *mapper) mapJSONSchema(
	spec *config.JSONSchemaExportSpec,
	tables []*schema.Table,
) *jsonschema.Schema {
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

	return sch
}

func (e *mapper) prepareDefaultValue(col *schema.Column) interface{} {
	if col.Default == nil || col.Default.Type != schema.ColumnDefaultTypeValue {
		return nil
	}

	return col.Default.Value
}

func (e *mapper) mapFormat(column *schema.Column) jsonschema.Format {
	switch {
	case column.Type.IsDatetime:
		return jsonschema.FormatDateTime
	case column.Type.IsUUID:
		return jsonschema.FormatUUID
	case column.Type.IsDate:
		return jsonschema.FormatDate
	case column.Type.IsStringable:
		if column.Name.Equal("email") || column.Name.Ends("_email") {
			return jsonschema.FormatEmail
		} else if column.Name.Equal("uri") || column.Name.Ends("_uri") {
			return jsonschema.FormatURI
		}
	}

	return jsonschema.FormatUnknown
}
