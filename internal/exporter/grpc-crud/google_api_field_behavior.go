package grpccrud

import (
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/tablemsg"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/db-exporter/internal/shared/proto/opts/googleapi"
)

type googleApiFieldBehaviorModifier struct{}

func (m *googleApiFieldBehaviorModifier) create() procedureModifierFactory {
	return func(file *proto.File, srv *service, tbl *tablemsg.Message) procedureModifier {
		return func(proc *procedure) {
			switch proc.Type {
			case procedureTypeGet:
				for _, field := range proc.Request.Fields {
					field.Options = append(field.Options, googleapi.FieldRequired())
				}
			case procedureTypeDelete:
				for _, field := range proc.Request.Fields {
					field.Options = append(field.Options, googleapi.FieldRequired())
				}
			case procedureTypeList:
				for _, field := range proc.Response.Fields {
					field.Options = append(field.Options, googleapi.FieldOutputOnly())
				}
			case procedureTypeCreate:
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
