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

var writeableDatabaseDrivers = []DatabaseDriver{
	DatabaseDriverPostgres,
}

type Database struct {
	Driver DatabaseDriver `yaml:"driver"`
	DSN    string         `yaml:"dsn"`
	Schema string         `yaml:"schema"`
}

func (d DatabaseDriver) Valid() bool {
	return slices.Contains(DatabaseDrivers, d)
}

func (d DatabaseDriver) CanWrite() bool {
	switch d {
	case DatabaseDriverPostgres:
		return true
	case DatabaseDriverDBML:
		return false
	}
	return true
}
