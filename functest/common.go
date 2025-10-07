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

func skipIfEnvNotFound(t *testing.T, envName string) {
	if os.Getenv(envName) == "" {
		t.Skip()
	}
}

func loadExpectedFiles(taskName string) map[string]string {
	dir := fmt.Sprintf("testdata/%s", taskName)

	files := loadFiles(dir, "")
	if len(files) == 0 {
		panic(fmt.Sprintf("expected files for test with task name %q not found", taskName))
	}

	return files
}

func loadFiles(dir, keyPrefix string) map[string]string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(fmt.Sprintf("failed to read directory %q: %s", dir, err))
	}

	files := map[string]string{}

	for _, entry := range entries {
		path := fmt.Sprintf("%s/%s", dir, entry.Name())

		if entry.IsDir() {
			kp := entry.Name()
			if keyPrefix != "" {
				kp = keyPrefix + "/" + entry.Name()
			}

			for k, v := range loadFiles(path, kp) {
				files[k] = v
			}
		} else {
			content, fileErr := os.ReadFile(path)

			if fileErr != nil {
				panic(fmt.Sprintf(
					"failed to load %s: %s",
					entry.Name(),
					fileErr,
				))
			}

			kp := ""
			if keyPrefix != "" {
				kp = keyPrefix + "/" + entry.Name()
			}

			files[kp] = string(content)
		}
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
			panic(fmt.Sprintf("failed to execute query %q: %s", query, err))
		}
	}
}
