package proto

import "github.com/artarts36/db-exporter/internal/shared/iox"

type Message struct {
	Name       string
	TopComment string
	Fields     []*Field
}

func (m *Message) write(buf iox.Writer) {
	if m.TopComment != "" {
		buf.WriteString("// ")
		buf.WriteString(m.TopComment)
		buf.WriteNewLine()
	}

	buf.WriteString("message " + m.Name + " {")

	if len(m.Fields) > 0 {
		buf.WriteNewLine()
	}

	for i, field := range m.Fields {
		if (i > 0 && len(m.Fields[i-1].Options) > 1) || field.TopComment != "" {
			buf.WriteString("\n")
		}

		field.write(buf.IncIndent())
	}

	buf.WriteString("}")
}
