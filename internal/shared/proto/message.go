package proto

type Message struct {
	Name   string
	Fields []*Field
}

func (m *Message) write(buf stringsBuffer) {
	buf.WriteString("message " + m.Name + " {")

	if len(m.Fields) > 0 {
		buf.WriteString("\n")
	}

	for _, field := range m.Fields {
		field.write(buf)
		buf.WriteString("\n")
	}

	buf.WriteString("}")
}
