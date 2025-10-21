package config

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"os"
	"path/filepath"

	"github.com/artarts36/db-exporter/internal/shared/fs"
)

type Loader struct {
	files   fs.Driver
	parsers map[string]Parser
}

func NewLoader(
	files fs.Driver,
	parsers map[string]Parser,
) *Loader {
	return &Loader{
		files:   files,
		parsers: parsers,
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

	if err = l.validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (l *Loader) validate(cfg *Config) error { //nolint:gocognit // todo
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

			if !db.Driver.CanRead() {
				return fmt.Errorf(
					"databases[%s] have driver %q which read schema unsupported. Available: %v",
					activity.Database,
					db.Driver,
					readableDatabaseDrivers,
				)
			}

			if activity.Spec != nil {
				expectingDatabaseDriver, isExpectingDatabaseDriver := activity.Spec.(ExpectingDatabaseDriver)
				if isExpectingDatabaseDriver {
					expectingDatabaseDriver.InjectDatabaseDriver(db.Driver)
				}

				validatableSpec, isValidatableSpec := activity.Spec.(ValidatableSpec)
				if isValidatableSpec {
					err := validatableSpec.Validate()
					if err != nil {
						return fmt.Errorf("task[%s][%d] have invalid spec: %w", tid, aid, err)
					}
				}
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
		if database.Schema != "" {
			continue
		}

		if database.Driver == DatabaseDriverMySQL {
			dsn, err := mysql.ParseDSN(database.DSN.Value)
			if err != nil {
				return fmt.Errorf("parse dsn %q: %w", database.DSN.Value, err)
			}
			database.Schema = dsn.DBName
		} else {
			database.Schema = DefaultDatabaseSchema
		}

		cfg.Databases[name] = database
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
