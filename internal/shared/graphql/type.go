package graphql

type Type interface {
	Name() string
	graphqlType()
}

type ScalarType int

type NamedType struct {
	name string
}

const (
	TypeInt ScalarType = iota
	TypeFloat
	TypeString
	TypeBoolean
	TypeID
)

func TypeOfName(name string) *NamedType {
	return &NamedType{name: name}
}

func (t ScalarType) Name() string {
	switch t {
	case TypeInt:
		return "Int"
	case TypeFloat:
		return "Float"
	case TypeString:
		return "String"
	case TypeBoolean:
		return "Boolean"
	case TypeID:
		return "ID"
	}
	return ""
}

func (t ScalarType) graphqlType() {}

func (t NamedType) Name() string {
	return t.name
}

func (t NamedType) graphqlType() {}
