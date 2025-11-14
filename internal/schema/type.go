package schema

import "fmt"

type DataType struct {
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
	IsJSON       bool
	IsInterval   bool
}

func (t *DataType) Clone() DataType {
	return DataType{
		Name:         t.Name,
		Length:       t.Length,
		IsNumeric:    t.IsNumeric,
		IsInteger:    t.IsInteger,
		IsFloat:      t.IsFloat,
		IsUUID:       t.IsUUID,
		IsStringable: t.IsStringable,
		IsDatetime:   t.IsDatetime,
		IsDate:       t.IsDate,
		IsBoolean:    t.IsBoolean,
		IsBinary:     t.IsBinary,
		IsJSON:       t.IsJSON,
		IsInterval:   t.IsInterval,
	}
}

func (t *DataType) WithLength(length string) DataType {
	newType := t.Clone()
	newType.Length = length

	return newType
}

func (t *DataType) MarkAsUUID() DataType {
	newType := t.Clone()
	newType.IsUUID = true

	return newType
}

func (t *DataType) String() string {
	if t.Length == "" {
		return t.Name
	}

	return fmt.Sprintf("%s(%s)", t.Name, t.Length)
}
