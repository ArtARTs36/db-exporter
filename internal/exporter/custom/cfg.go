package custom

import "fmt"

type Specification struct {
	Template string `yaml:"template" json:"template"`
	Output   struct {
		Extension string `yaml:"extension" json:"extension"`
	} `yaml:"output" json:"output"`
}

func (s *Specification) Validate() error {
	if s.Template == "" {
		return fmt.Errorf("custom export template is required")
	}

	return nil
}
