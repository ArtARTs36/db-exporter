package env

import (
	"fmt"
	"github.com/buildkite/interpolate"
)

type Injector struct {
	env interpolate.Env
}

func NewInjector() *Injector {
	return &Injector{
		env: &interpolateEnv{},
	}
}

func (i *Injector) Inject(expression string) (string, error) {
	val, err := i.inject(expression)
	if err != nil {
		return "", err
	}
	if val == "" {
		return "", fmt.Errorf("failed to interpolate expression %q", expression)
	}
	return val, nil
}

func (i *Injector) inject(expression string) (string, error) {
	return interpolate.Interpolate(i.env, expression)
}
