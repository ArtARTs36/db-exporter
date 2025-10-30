package markdown

import (
	"fmt"

	"github.com/artarts36/db-exporter/internal/exporter/diagram"
)

type Specification struct {
	WithDiagram bool                  `yaml:"with_diagram" json:"with_diagram"`
	Diagram     diagram.Specification `yaml:"diagram" json:"diagram"`
}

func (s *Specification) Validate() error {
	if err := s.Diagram.Validate(); err != nil {
		return fmt.Errorf("diagram: %w", err)
	}
	return nil
}
