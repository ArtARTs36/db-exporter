package dbml

import "strings"

type File struct {
	Tables []*Table
	Refs   []*Ref
	Enums  []*Enum
}

func (f *File) Render() string {
	w := &strings.Builder{}

	for i, table := range f.Tables {
		table.Render(w)

		if i < len(f.Tables)-1 {
			w.WriteString("\n")
		}
	}

	if len(f.Refs) > 0 {
		w.WriteString("\n")
	}

	for _, ref := range f.Refs {
		ref.Render(w)
	}

	if w.Len() > 0 && len(f.Enums) > 0 {
		w.WriteString("\n")
	}

	for _, enum := range f.Enums {
		enum.Render(w)
	}

	return w.String()
}
