package tablemsg

import (
	"fmt"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type Message struct {
	Table *schema.Table

	Proto      *proto.Message
	Fields     map[string]*proto.Field
	PrimaryKey []*proto.Field
}

func (m *Message) CloneField(columnName string) (*proto.Field, error) {
	field, ok := m.Fields[columnName]
	if !ok {
		return nil, fmt.Errorf("field %s not found", columnName)
	}

	return field.Clone(), nil
}
