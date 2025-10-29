package proto

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/iox"
	"strconv"
)

type Field struct {
	Repeated   bool
	TopComment string
	Type       string
	Name       string
	ID         int
	Options    []*FieldOption
}

type FieldOption struct {
	Name  string
	Value interface{}
}

type ConstValue string

func (f *Field) write(buf iox.Writer) {
	if f.TopComment != "" {
		buf.WriteString("// " + f.TopComment)
		buf.WriteNewLine()
	}

	buf.WriteString("")

	if f.Repeated {
		buf.WriteInline("repeated ")
	}

	buf.WriteInline(f.Type + " " + f.Name + " = " + strconv.Itoa(f.ID))

	if len(f.Options) > 0 {
		buf.WriteInline(" [")

		if len(f.Options) == 1 {
			f.Options[0].write(buf.WithoutIndent())
		} else {
			buf.WriteNewLine()
			for i, opt := range f.Options {
				opt.write(buf.IncIndent())

				if i < len(f.Options)-1 {
					buf.WriteInline(",")
				}
				buf.WriteNewLine()
			}
			buf.WriteString("")
		}

		buf.WriteInline("]")
	}

	buf.WriteInline(";\n")
}

func (f *FieldOption) write(buf iox.Writer) {
	buf.WriteString(f.Name + " = " + f.resolveValue())
}

func (f *FieldOption) resolveValue() string {
	switch val := f.Value.(type) {
	case ConstValue:
		return string(val)
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
