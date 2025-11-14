package sqltype

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

// DBML([a-zA-Z1-9]+)\s+= schema.DataType\{Name: "(.*)"(.*)

var (
	DBMLChar    = schema.DataType{Name: "char", IsStringable: true}
	DBMLVarchar = schema.DataType{Name: "varchar", IsStringable: true}
	DBMLBinary  = schema.DataType{Name: "binary", IsBinary: true}
	DBMLText    = schema.DataType{Name: "text", IsStringable: true}

	DBMLInt       = schema.DataType{Name: "int", IsInteger: true, IsNumeric: true}
	DBMLInteger   = schema.DataType{Name: "integer", IsInteger: true, IsNumeric: true}
	DBMLTimestamp = schema.DataType{Name: "timestamp", IsDatetime: true}
	DBMLUUID      = schema.DataType{Name: "uuid", IsUUID: true, IsStringable: true}

	DBMLFloat  = schema.DataType{Name: "float", IsFloat: true, IsNumeric: true}
	DBMLFloat8 = schema.DataType{Name: "float8", IsFloat: true, IsNumeric: true}

	DBMLBoolean = schema.DataType{Name: "boolean", IsBoolean: true}
	DBMLBool    = schema.DataType{Name: "bool", IsBoolean: true}
)

var dbmlTypeMap = map[string]schema.DataType{
	"char":    DBMLChar,
	"varchar": DBMLVarchar,
	"binary":  DBMLBinary,
	"text":    DBMLText,

	"int":       DBMLInt,
	"integer":   DBMLInteger,
	"timestamp": DBMLTimestamp,
	"uuid":      DBMLUUID,

	"float":  DBMLFloat,
	"float8": DBMLFloat8,

	"boolean": DBMLBoolean,
	"bool":    DBMLBool,
}

func MapDBMLType(name string) schema.DataType {
	return mapType(dbmlTypeMap, name)
}

func mapGoTypeFromDBML(t schema.DataType) golang.Type {
	if t.IsStringable {
		return golang.TypeString
	}

	if t.IsBoolean {
		return golang.TypeBool
	}

	if t.IsInteger {
		return golang.TypeInt
	}

	if t.IsFloat {
		return golang.TypeFloat64
	}

	return golang.TypeString
}
