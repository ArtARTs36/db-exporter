package sqltype

import "github.com/artarts36/db-exporter/internal/schema"

var (
	MySQLChar     = schema.DataType{Name: "char", IsStringable: true}
	MySQLVarchar  = schema.DataType{Name: "varchar", Length: "256", IsStringable: true}
	MySQLBinary   = schema.DataType{Name: "binary"}
	MySQLText     = schema.DataType{Name: "text", IsStringable: true}
	MySQLLongText = schema.DataType{Name: "longtext", IsStringable: true}

	MySQLInt       = schema.DataType{Name: "int", IsInteger: true, IsNumeric: true}
	MySQLInteger   = schema.DataType{Name: "integer", IsInteger: true, IsNumeric: true}
	MySQLSmallInt  = schema.DataType{Name: "smallint", IsInteger: true, IsNumeric: true}
	MySQLBigInt    = schema.DataType{Name: "bigint", IsInteger: true, IsNumeric: true}
	MySQLTinyint   = schema.DataType{Name: "tinyint", IsInteger: true, IsNumeric: true}
	MySQLMediumint = schema.DataType{Name: "mediumint", IsInteger: true, IsNumeric: true}

	MySQLFloat   = schema.DataType{Name: "float", IsFloat: true, IsNumeric: true}
	MySQLDouble  = schema.DataType{Name: "double", IsFloat: true, IsNumeric: true}
	MySQLDecimal = schema.DataType{Name: "decimal", IsFloat: true, IsNumeric: true}

	MySQLBit = schema.DataType{Name: "bit"}

	MySQLTimestamp = schema.DataType{Name: "timestamp"}
	MySQLDate      = schema.DataType{Name: "date", IsDate: true}
	MySQLTime      = schema.DataType{Name: "time"}
	MySQLDateTime  = schema.DataType{Name: "datetime", IsDatetime: true}

	MySQLLongBlob = schema.DataType{Name: "longblob"}

	MySQLLineString = schema.DataType{Name: "linestring"}
	MySQLPolygon    = schema.DataType{Name: "polygon"}

	MySQLPseudoUUID = MySQLVarchar.WithLength("36")
)

var mysqlTypeMap = map[string]schema.DataType{
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

func MapMySQLType(name string) schema.DataType {
	return mapType(mysqlTypeMap, name)
}
