package schema

type ColumnType int

const (
	ColumnTypeInteger   ColumnType = iota
	ColumnTypeInteger64 ColumnType = iota
	ColumnTypeInteger16 ColumnType = iota
	ColumnTypeString    ColumnType = iota
	ColumnTypeTimestamp ColumnType = iota
	ColumnTypeBoolean   ColumnType = iota
	ColumnTypeFloat32   ColumnType = iota
	ColumnTypeFloat64   ColumnType = iota
	ColumnTypeBytes     ColumnType = iota
)

func (t ColumnType) String() string {
	switch t {
	case ColumnTypeInteger:
		return "integer"
	case ColumnTypeInteger64:
		return "integer64"
	case ColumnTypeInteger16:
		return "integer16"
	case ColumnTypeString:
		return "string"
	case ColumnTypeTimestamp:
		return "timestamp"
	case ColumnTypeBoolean:
		return "boolean"
	case ColumnTypeFloat32:
		return "float32"
	case ColumnTypeFloat64:
		return "float64"
	case ColumnTypeBytes:
		return "[]byte"
	default:
		return "string"
	}
}
