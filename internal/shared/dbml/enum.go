package dbml

import (
	"fmt"
	"strings"
)

type Enum struct {
	Name   string
	Values []EnumValue
}

type EnumValue struct {
	Name string

	Settings EnumValueSettings
}

type EnumValueSettings struct {
	Note string
}

func (e *Enum) Render() string {
	const strsMinLen = 2

	strs := make([]string, 0, strsMinLen+len(e.Values))

	strs = append(strs, fmt.Sprintf("Enum %s {", e.Name))

	for _, value := range e.Values {
		strs = append(strs, fmt.Sprintf("  %s", value.Render()))
	}

	strs = append(strs, "}")

	return strings.Join(strs, "\n")
}

func (v *EnumValue) Render() string {
	st := v.Settings.Render()
	if st != "" {
		st = " " + st
	}
	return fmt.Sprintf("%s%s", v.Name, st)
}

func (s *EnumValueSettings) Render() string {
	strs := make([]string, 0)

	if s.Note != "" {
		strs = append(strs, fmt.Sprintf("note: '%s'", s.Note))
	}

	if len(strs) == 0 {
		return ""
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}
