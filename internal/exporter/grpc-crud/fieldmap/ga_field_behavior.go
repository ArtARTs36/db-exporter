package fieldmap

import (
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/db-exporter/internal/shared/proto/opts/googleapi"
)

type GoogleAPIFieldBehaviorModifier struct{}

func (m *GoogleAPIFieldBehaviorModifier) ModifyTableField(
	file *presentation.File,
	col *schema.Column,
	field *proto.Field,
) {
	file.Imports.Add("google/api/annotations.proto")
	field.Options = append(field.Options, m.required(col), googleapi.FieldOutputOnly())
}

func (m *GoogleAPIFieldBehaviorModifier) required(col *schema.Column) *proto.FieldOption {
	if col.Nullable {
		return googleapi.FieldOptional()
	}
	return googleapi.FieldRequired()
}
