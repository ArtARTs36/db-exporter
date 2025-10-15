package fieldmap

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type Nop struct{}

func (Nop) ModifyTableField(file *proto.File, col *schema.Column, field *proto.Field) {}
