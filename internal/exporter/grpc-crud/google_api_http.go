package grpccrud

import (
	"fmt"
	"strings"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/db-exporter/internal/shared/proto/opts/googleapihttp"
)

type googleApiHTTPProcedureModifier struct {
	pathPrefix string
}

func (m *googleApiHTTPProcedureModifier) create() procedureModifierFactory {
	return func(file *proto.File, srv *service, table *schema.Table, tableMessage *tableMessage) procedureModifier {
		basePath := fmt.Sprintf("%s/%s", m.pathPrefix, table.Name.Snake().Lower())

		return func(proc *procedure) {
			var opt *proto.ServiceProcedureOption

			switch proc.Type {
			case procedureTypeList:
				opt = googleapihttp.Get(basePath)
			case procedureTypeGet:
				opt = googleapihttp.Get(m.pathTo(basePath, tableMessage))
			case procedureTypeCreate:
				opt = googleapihttp.Post(basePath)
			case procedureTypePatch:
				opt = googleapihttp.Patch(m.pathTo(basePath, tableMessage))
			case procedureTypeDelete:
				opt = googleapihttp.Delete(m.pathTo(basePath, tableMessage))
			default:
				return
			}

			file.Imports.Add("google/api/annotations.proto")

			proc.Options = append(proc.Options, opt)
		}
	}
}

func (m *googleApiHTTPProcedureModifier) pathTo(basePath string, msg *tableMessage) string {
	return fmt.Sprintf("%s/%s", basePath, m.fieldsToPath(msg))
}

func (m *googleApiHTTPProcedureModifier) fieldsToPath(msg *tableMessage) string {
	path := strings.Builder{}

	for i, field := range msg.PrimaryKey {
		path.WriteString("{" + field.Name + "}")

		if i < len(msg.PrimaryKey)-1 {
			path.WriteString("/")
		}
	}

	return path.String()
}
