package golang

import "fmt"

const (
	TypeString    = "string"
	TypeByteSlice = "[]byte"
	TypeFloat32   = "float32"
	TypeFloat64   = "float64"
	TypeBool      = "bool"
	TypeInt16     = "int16"
	TypeInt64     = "int64"
)

const (
	TypeSQLNullInt64   = "sql.NullInt64"
	TypeSQLNullInt16   = "sql.NullInt16"
	TypeSQLNullFloat64 = "sql.NullFloat64"
	TypeSQLNullBool    = "sql.NullBool"
	TypeSQLNullString  = "sql.NullString"
	TypeSQLNullTime    = "sql.NullTime"

	TypeTimeTime = "time.Time"
)

func Ptr(t string) string {
	return fmt.Sprintf("*%s", t)
}
