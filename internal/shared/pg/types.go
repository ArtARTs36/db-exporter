package pg

const (
	TypeText             = "text"
	TypeUUID             = "uuid"
	TypeCharacter        = "character"
	TypeCharacterVarying = "character varying"

	TypeTimestampWithoutTZ = "timestamp without time zone"
	TypeTimestampWithTZ    = "timestamp with time zone"

	TypeBoolean = "boolean"
	TypeBit     = "bit"

	TypeBytea = "bytea"

	TypeInteger     = "integer"
	TypeBigint      = "bigint"
	TypeInt4        = "int4"
	TypeInt8        = "int8"
	TypeSmallInt    = "smallint"
	TypeSmallSerial = "smallserial"
	TypeSerial      = "serial"

	TypeMoney   = "money"
	TypeNumeric = "numeric"

	TypeReal            = "real"
	TypeDoublePrecision = "double precision"
	TypeFloat8          = "float8"
	TypeDecimal         = "decimal"
)
