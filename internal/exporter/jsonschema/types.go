package jsonschema

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/jsonschema"
)

func mapJSONType(col *schema.Column) jsonschema.Type {
	switch {
	case col.DataType.IsNumeric:
		return jsonschema.TypeNumber
	case col.DataType.IsStringable:
		return jsonschema.TypeString
	case col.DataType.IsBoolean:
		return jsonschema.TypeBoolean
	default:
		return jsonschema.TypeString
	}
}
