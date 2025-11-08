package proto

import (
	"github.com/artarts36/db-exporter/internal/shared/iox"
	"github.com/artarts36/gds"
	"strconv"
)

type Enum struct {
	Name        gds.String
	Values      []string
	valuePrefix *gds.String
}

func NewEnum(name gds.String, valuesCount int) *Enum {
	enum := &Enum{
		Name:        *name.Pascal(),
		Values:      make([]string, 0, 1+valuesCount),
		valuePrefix: name.Upper().Append("_"),
	}

	enum.AddValue("UNSPECIFIED")

	return enum
}

func NewEnumWithValues(name gds.String, values []string) *Enum {
	enum := NewEnum(name, len(values))
	enum.AddValue(values...)
	return enum
}

func (e *Enum) AddValue(value ...string) {
	for _, v := range value {
		e.Values = append(e.Values, e.valuePrefix.Append(gds.NewString(v).Upper().Value).Value)
	}
}

func (e *Enum) write(buf iox.Writer) {
	buf.WriteString(e.Name.Prepend("enum ").Append(" {").Value)
	buf.WriteNewLine()

	valuesBuf := buf.IncIndent()

	for i, v := range e.Values {
		valuesBuf.WriteString(v + " = " + strconv.Itoa(i) + ";")
		valuesBuf.WriteNewLine()
	}

	buf.WriteString("}")
}
