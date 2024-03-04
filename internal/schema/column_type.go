package schema

type ColumnType int

const (
	ColumnTypeInteger   ColumnType = iota
	ColumnTypeString    ColumnType = iota
	ColumnTypeTimestamp ColumnType = iota
	ColumnTypeBoolean   ColumnType = iota
)

func (t ColumnType) String() string {
	switch t {
	case ColumnTypeInteger:
		return "integer"
	case ColumnTypeString:
		return "string"
	case ColumnTypeTimestamp:
		return "timestamp"
	case ColumnTypeBoolean:
		return "boolean"
	default:
		return "string"
	}
}
