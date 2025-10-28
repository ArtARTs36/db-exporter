package functest

import (
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/cmd"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for pg driver
	"github.com/stretchr/testify/assert"
)

type pgTestEnvironment struct {
	DSN        string
	BinaryPath string
	db         *sqlx.DB
}

func initPgTestEnvironment() *pgTestEnvironment {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		panic("PG_DSN not found")
	}

	binaryPath := os.Getenv("DB_EXPORTER_BIN")
	if binaryPath == "" {
		panic("DB_EXPORTER_BIN not found")
	}

	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed connect to db: %s", err))
	}

	return &pgTestEnvironment{
		DSN:        dsn,
		BinaryPath: binaryPath,
		db:         conn,
	}
}

func TestPGExport(t *testing.T) {
	skipIfRunningShortTests(t)
	skipIfEnvNotFound(t, "PG_DSN")

	env := initPgTestEnvironment()

	mustExecQueries(env.db, []string{
		`CREATE TYPE mood AS ENUM ('ok', 'happy');`,
		`CREATE TABLE users
		(
		    id   integer NOT NULL,
		    name character varying(64) NOT NULL,
		    balance real NOT NULL,
		    prev_balance real,
		    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		    current_mood mood NOT NULL,
		    updated_at timestamp,
		
		    CONSTRAINT users_pk PRIMARY KEY (id)
		);`,
		`CREATE TABLE phones
		(
		    user_id   integer NOT NULL,
		    number character varying NOT NULL,
		    
		    CONSTRAINT phones_pk PRIMARY KEY (user_id, number)
		);`,
		`ALTER TABLE phones ADD CONSTRAINT phone_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id);`,
		`INSERT INTO users (id, name, balance, prev_balance, created_at, current_mood) VALUES
		(1, 'Artem', 999999999, null, '2025-10-26 21:21:27.699806', 'ok'),
		(2, 'Ivan', 88888888, null, '2025-10-26 21:21:27.699806', 'happy');`,
	})

	t.Cleanup(func() {
		mustExecQueries(env.db, []string{
			`DROP TABLE phones;`,
			`DROP TABLE users;`,
			`DROP TYPE mood;`,
		})
	})

	cases := []struct {
		Title               string
		ConfigPath          string
		TaskName            string
		CheckOnlyFileExists bool
	}{
		{
			Title:      "test pg with csv",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/csv",
		},
		{
			Title:      "test pg with diagram",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/diagram",
		},
		{
			Title:      "test pg with go-entities",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/go-entities",
		},
		{
			Title:      "test pg with go-entity-repository",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/go-entity-repository",
		},
		{
			Title:      "test pg with go-entity-repository with external interfaces",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/go-entity-repository_interfaces_external",
		},
		{
			Title:      "test pg with go-entity-repository with internal interfaces",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/go-entity-repository_interfaces_internal",
		},
		{
			Title:      "test pg with laravel-models",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/laravel-models/per-table",
		},
		{
			Title:      "test pg with grpc-crud",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/grpc-crud/all",
		},
		{
			Title:      "test pg with grpc-crud with all options",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/grpc-crud-with-options/all",
		},
		{
			Title:      "test pg with custom (all in one file)",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/custom/all",
		},
		{
			Title:      "test pg with custom (per-table)",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/custom/per-table",
		},
		{
			Title:      "test pg with ddl (all)",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/ddl/all",
		},
		{
			Title:      "test pg with ddl (per-table)",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/ddl/per-table",
		},
		{
			Title:      "test pg with dbml (all in one file)",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/dbml/all",
		},
		{
			Title:      "test pg with dbml (per table)",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/dbml/per-table",
		},
		{
			Title:      "test pg graphql (all)",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/graphql/all",
		},
		{
			Title:      "test pg graphql (per-table)",
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/graphql/per-table",
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Title, func(t *testing.T) {
			t.Cleanup(func() {
				removeDir("./out")
			})

			res, cmdErr := cmd.NewCommand(env.BinaryPath).Run(
				t.Context(),
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

				if !tCase.CheckOnlyFileExists {
					assert.Equal(t, expFileContent, outFileContent, "file %q: not equal expected content", expFileName)
				}
			}
		})
	}
}
