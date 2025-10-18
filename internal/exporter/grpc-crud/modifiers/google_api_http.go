package modifiers

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
	"github.com/artarts36/db-exporter/internal/shared/proto"
	"github.com/artarts36/db-exporter/internal/shared/proto/opts/googleapi"
	"strings"
)

type GoogleApiHttp struct {
	PathPrefix string
}

func (m *GoogleApiHttp) ModifyProcedure(proc *presentation.Procedure) {
	basePath := fmt.Sprintf(
		"%s/%s",
		m.PathPrefix,
		proc.Service().TableMessage().GetTable().Name.Snake().Lower().Value,
	)

	var opt *proto.ServiceProcedureOption

	switch proc.Type {
	case presentation.ProcedureTypeList:
		opt = googleapi.Get(basePath)
	case presentation.ProcedureTypeGet:
		opt = googleapi.Get(m.pathTo(basePath, proc.Service().TableMessage()))
	case presentation.ProcedureTypeCreate:
		opt = googleapi.Post(basePath)
	case presentation.ProcedureTypePatch:
		opt = googleapi.Patch(m.pathTo(basePath, proc.Service().TableMessage()))
	case presentation.ProcedureTypeDelete:
		opt = googleapi.Delete(m.pathTo(basePath, proc.Service().TableMessage()))
	default:
		return
	}

	proc.Service().File().AddImport("google/api/annotations.proto")

	proc.Options = append(proc.Options, opt)
}

func (m *GoogleApiHttp) pathTo(basePath string, msg *presentation.TableMessage) string {
	return fmt.Sprintf("%s/%s", basePath, m.fieldsToPath(msg))
}

func (m *GoogleApiHttp) fieldsToPath(msg *presentation.TableMessage) string {
	path := strings.Builder{}

	for i, field := range msg.PrimaryKey {
		path.WriteString("{" + field.Name() + "}")

		if i < len(msg.PrimaryKey)-1 {
			path.WriteString("/")
		}
	}

	return path.String()
}
