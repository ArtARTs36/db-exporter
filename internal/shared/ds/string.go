package ds

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/camelcase"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type String struct {
	Val string
}

type SplitWord struct {
	Word           string
	SeparatorAfter string
}

func NewString(val string) *String {
	return &String{
		Val: val,
	}
}

func (s *String) Scan(val any) error {
	switch v := val.(type) {
	case string:
		s.Val = v

		return nil
	case []byte:
		s.Val = string(v)

		return nil
	case nil:
		s.Val = ""

		return nil
	default:
		return fmt.Errorf("unexpected type %q", reflect.TypeOf(val).String())
	}
}

func (s String) Value() (driver.Value, error) {
	return s.Val, nil
}

func (s *String) String() string {
	return s.Val
}

func (s *String) Replace(old, new string) string {
	return strings.ReplaceAll(s.Val, old, new)
}

func (s *String) Pascal() *String {
	return NewString(strcase.ToCamel(s.Val))
}

func (s *String) Len() int {
	return len(s.Val)
}

func (s *String) IsNotEmpty() bool {
	return s.Len() != 0
}

func (s *String) Singular() *String {
	return NewString(inflection.Singular(s.Val))
}

func (s *String) Ends(suffix string) bool {
	return strings.HasSuffix(s.Val, suffix)
}

func (s *String) SplitCamel() []string {
	return camelcase.Split(s.Val)
}

func (s *String) SplitWords() []*SplitWord {
	if len(s.Val) == 0 {
		return []*SplitWord{}
	}

	srcBytes := []byte(s.Val)

	var words []*SplitWord
	currWordBytes := []byte{}

	prevCharIsLower := strings.ToLower(string(srcBytes[0])) == string(srcBytes[0])
	wordPos := 0

	for i, b := range srcBytes {
		currChar := string(b)
		currCharIsLower := strings.ToLower(currChar) == currChar

		if b == '_' || b == '-' || b == ' ' { //nolint:gocritic // not required
			words = append(words, &SplitWord{
				Word:           string(currWordBytes),
				SeparatorAfter: currChar,
			})
			wordPos = 0
			currWordBytes = []byte{}
		} else if prevCharIsLower != currCharIsLower && wordPos > 1 { // currWord: Aaa, currChar: B
			words = append(words, &SplitWord{
				Word: string(currWordBytes),
			})
			wordPos = 1
			currWordBytes = []byte{
				b,
			}
		} else {
			currWordBytes = append(currWordBytes, b)

			if i == len(srcBytes)-1 {
				words = append(words, &SplitWord{
					Word: string(currWordBytes),
				})
				break
			}

			wordPos++
		}

		prevCharIsLower = currCharIsLower
	}

	return words
}

func (s *String) FixAbbreviations(abbrSet map[string]bool) *String {
	split := s.SplitWords()
	words := make([]string, 0, len(split))

	for _, word := range split {
		w := strings.ToLower(word.Word)
		_, exists := abbrSet[w]
		if exists {
			words = append(words, strings.ToUpper(w), word.SeparatorAfter)
		} else {
			words = append(words, word.Word, word.SeparatorAfter)
		}
	}

	return NewString(strings.Join(words, ""))
}

func (s *String) Lower() *String {
	return NewString(strings.ToLower(s.Val))
}

func (s *String) Equal(str string) bool {
	return s.Val == str
}
