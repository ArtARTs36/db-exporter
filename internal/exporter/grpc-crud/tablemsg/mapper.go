package tablemsg

import (
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type Mapper struct {
	fieldModifier fieldModifier
}

type fieldModifier interface {
	ModifyTableField(file *proto.File, col *schema.Column, field *proto.Field)
}

func NewMapper(fieldmod fieldModifier) *Mapper {
	return &Mapper{
		fieldModifier: fieldmod,
	}
}

func (m *Mapper) MapTable(
	file *proto.File,
	table *schema.Table,
	fieldTypeMapper func(col *schema.Column) string,
) *presentation.TableMessage {
	msg := &presentation.TableMessage{
		Table: table,
		Proto: &proto.Message{
			Name:   table.Name.Pascal().Singular().Value,
			Fields: make([]*proto.Field, 0, len(table.Columns)),
		},
		Fields: make(map[string]*proto.Field),
	}

	if table.PrimaryKey != nil {
		msg.PrimaryKey = make([]*proto.Field, 0, table.PrimaryKey.ColumnsNames.Len())
	}

	for i, column := range table.Columns {
		field := &proto.Field{
			Name: column.Name.Snake().Lower().Value,
			Type: fieldTypeMapper(column),
			ID:   i + 1,
		}

		m.fieldModifier.ModifyTableField(file, column, field)

		msg.Proto.Fields = append(msg.Proto.Fields, field)
		msg.Fields[column.Name.Value] = field

		if column.PrimaryKey != nil {
			msg.PrimaryKey = append(msg.PrimaryKey, field)
		}
	}

	return msg
}
