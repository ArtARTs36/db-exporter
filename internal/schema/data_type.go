package schema

type DataType int

const (
	DataTypeInteger   DataType = iota
	DataTypeInteger64 DataType = iota
	DataTypeInteger16 DataType = iota
	DataTypeString    DataType = iota
	DataTypeTimestamp DataType = iota
	DataTypeBoolean   DataType = iota
	DataTypeFloat32   DataType = iota
	DataTypeFloat64   DataType = iota
	DataTypeBytes     DataType = iota
)

func (t DataType) String() string {
	switch t {
	case DataTypeInteger:
		return "integer"
	case DataTypeInteger64:
		return "integer64"
	case DataTypeInteger16:
		return "integer16"
	case DataTypeString:
		return "string"
	case DataTypeTimestamp:
		return "timestamp"
	case DataTypeBoolean:
		return "boolean"
	case DataTypeFloat32:
		return "float32"
	case DataTypeFloat64:
		return "float64"
	case DataTypeBytes:
		return "[]byte"
	default:
		return "string"
	}
}

func (t DataType) IsInteger() bool {
	return t == DataTypeInteger || t == DataTypeInteger64 || t == DataTypeInteger16
}
