package functest

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func skipIfRunningShortTests(t *testing.T) {
	if os.Getenv("FUNCTEST") != "on" {
		t.Skip()
	}
}

func loadExpectedFiles(testName string, i int) map[string]string {
	dir := fmt.Sprintf("expected_files/%s/%d", testName, i)

	return loadFiles(dir)
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
