package migrations

import (
	"slices"

	"github.com/artarts36/db-exporter/internal/shared/goose"
	"github.com/artarts36/db-exporter/internal/shared/gosqlmigrate"
	"github.com/artarts36/db-exporter/internal/shared/laravel"
)

type TableDetector struct {
	migrationsTables map[string][]*Table
}

type Table struct {
	Name         string
	ColumnsNames []string
}

func NewTableDetector() *TableDetector {
	return newTableDetector(map[string][]*Table{
		goose.MigrationsTable: {
			{
				Name:         goose.MigrationsTable,
				ColumnsNames: goose.MigrationsTableColumns,
			},
		},
		laravel.MigrationsTable: {
			{
				Name:         laravel.MigrationsTable,
				ColumnsNames: laravel.MigrationsTableColumns,
			},
			{
				Name:         gosqlmigrate.Table,
				ColumnsNames: gosqlmigrate.TableColumns,
			},
		},
	})
}

func newTableDetector(migrationsTables map[string][]*Table) *TableDetector {
	for _, tables := range migrationsTables {
		for _, table := range tables {
			slices.Sort(table.ColumnsNames)
		}
	}

	return &TableDetector{migrationsTables: migrationsTables}
}

func (d *TableDetector) IsMigrationsTable(tableName string, columnsNames []string) bool {
	tables, exists := d.migrationsTables[tableName]
	if !exists {
		return false
	}

	for _, table := range tables {
		if len(table.ColumnsNames) != len(columnsNames) {
			continue
		}

		slices.Sort(columnsNames)

		if slices.Equal(table.ColumnsNames, columnsNames) {
			return true
		}
	}

	return false
}
