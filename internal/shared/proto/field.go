package proto

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/indentx"
	"strconv"
)

type Field struct {
	Repeated bool
	Type     string
	Name     string
	ID       int
	Options  []*FieldOption
}

type FieldOption struct {
	Name  string
	Value interface{}
}

func (f *Field) write(buf stringsBuffer, indent *indentx.Indent) {
	buf.WriteString(indent.Curr())

	if f.Repeated {
		buf.WriteString("repeated ")
	}

	buf.WriteString(f.Type + " " + f.Name + " = " + strconv.Itoa(f.ID))

	if len(f.Options) > 0 {
		buf.WriteString(" [")

		if len(f.Options) == 1 {
			f.Options[0].write(buf, indentx.Zero())
		} else {
			buf.WriteString("\n")
			for i, opt := range f.Options {
				opt.write(buf, indent.Next())

				if i < len(f.Options)-1 {
					buf.WriteString(",")
				}
				buf.WriteString("\n")
			}
			buf.WriteString(indent.Curr())
		}

		buf.WriteString("]")
	}

	buf.WriteString(";")
}

func (f *FieldOption) write(buf stringsBuffer, indent *indentx.Indent) {
	buf.WriteString(indent.Curr())
	buf.WriteString(f.Name + " = " + f.resolveValue())
}

func (f *FieldOption) resolveValue() string {
	switch val := f.Value.(type) {
	case string:
		return strconv.Quote(val)
	case int:
		return strconv.Itoa(val)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprint(val)
	}
}

func (f *Field) Clone() *Field {
	return &Field{
		Repeated: f.Repeated,
		Type:     f.Type,
		Name:     f.Name,
		ID:       f.ID,
	}
}
