package db

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Connection struct {
	db         *sqlx.DB
	driverName DriverName
	dsn        string
}

func NewConnection(driverName DriverName, dsn string) *Connection {
	return &Connection{driverName: driverName, dsn: dsn}
}

func (c *Connection) Connect(ctx context.Context) (*sqlx.DB, error) {
	if c.db == nil {
		slog.DebugContext(ctx, "[db-connection] connecting to database")

		db, err := sqlx.Connect(c.driverName.String(), c.dsn)
		if err != nil {
			return nil, err
		}

		slog.InfoContext(ctx, "[db-connection] connected to database")

		c.db = db
	}

	return c.db, nil
}

func (c *Connection) Close() error {
	if c.db == nil {
		return nil
	}

	return c.db.Close()
}
