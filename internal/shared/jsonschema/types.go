package jsonschema

type Type string

const (
	TypeString  Type = "string"
	TypeBoolean Type = "boolean"
	TypeObject  Type = "object"
	TypeArray   Type = "array"
	TypeInteger Type = "integer"
	TypeNumber  Type = "number"
)
