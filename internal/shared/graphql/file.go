package graphql

import "github.com/artarts36/db-exporter/internal/shared/iox"

type File struct {
	Types []*Object
}

func (f *File) Build(w iox.Writer) {
	for i, t := range f.Types {
		t.Build(w)

		if i < len(f.Types)-1 {
			w.WriteNewLine()
			w.WriteNewLine()
		}
	}
}
