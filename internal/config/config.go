package config

import "fmt"

type Config struct {
	Databases map[string]Database `yaml:"databases"`
	Tasks     map[string]Task     `yaml:"tasks"`
	Options   struct {
		WithMigrationsTable bool `yaml:"with_migrations_table"`
		PrintStat           bool `yaml:"print_stat"`
		Debug               bool `yaml:"debug"`
	} `yaml:"options"`
}

type Task struct {
	Activities []Activity `yaml:"activities"`
	Commit     Commit     `yaml:"commit"`
}

type Database struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type Commit struct {
	Message string `yaml:"message"`
	Author  string `yaml:"author"`
	Push    bool   `yaml:"push"`
}

func (c *Config) OnlyTasks(taskNames []string) error {
	tasks := map[string]Task{}
	for _, name := range taskNames {
		task, ok := c.Tasks[name]
		if !ok {
			return fmt.Errorf("task %q not found", name)
		}

		tasks[name] = task
	}

	c.Tasks = tasks

	return nil
}

func (c *Commit) Valid() bool {
	return c.Message != "" || c.Author != "" || c.Push
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
