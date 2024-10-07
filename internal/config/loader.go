package config

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/env"
	"os"

	"gopkg.in/yaml.v3"
)

type Loader struct {
	envInjector *env.Injector
}

func NewLoader(envInjector *env.Injector) *Loader {
	return &Loader{
		envInjector: envInjector,
	}
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

	err = l.injectEnvVars(&cfg)
	if err != nil {
		return nil, err
	}

	if err = l.validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (l *Loader) injectEnvVars(cfg *Config) error {
	for dbName, database := range cfg.Databases {
		val, err := l.envInjector.Inject(database.DSN)
		if err != nil {
			return fmt.Errorf("database[%s]: failed to inject environment variable: %w", dbName, err)
		}

		database.DSN = val
		cfg.Databases[dbName] = database
	}

	for taskName, task := range cfg.Tasks {
		if task.Commit.Author == "" {
			continue
		}

		val, err := l.envInjector.Inject(task.Commit.Author)
		if err != nil {
			return fmt.Errorf("tasks[%s]: failed to inject environment variable into author: %w", taskName, err)
		}

		task.Commit.Author = val
	}

	return nil
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
