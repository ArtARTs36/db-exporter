package service

import (
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/shared/proto/opts/googleapi"
)

type googleApiFieldBehaviorModifier struct{}

func (m *googleApiFieldBehaviorModifier) create() ProcedureModifierFactory {
	return func(file *presentation.File, srv *presentation.Service, tbl *presentation.TableMessage) ProcedureModifier {
		return func(proc *presentation.Procedure) {
			switch proc.Type {
			case presentation.ProcedureTypeGet:
				for _, field := range proc.Request.Fields {
					field.Options = append(field.Options, googleapi.FieldRequired())
				}
			case presentation.ProcedureTypeDelete:
				for _, field := range proc.Request.Fields {
					field.Options = append(field.Options, googleapi.FieldRequired())
				}
			case presentation.ProcedureTypeList:
				for _, field := range proc.Response.Fields {
					field.Options = append(field.Options, googleapi.FieldOutputOnly())
				}
			case presentation.ProcedureTypeCreate:
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
