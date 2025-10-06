package config

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/env"
	"github.com/artarts36/db-exporter/internal/shared/fs"
	"github.com/artarts36/gds"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"os"
	"path/filepath"
	"strings"
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
		val, err := l.envInjector.Inject(database.DSN, nil)
		if err != nil {
			return fmt.Errorf("database[%s]: failed to inject environment variable: %w", dbName, err)
		}

		database.DSN = val
		cfg.Databases[dbName] = database
	}

	err := l.injectEnvVarsToTasks(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (l *Loader) injectEnvVarsToTasks(cfg *Config) error {
	for taskName, task := range cfg.Tasks {
		if task.Commit.Author != "" {
			val, err := l.envInjector.Inject(task.Commit.Author, nil)
			if err != nil {
				return fmt.Errorf("tasks[%s]: failed to inject environment variable into author: %w", taskName, err)
			}

			task.Commit.Author = val
		}

		for actID, activity := range task.Activities {
			if activity.Tables.List.IsNotEmpty() && activity.Tables.List.IsString() {
				val, err := l.envInjector.Inject(activity.Tables.List.First(), &env.InjectOpts{
					AllowEmptyVars: true,
				})
				if err != nil {
					return fmt.Errorf(
						"tasks[%s]: failed to inject environment variable into activity tables list: %w",
						taskName,
						err,
					)
				}

				envTablesSet := gds.NewSet[string]()

				if val != "" {
					for _, table := range strings.Split(val, ",") {
						envTablesSet.Add(strings.Trim(table, " "))
					}
				}

				cfg.Tasks[taskName].Activities[actID].Tables.List.Set = *envTablesSet
			}
		}
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

			if !db.Driver.CanRead() {
				return fmt.Errorf(
					"databases[%s] have driver %q which read schema unsupported. Available: %v",
					activity.Database,
					db.Driver,
					readableDatabaseDrivers,
				)
			}

			if activity.Export.Spec != nil {
				validatableSpec, isValidatableSpec := activity.Export.Spec.(ValidatableSpec)
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
		if database.Schema == "" {
			if database.Driver == DatabaseDriverMySQL {
				dsn, err := mysql.ParseDSN(database.DSN)
				if err != nil {
					return fmt.Errorf("parse dsn %q: %w", database.DSN, err)
				}
				database.Schema = dsn.DBName
			} else {
				database.Schema = DefaultDatabaseSchema
			}

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
