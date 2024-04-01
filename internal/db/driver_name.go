package db

import "fmt"

type DriverName int

const (
	DriverNameUndefined DriverName = iota
	DriverNamePostgres
)

func CreateDriverName(name string) (DriverName, error) {
	if name == "postgres" || name == "pg" {
		return DriverNamePostgres, nil
	}

	return DriverNameUndefined, fmt.Errorf("driver name %q unsupported", name)
}

func (n DriverName) String() string {
	switch n {
	case DriverNameUndefined:
		return "undefined"
	case DriverNamePostgres:
		return "postgres"
	}

	return "undefined"
}
