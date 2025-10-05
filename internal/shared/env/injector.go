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

type InjectOpts struct {
	AllowEmptyVars bool
}

func (i *Injector) Inject(expression string, opts *InjectOpts) (string, error) {
	val, err := i.inject(expression)
	if err != nil {
		return "", err
	}

	if val == "" {
		if opts != nil && opts.AllowEmptyVars {
			return "", nil
		}

		return "", fmt.Errorf("failed to interpolate expression %q", expression)
	}

	return val, nil
}

func (i *Injector) inject(expression string) (string, error) {
	return interpolate.Interpolate(i.env, expression)
}
