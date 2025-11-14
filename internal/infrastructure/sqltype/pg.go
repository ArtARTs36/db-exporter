package sqltype

import "github.com/artarts36/db-exporter/internal/schema"

// PG([a-zA-Z1-9]+)\s+= schema.DataType\{Name: "(.*)"

var (
	PGText             = schema.DataType{Name: "text", IsStringable: false}
	PGUUID             = schema.DataType{Name: "uuid", IsStringable: true, IsUUID: true}
	PGCharacter        = schema.DataType{Name: "character", IsStringable: true}
	PGChar             = schema.DataType{Name: "char", IsStringable: true}
	PGCharacterVarying = schema.DataType{Name: "character varying", IsStringable: true}
	PGBpchar           = schema.DataType{Name: "bpchar", IsStringable: true}

	PGTimestampWithoutTZ = schema.DataType{Name: "timestamp without time zone", IsDatetime: true}
	PGTimestampWithTZ    = schema.DataType{Name: "timestamp with time zone", IsDatetime: true}
	PGTimeWithoutTZ      = schema.DataType{Name: "time without time zone", IsDatetime: true}
	PGTimeWithTZ         = schema.DataType{Name: "time with time zone", IsDatetime: true}
	PGDate               = schema.DataType{Name: "date", IsDate: true}
	PGInterval           = schema.DataType{Name: "interval", IsInterval: true}

	PGBoolean = schema.DataType{Name: "boolean", IsBoolean: true}
	PGBit     = schema.DataType{Name: "bit"}

	PGBytea = schema.DataType{Name: "bytea"}

	PGInteger     = schema.DataType{Name: "integer", IsInteger: true, IsNumeric: true}
	PGBigint      = schema.DataType{Name: "bigint", IsInteger: true, IsNumeric: true}
	PGInt         = schema.DataType{Name: "int", IsInteger: true, IsNumeric: true}
	PGInt4        = schema.DataType{Name: "int4", IsInteger: true, IsNumeric: true}
	PGInt8        = schema.DataType{Name: "int8", IsInteger: true, IsNumeric: true}
	PGSmallInt    = schema.DataType{Name: "smallint", IsInteger: true, IsNumeric: true}
	PGSmallSerial = schema.DataType{Name: "smallserial", IsInteger: true, IsNumeric: true}
	PGSerial      = schema.DataType{Name: "serial", IsInteger: true, IsNumeric: true}
	PGBigSerial   = schema.DataType{Name: "bigserial", IsInteger: true, IsNumeric: true}

	PGMoney   = schema.DataType{Name: "money"}
	PGNumeric = schema.DataType{Name: "numeric", IsNumeric: true}

	PGReal            = schema.DataType{Name: "real", IsFloat: true, IsNumeric: true}
	PGDoublePrecision = schema.DataType{Name: "double precision", IsFloat: true, IsNumeric: true}
	PGFloat8          = schema.DataType{Name: "float8", IsFloat: true, IsNumeric: true}
	PGDecimal         = schema.DataType{Name: "decimal", IsFloat: true, IsNumeric: true}

	PGCidr    = schema.DataType{Name: "cidr"}
	PGInet    = schema.DataType{Name: "inet"}
	PGMacaddr = schema.DataType{Name: "macaddr"}

	PGXML   = schema.DataType{Name: "xml"}
	PGJSON  = schema.DataType{Name: "json", IsJSON: true}
	PGJSONB = schema.DataType{Name: "jsonb", IsJSON: true}

	PGTSVector = schema.DataType{Name: "tsvector"}
	PGTSQuery  = schema.DataType{Name: "tsquery"}

	PGArray = schema.DataType{Name: "array"}

	PGPoint        = schema.DataType{Name: "point"}
	PGLine         = schema.DataType{Name: "line"}
	PGLseg         = schema.DataType{Name: "lseg"}
	PGBox          = schema.DataType{Name: "box"}
	PGPath         = schema.DataType{Name: "path"}
	PGPolygon      = schema.DataType{Name: "polygon"}
	PGCircle       = schema.DataType{Name: "circle"}
	PGTxidSnapshot = schema.DataType{Name: "txid_snapshot"}
)

var pgTypeMap = map[string]schema.DataType{
	"text":              PGText,
	"uuid":              PGUUID,
	"character":         PGCharacter,
	"char":              PGChar,
	"character varying": PGCharacterVarying,
	"bpchar":            PGBpchar,

	"timestamp without time zone": PGTimestampWithoutTZ,
	"timestamp with time zone":    PGTimestampWithTZ,
	"time without time zone":      PGTimeWithoutTZ,
	"time with time zone":         PGTimeWithTZ,
	"date":                        PGDate,
	"interval":                    PGInterval,

	"boolean": PGBoolean,
	"bit":     PGBit,

	"bytea": PGBytea,

	"integer":     PGInteger,
	"bigint":      PGBigint,
	"int":         PGInt,
	"int4":        PGInt4,
	"int8":        PGInt8,
	"smallint":    PGSmallInt,
	"smallserial": PGSmallSerial,
	"serial":      PGSerial,
	"bigserial":   PGBigSerial,

	"money":   PGMoney,
	"numeric": PGNumeric,

	"real":             PGReal,
	"double precision": PGDoublePrecision,
	"float8":           PGFloat8,
	"decimal":          PGDecimal,

	"cidr":    PGCidr,
	"inet":    PGInet,
	"macaddr": PGMacaddr,

	"xml":   PGXML,
	"json":  PGJSON,
	"jsonb": PGJSONB,

	"tsvector": PGTSVector,
	"tsquery":  PGTSQuery,

	"array": PGArray,

	"point":         PGPoint,
	"lint":          PGLine,
	"lseg":          PGLseg,
	"box":           PGBox,
	"path":          PGPath,
	"polygon":       PGPolygon,
	"circle":        PGCircle,
	"txid_snapshot": PGTxidSnapshot,
}

func MapPGType(name string) schema.DataType {
	return mapType(pgTypeMap, name)
}
