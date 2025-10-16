package service

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/shared/proto/opts/googleapi"
	"strings"

	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type googleApiHTTPProcedureModifier struct {
	pathPrefix string
}

func (m *googleApiHTTPProcedureModifier) create() ProcedureModifierFactory {
	return func(file *presentation.File, srv *presentation.Service, tbl *presentation.TableMessage) ProcedureModifier {
		basePath := fmt.Sprintf("%s/%s", m.pathPrefix, tbl.Table.Name.Snake().Lower())

		return func(proc *presentation.Procedure) {
			var opt *proto.ServiceProcedureOption

			switch proc.Type {
			case presentation.ProcedureTypeList:
				opt = googleapi.Get(basePath)
			case presentation.ProcedureTypeGet:
				opt = googleapi.Get(m.pathTo(basePath, tbl))
			case presentation.ProcedureTypeCreate:
				opt = googleapi.Post(basePath)
			case presentation.ProcedureTypePatch:
				opt = googleapi.Patch(m.pathTo(basePath, tbl))
			case presentation.ProcedureTypeDelete:
				opt = googleapi.Delete(m.pathTo(basePath, tbl))
			default:
				return
			}

			file.Imports.Add("google/api/annotations.proto")

			proc.Options = append(proc.Options, opt)
		}
	}
}

func (m *googleApiHTTPProcedureModifier) pathTo(basePath string, msg *presentation.TableMessage) string {
	return fmt.Sprintf("%s/%s", basePath, m.fieldsToPath(msg))
}

func (m *googleApiHTTPProcedureModifier) fieldsToPath(msg *presentation.TableMessage) string {
	path := strings.Builder{}

	for i, field := range msg.PrimaryKey {
		path.WriteString("{" + field.Name + "}")

		if i < len(msg.PrimaryKey)-1 {
			path.WriteString("/")
		}
	}

	return path.String()
}
