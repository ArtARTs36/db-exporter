package config

type ValidatableSpec interface {
	Validate() error
}
