package config

import (
	"slices"
)

type DatabaseDriver string

const (
	DatabaseDriverPostgres DatabaseDriver = "postgres"
	DatabaseDriverDBML     DatabaseDriver = "dbml"
)

const (
	DefaultDatabaseSchema = "public"
)

var DatabaseDrivers = []DatabaseDriver{
	DatabaseDriverPostgres,
	DatabaseDriverDBML,
}

type Database struct {
	Driver DatabaseDriver `yaml:"driver"`
	DSN    string         `yaml:"dsn"`
	Schema string         `yaml:"schema"`
}

func (d DatabaseDriver) Valid() bool {
	return slices.Contains(DatabaseDrivers, d)
}
