package jsonschema

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/jsonschema"
)

func mapJSONType(col *schema.Column) jsonschema.Type {
	switch col.PreparedType {
	case schema.DataTypeInteger, schema.DataTypeInteger16, schema.DataTypeInteger64:
		return jsonschema.TypeInteger
	case schema.DataTypeString, schema.DataTypeTimestamp:
		return jsonschema.TypeString
	case schema.DataTypeBoolean:
		return jsonschema.TypeBoolean
	case schema.DataTypeFloat32, schema.DataTypeFloat64:
		return jsonschema.TypeNumber
	case schema.DataTypeBytes:
		return jsonschema.TypeString
	default:
		return jsonschema.TypeString
	}
}
