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

	fieldLine := buf.Line()

	if f.Repeated {
		fieldLine.WriteString("repeated ")
	}

	fieldLine.WriteString(f.Type + " " + f.Name + " = " + strconv.Itoa(f.ID))

	if len(f.Options) > 0 {
		fieldLine.WriteString(" [")
		buf.WriteNewLine()

		optsBuf := buf.IncIndent()

		// write options
		for i, opt := range f.Options {
			optLine := optsBuf.Line()
			opt.write(optLine)

			if i < len(f.Options)-1 {
				optLine.WriteString(",")
			}
			buf.WriteNewLine()
		}
		// close options
		buf.WriteString("];")
		buf.WriteNewLine()
	} else {
		fieldLine.WriteString(";")
		buf.WriteNewLine()
	}
}

func (f *Field) hasOptions() bool {
	return len(f.Options) > 0
}

func (f *FieldOption) write(buf iox.StringWriter) {
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
