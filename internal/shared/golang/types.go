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

func Ptr(t string) string {
	return fmt.Sprintf("*%s", t)
}
