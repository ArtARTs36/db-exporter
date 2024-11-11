package sql

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/gds"
)

type DDL struct {
	Name        gds.String
	UpQueries   []string
	DownQueries []string
}

type DDLBuilder interface {
	Build(schema *schema.Schema, opts BuildDDLOpts) (*DDL, error)
	BuildPerTable(schema *schema.Schema, opts BuildDDLOpts) ([]*DDL, error)
}

type DDLBuilderManager struct {
	builders map[config.DatabaseDriver]DDLBuilder
}

func NewDDLBuilderManager() *DDLBuilderManager {
	return &DDLBuilderManager{
		builders: map[config.DatabaseDriver]DDLBuilder{ //nolint:exhaustive // no all drivers unsupported ddl
			config.DatabaseDriverPostgres: NewPostgresDDLBuilder(),
			config.DatabaseDriverMySQL:    NewMySQLDDLBuilder(),
		},
	}
}

func (m *DDLBuilderManager) For(driver config.DatabaseDriver) DDLBuilder {
	builder, ok := m.builders[driver]
	if ok {
		return builder
	}
	return NewPostgresDDLBuilder()
}

func (d *DDL) filled() bool {
	return len(d.UpQueries) > 0
}
