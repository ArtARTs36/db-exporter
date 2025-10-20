package graphql

import (
	"strings"
)

type Object struct {
	kind   string
	name   string
	fields []*Field
}

func NewType(name string) *Object {
	return &Object{
		kind:   "type",
		name:   name,
		fields: []*Field{},
	}
}

func NewInput(name string) *Object {
	return &Object{
		kind:   "input",
		name:   name,
		fields: []*Field{},
	}
}

func (t *Object) AddField(name string) *Field {
	prop := &Field{
		name: name,
	}

	t.AttachField(prop)

	return prop
}

func (t *Object) AttachField(prop *Field) {
	t.fields = append(t.fields, prop)
}

func (t *Object) Build() string {
	stringsBuilder := &strings.Builder{}

	stringsBuilder.WriteString(t.kind + " " + t.name + " {\n")

	for _, property := range t.fields {
		property.Write(stringsBuilder)
	}

	stringsBuilder.WriteByte('}')

	return stringsBuilder.String()
}

func (t *Object) Name() string {
	return t.name
}

func (t *Object) graphqlType() {}
