package graphql

import (
	"fmt"
	"strings"
)

type Field struct {
	name     string
	required bool
	comment  string

	typ  Type
	list bool
}

func (p *Field) Of(typ Type) *Field {
	p.typ = typ
	return p
}

func (p *Field) ListOf(typ Type) *Field {
	p.Of(typ)
	p.list = true
	return p
}

func (p *Field) Require() *Field {
	p.required = true
	return p
}

func (p *Field) Comment(comment string) *Field {
	p.comment = comment
	return p
}

func (p *Field) Write(w *strings.Builder) {
	kind := p.typ.Name()
	if p.list {
		kind = fmt.Sprintf("[%s!]", kind)
	}

	required := ""
	if p.required {
		required = "!"
	}

	if p.comment != "" {
		w.WriteString("  # " + p.comment + "\n")
	}

	w.WriteString("  " + p.name + ": " + kind + required + "\n")
}
