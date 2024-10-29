package graphql

import (
	"fmt"
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
	const minObjectLinesCount = 2

	strs := make([]string, 0, minObjectLinesCount+len(t.fields))

	strs = append(strs, fmt.Sprintf("%s %s {", t.kind, t.name))

	for _, property := range t.fields {
		strs = append(strs, property.Build())
	}

	strs = append(strs, "}")

	return strings.Join(strs, "\n")
}

func (t *Object) Name() string {
	return t.name
}

func (t *Object) graphqlType() {}
