package schemaloader

import "fmt"

func CreateLoader(driverName string) (Loader, error) {
	if driverName == "pg" {
		return &PGLoader{}, nil
	}

	return nil, fmt.Errorf("driver %q unsupported", driverName)
}
