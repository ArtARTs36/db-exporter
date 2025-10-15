package service

import "github.com/artarts36/db-exporter/internal/shared/proto"

type ProcedureType string

const (
	ProcedureTypeList   ProcedureType = "List"
	ProcedureTypeGet    ProcedureType = "Get"
	ProcedureTypeCreate ProcedureType = "Create"
	ProcedureTypePatch  ProcedureType = "Patch"
	ProcedureTypeDelete ProcedureType = "Delete"
)

type Procedure struct {
	Name string

	Type ProcedureType

	Request  *proto.Message
	Response *proto.Message

	Options []*proto.ServiceProcedureOption
}

func (p *Procedure) ToProto() *proto.ServiceProcedure {
	return &proto.ServiceProcedure{
		Name:    p.Name,
		Param:   p.Request.Name,
		Returns: p.Response.Name,
		Options: p.Options,
	}
}
