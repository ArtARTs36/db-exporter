package dbml

import "strings"

type File struct {
	Tables []*Table
	Refs   []*Ref
	Enums  []*Enum
}

func (f *File) Render() string {
	strs := make([]string, 0, len(f.Tables)+len(f.Refs)+len(f.Enums))

	for _, table := range f.Tables {
		strs = append(strs, table.Render())
	}

	for _, ref := range f.Refs {
		strs = append(strs, ref.Render())
	}

	for _, enum := range f.Enums {
		strs = append(strs, enum.Render())
	}

	return strings.Join(strs, "\n\n")
}
