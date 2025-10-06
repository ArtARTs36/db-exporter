package functest

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type mysqlTestEnvironment struct {
	DSN        string
	BinaryPath string
	db         *sqlx.DB
}

func initMysqlTestEnvironment() *mysqlTestEnvironment {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		panic("MYSQL_DSN not found")
	}

	binaryPath := os.Getenv("DB_EXPORTER_BIN")
	if binaryPath == "" {
		panic("DB_EXPORTER_BIN not found")
	}

	conn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed connect to db: %s", err))
	}

	return &mysqlTestEnvironment{
		DSN:        dsn,
		BinaryPath: binaryPath,
		db:         conn,
	}
}

func TestMySQLExport(t *testing.T) {
	skipIfRunningShortTests(t)
	skipIfEnvNotFound(t, "MYSQL_DSN")

	env := initMysqlTestEnvironment()

	cases := []struct {
		Title       string
		InitQueries []string
		DownQueries []string

		ConfigPath string
		TaskName   string
	}{
		{
			Title: "test mysql with csv",
			InitQueries: []string{
				`CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`,
			},
			DownQueries: []string{
				"DROP TABLE users",
			},
			ConfigPath: "mysql_test.yml",
			TaskName:   "mysql/custom_export",
		},
		{
			Title: "test mysql with ddl",
			InitQueries: []string{
				`CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`,
			},
			DownQueries: []string{
				"DROP TABLE users",
			},
			ConfigPath: "mysql_test.yml",
			TaskName:   "mysql/ddl",
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Title, func(t *testing.T) {
			mustExecQueries(env.db, tCase.InitQueries)

			defer func() {
				mustExecQueries(env.db, tCase.DownQueries)
				removeDir("./out")
			}()

			res, cmdErr := cmd.NewCommand(env.BinaryPath).Run(
				context.Background(),
				fmt.Sprintf("--config=%s", tCase.ConfigPath),
				fmt.Sprintf("--tasks=%s", tCase.TaskName),
			)
			if cmdErr != nil {
				t.Fatalf("failed to exec command: %s: %s: %s", cmdErr, res.Stdout, res.Stderr)
			}
			assert.NoError(t, cmdErr)

			expectedFiles := loadExpectedFiles(tCase.TaskName)

			outFiles := loadFiles("./out", "")
			require.NotEmpty(t, outFiles)

			for expFileName, expFileContent := range expectedFiles {
				outFileContent, outFileExists := outFiles[expFileName]

				require.True(t, outFileExists, "file %q: not exists", expFileName)
				assert.Equal(t, expFileContent, outFileContent, "file %q: not equal expected content", expFileName)
			}
		})
	}
}
