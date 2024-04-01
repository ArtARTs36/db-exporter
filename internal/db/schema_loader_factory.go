package db

func CreateSchemaLoader(conn *Connection) (SchemaLoader, error) {
	return NewPGLoader(conn), nil
}
