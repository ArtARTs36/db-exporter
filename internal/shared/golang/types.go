package golang

import (
	"fmt"
)

type Type struct {
	Name        string
	PackageName string
	PackagePath string

	Null *Type
}

var (
	TypeString    = Type{Name: "string", Null: &TypeSQLNullString}
	TypeByteSlice = Type{Name: "[]byte"}
	TypeFloat32   = Type{Name: "float32"}
	TypeFloat64   = Type{Name: "float64", Null: &TypeSQLNullFloat64}
	TypeBool      = Type{Name: "bool", Null: &TypeSQLNullBool}
	TypeInt       = Type{Name: "int", Null: &TypeSQLNullInt64}
	TypeInt16     = Type{Name: "int16", Null: &TypeSQLNullInt16}
	TypeInt64     = Type{Name: "int64", Null: &TypeSQLNullInt64}

	TypeTimeTime     = Type{Name: "Time", PackageName: "time", PackagePath: "time", Null: &TypeSQLNullTime}
	TypeTimeDuration = Type{Name: "Duration", PackageName: "time", PackagePath: "time"}

	TypeSQLNullInt64   = Type{Name: "NullInt64", PackageName: "sql", PackagePath: "database/sql"}
	TypeSQLNullInt16   = Type{Name: "NullInt16", PackageName: "sql", PackagePath: "database/sql"}
	TypeSQLNullFloat64 = Type{Name: "NullFloat64", PackageName: "sql", PackagePath: "database/sql"}
	TypeSQLNullBool    = Type{Name: "NullBool", PackageName: "sql", PackagePath: "database/sql"}
	TypeSQLNullString  = Type{Name: "NullString", PackageName: "sql", PackagePath: "database/sql"}
	TypeSQLNullTime    = Type{Name: "NullTime", PackageName: "sql", PackagePath: "database/sql"}
)

func (t *Type) Call() string {
	if t.PackageName == "" {
		return t.Name
	}

	return fmt.Sprintf("%s.%s", t.PackageName, t.Name)
}

func Ptr(t string) string {
	return fmt.Sprintf("*%s", t)
}
