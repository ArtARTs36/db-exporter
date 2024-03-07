package migrations

import (
	"slices"

	"github.com/artarts36/db-exporter/internal/shared/golangmigrate"
	"github.com/artarts36/db-exporter/internal/shared/goose"
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
	return &TableDetector{
		migrationsTables: map[string][]*Table{
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
					Name:         golangmigrate.Table,
					ColumnsNames: golangmigrate.TableColumns,
				},
			},
		},
	}
}

func (d *TableDetector) IsMigrationsTable(tableName string, columnsNames []string) bool {
	tables, exists := d.migrationsTables[tableName]
	if !exists {
		return false
	}

	for _, table := range tables {
		if len(table.ColumnsNames) != len(columnsNames) {
			return false
		}

		slices.Sort(table.ColumnsNames)
		slices.Sort(columnsNames)

		if slices.Equal(table.ColumnsNames, columnsNames) {
			return true
		}
	}

	return false
}
