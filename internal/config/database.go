package config

import (
	"slices"
)

type DatabaseDriver string

const (
	DatabaseDriverPostgres DatabaseDriver = "postgres"
	DatabaseDriverDBML     DatabaseDriver = "dbml"
	DatabaseDriverMySQL    DatabaseDriver = "mysql"
)

const (
	DefaultDatabaseSchema = "public"
)

var DatabaseDrivers = []DatabaseDriver{
	DatabaseDriverPostgres,
	DatabaseDriverDBML,
	DatabaseDriverMySQL,
}

var readableDatabaseDrivers = []DatabaseDriver{
	DatabaseDriverPostgres,
	DatabaseDriverDBML,
}

var writeableDatabaseDrivers = []DatabaseDriver{
	DatabaseDriverPostgres,
}

var migrateableDatabaseDrivers = []DatabaseDriver{
	DatabaseDriverPostgres,
	DatabaseDriverMySQL,
}

type Database struct {
	Driver DatabaseDriver `yaml:"driver"`
	DSN    string         `yaml:"dsn"`
	Schema string         `yaml:"schema"`
}

func (d DatabaseDriver) Valid() bool {
	return slices.Contains(DatabaseDrivers, d)
}

func (d DatabaseDriver) CanRead() bool {
	return slices.Contains(readableDatabaseDrivers, d)
}

func (d DatabaseDriver) CanWrite() bool {
	return slices.Contains(writeableDatabaseDrivers, d)
}

func (d DatabaseDriver) CanMigrate() bool {
	return slices.Contains(migrateableDatabaseDrivers, d)
}
