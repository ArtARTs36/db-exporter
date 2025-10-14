package conn

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
	"log/slog"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
)

type Connection struct {
	db                 *sqlx.DB
	cfg                config.Database
	transactionManager *manager.Manager
}

type Transactioner func(context.Context, func(ctx context.Context) error) error

func NewOpenedConnection(db *sqlx.DB) (*Connection, error) {
	transactionManager, err := manager.New(trmsqlx.NewDefaultFactory(db))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize transaction manager: %w", err)
	}

	return &Connection{
		transactionManager: transactionManager,
		db:                 db,
	}, nil
}

func NewConnection(cfg config.Database) *Connection {
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

		c.transactionManager, err = manager.New(trmsqlx.NewDefaultFactory(c.db))
		if err != nil {
			return nil, fmt.Errorf("failed to initialize transaction manager: %w", err)
		}
	}

	return c.db, nil
}

func (c *Connection) Transact(ctx context.Context, fn func(context.Context) error) error {
	return c.transactionManager.Do(ctx, fn)
}

func (c *Connection) Close() error {
	if c.db == nil {
		return nil
	}

	return c.db.Close()
}

func (c *Connection) ExtContext(ctx context.Context) (sqlx.ExtContext, error) {
	if _, err := c.Connect(ctx); err != nil {
		return nil, err
	}

	return trmsqlx.DefaultCtxGetter.DefaultTrOrDB(ctx, c.db), nil
}

func (c *Connection) Database() config.Database {
	return c.cfg
}
