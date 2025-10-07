package config

type ValidatableSpec interface {
	Validate() error
}

type ExpectingDatabaseDriver interface {
	InjectDatabaseDriver(driver DatabaseDriver)
}
