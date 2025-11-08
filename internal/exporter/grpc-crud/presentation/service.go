package presentation

import (
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type Service struct {
	table *TableMessage

	file *File

	proto *proto.Service
}

func (s *Service) AddProcedure(name string, typ ProcedureType, reqBuild func(message *Message)) *Service {
	req := newMessage(MessageTypeRequest, s)

	reqBuild(req)

	p := &Procedure{
		proto: &proto.ServiceProcedure{
			Name:    name,
			Param:   req.proto.Name,
			Returns: s.table.Name(),
			Options: make([]*proto.ServiceProcedureOption, 0),
		},
		typ:     typ,
		service: s,
	}

	s.file.AddMessage(req.proto)

	s.file.cfg.modifyProcedure(p)

	s.proto.Procedures = append(s.proto.Procedures, p.proto)

	return s
}

func (s *Service) AddProcedureFn(
	name string,
	typ ProcedureType,
	reqBuild func(message *Message),
	respBuild func(message *Message),
) *Service {
	req := newMessage(MessageTypeRequest, s)
	resp := newMessage(MessageTypeResponse, s)

	reqBuild(req)
	respBuild(resp)

	p := &Procedure{
		proto: &proto.ServiceProcedure{
			Name:    name,
			Param:   req.proto.Name,
			Returns: resp.proto.Name,
			Options: make([]*proto.ServiceProcedureOption, 0),
		},
		typ:     typ,
		service: s,
	}

	s.file.AddMessage(req.proto)
	s.file.AddMessage(resp.proto)

	s.file.cfg.modifyProcedure(p)

	s.proto.Procedures = append(s.proto.Procedures, p.proto)

	return s
}

func (s *Service) File() *File {
	return s.file
}

func (s *Service) HasProcedures() bool {
	return len(s.proto.Procedures) > 0
}

func (s *Service) TableMessage() *TableMessage {
	return s.table
}

func (s *Service) SetCommentTop(comment string) *Service {
	s.proto.CommentTop = comment

	return s
}
