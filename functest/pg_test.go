package functest

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/shared/cmd"
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

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed connect to db: %s", err))
	}

	return &pgTestEnvironment{
		DSN:        dsn,
		BinaryPath: binaryPath,
		db:         db,
	}
}

func TestPG(t *testing.T) {
	skipIfRunningShortTests(t)

	env := initPgTestEnvironment()

	cases := []struct {
		Name        string
		InitQueries []string
		DownQueries []string

		ConfigPath string
		TaskName   string
	}{
		{
			Name: "test pg with go-structs",
			InitQueries: []string{
				`CREATE TABLE users
(
    id   integer NOT NULL,
    name character varying NOT NULL,
    country_id integer,
    balance real NOT NULL,
    prev_balance real,
    phone character varying,
    created_at timestamp NOT NULL,
    updated_at timestamp,

    CONSTRAINT users_pk PRIMARY KEY (id)
);`,
				`CREATE TABLE countries
(
    id integer NOT NULL,
    name character varying NOT NULL,
    
    CONSTRAINT countries_pk PRIMARY KEY (id)
)`,
				`ALTER TABLE users ADD CONSTRAINT user_country_fk FOREIGN KEY (country_id) REFERENCES countries(id);`,
			},
			DownQueries: []string{
				"DROP TABLE users",
				"DROP TABLE countries",
			},
			ConfigPath: "config.yml",
			TaskName:   "pg_go_structs",
		},
		{
			Name: "test pg with diagram",
			InitQueries: []string{
				`CREATE TABLE users
(
    id   integer NOT NULL,
    name character varying NOT NULL,
    country_id integer,
    balance real NOT NULL,
    prev_balance real,
    phone character varying,
    created_at timestamp NOT NULL,
    updated_at timestamp,

    CONSTRAINT users_pk PRIMARY KEY (id)
);`,
				`CREATE TABLE countries
(
    id integer NOT NULL,
    name character varying NOT NULL,
    
    CONSTRAINT countries_pk PRIMARY KEY (id)
)`,
				`ALTER TABLE users ADD CONSTRAINT user_country_fk FOREIGN KEY (country_id) REFERENCES countries(id);`,
			},
			DownQueries: []string{
				"DROP TABLE users",
				"DROP TABLE countries",
			},
			ConfigPath: "config.yml",
			TaskName:   "pg_diagram",
		},
	}

	for i, tCase := range cases {
		t.Run(tCase.Name, func(t *testing.T) {
			expectedFiles := loadExpectedFiles("pg_test", i)

			mustExecQueries(env.db, tCase.InitQueries)

			defer func() {
				mustExecQueries(env.db, tCase.DownQueries)
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

			outFiles := loadFiles("./out")

			for expFileName, expFileContent := range expectedFiles {
				outFileContent, outFileExists := outFiles[expFileName]

				assert.True(t, outFileExists)
				assert.Equal(t, expFileContent, outFileContent)
			}

			removeDir("./out")
		})
	}
}
