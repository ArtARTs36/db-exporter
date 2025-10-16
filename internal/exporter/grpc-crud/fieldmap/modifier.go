package fieldmap

import (
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type Modifier interface {
	ModifyTableField(file *presentation.File, col *schema.Column, field *proto.Field)
}

type compositeModifier struct {
	modifiers []Modifier
}

func Compose(modifiers []Modifier) Modifier {
	return &compositeModifier{
		modifiers: modifiers,
	}
}

func (m *compositeModifier) ModifyTableField(file *presentation.File, col *schema.Column, field *proto.Field) {
	for _, modifier := range m.modifiers {
		modifier.ModifyTableField(file, col, field)
	}
}
