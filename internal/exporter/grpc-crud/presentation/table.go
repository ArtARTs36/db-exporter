package presentation

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type TableMessage struct {
	Table *schema.Table

	// map[column.name]*Field
	fields     map[string]*Field
	PrimaryKey []*Field

	message *Message
}

func newTableMessage(table *schema.Table, srv *Service) *TableMessage {
	msg := &TableMessage{
		Table: table,
		message: &Message{
			proto: &proto.Message{
				Name:   table.Name.Pascal().Singular().Value,
				Fields: make([]*proto.Field, 0, len(table.Columns)),
			},
			srv: srv,
		},
		fields: make(map[string]*Field),
	}

	if table.PrimaryKey != nil {
		msg.PrimaryKey = make([]*Field, 0, table.PrimaryKey.ColumnsNames.Len())
	}

	return msg
}

func (t *TableMessage) Name() string {
	return t.message.proto.Name
}

func (t *TableMessage) GetField(name string) (*Field, bool) {
	f, ok := t.fields[name]
	return f, ok
}

func (t *TableMessage) GetTable() *schema.Table {
	return t.Table
}

func (t *TableMessage) CreateField(name string, columnName string, creator func(*Field)) *TableMessage {
	t.createField(name, columnName, creator)

	return t
}

func (t *TableMessage) CreatePrimaryKeyField(name string, columnName string, creator func(field *Field)) *TableMessage {
	field := t.createField(name, columnName, creator)
	t.PrimaryKey = append(t.PrimaryKey, field)

	return t
}

func (t *TableMessage) createField(name string, columnName string, creator func(*Field)) *Field {
	field := t.message.createField(name, creator)
	t.fields[columnName] = field

	return field
}
