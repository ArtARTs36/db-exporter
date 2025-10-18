package presentation

import "github.com/artarts36/db-exporter/internal/shared/proto"

type Field struct {
	required bool

	proto   *proto.Field
	message *Message
}

func (f *Field) AddOption(option *proto.FieldOption) *Field {
	f.proto.Options = append(f.proto.Options, option)

	return f
}

func (f *Field) AsRepeated() *Field {
	f.proto.Repeated = true

	return f
}

func (f *Field) SetType(typ string) *Field {
	f.proto.Type = typ

	return f
}

func (f *Field) AsRequired() *Field {
	f.required = true
	return f
}

func (f *Field) IsRequired() bool {
	return f.required
}

func (f *Field) Message() *Message {
	return f.message
}
