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
		if field.TopComment != "" {
			buf.WriteString("\n")
		}

		field.write(buf, indent.Next())
		buf.WriteString("\n")

		if len(field.Options) > 1 && i < len(m.Fields)-1 {
			buf.WriteString("\n")
		}
	}

	buf.WriteString("}")
}
