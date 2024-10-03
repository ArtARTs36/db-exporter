package config

type Config struct {
	Databases map[string]Database `yaml:"databases"`
	Tasks     map[string]Task     `yaml:"tasks"`
	Options   struct {
		Commit              Commit `yaml:"commit"`
		WithMigrationsTable bool   `yaml:"with_migrations_table"`
		PrintStat           bool   `yaml:"print_stat"`
		Debug               bool   `yaml:"debug"`
	} `yaml:"options"`
}

type Task []Activity

type Database struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type Commit struct {
	CommitMessage string
	CommitAuthor  string
	CommitPush    bool
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
