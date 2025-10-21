package dbml

import (
	"strings"
)

type Ref struct {
	From string
	Type string
	To   string
}

func (r *Ref) Render(w *strings.Builder) {
	w.WriteString("Ref: " + r.From + " " + r.Type + " " + r.To + "\n")
}
