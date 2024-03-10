package goose

import (
	"fmt"
	"time"
)

const MigrationsTable = "goose_db_version"
const migrationTimeFormat = "20060102150405"

var MigrationsTableColumns = []string{
	"id",
	"version_id",
	"is_applied",
	"tstamp",
}

func CreateMigrationFilename(migrationName string, offset int) string {
	// 20240229220526_create_cars_table.sql
	return fmt.Sprintf(
		"%s_%s.sql",
		time.Now().Add(time.Duration(offset)*time.Second).Format(migrationTimeFormat),
		migrationName,
	)
}
