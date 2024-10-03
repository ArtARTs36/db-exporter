package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Loader struct {
}

func (l *Loader) Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", path, err)
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	err = l.fillDefaults(&cfg)
	if err != nil {
		return nil, err
	}

	if err = l.validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (l *Loader) validate(cfg *Config) error {
	for tid, task := range cfg.Tasks {
		for aid, activity := range task.Activities {
			if _, ok := cfg.Databases[activity.Database]; !ok {
				return fmt.Errorf("task[%s][%d] have invalid database name %q", tid, aid, activity.Database)
			}
		}
	}

	return nil
}

func (l *Loader) fillDefaults(cfg *Config) error {
	defaultDB, exists := cfg.GetDefaultDatabaseName()
	if !exists {
		return fmt.Errorf("databases not filled")
	}

	for tid, task := range cfg.Tasks {
		for aid, activity := range task.Activities {
			if activity.Database == "" {
				cfg.Tasks[tid].Activities[aid].Database = defaultDB
			}
		}
	}

	return nil
}
