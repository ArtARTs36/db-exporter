package php

type Type int

const (
	TypeUndefined Type = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeString
	TypeObject
)

func (t Type) String() string {
	switch t {
	case TypeUndefined:
		return ""
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeBool:
		return "bool"
	case TypeString:
		return "string"
	case TypeObject:
		return "object"
	}
	return ""
}
