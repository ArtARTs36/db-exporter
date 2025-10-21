package dbml

import (
	"fmt"
	"strings"
)

type Table struct {
	Name    string
	Columns []*Column

	Note string
}

type Column struct {
	Name string
	Type string

	Settings ColumnSettings
}

type ColumnSettings struct {
	PrimaryKey bool
	Increment  bool
	Note       string
	Unique     bool
	Default    ColumnDefault

	null bool
}

func (t *Table) Render(w *strings.Builder) {
	w.WriteString("Table " + t.Name + " {\n")

	for _, column := range t.Columns {
		w.WriteString("  " + column.Render())
	}

	if t.Note != "" {
		w.WriteString("Note: '" + t.Note + "'")
	}

	w.WriteString("}\n")
}

func (c *Column) Render() string {
	settingsStr := c.Settings.Render()
	if settingsStr != "" {
		settingsStr = fmt.Sprintf(" [%s]", settingsStr)
	}

	return c.Name + " " + c.Type + "" + settingsStr + "\n"
}

func (c *Column) AsNullable() {
	c.Settings.null = true
}

func (s *ColumnSettings) Render() string {
	strs := make([]string, 0)

	renderStringSetting := func(name, val string) string {
		return fmt.Sprintf("%s: '%s'", name, val)
	}

	if s.PrimaryKey {
		strs = append(strs, "primary key")
	}

	if !s.null {
		strs = append(strs, "not null")
	}

	if s.Increment {
		strs = append(strs, "increment")
	}

	if s.Note != "" {
		strs = append(strs, renderStringSetting("note", s.Note))
	}

	if s.Unique {
		strs = append(strs, "unique")
	}

	if s.Default.Value != "" {
		strs = append(strs, s.Default.Render())
	}

	return strings.Join(strs, ", ")
}
