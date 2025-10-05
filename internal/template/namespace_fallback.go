package template

import (
	"errors"

	"github.com/tyler-sommer/stick"
)

type NamespaceFallbackLoader struct {
	namespaceLoader stick.Loader
	fallbackLoader  stick.Loader
}

func NewNamespaceFallbackLoader(
	namespaceLoader stick.Loader,
	fallbackLoader stick.Loader,
) *NamespaceFallbackLoader {
	return &NamespaceFallbackLoader{
		namespaceLoader: namespaceLoader,
		fallbackLoader:  fallbackLoader,
	}
}

func (l *NamespaceFallbackLoader) Load(name string) (stick.Template, error) {
	template, err := l.namespaceLoader.Load(name)
	if err != nil {
		if errors.Is(err, errNameNotContainsNamespace) {
			return l.fallbackLoader.Load(name)
		}

		return nil, err
	}

	return template, nil
}
