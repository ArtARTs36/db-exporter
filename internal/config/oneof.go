package config

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/artarts36/gds"
)

type stringOrStringSlice struct {
	gds.Set[string]

	isString bool
}

func (s *stringOrStringSlice) IsString() bool {
	return s.isString
}

func (s *stringOrStringSlice) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind { //nolint:exhaustive//not need
	case yaml.ScalarNode:
		s.Set = *gds.NewSet[string](n.Value)

		s.isString = true

		return nil
	case yaml.SequenceNode:
		err := n.Decode(&s.Set)
		if err != nil {
			return err
		}

		s.isString = false

		return nil
	}

	return fmt.Errorf("expected scalar or sequence node, got %v", n.Kind)
}
