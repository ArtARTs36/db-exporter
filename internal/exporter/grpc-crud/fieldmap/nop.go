package fieldmap

import (
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type Nop struct{}

func (Nop) ModifyTableField(file *presentation.File, col *schema.Column, field *proto.Field) {}
