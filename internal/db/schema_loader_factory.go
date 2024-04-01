package db

import "fmt"

func CreateSchemaLoader(driverName string, conn *Connection) (SchemaLoader, error) {
	if driverName == "pg" {
		return &PGLoader{
			conn: conn,
		}, nil
	}

	return nil, fmt.Errorf("driver %q unsupported", driverName)
}
