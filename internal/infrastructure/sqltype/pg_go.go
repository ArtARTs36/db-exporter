package sqltype

import (
	"github.com/artarts36/db-exporter/internal/schema"
)

var pgGoTypeMap = map[schema.Type]schema.DataType{
	PGText:             schema.DataTypeString,
	PGUUID:             schema.DataTypeString,
	PGCharacter:        schema.DataTypeString,
	PGCharacterVarying: schema.DataTypeString,

	PGTimestampWithoutTZ: schema.DataTypeTimestamp,
	PGTimestampWithTZ:    schema.DataTypeTimestamp,

	PGInteger: schema.DataTypeInteger,
	PGInt4:    schema.DataTypeInteger,
	PGInt8:    schema.DataTypeInteger,
	PGSerial:  schema.DataTypeInteger,

	PGSmallInt:    schema.DataTypeInteger16,
	PGSmallSerial: schema.DataTypeInteger16,

	PGBigint: schema.DataTypeInteger64,

	PGBoolean: schema.DataTypeBoolean,
	PGBit:     schema.DataTypeBoolean,

	PGDoublePrecision: schema.DataTypeFloat32,
	PGFloat8:          schema.DataTypeFloat32,
	PGDecimal:         schema.DataTypeFloat32,

	PGMoney:   schema.DataTypeFloat64,
	PGReal:    schema.DataTypeFloat64,
	PGNumeric: schema.DataTypeFloat64,

	PGBytea: schema.DataTypeBytes,
}
