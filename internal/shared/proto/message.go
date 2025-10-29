package proto

import "github.com/artarts36/db-exporter/internal/shared/iox"

type Message struct {
	Name   string
	Fields []*Field
}

func (m *Message) write(buf iox.Writer) {
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
