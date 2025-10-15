package grpccrud

import "github.com/artarts36/db-exporter/internal/shared/proto"

type procedureType string

const (
	procedureTypeList   procedureType = "List"
	procedureTypeGet    procedureType = "Get"
	procedureTypeCreate procedureType = "Create"
	procedureTypePatch  procedureType = "Patch"
	procedureTypeDelete procedureType = "Delete"
)

type procedure struct {
	Name string

	Type procedureType

	Request  *proto.Message
	Response *proto.Message

	Options []*proto.ServiceProcedureOption
}

func (p *procedure) ToProto() *proto.ServiceProcedure {
	return &proto.ServiceProcedure{
		Name:    p.Name,
		Param:   p.Request.Name,
		Returns: p.Response.Name,
		Options: p.Options,
	}
}
