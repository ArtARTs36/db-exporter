package goose

import (
	"fmt"
	"time"
)

func CreateMigrationFilename(migrationName string) string {
	// 20240229220526_create_cars_table.sql
	return fmt.Sprintf(
		"%s_%s.sql",
		time.Now().Format("20060102150405"),
		migrationName,
	)
}
