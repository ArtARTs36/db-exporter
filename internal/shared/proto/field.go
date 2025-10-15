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

func (f *Field) write(buf stringsBuffer) {
	buf.WriteString("    ")

	if f.Repeated {
		buf.WriteString("repeated ")
	}

	buf.WriteString(f.Type + " " + f.Name + " = " + strconv.Itoa(f.ID) + ";")
}
