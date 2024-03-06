package migrations

import (
	"github.com/artarts36/db-exporter/internal/shared/laravel"
	"slices"

	"github.com/artarts36/db-exporter/internal/shared/goose"
)

type TableDetector struct {
	migrationsTables map[string]*Table
}

type Table struct {
	Name         string
	ColumnsNames []string
}

func NewTableDetector() *TableDetector {
	return &TableDetector{
		migrationsTables: map[string]*Table{
			goose.MigrationsTable: {
				Name:         goose.MigrationsTable,
				ColumnsNames: goose.MigrationsTableColumns,
			},
			laravel.MigrationsTable: {
				Name:         laravel.MigrationsTable,
				ColumnsNames: laravel.MigrationsTableColumns,
			},
		},
	}
}

func (d *TableDetector) IsMigrationsTable(tableName string, columnsNames []string) bool {
	table, exists := d.migrationsTables[tableName]
	if !exists {
		return false
	}

	if len(table.ColumnsNames) != len(columnsNames) {
		return false
	}

	slices.Sort(table.ColumnsNames)
	slices.Sort(columnsNames)

	return slices.Equal(table.ColumnsNames, columnsNames)
}
