package jsonschema

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/jsonschema"
)

func mapJSONType(col *schema.Column) jsonschema.Type {
	switch col.PreparedType {
	case schema.ColumnTypeInteger, schema.ColumnTypeInteger16, schema.ColumnTypeInteger64:
		return jsonschema.TypeInteger
	case schema.ColumnTypeString, schema.ColumnTypeTimestamp:
		return jsonschema.TypeString
	case schema.ColumnTypeBoolean:
		return jsonschema.TypeBoolean
	case schema.ColumnTypeFloat32, schema.ColumnTypeFloat64:
		return jsonschema.TypeNumber
	default:
		return jsonschema.TypeString
	}
}
