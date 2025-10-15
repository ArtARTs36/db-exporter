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

	for _, field := range m.Fields {
		field.write(buf, indent.Next())
		buf.WriteString("\n")
	}

	buf.WriteString("}")
}
