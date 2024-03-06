package goose

import (
	"fmt"
	"time"
)

const MigrationsTable = "goose_db_version"

var MigrationsTableColumns = []string{
	"id",
	"version_id",
	"is_applied",
	"tstamp",
}

func IsMigrationsTable(table string) bool {
	return table == MigrationsTable
}

func CreateMigrationFilename(migrationName string) string {
	// 20240229220526_create_cars_table.sql
	return fmt.Sprintf(
		"%s_%s.sql",
		time.Now().Format("20060102150405"),
		migrationName,
	)
}
