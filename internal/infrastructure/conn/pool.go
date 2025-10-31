package conn

import (
	"errors"
	"fmt"
	"github.com/artarts36/db-exporter/internal/schema"
)

type Pool struct {
	connections map[string]*Connection
}

func NewPool() *Pool {
	return &Pool{
		connections: map[string]*Connection{},
	}
}

func (p *Pool) Setup(ds map[string]schema.Database) {
	for name, d := range ds {
		p.connections[name] = NewConnection(d)
	}
}

func (p *Pool) Get(name string) (*Connection, bool) {
	conn, ok := p.connections[name]
	return conn, ok
}

func (p *Pool) All() map[string]*Connection {
	return p.connections
}

func (p *Pool) Close() error {
	errs := make([]error, 0)

	for conn, connection := range p.connections {
		err := connection.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection to db %q: %w", conn, err))
		}
	}

	return errors.Join(errs...)
}
