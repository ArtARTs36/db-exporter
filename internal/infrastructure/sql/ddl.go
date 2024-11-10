package sql

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
)

type DDL struct {
	Name        string
	UpQueries   []string
	DownQueries []string
}

type DDLBuilder interface {
	Build(schema *schema.Schema, opts BuildDDLOpts) (*DDL, error)
	BuildPerTable(schema *schema.Schema, opts BuildDDLOpts) ([]*DDL, error)
	CreateSequence(seq *schema.Sequence, params CreateSequenceParams) (string, error)
	CreateEnum(enum *schema.Enum) string
	DropType(name string, ifExists bool) string
	DropSequence(seq *schema.Sequence, ifExists bool) string
}

type DDLBuilderManager struct {
	builders map[config.DatabaseDriver]DDLBuilder
}

func NewDDLBuilderManager() *DDLBuilderManager {
	return &DDLBuilderManager{
		builders: map[config.DatabaseDriver]DDLBuilder{
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
