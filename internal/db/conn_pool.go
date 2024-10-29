package db

import (
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/config"
)

type ConnectionPool struct {
	connections map[string]*Connection
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		connections: map[string]*Connection{},
	}
}

func (p *ConnectionPool) Setup(ds map[string]config.Database) {
	for name, db := range ds {
		p.connections[name] = NewConnection(db)
	}
}

func (p *ConnectionPool) Get(name string) (*Connection, bool) {
	conn, ok := p.connections[name]
	return conn, ok
}

func (p *ConnectionPool) Close() error {
	errs := make([]error, 0)

	for db, connection := range p.connections {
		err := connection.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection to db %q: %w", db, err))
		}
	}

	return errors.Join(errs...)
}
