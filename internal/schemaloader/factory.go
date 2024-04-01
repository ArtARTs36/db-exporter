package schemaloader

import "fmt"

func CreateLoader(driverName string, conn *Connection) (Loader, error) {
	if driverName == "pg" {
		return &PGLoader{
			conn: conn,
		}, nil
	}

	return nil, fmt.Errorf("driver %q unsupported", driverName)
}
