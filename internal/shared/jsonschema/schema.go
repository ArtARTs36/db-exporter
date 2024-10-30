package jsonschema

import (
	"encoding/json"
	"fmt"
)

type Schema struct {
	Title                string              `json:"title,omitempty"`
	Description          string              `json:"description,omitempty"`
	Schema               string              `json:"$schema,omitempty"`
	Type                 string              `json:"type,omitempty"`
	Properties           map[string]Property `json:"properties,omitempty"`
	AdditionalProperties bool                `json:"additional_properties,omitempty"`
	Definitions          map[string]Property `json:"definitions,omitempty"`
}

type Property struct {
	Title                string              `json:"title,omitempty"`
	Description          string              `json:"description,omitempty"`
	Type                 Type                `json:"type,omitempty"`
	Properties           map[string]Property `json:"properties,omitempty"`
	Required             []string            `json:"required,omitempty"`
	AdditionalProperties *bool               `json:"additional_properties,omitempty"`
	Ref                  string              `json:"$ref,omitempty"`
	Format               Format              `json:"format,omitempty"`
	Default              interface{}         `json:"default,omitempty"`
	Enum                 []string            `json:"enum,omitempty"`
}

func Draft04() *Schema {
	return &Schema{
		Schema:      "http://json-schema.org/draft-04/schema#",
		Type:        "object",
		Properties:  map[string]Property{},
		Definitions: map[string]Property{},
	}
}

func NewProperty(typ Type) Property {
	return Property{Type: typ, Properties: map[string]Property{}}
}

func (s *Schema) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Schema) MarshallPretty() ([]byte, error) {
	return json.MarshalIndent(s, "", "    ")
}

func (s *Schema) AddDefinition(name string, prop Property) string {
	s.Definitions[name] = prop

	return fmt.Sprintf("#/definitions/%s", name)
}

func (p *Property) DisableAdditionalProperties() {
	v := false

	p.AdditionalProperties = &v
}
