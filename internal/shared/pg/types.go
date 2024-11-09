package pg

const (
	TypeText             = "text"
	TypeUUID             = "uuid"
	TypeCharacter        = "character"
	TypeChar             = "char"
	TypeCharacterVarying = "character varying"

	TypeTimestampWithoutTZ = "timestamp without time zone"
	TypeTimestampWithTZ    = "timestamp with time zone"
	TypeTimeWithoutTZ      = "time without time zone"
	TypeTimeWithTZ         = "time with time zone"
	TypeDate               = "date"
	TypeInterval           = "interval"

	TypeBoolean = "boolean"
	TypeBit     = "bit"

	TypeBytea = "bytea"

	TypeInteger     = "integer"
	TypeBigint      = "bigint"
	TypeInt         = "int"
	TypeInt4        = "int4"
	TypeInt8        = "int8"
	TypeSmallInt    = "smallint"
	TypeSmallSerial = "smallserial"
	TypeSerial      = "serial"
	TypeBigSerial   = "bigserial"

	TypeMoney   = "money"
	TypeNumeric = "numeric"

	TypeReal            = "real"
	TypeDoublePrecision = "double precision"
	TypeFloat8          = "float8"
	TypeDecimal         = "decimal"

	TypeCidr    = "cidr"
	TypeInet    = "inet"
	TypeMacaddr = "macaddr"

	TypeXML  = "xml"
	TypeJSON = "json"

	TypeTsVector = "tsvector"
	TypeTsQuery  = "tsquery"

	TypeArray = "array"

	TypePoint    = "point"
	TypeLine     = "lint"
	TypeLseq     = "lseq"
	TypeBox      = "box"
	TypePath     = "path"
	TypePolygon  = "polygon"
	TypeCircle   = "circle"
	TxidSnapshot = "txid_snapshot"
)
