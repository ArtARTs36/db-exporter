package dbml

import (
	"fmt"
	"regexp"
)

type ColumnDefaultType int

const (
	ColumnDefaultTypeNumber ColumnDefaultType = iota
	ColumnDefaultTypeString
	ColumnDefaultTypeExpression
	ColumnDefaultTypeBoolean
)

var (
	ColumnDefaultValueStringRegexp     = regexp.MustCompile(`^'(.*)'$`)
	ColumnDefaultValueExpressionRegexp = regexp.MustCompile("^`(.*)`$")
)

type ColumnDefault struct {
	Type  ColumnDefaultType
	Value string
}

func (d *ColumnDefault) Render() string {
	return fmt.Sprintf("default: %s", d.renderValue())
}

func (d *ColumnDefault) renderValue() string {
	switch d.Type {
	case ColumnDefaultTypeNumber, ColumnDefaultTypeBoolean:
		return d.Value
	case ColumnDefaultTypeString:
		return fmt.Sprintf("'%s'", d.Value)
	case ColumnDefaultTypeExpression:
		return fmt.Sprintf("`%s`", d.Value)
	}
	return d.Value
}
