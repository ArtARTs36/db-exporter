package schema

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
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

func (s *String) Replace(old, new string) string {
	return strings.ReplaceAll(s.Value, old, new)
}

func (s *String) Pascal() *String {
	return &String{
		Value: strcase.ToCamel(s.Value),
	}
}

func (s *String) Len() int {
	return len(s.Value)
}

func (s *String) IsNotEmpty() bool {
	return s.Len() != 0
}

func (s *String) Singular() *String {
	return &String{
		Value: inflection.Singular(s.Value),
	}
}

func (s *String) Ends(suffix string) bool {
	return strings.HasSuffix(s.Value, suffix)
}

func (s *String) SplitCamel() []string {
	return camelcase.Split(s.Value)
}

func (s *String) FixAbbreviations(abbrSet map[string]bool) *String {
	words := s.SplitCamel()
	for i, word := range words {
		w := strings.ToLower(word)
		_, exists := abbrSet[w]
		if !exists {
			continue
		}

		words[i] = strings.ToUpper(w)
	}

	return &String{
		Value: strings.Join(words, ""),
	}
}

func (s *String) Lower() *String {
	return &String{
		Value: strings.ToLower(s.Value),
	}
}
