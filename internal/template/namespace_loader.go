package template

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tyler-sommer/stick"
)

type NamespaceLoader struct {
	namespaces map[string]stick.Loader
}

var errNameNotContainsNamespace = errors.New("name not contains namespace")

func NewNamespaceLoader(namespaces map[string]stick.Loader) *NamespaceLoader {
	return &NamespaceLoader{
		namespaces: namespaces,
	}
}

func (l *NamespaceLoader) Load(fullName string) (stick.Template, error) {
	namespace, name, err := l.splitName(fullName)
	if err != nil {
		return nil, err
	}

	loader, ok := l.namespaces[namespace]
	if !ok {
		return nil, fmt.Errorf("namespace %q not found", namespace)
	}

	return loader.Load(name)
}

func (*NamespaceLoader) splitName(fullName string) (namespace string, name string, err error) {
	if len(fullName) == 0 {
		return "", "", errors.New("empty file name")
	}

	if !strings.HasPrefix(fullName, "@") {
		return "", "", errNameNotContainsNamespace
	}

	fullName = fullName[1:]
	sepIndex := strings.Index(fullName, "/")

	return fullName[0:sepIndex], fullName[sepIndex+1:], nil
}
