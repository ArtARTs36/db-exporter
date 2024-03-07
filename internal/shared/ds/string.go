package ds

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

func NewString(val string) *String {
	return &String{
		Value: val,
	}
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
	return NewString(strcase.ToCamel(s.Value))
}

func (s *String) Len() int {
	return len(s.Value)
}

func (s *String) IsNotEmpty() bool {
	return s.Len() != 0
}

func (s *String) Singular() *String {
	return NewString(inflection.Singular(s.Value))
}

func (s *String) Ends(suffix string) bool {
	return strings.HasSuffix(s.Value, suffix)
}

func (s *String) SplitCamel() []string {
	return camelcase.Split(s.Value)
}

func (s *String) SplitWords() []string {
	srcBytes := []byte(s.Value)

	words := []string{}
	currWordBytes := []byte{}

	for i, b := range srcBytes {
		if b == '_' || b == '-' || b == ' ' {
			words = append(words, string(currWordBytes))
		} else {
			currWordBytes = append(currWordBytes, b)

			if i == len(srcBytes)-1 {
				words = append(words, string(currWordBytes))
			}
		}
	}

	return words
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

	return NewString(strings.Join(words, ""))
}

func (s *String) Lower() *String {
	return NewString(strings.ToLower(s.Value))
}
