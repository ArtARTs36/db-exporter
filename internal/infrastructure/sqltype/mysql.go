package sqltype

import "github.com/artarts36/db-exporter/internal/schema"

var (
	MySQLChar     = schema.Type{Name: "char", IsStringable: true}
	MySQLVarchar  = schema.Type{Name: "varchar", Length: "256", IsStringable: true}
	MySQLBinary   = schema.Type{Name: "binary"}
	MySQLText     = schema.Type{Name: "text", IsStringable: true}
	MySQLLongText = schema.Type{Name: "longtext", IsStringable: true}

	MySQLInt       = schema.Type{Name: "int", IsInteger: true, IsNumeric: true}
	MySQLInteger   = schema.Type{Name: "integer", IsInteger: true, IsNumeric: true}
	MySQLSmallInt  = schema.Type{Name: "smallint", IsInteger: true, IsNumeric: true}
	MySQLBigInt    = schema.Type{Name: "bigint", IsInteger: true, IsNumeric: true}
	MySQLTinyint   = schema.Type{Name: "tinyint", IsInteger: true, IsNumeric: true}
	MySQLMediumint = schema.Type{Name: "mediumint", IsInteger: true, IsNumeric: true}

	MySQLFloat   = schema.Type{Name: "float", IsFloat: true, IsNumeric: true}
	MySQLDouble  = schema.Type{Name: "double", IsFloat: true, IsNumeric: true}
	MySQLDecimal = schema.Type{Name: "decimal", IsFloat: true, IsNumeric: true}

	MySQLBit = schema.Type{Name: "bit"}

	MySQLTimestamp = schema.Type{Name: "timestamp"}
	MySQLDate      = schema.Type{Name: "date", IsDate: true}
	MySQLTime      = schema.Type{Name: "time"}
	MySQLDateTime  = schema.Type{Name: "datetime", IsDatetime: true}

	MySQLLongBlob = schema.Type{Name: "longblob"}

	MySQLLineString = schema.Type{Name: "linestring"}
	MySQLPolygon    = schema.Type{Name: "polygon"}

	MySQLPseudoUUID = MySQLVarchar.WithLength("36")
)

var mysqlTypeMap = map[string]schema.Type{
	"char":     MySQLChar,
	"varchar":  MySQLVarchar,
	"binary":   MySQLBinary,
	"text":     MySQLText,
	"longtext": MySQLLongText,

	"int":       MySQLInt,
	"integer":   MySQLInteger,
	"smallint":  MySQLSmallInt,
	"bigint":    MySQLBigInt,
	"tinyint":   MySQLTinyint,
	"mediumint": MySQLMediumint,

	"float":   MySQLFloat,
	"double":  MySQLDouble,
	"decimal": MySQLDecimal,

	"bit": MySQLBit,

	"timestamp": MySQLTimestamp,
	"date":      MySQLDate,
	"time":      MySQLTime,
	"datetime":  MySQLDateTime,

	"longblob": MySQLLongBlob,

	"linestring": MySQLLineString,
	"polygon":    MySQLPolygon,
}

func MapMySQLType(name string) schema.Type {
	return mapType(mysqlTypeMap, name)
}
