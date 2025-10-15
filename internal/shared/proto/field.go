package proto

import (
	"strconv"
)

type Field struct {
	Repeated bool
	Type     string
	Name     string
	ID       int
}

func (f *Field) write(buf stringsBuffer, indent *Indent) {
	buf.WriteString(indent.curr)

	if f.Repeated {
		buf.WriteString("repeated ")
	}

	buf.WriteString(f.Type + " " + f.Name + " = " + strconv.Itoa(f.ID) + ";")
}

func (f *Field) Clone() *Field {
	return &Field{
		Repeated: f.Repeated,
		Type:     f.Type,
		Name:     f.Name,
		ID:       f.ID,
	}
}
