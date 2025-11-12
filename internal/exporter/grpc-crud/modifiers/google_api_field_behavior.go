package modifiers

import (
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/db-exporter/internal/shared/proto/opts/googleapi"
)

type GoogleAPIFieldBehavior struct{}

func (m *GoogleAPIFieldBehavior) ModifyField(field *presentation.Field) {
	field.Message().Service().File().AddImport("google/api/field_behavior.proto")
	field.AddOption(m.option(field))
}

func (m *GoogleAPIFieldBehavior) option(field *presentation.Field) *proto.FieldOption {
	if field.Autofilled() {
		return googleapi.FieldOutputOnly()
	}

	if field.IsRequired() {
		return googleapi.FieldRequired()
	}

	return googleapi.FieldOptional()
}
