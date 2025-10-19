package presentation

import (
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type MessageType uint8

const (
	MessageTypeRequest MessageType = iota
	MessageTypeResponse
	MessageTypeTable
)

type Message struct {
	proto *proto.Message
	srv   *Service
	typ   MessageType
}

func newMessage(typ MessageType, srv *Service) *Message {
	return &Message{
		proto: &proto.Message{},
		srv:   srv,
		typ:   typ,
	}
}

func (msg *Message) SetName(name string) *Message {
	msg.proto.Name = name

	return msg
}

func (msg *Message) CreateField(name string, creator func(*Field)) *Message {
	msg.createField(name, creator)

	return msg
}

func (msg *Message) Type() MessageType {
	return msg.typ
}

func (msg *Message) createField(name string, creator func(*Field)) *Field {
	field := &Field{
		proto: &proto.Field{
			Name:    name,
			Options: make([]*proto.FieldOption, 0),
			ID:      len(msg.proto.Fields) + 1,
			Type:    "string",
		},
		message: msg,
	}

	creator(field)

	msg.srv.file.cfg.modifyField(field)

	msg.proto.Fields = append(msg.proto.Fields, field.proto)

	return field
}

func (msg *Message) Service() *Service {
	return msg.srv
}
