package functest

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

func skipIfRunningShortTests(t *testing.T) {
	if os.Getenv("FUNCTEST") != "on" {
		t.Skip()
	}
}

func loadExpectedFiles(taskName string) map[string]string {
	dir := fmt.Sprintf("data/%s", taskName)

	files := loadFiles(dir)
	if len(files) == 0 {
		panic(fmt.Sprintf("expected files for test with task name %q not found", taskName))
	}

	return files
}

func loadFiles(dir string) map[string]string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(fmt.Sprintf("failed to read directory %q: %s", dir, err))
	}

	files := map[string]string{}

	for _, entry := range entries {
		content, fileErr := os.ReadFile(fmt.Sprintf("%s/%s", dir, entry.Name()))

		if fileErr != nil {
			panic(fmt.Sprintf(
				"failed to load %s: %s",
				entry.Name(),
				fileErr,
			))
		}

		files[entry.Name()] = string(content)
	}

	return files
}

func removeDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Printf("failed to remove %q: %s", dir, err)
	}
}

func mustExecQueries(db *sqlx.DB, queries []string) {
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			panic(err)
		}
	}
}
