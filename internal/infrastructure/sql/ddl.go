package sql

import (
	"github.com/artarts36/db-exporter/internal/config"
	"github.com/artarts36/db-exporter/internal/schema"
)

type DDLBuilder interface {
	BuildDDL(table *schema.Table, params BuildDDLParams) ([]string, error)
	CreateSequence(seq *schema.Sequence, params CreateSequenceParams) (string, error)
	DropTable(table *schema.Table, useIfExists bool) string
	CreateEnum(enum *schema.Enum) string
	DropType(name string, ifExists bool) string
	DropSequence(seq *schema.Sequence, ifExists bool) string
	CommentOnColumn(col *schema.Column) string
}

type DDLBuilderManager struct {
	builders map[config.DatabaseDriver]DDLBuilder
}

func NewDDLBuilderManager() *DDLBuilderManager {
	return &DDLBuilderManager{
		builders: map[config.DatabaseDriver]DDLBuilder{
			config.DatabaseDriverPostgres: NewPostgresDDLBuilder(),
		},
	}
}

func (m *DDLBuilderManager) For(_ config.DatabaseDriver) DDLBuilder {
	return NewPostgresDDLBuilder()
}
