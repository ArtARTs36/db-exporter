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

type SplitWord struct {
	Word           string
	SeparatorAfter string
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

func (s *String) Camel() *String {
	return NewString(strcase.ToLowerCamel(s.Value))
}

func (s *String) Len() int {
	return len(s.Value)
}

func (s *String) IsEmpty() bool {
	return s.Len() == 0
}

func (s *String) IsNotEmpty() bool {
	return s.Len() != 0
}

func (s *String) Singular() *String {
	return NewString(inflection.Singular(s.Value))
}

func (s *String) Plural() *String {
	return NewString(inflection.Plural(s.Value))
}

func (s *String) Starts(prefix string) bool {
	return strings.HasPrefix(s.Value, prefix)
}

func (s *String) Ends(suffix string) bool {
	return strings.HasSuffix(s.Value, suffix)
}

func (s *String) SplitCamel() []string {
	return camelcase.Split(s.Value)
}

func (s *String) SplitWords() []*SplitWord {
	if len(s.Value) == 0 {
		return []*SplitWord{}
	}

	srcBytes := []byte(s.Value)

	var words []*SplitWord
	currWordBytes := []byte{}

	prevCharIsLower := strings.ToLower(string(srcBytes[0])) == string(srcBytes[0])
	wordPos := 0

	for i, b := range srcBytes {
		currChar := string(b)
		currCharIsLower := strings.ToLower(currChar) == currChar

		if b == '_' || b == '-' || b == ' ' || b == '.' || b == '/' { //nolint:gocritic // not required
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

func (s *String) PluralFixAbbreviations(abbrSet map[string]string) *String {
	split := s.SplitWords()
	words := make([]string, 0, len(split))

	for i, word := range split {
		w := strings.ToLower(word.Word)
		newWord, exists := abbrSet[w]
		if exists {
			if i < len(split)-1 {
				newWord = strings.ToUpper(w)
			}
		} else {
			if i == len(split)-1 {
				newWord = inflection.Plural(word.Word)
			} else {
				newWord = word.Word
			}
		}

		words = append(words, newWord, word.SeparatorAfter)
	}

	return NewString(strings.Join(words, ""))
}

func (s *String) Lower() *String {
	return NewString(strings.ToLower(s.Value))
}

func (s *String) Upper() *String {
	return NewString(strings.ToUpper(s.Value))
}

func (s *String) Equal(strs ...string) bool {
	for _, str := range strs {
		if s.Value == str {
			return true
		}
	}

	return true
}

func (s *String) FirstLine() *String {
	lines := strings.Split(s.Value, "\n")
	if len(lines) == 0 {
		return &String{Value: ""}
	}

	return &String{Value: lines[0]}
}

func (s *String) TrimPrefix(cutset string) *String {
	return &String{Value: strings.TrimPrefix(s.Value, cutset)}
}

func (s *String) TrimSpaces() *String {
	return &String{
		Value: strings.Trim(s.Value, " "),
	}
}

func (s *String) Prepend(prefix string) *String {
	return NewString(fmt.Sprintf("%s%s", prefix, s.Value))
}

func (s *String) Append(suffix string) *String {
	return NewString(fmt.Sprintf("%s%s", s.Value, suffix))
}
