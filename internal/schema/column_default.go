package schema

type ColumnDefaultType int

const (
	ColumnDefaultTypeUnknown ColumnDefaultType = iota
	ColumnDefaultTypeFunc
	ColumnDefaultTypeValue
	ColumnDefaultTypeAutoincrement
)

type ColumnDefault struct {
	Type  ColumnDefaultType
	Value interface{}
}
