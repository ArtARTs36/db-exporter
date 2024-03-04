package schema

import (
	"fmt"
	"reflect"
	"strings"
)

type String struct {
	Value string
}

func (s *String) Scan(val any) error {
	switch v := val.(type) {
	case string:
		s.Value = v

		return nil
	case []byte:
		s.Value = string(v)

		return nil
	case nil:
		s.Value = ""

		return nil
	default:
		return fmt.Errorf("unexpected type %q", reflect.TypeOf(val).String())
	}
}

func (s *String) String() string {
	return s.Value
}

func (s *String) Replace(old string, new string) string {
	return strings.ReplaceAll(s.Value, old, new)
}
