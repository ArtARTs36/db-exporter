package proto

import "github.com/artarts36/db-exporter/internal/shared/ds"

type Enum struct {
	Name        *ds.String
	Values      []string
	valuePrefix *ds.String
}

func NewEnum(name *ds.String, valuesCount int) *Enum {
	enum := &Enum{
		Name:        name.Pascal(),
		Values:      make([]string, 0, 1+valuesCount),
		valuePrefix: name.Upper().Append("_"),
	}

	enum.AddValue("UNDEFINED")

	return enum
}

func NewEnumWithValues(name *ds.String, values []string) *Enum {
	enum := NewEnum(name, len(values))
	enum.AddValue(values...)
	return enum
}

func (e *Enum) AddValue(value ...string) {
	for _, v := range value {
		e.Values = append(e.Values, e.valuePrefix.Append(ds.NewString(v).Upper().Value).Value)
	}
}
