package sqltype

import "github.com/artarts36/db-exporter/internal/schema"

// PG([a-zA-Z1-9]+)\s+= schema.Type\{Name: "(.*)"

var (
	PGText             = schema.Type{Name: "text", IsStringable: false}
	PGUUID             = schema.Type{Name: "uuid", IsStringable: true, IsUUID: true}
	PGCharacter        = schema.Type{Name: "character", IsStringable: true}
	PGChar             = schema.Type{Name: "char", IsStringable: true}
	PGCharacterVarying = schema.Type{Name: "character varying", IsStringable: true}

	PGTimestampWithoutTZ = schema.Type{Name: "timestamp without time zone", IsDatetime: true}
	PGTimestampWithTZ    = schema.Type{Name: "timestamp with time zone", IsDatetime: true}
	PGTimeWithoutTZ      = schema.Type{Name: "time without time zone", IsDatetime: true}
	PGTimeWithTZ         = schema.Type{Name: "time with time zone", IsDatetime: true}
	PGDate               = schema.Type{Name: "date", IsDate: true}
	PGInterval           = schema.Type{Name: "interval"}

	PGBoolean = schema.Type{Name: "boolean", IsBoolean: true}
	PGBit     = schema.Type{Name: "bit"}

	PGBytea = schema.Type{Name: "bytea"}

	PGInteger     = schema.Type{Name: "integer", IsInteger: true, IsNumeric: true}
	PGBigint      = schema.Type{Name: "bigint", IsInteger: true, IsNumeric: true}
	PGInt         = schema.Type{Name: "int", IsInteger: true, IsNumeric: true}
	PGInt4        = schema.Type{Name: "int4", IsInteger: true, IsNumeric: true}
	PGInt8        = schema.Type{Name: "int8", IsInteger: true, IsNumeric: true}
	PGSmallInt    = schema.Type{Name: "smallint", IsInteger: true, IsNumeric: true}
	PGSmallSerial = schema.Type{Name: "smallserial", IsInteger: true, IsNumeric: true}
	PGSerial      = schema.Type{Name: "serial", IsInteger: true, IsNumeric: true}
	PGBigSerial   = schema.Type{Name: "bigserial", IsInteger: true, IsNumeric: true}

	PGMoney   = schema.Type{Name: "money"}
	PGNumeric = schema.Type{Name: "numeric", IsNumeric: true}

	PGReal            = schema.Type{Name: "real", IsFloat: true, IsNumeric: true}
	PGDoublePrecision = schema.Type{Name: "double precision", IsFloat: true, IsNumeric: true}
	PGFloat8          = schema.Type{Name: "float8", IsFloat: true, IsNumeric: true}
	PGDecimal         = schema.Type{Name: "decimal", IsFloat: true, IsNumeric: true}

	PGCidr    = schema.Type{Name: "cidr"}
	PGInet    = schema.Type{Name: "inet"}
	PGMacaddr = schema.Type{Name: "macaddr"}

	PGXML  = schema.Type{Name: "xml"}
	PGJSON = schema.Type{Name: "json"}

	PGTSVector = schema.Type{Name: "tsvector"}
	PGTSQuery  = schema.Type{Name: "tsquery"}

	PGArray = schema.Type{Name: "array"}

	PGPoint        = schema.Type{Name: "point"}
	PGLine         = schema.Type{Name: "lint"}
	PGLseq         = schema.Type{Name: "lseq"}
	PGBox          = schema.Type{Name: "box"}
	PGPath         = schema.Type{Name: "path"}
	PGPolygon      = schema.Type{Name: "polygon"}
	PGCircle       = schema.Type{Name: "circle"}
	PGTxidSnapshot = schema.Type{Name: "txid_snapshot"}
)

var pgTypeMap = map[string]schema.Type{
	"text":              PGText,
	"uuid":              PGUUID,
	"character":         PGCharacter,
	"char":              PGChar,
	"character varying": PGCharacterVarying,

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

	"xml":  PGXML,
	"json": PGJSON,

	"tsvector": PGTSVector,
	"tsquery":  PGTSQuery,

	"array": PGArray,

	"point":         PGPoint,
	"lint":          PGLine,
	"lseq":          PGLseq,
	"box":           PGBox,
	"path":          PGPath,
	"polygon":       PGPolygon,
	"circle":        PGCircle,
	"txid_snapshot": PGTxidSnapshot,
}

func MapPGType(name string) schema.Type {
	return mapType(pgTypeMap, name)
}
