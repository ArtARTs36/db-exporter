package config

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/specw"
)

type Config struct {
	Databases map[string]schema.Database `yaml:"databases" json:"databases"`
	Tasks     map[string]Task            `yaml:"tasks" json:"tasks"`
	Options   struct {
		WithMigrationsTable bool `yaml:"with_migrations_table" json:"with_migrations_table"`
		NoPrintStat         bool `yaml:"no_print_stat" json:"no_print_stat"`
		Debug               bool `yaml:"debug" json:"debug"`
	} `yaml:"options"`
}

type Task struct {
	Activities []Activity `yaml:"activities" json:"activities"`
	Commit     Commit     `yaml:"commit" json:"commit"`
}

type Commit struct {
	Message string                         `yaml:"message" json:"message"`
	Author  *specw.Env[specw.GitCommitter] `yaml:"author" json:"author"`
	Push    bool                           `yaml:"push" json:"push"`
}

func (c *Commit) Valid() bool {
	return c.Message != "" || c.Author != nil || c.Push
}

func (c *Config) GetDefaultDatabaseName() (string, bool) {
	const defaultDBKey = "default"

	if len(c.Databases) == 0 {
		return "", false
	}

	if _, exists := c.Databases[defaultDBKey]; exists {
		return defaultDBKey, true
	}

	for dbName := range c.Databases {
		return dbName, true
	}

	return "", false
}

func (c *Config) GetDefaultDatabase() (schema.Database, bool) {
	const defaultDBKey = "default"

	if len(c.Databases) == 0 {
		return schema.Database{}, false
	}

	if db, exists := c.Databases[defaultDBKey]; exists {
		return db, true
	}

	for _, db := range c.Databases {
		return db, true
	}

	return schema.Database{}, false
}

func (c *Config) UsingDatabases() map[string]schema.Database {
	dbs := map[string]schema.Database{}

	for _, task := range c.Tasks {
		for _, activity := range task.Activities {
			if _, ok := dbs[activity.Database]; !ok {
				dbs[activity.Database] = c.Databases[activity.Database]
			}
		}
	}

	return dbs
}
