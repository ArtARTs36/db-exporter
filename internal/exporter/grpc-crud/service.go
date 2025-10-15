package grpccrud

import "github.com/artarts36/db-exporter/internal/shared/proto"

type service struct {
	Name string

	Procedures []*procedure
	Messages   []*proto.Message
}

func (s *service) ToProto() *proto.Service {
	procs := make([]*proto.ServiceProcedure, len(s.Procedures))
	for i, proc := range s.Procedures {
		procs[i] = proc.ToProto() // @todo need to optimize
	}

	return &proto.Service{
		Name:       s.Name,
		Procedures: procs,
	}
}
