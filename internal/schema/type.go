package schema

import "fmt"

type Type struct {
	Name string

	Length string

	IsNumeric    bool
	IsInteger    bool
	IsFloat      bool
	IsUUID       bool
	IsStringable bool
	IsDatetime   bool
	IsDate       bool
	IsBoolean    bool
	IsBinary     bool
}

func (t *Type) Clone() Type {
	return Type{
		Name:         t.Name,
		Length:       t.Length,
		IsNumeric:    t.IsNumeric,
		IsInteger:    t.IsInteger,
		IsFloat:      t.IsFloat,
		IsUUID:       t.IsUUID,
		IsStringable: t.IsStringable,
		IsDatetime:   t.IsDatetime,
		IsDate:       t.IsDate,
	}
}

func (t *Type) WithLength(length string) Type {
	newType := t.Clone()
	newType.Length = length

	return newType
}

func (t *Type) MarkAsUUID() Type {
	newType := t.Clone()
	newType.IsUUID = true

	return newType
}

func (t *Type) String() string {
	if t.Length == "" {
		return t.Name
	}

	return fmt.Sprintf("%s(%s)", t.Name, t.Length)
}
