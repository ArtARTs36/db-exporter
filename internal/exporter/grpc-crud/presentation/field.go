package presentation

import "github.com/artarts36/db-exporter/internal/shared/proto"

type Field struct {
	proto *proto.Field
}

func (f *Field) AddOption(option *proto.FieldOption) *Field {
	f.proto.Options = append(f.proto.Options, option)

	return f
}

func (f *Field) AsRepeated() *Field {
	f.proto.Repeated = true

	return f
}

func (f *Field) SetType(typ string) {
	f.proto.Type = typ
}
