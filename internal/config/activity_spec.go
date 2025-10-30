package config

import "github.com/artarts36/db-exporter/internal/schema"

type ValidatableSpec interface {
	Validate() error
}

type ExpectingDatabaseDriver interface {
	InjectDatabaseDriver(driver schema.DatabaseDriver)
}
