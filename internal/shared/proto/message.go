package proto

import "github.com/artarts36/db-exporter/internal/shared/indentx"

type Message struct {
	Name   string
	Fields []*Field
}

func (m *Message) write(buf stringsBuffer, indent *indentx.Indent) {
	buf.WriteString("message " + m.Name + " {")

	if len(m.Fields) > 0 {
		buf.WriteString("\n")
	}

	for i, field := range m.Fields {
		if (i > 0 && len(m.Fields[i-1].Options) > 1) || field.TopComment != "" {
			buf.WriteString("\n")
		}

		field.write(buf, indent.Next())
	}

	buf.WriteString("}")
}
