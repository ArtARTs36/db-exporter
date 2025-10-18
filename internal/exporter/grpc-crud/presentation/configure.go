package presentation

type Configurator func(cfg *config)

type config struct {
	modifyProcedure func(procedure *Procedure)
	modifyField     func(field *Field)
}

func newConfig(configurators []Configurator) *config {
	cfg := &config{}

	for _, configure := range configurators {
		configure(cfg)
	}

	if cfg.modifyProcedure == nil {
		cfg.modifyProcedure = func(procedure *Procedure) {}
	}

	if cfg.modifyField == nil {
		cfg.modifyField = func(field *Field) {}
	}

	return cfg
}

func WithModifyProcedure(modifier func(procedure *Procedure)) Configurator {
	return func(cfg *config) {
		cfg.modifyProcedure = modifier
	}
}

func WithModifyField(modifier func(field *Field)) Configurator {
	return func(cfg *config) {
		cfg.modifyField = modifier
	}
}
