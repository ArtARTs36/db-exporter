package sqltype

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/golang"
)

// DBML([a-zA-Z1-9]+)\s+= schema.Type\{Name: "(.*)"(.*)

var (
	DBMLChar    = schema.Type{Name: "char", IsStringable: true}
	DBMLVarchar = schema.Type{Name: "varchar", IsStringable: true}
	DBMLBinary  = schema.Type{Name: "binary", IsBinary: true}
	DBMLText    = schema.Type{Name: "text", IsStringable: true}

	DBMLInt       = schema.Type{Name: "int", IsInteger: true, IsNumeric: true}
	DBMLInteger   = schema.Type{Name: "integer", IsInteger: true, IsNumeric: true}
	DBMLTimestamp = schema.Type{Name: "timestamp", IsDatetime: true}
	DBMLUUID      = schema.Type{Name: "uuid", IsUUID: true, IsStringable: true}

	DBMLFloat  = schema.Type{Name: "float", IsFloat: true, IsNumeric: true}
	DBMLFloat8 = schema.Type{Name: "float8", IsFloat: true, IsNumeric: true}

	DBMLBoolean = schema.Type{Name: "boolean", IsBoolean: true}
	DBMLBool    = schema.Type{Name: "bool", IsBoolean: true}
)

var dbmlTypeMap = map[string]schema.Type{
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

func MapDBMLType(name string) schema.Type {
	return mapType(dbmlTypeMap, name)
}

func MapGoTypeFromDBML(t schema.Type) golang.Type {
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
