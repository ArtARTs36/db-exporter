package modifiers

import (
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/shared/proto/opts/bufvalidate"
	"strconv"
)

type BufValidate struct {
}

func (p *BufValidate) ModifyField(field *presentation.Field) {
	if field.Message().Type() != presentation.MessageTypeRequest {
		return
	}

	added := p.addOption(field)
	if added {
		field.Message().Service().File().AddImport("buf/validate/validate.proto")
	}
}

func (p *BufValidate) addOption(field *presentation.Field) bool {
	if field.IsRequired() {
		field.AddOption(bufvalidate.Required())
	}

	if field.Column() == nil {
		return false
	}

	colType := field.Column().DataType

	switch {
	case colType.IsUUID:
		field.AddOption(bufvalidate.UUID())
		return true
	case colType.IsStringable:
		added := false

		if field.Column().Name.Lower().Ends("email") {
			field.AddOption(bufvalidate.Email())
			added = true
		}

		if field.Column().CharacterLength > 0 {
			field.AddOption(bufvalidate.MaxLen(strconv.Itoa(int(field.Column().CharacterLength))))
			added = true
		} else if field.Column().DataType.Length != "" {
			field.AddOption(bufvalidate.MaxLen(field.Column().DataType.Length))
			added = true
		}

		return added
	}

	return true
}
