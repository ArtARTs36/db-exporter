package golangmigrate

import (
	"fmt"
	"time"
)

const Table = "migrations"
const migrationTimeFormat = "20060102150405"

var TableColumns = []string{
	"id",
	"applied_at",
}

func CreateMigrationFilename(migrationName string) string {
	// 20240229220526_create_cars_table.sql
	return fmt.Sprintf(
		"%s_%s.sql",
		time.Now().Format(migrationTimeFormat),
		migrationName,
	)
}
