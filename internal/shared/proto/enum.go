package proto

import (
	"github.com/artarts36/gds"
	"strconv"
)

type Enum struct {
	Name        *gds.String
	Values      []string
	valuePrefix *gds.String
}

func NewEnum(name *gds.String, valuesCount int) *Enum {
	enum := &Enum{
		Name:        name.Pascal(),
		Values:      make([]string, 0, 1+valuesCount),
		valuePrefix: name.Upper().Append("_"),
	}

	enum.AddValue("UNDEFINED")

	return enum
}

func NewEnumWithValues(name *gds.String, values []string) *Enum {
	enum := NewEnum(name, len(values))
	enum.AddValue(values...)
	return enum
}

func (e *Enum) AddValue(value ...string) {
	for _, v := range value {
		e.Values = append(e.Values, e.valuePrefix.Append(gds.NewString(v).Upper().Value).Value)
	}
}

func (e *Enum) write(buf stringsBuffer) {
	buf.WriteString(e.Name.Prepend("enum ").Append(" {").Value)

	buf.WriteString("\n")

	for i, v := range e.Values {
		buf.WriteString("    " + v + " = " + strconv.Itoa(i) + ";")
		buf.WriteString("\n")
	}

	buf.WriteString("}")
}
