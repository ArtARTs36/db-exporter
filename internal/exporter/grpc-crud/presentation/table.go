package presentation

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type TableMessage struct {
	Table *schema.Table

	Proto      *proto.Message
	Fields     map[string]*proto.Field
	PrimaryKey []*proto.Field
}
