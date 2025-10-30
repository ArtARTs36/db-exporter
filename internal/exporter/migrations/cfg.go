package migrations

import (
	"errors"
	"fmt"

	"github.com/artarts36/db-exporter/internal/schema"
)

type Specification struct {
	Use struct {
		IfNotExists bool `yaml:"if_not_exists" json:"if_not_exists"`
		IfExists    bool `yaml:"if_exists" json:"if_exists"`
	} `yaml:"use"`
	Target schema.DatabaseDriver `yaml:"target" json:"target"`
}

func (m *Specification) InjectDatabaseDriver(driver schema.DatabaseDriver) {
	if m.Target != "" {
		return
	}

	m.Target = driver
}

func (m *Specification) Validate() error {
	if m.Target == "" {
		return errors.New("target is required")
	}

	if !m.Target.Valid() {
		return fmt.Errorf(
			"target have unsupported driver %q. Available: %v",
			m.Target,
			schema.GetWriteableDatabaseDrivers(),
		)
	}

	if !m.Target.CanMigrate() {
		return fmt.Errorf("target have driver %q, which unsupported migrate queries", m.Target)
	}

	return nil
}
