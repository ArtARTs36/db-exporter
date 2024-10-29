package config

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

type Commit struct {
	Message string `yaml:"message"`
	Author  string `yaml:"author"`
	Push    bool   `yaml:"push"`
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
