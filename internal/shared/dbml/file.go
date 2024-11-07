package dbml

import "strings"

type File struct {
	Tables []*Table
	Refs   []*Ref
}

func (f *File) Render() string {
	strs := make([]string, 0, len(f.Tables))

	for _, table := range f.Tables {
		strs = append(strs, table.Render())
	}

	for _, ref := range f.Refs {
		strs = append(strs, ref.Render())
	}

	return strings.Join(strs, "\n\n")
}
