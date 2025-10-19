package modifiers

import (
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/shared/proto/opts/bufvalidate"
)

type BufValidate struct {
}

func (p *BufValidate) ModifyField(field *presentation.Field) {
	if field.Message().Type() != presentation.MessageTypeRequest || field.Column() == nil {
		return
	}

	field.Message().Service().File().AddImport("buf/validate/validate.proto")

	if field.IsRequired() {
		field.AddOption(bufvalidate.Required())
	}

	colType := field.Column().Type

	switch {
	case colType.IsUUID:
		field.AddOption(bufvalidate.UUID())
	case colType.IsStringable:
		if field.Column().Name.Lower().Ends("email") {
			field.AddOption(bufvalidate.Email())
		}

		if field.Column().Type.Length != "" {
			field.AddOption(bufvalidate.MaxLen(field.Column().Type.Length))
		}
	}
}
