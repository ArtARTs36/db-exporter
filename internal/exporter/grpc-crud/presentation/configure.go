package presentation

type Configurator func(cfg *config)

type config struct {
	modifyProcedure func(procedure *Procedure)
}

func newConfig(configurators []Configurator) *config {
	cfg := &config{}

	for _, configure := range configurators {
		configure(cfg)
	}

	if cfg.modifyProcedure == nil {
		cfg.modifyProcedure = func(procedure *Procedure) {}
	}

	return cfg
}

func WithModifyProcedure(modifier func(procedure *Procedure)) Configurator {
	return func(cfg *config) {
		cfg.modifyProcedure = modifier
	}
}
