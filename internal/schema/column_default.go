package schema

type ColumnDefaultType int

const (
	ColumnDefaultTypeUnknown ColumnDefaultType = iota
	ColumnDefaultTypeFunc                      // @todo need refactor to expression
	ColumnDefaultTypeValue
	ColumnDefaultTypeAutoincrement
)

type ColumnDefault struct {
	Type  ColumnDefaultType
	Value interface{}
}
