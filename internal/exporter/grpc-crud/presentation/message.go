package presentation

import "github.com/artarts36/db-exporter/internal/shared/proto"

type Message struct {
	proto *proto.Message
}

func newMessage() *Message {
	return &Message{
		proto: &proto.Message{},
	}
}

func (msg *Message) SetName(name string) *Message {
	msg.proto.Name = name

	return msg
}

func (msg *Message) CreateField(name string, creator func(*Field)) *Message {
	field := &Field{
		proto: &proto.Field{
			Name:    name,
			Options: make([]*proto.FieldOption, 0),
			ID:      len(msg.proto.Fields) + 1,
			Type:    "string",
		},
	}

	creator(field)

	msg.proto.Fields = append(msg.proto.Fields, field.proto)

	return msg
}
