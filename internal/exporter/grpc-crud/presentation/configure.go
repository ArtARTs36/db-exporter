package presentation

type Configurator func(cfg *config)

type config struct {
	modifyProcedure func(procedure *Procedure)
	modifyField     func(field *Field)
	modifyService   func(service *Service)
}

func newConfig(configurators []Configurator) *config {
	cfg := &config{}

	for _, configure := range configurators {
		configure(cfg)
	}

	if cfg.modifyProcedure == nil {
		cfg.modifyProcedure = func(*Procedure) {}
	}

	if cfg.modifyField == nil {
		cfg.modifyField = func(*Field) {}
	}

	if cfg.modifyService == nil {
		cfg.modifyService = func(*Service) {}
	}

	return cfg
}

func WithModifyProcedure(modifier func(procedure *Procedure)) Configurator {
	return func(cfg *config) {
		if cfg.modifyProcedure == nil {
			cfg.modifyProcedure = modifier
		} else {
			prevModifier := cfg.modifyProcedure

			cfg.modifyProcedure = func(procedure *Procedure) {
				modifier(procedure)
				prevModifier(procedure)
			}
		}
	}
}

func WithModifyField(modifier func(field *Field)) Configurator {
	return func(cfg *config) {
		if cfg.modifyField == nil {
			cfg.modifyField = modifier
		} else {
			prevModifier := cfg.modifyField

			cfg.modifyField = func(field *Field) {
				modifier(field)
				prevModifier(field)
			}
		}
	}
}

func WithModifyService(modifier func(*Service)) Configurator {
	return func(cfg *config) {
		cfg.modifyService = modifier
	}
}
