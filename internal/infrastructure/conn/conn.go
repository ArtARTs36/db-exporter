package conn

import (
	"context"
	"github.com/artarts36/db-exporter/internal/schema"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Connection struct {
	db  *sqlx.DB
	cfg schema.Database
}

func NewOpenedConnection(db *sqlx.DB) (*Connection, error) {
	return &Connection{
		db: db,
	}, nil
}

func NewConnection(cfg schema.Database) *Connection {
	return &Connection{cfg: cfg}
}

func (c *Connection) Connect(ctx context.Context) (*sqlx.DB, error) {
	if c.db == nil {
		slog.DebugContext(ctx, "[db-connection] connecting to database")

		db, err := sqlx.Connect(string(c.cfg.Driver), c.cfg.DSN.Value)
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

func (c *Connection) Database() schema.Database {
	return c.cfg
}
