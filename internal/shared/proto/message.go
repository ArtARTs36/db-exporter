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
	buf.WriteNewLine()

	m.writeFields(buf)

	buf.WriteString("}")
}

func (m *Message) writeFields(buf iox.Writer) {
	for i, field := range m.Fields {
		if (i > 0 && m.Fields[i-1].hasOptions()) || field.TopComment != "" {
			buf.WriteNewLine()
		}

		field.write(buf.IncIndent())
	}
}
