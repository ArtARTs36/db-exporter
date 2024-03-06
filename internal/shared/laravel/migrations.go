package laravel

import (
	"fmt"
	"time"
)

const (
	MigrationsTable     = "migrations"
	migrationTimeFormat = "2006_01_02_150405"
)

var MigrationsTableColumns = []string{
	"id",
	"migration",
	"batch",
}

func CreateMigrationFilename(migrationName string, offset int) string {
	return fmt.Sprintf(
		"%s_%s.php",
		time.Now().Add(time.Duration(offset)*time.Second).Format(migrationTimeFormat),
		migrationName,
	)
}
