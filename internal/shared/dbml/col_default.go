package dbml

import "fmt"

type ColumnDefaultType int

const (
	ColumnDefaultTypeNumber ColumnDefaultType = iota
	ColumnDefaultTypeString
	ColumnDefaultTypeExpression
	ColumnDefaultTypeBoolean
)

type ColumnDefault struct {
	Type  ColumnDefaultType
	Value string
}

func (d *ColumnDefault) Render() string {
	switch d.Type {
	case ColumnDefaultTypeNumber:
		return d.Value
	case ColumnDefaultTypeString:
		return fmt.Sprintf("'%s'", d.Value)
	case ColumnDefaultTypeExpression:
		return fmt.Sprintf("`%s`", d.Value)
	case ColumnDefaultTypeBoolean:
		return d.Value
	}
	return d.Value
}
