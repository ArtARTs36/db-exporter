package schema

import "fmt"

func CreateLoader(driverName string) (Loader, error) {
	if driverName == "pg" || driverName == "postgres" {
		return &PGLoader{}, nil
	}

	return nil, fmt.Errorf("driver %q unsupported", driverName)
}
