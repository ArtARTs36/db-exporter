package dbml

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/iox"
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

func (e *Enum) Render(w iox.Writer) {
	w.WriteString("Enum " + e.Name + " {\n")

	for _, value := range e.Values {
		w.WriteString("  ")
		value.Render(w)
		w.WriteString("\n")
	}

	w.WriteByte('}')
}

func (v *EnumValue) Render(w iox.Writer) {
	w.WriteByte('"')
	w.WriteString(v.Name)
	w.WriteByte('"')

	st := v.Settings.Render()
	if st != "" {
		st = " " + st
	}

	w.WriteString(st)
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
