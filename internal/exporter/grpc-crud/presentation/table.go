package presentation

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type TableMessage struct {
	Table *schema.Table

	Proto *proto.Message

	// map[column.name]proto.Field
	Fields     map[string]*proto.Field
	PrimaryKey []*proto.Field
}

func (t *TableMessage) Name() string {
	return t.Proto.Name
}

func (t *TableMessage) GetField(name string) (*proto.Field, bool) {
	f, ok := t.Fields[name]
	return f, ok
}

func (t *TableMessage) GetTable() *schema.Table {
	return t.Table
}
