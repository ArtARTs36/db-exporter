package jsonschema

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/jsonschema"
)

func mapJSONType(col *schema.Column) jsonschema.Type {
	switch true {
	case col.Type.IsNumeric:
		return jsonschema.TypeNumber
	case col.Type.IsStringable:
		return jsonschema.TypeString
	case col.Type.IsBoolean:
		return jsonschema.TypeBoolean
	default:
		return jsonschema.TypeString
	}
}
