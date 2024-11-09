package sqltype

import (
	"github.com/artarts36/db-exporter/internal/schema"
)

var pgGoTypeMap = map[schema.Type]schema.DataType{
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

func MapGoTypeFromPG(t schema.Type) schema.DataType {
	if t.IsStringable {
		return schema.DataTypeString
	}

	if t.IsDatetime {
		return schema.DataTypeTimestamp
	}

	if t.IsBoolean {
		return schema.DataTypeBoolean
	}

	dt, ok := pgGoTypeMap[t]
	if ok {
		return dt
	}

	if t.IsInteger {
		return schema.DataTypeInteger
	}

	if t.IsFloat {
		return schema.DataTypeFloat64
	}

	return schema.DataTypeString
}
