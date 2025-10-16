package service

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
)

type (
	ProcedureModifierFactory func(
		file *presentation.File,
		srv *presentation.Service,
		tableMessage *presentation.TableMessage,
	) ProcedureModifier

	ProcedureModifier func(proc *presentation.Procedure)
)

func CompositeProcedureModifier(modifiers []ProcedureModifier) ProcedureModifier {
	return func(proc *presentation.Procedure) {
		for _, modifier := range modifiers {
			modifier(proc)
		}
	}
}

func CompositeProcedureModifierFactory(factories []ProcedureModifierFactory) ProcedureModifierFactory {
	return func(file *presentation.File, srv *presentation.Service, tableMessage *presentation.TableMessage) ProcedureModifier {
		modifiers := make([]ProcedureModifier, len(factories))

		for i, f := range factories {
			modifiers[i] = f(file, srv, tableMessage)
		}

		return CompositeProcedureModifier(modifiers)
	}
}

func SelectProcedureModifier(spec *config.GRPCCrudExportSpec) ProcedureModifierFactory {
	if spec.With.Object == nil {
		return nil
	}

	modifierFactories := []ProcedureModifierFactory{}

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

	return CompositeProcedureModifierFactory(modifierFactories)
}

func nopProcedureModifier() ProcedureModifierFactory {
	return ProcedureModifierFactory(nil)
}
