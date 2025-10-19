package presentation

import (
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type Service struct {
	table *TableMessage

	file *File

	proto *proto.Service
}

func (s *Service) AddProcedureFn(
	name string,
	typ ProcedureType,
	reqBuild func(message *Message),
	respBuild func(message *Message),
) *Procedure {
	req := newMessage(s)
	resp := newMessage(s)

	reqBuild(req)
	respBuild(resp)

	p := &Procedure{
		Name:     name,
		Type:     typ,
		Request:  req.proto,
		Response: resp.proto,
		Options:  make([]*proto.ServiceProcedureOption, 0),
		service:  s,
	}

	s.file.AddMessage(req.proto)
	s.file.AddMessage(resp.proto)

	s.file.cfg.modifyProcedure(p)

	s.proto.Procedures = append(s.proto.Procedures, p.ToProto())

	return p
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
