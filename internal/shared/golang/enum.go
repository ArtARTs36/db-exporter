package golang

import (
	"github.com/artarts36/db-exporter/internal/shared/ds"
)

type StringEnum struct {
	Name   *ds.String
	Values []*StringEnumValue
}

type StringEnumValue struct {
	Name  string
	Value string
}

func NewStringEnumOfValues(name *ds.String, values []string) *StringEnum {
	enum := &StringEnum{
		Name:   name.Pascal(),
		Values: make([]*StringEnumValue, 0, 1+len(values)),
	}

	enum.AddNamedValue("UNDEFINED", "")
	enum.AddValue(values...)

	return enum
}

func (e *StringEnum) AddNamedValue(name string, value string) {
	v := &StringEnumValue{
		Name:  e.Name.Append(ds.NewString(name).Pascal().Value).Value,
		Value: value,
	}

	e.Values = append(e.Values, v)
}

func (e *StringEnum) AddValue(values ...string) {
	for _, value := range values {
		e.AddNamedValue(value, value)
	}
}
