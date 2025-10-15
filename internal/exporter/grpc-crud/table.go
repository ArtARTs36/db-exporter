package grpccrud

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type tableMessage struct {
	Proto      *proto.Message
	Fields     map[string]*proto.Field
	PrimaryKey []*proto.Field
}

func (m *tableMessage) CloneField(columnName string) (*proto.Field, error) {
	field, ok := m.Fields[columnName]
	if !ok {
		return nil, fmt.Errorf("field %s not found", columnName)
	}

	return field.Clone(), nil
}

func newTableMessage(table *schema.Table, fieldTypeMapper func(col *schema.Column) string) *tableMessage {
	msg := &tableMessage{
		Proto: &proto.Message{
			Name:   table.Name.Pascal().Singular().Value,
			Fields: make([]*proto.Field, 0, len(table.Columns)),
		},
		Fields:     make(map[string]*proto.Field),
		PrimaryKey: make([]*proto.Field, 0, len(table.Columns)),
	}

	for i, column := range table.Columns {
		field := &proto.Field{
			Name: column.Name.Snake().Lower().Value,
			Type: fieldTypeMapper(column),
			ID:   i + 1,
		}

		msg.Proto.Fields = append(msg.Proto.Fields, field)
		msg.Fields[column.Name.Value] = field

		if column.PrimaryKey != nil {
			msg.PrimaryKey = append(msg.PrimaryKey, field)
		}
	}

	return msg
}
