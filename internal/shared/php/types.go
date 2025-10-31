package php

type Type string

const (
	TypeUndefined Type = ""
	TypeInt       Type = "int"
	TypeFloat     Type = "float"
	TypeBool      Type = "bool"
	TypeString    Type = "string"
	TypeObject    Type = "object"
)
