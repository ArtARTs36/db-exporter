package env

import "github.com/buildkite/interpolate"

type Injector struct {
	env interpolate.Env
}

func NewInjector() *Injector {
	return &Injector{
		env: &interpolateEnv{},
	}
}

func (i *Injector) Inject(expression string) (string, error) {
	return interpolate.Interpolate(i.env, expression)
}
