package grpccrud

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/tablemsg"
	"github.com/artarts36/db-exporter/internal/shared/proto"
)

type (
	procedureModifierFactory func(
		file *proto.File,
		srv *service,
		tableMessage *tablemsg.Message,
	) procedureModifier

	procedureModifier func(proc *procedure)
)

func compositeProcedureModifier(modifiers []procedureModifier) procedureModifier {
	return func(proc *procedure) {
		for _, modifier := range modifiers {
			modifier(proc)
		}
	}
}

func compositeProcedureModifierFactory(factories []procedureModifierFactory) procedureModifierFactory {
	return func(file *proto.File, srv *service, tableMessage *tablemsg.Message) procedureModifier {
		modifiers := make([]procedureModifier, len(factories))

		for i, f := range factories {
			modifiers[i] = f(file, srv, tableMessage)
		}

		return compositeProcedureModifier(modifiers)
	}
}

func selectProcedureModifier(spec *config.GRPCCrudExportSpec) procedureModifierFactory {
	if spec.With.Object == nil {
		return nil
	}

	modifierFactories := []procedureModifierFactory{}

	if spec.With.Object.GoogleApiHTTP.Object != nil {
		m := &googleApiHTTPProcedureModifier{}
		modifierFactories = append(modifierFactories, m.create())
	}

	if spec.With.Object.GoogleAPIFieldBehavior.Object != nil {
		m := &googleApiFieldBehaviorModifier{}
		modifierFactories = append(modifierFactories, m.create())
	}

	if len(modifierFactories) == 0 {
		return nopProcedureModifier()
	}

	if len(modifierFactories) == 1 {
		return modifierFactories[0]
	}

	return compositeProcedureModifierFactory(modifierFactories)
}

func nopProcedureModifier() procedureModifierFactory {
	return procedureModifierFactory(nil)
}
