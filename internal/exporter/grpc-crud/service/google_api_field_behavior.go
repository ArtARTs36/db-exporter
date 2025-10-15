package service

import (
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/tablemsg"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/db-exporter/internal/shared/proto/opts/googleapi"
)

type googleApiFieldBehaviorModifier struct{}

func (m *googleApiFieldBehaviorModifier) create() ProcedureModifierFactory {
	return func(file *proto.File, srv *Service, tbl *tablemsg.Message) ProcedureModifier {
		return func(proc *Procedure) {
			switch proc.Type {
			case ProcedureTypeGet:
				for _, field := range proc.Request.Fields {
					field.Options = append(field.Options, googleapi.FieldRequired())
				}
			case ProcedureTypeDelete:
				for _, field := range proc.Request.Fields {
					field.Options = append(field.Options, googleapi.FieldRequired())
				}
			case ProcedureTypeList:
				for _, field := range proc.Response.Fields {
					field.Options = append(field.Options, googleapi.FieldOutputOnly())
				}
			case ProcedureTypeCreate:
				for _, field := range proc.Request.Fields {
					col := tbl.Table.GetColumn(field.Name)
					if col == nil {
						continue
					}

					if col.Nullable {
						field.Options = append(field.Options, googleapi.FieldOptional())
					} else {
						field.Options = append(field.Options, googleapi.FieldRequired())
					}
				}
			}
		}
	}
}
