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
	TypeSqlNullInt64   = "sql.NullInt64"
	TypeSqlNullInt16   = "sql.NullInt16"
	TypeSqlNullFloat64 = "sql.NullFloat64"
	TypeSqlNullBool    = "sql.NullBool"
	TypeSqlNullString  = "sql.NullString"
	TypeSqlNullTime    = "sql.NullTime"

	TypeTimeTime = "time.Time"
)

func Ptr(t string) string {
	return fmt.Sprintf("*%s", t)
}
