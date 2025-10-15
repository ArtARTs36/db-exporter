package grpccrud

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type (
	procedureModifierFactory func(file *proto.File, srv *service, table *schema.Table) procedureModifier

	procedureModifier func(proc *procedure)
)

func selectProcedureModifier(spec *config.GRPCCrudExportSpec) procedureModifierFactory {
	if spec.With.Object == nil {
		return nil
	}

	if spec.With.Object.GoogleApiHTTP.Object != nil {
		m := &googleApiHTTPProcedureModifier{}

		return m.create()
	}

	return nopProcedureModifier()
}

func nopProcedureModifier() procedureModifierFactory {
	return procedureModifierFactory(nil)
}
