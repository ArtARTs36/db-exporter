package config

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/env"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"os"
	"path/filepath"
)

type Loader struct {
	files       fs.Driver
	envInjector *env.Injector
	parsers     map[string]Parser
}

func NewLoader(
	files fs.Driver,
	envInjector *env.Injector,
	parsers map[string]Parser,
) *Loader {
	return &Loader{
		files:       files,
		envInjector: envInjector,
		parsers:     parsers,
	}
}

func (l *Loader) Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", path, err)
	}

	pathExt := filepath.Ext(path)
	parser, ok := l.parsers[pathExt]
	if !ok {
		return nil, fmt.Errorf("parser for exension %q not found", pathExt)
	}

	cfg, err := parser(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s config: %w", pathExt, err)
	}

	err = l.fillDefaults(cfg)
	if err != nil {
		return nil, err
	}

	err = l.injectEnvVars(cfg)
	if err != nil {
		return nil, err
	}

	if err = l.validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
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
			db, ok := cfg.Databases[activity.Database]
			if !ok {
				return fmt.Errorf("task[%s][%d] have invalid database name %q", tid, aid, activity.Database)
			}

			if !db.Driver.Valid() {
				return fmt.Errorf(
					"databases[%s] have unsupported driver %q. Available: %v",
					activity.Database,
					db.Driver,
					DatabaseDrivers,
				)
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

	for name, database := range cfg.Databases {
		if database.Schema == "" {
			database.Schema = DefaultDatabaseSchema
			cfg.Databases[name] = database
		}
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
