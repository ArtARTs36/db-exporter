package dbml

import (
	"github.com/artarts36/db-exporter/internal/shared/iox"
)

type Ref struct {
	From string
	Type string
	To   string
}

func (r *Ref) Render(w iox.Writer) {
	w.WriteString("Ref: " + r.From + " " + r.Type + " " + r.To + "\n")
}
