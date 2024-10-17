package schema

type ColumnDefaultValueType int

const (
	ColumnDefaultValueTypeUnknown ColumnDefaultValueType = iota
	ColumnDefaultValueTypeFunc
	ColumnDefaultValueTypeInteger
	ColumnDefaultValueTypeString
)

type ColumnDefault struct {
	Type  ColumnDefaultValueType
	Value string
}
