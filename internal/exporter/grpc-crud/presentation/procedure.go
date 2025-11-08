package presentation

import "github.com/artarts36/db-exporter/internal/shared/proto"

type ProcedureType string

const (
	ProcedureTypeList     ProcedureType = "List"
	ProcedureTypeGet      ProcedureType = "Get"
	ProcedureTypeCreate   ProcedureType = "Create"
	ProcedureTypeUpdate   ProcedureType = "Update"
	ProcedureTypeDelete   ProcedureType = "Delete"
	ProcedureTypeUndelete ProcedureType = "Undelete"
)

type Procedure struct {
	typ ProcedureType

	service *Service

	proto *proto.ServiceProcedure
}

func (p *Procedure) ToProto() *proto.ServiceProcedure {
	return &proto.ServiceProcedure{}
}

func (p *Procedure) AddOption(option *proto.ServiceProcedureOption) {
	p.proto.Options = append(p.proto.Options, option)
}

func (p *Procedure) Service() *Service {
	return p.service
}

func (p *Procedure) SetCommentTop(comment string) *Procedure {
	p.proto.CommentTop = comment
	return p
}

func (p *Procedure) Type() ProcedureType {
	return p.typ
}
