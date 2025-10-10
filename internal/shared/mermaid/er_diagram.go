package mermaid

import "strings"

type KeyType string

const (
	KeyTypeUnspecified KeyType = ""
	KeyTypePK          KeyType = "PK"
	KeyTypeFK          KeyType = "FK"
	KeyTypeUK          KeyType = "UK"
)

type ErDiagram struct {
	entities  []*Entity
	relations []*Relation
}

func NewErDiagram() *ErDiagram {
	return &ErDiagram{
		entities:  make([]*Entity, 0),
		relations: make([]*Relation, 0),
	}
}

type Entity struct {
	Name   string
	Fields []*EntityField
}

type EntityField struct {
	Name     string
	DataType string
	KeyType  KeyType
}

type Relation struct {
	Owner   string
	Related string
	Action  string
}

func (d *ErDiagram) AddEntity(entity *Entity) {
	d.entities = append(d.entities, entity)
}

func (d *ErDiagram) AddRelation(relation *Relation) {
	d.relations = append(d.relations, relation)
}

func (d *ErDiagram) Build() string {
	result := &strings.Builder{}

	_, _ = result.WriteString("erDiagram")
	_, _ = result.WriteString("\n")

	for _, relation := range d.relations {
		relation.build(result, "  ")
		write(result, "\n")
	}

	for _, entity := range d.entities {
		entity.build(result, "  ")
		write(result, "\n")
	}

	return result.String()
}

func (e *Entity) build(builder *strings.Builder, indent string) {
	write(builder, indent)
	write(builder, e.Name)
	write(builder, " {")

	if len(e.Fields) > 0 {
		write(builder, "\n")
	}

	fieldIndent := indent + "  "

	for _, field := range e.Fields {
		field.build(builder, fieldIndent)
		write(builder, "\n")
	}

	write(builder, indent)
	write(builder, "}")
}

func (f *EntityField) build(builder *strings.Builder, indent string) {
	write(builder, indent)
	write(builder, f.DataType)
	write(builder, " ")
	write(builder, f.Name)

	if f.KeyType != KeyTypeUnspecified {
		write(builder, " ")
		write(builder, string(f.KeyType))
	}
}

func (r *Relation) build(builder *strings.Builder, indent string) {
	builder.WriteString(indent)
	write(builder, r.Related)
	write(builder, " ")
	write(builder, "||--o{ ")
	write(builder, r.Owner)
	write(builder, " : ")
	write(builder, r.Action)
}

func write(builder *strings.Builder, arg string) {
	_, _ = builder.WriteString(arg)
}
