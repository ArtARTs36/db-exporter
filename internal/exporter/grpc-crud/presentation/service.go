package presentation

import (
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type Service struct {
	Name string

	Procedures []*Procedure
	Messages   []*proto.Message
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
