package presentation

import (
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type Service struct {
	Name string

	TableMessage *TableMessage

	Procedures []*Procedure
	Messages   []*proto.Message

	file *File
}

func (s *Service) ToProto() *proto.Service {
	procs := make([]*proto.ServiceProcedure, len(s.Procedures))
	for i, proc := range s.Procedures {
		procs[i] = proc.ToProto() // @todo need to optimize
	}

	return &proto.Service{
		Name:       s.Name,
		Procedures: procs,
	}
}

func (s *Service) AddProcedure(
	name string,
	typ ProcedureType,
	req *proto.Message,
	resp *proto.Message,
) *Procedure {
	p := &Procedure{
		Name:     name,
		Type:     typ,
		Request:  req,
		Response: resp,
		Options:  make([]*proto.ServiceProcedureOption, 0),
		service:  s,
	}

	s.Procedures = append(s.Procedures, p)
	s.Messages = append(s.Messages, req, resp)
	//s.file.AddMessage(req)
	//s.file.AddMessage(resp)

	s.file.cfg.modifyProcedure(p)

	return p
}

func (s *Service) File() *File {
	return s.file
}
