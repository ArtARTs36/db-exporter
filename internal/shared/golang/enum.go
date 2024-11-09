package golang

import "github.com/artarts36/gds"

type StringEnum struct {
	Name   *gds.String
	Values []*StringEnumValue
}

type StringEnumValue struct {
	Name  string
	Value string
}

func NewStringEnumOfValues(name *gds.String, values []string) *StringEnum {
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
		Name:  e.Name.Append(gds.NewString(name).Pascal().Value).Value,
		Value: value,
	}

	e.Values = append(e.Values, v)
}

func (e *StringEnum) AddValue(values ...string) {
	for _, value := range values {
		e.AddNamedValue(value, value)
	}
}
