package sqltype

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

var pgGoTypeMap = map[schema.Type]golang.Type{
	PGInteger: golang.TypeInt,
	PGInt4:    golang.TypeInt,
	PGInt8:    golang.TypeInt,
	PGSerial:  golang.TypeInt,

	PGSmallInt:    golang.TypeInt16,
	PGSmallSerial: golang.TypeInt16,

	PGBigint: golang.TypeInt64,

	PGBoolean: golang.TypeBool,
	PGBit:     golang.TypeBool,

	PGDoublePrecision: golang.TypeFloat32,
	PGFloat8:          golang.TypeFloat32,
	PGDecimal:         golang.TypeFloat32,

	PGMoney:   golang.TypeFloat64,
	PGReal:    golang.TypeFloat64,
	PGNumeric: golang.TypeFloat64,

	PGBytea: golang.TypeByteSlice,
}

func mapGoTypeFromPG(t schema.Type) golang.Type {
	if t.IsStringable {
		return golang.TypeString
	}

	if t.IsDatetime {
		return golang.TypeTimeTime
	}

	if t.IsBoolean {
		return golang.TypeBool
	}

	dt, ok := pgGoTypeMap[t]
	if ok {
		return dt
	}

	if t.IsInteger {
		return golang.TypeInt
	}

	if t.IsFloat {
		return golang.TypeFloat64
	}

	return golang.TypeString
}
