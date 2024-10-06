package functest

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/db"
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

	env := initPgTestEnvironment()

	cases := []struct {
		Title       string
		InitQueries []string
		DownQueries []string

		ConfigPath string
		TaskName   string
	}{
		{
			Title: "test pg with csv",
			InitQueries: []string{
				`CREATE TABLE users
(
    id   integer NOT NULL,
    name character varying NOT NULL,

    CONSTRAINT users_pk PRIMARY KEY (id)
);`,
				`INSERT INTO users (id, name) VALUES
				(1, 'a'),
				(2, 'b')
				`,
			},
			DownQueries: []string{
				"DROP TABLE users",
			},
			ConfigPath: "config.yml",
			TaskName:   "pg_csv_export",
		},
		{
			Title: "test pg with diagram",
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
		{
			Title: "test pg with go-structs",
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
			TaskName:   "pg_go-structs",
		},
		{
			Title: "test pg with yaml-fixtures",
			InitQueries: []string{
				`CREATE TABLE users
(
    id   integer NOT NULL,
    name character varying NOT NULL,

    CONSTRAINT users_pk PRIMARY KEY (id)
);`,
				`INSERT INTO users (id, name) VALUES
(1, 'a'),
(2, 'b')
`,
			},
			DownQueries: []string{
				"DROP TABLE users",
			},
			ConfigPath: "config.yml",
			TaskName:   "pg_yaml-fixtures_export",
		},
		{
			Title: "test pg with grpc-crud",
			InitQueries: []string{
				`CREATE TABLE users
(
    id   integer NOT NULL,
    name character varying NOT NULL,

    CONSTRAINT users_pk PRIMARY KEY (id)
);`,
				`INSERT INTO users (id, name) VALUES
(1, 'a'),
(2, 'b')
`,
			},
			DownQueries: []string{
				"DROP TABLE users",
			},
			ConfigPath: "config.yml",
			TaskName:   "pg_grpc-crud",
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Title, func(t *testing.T) {
			expectedFiles := loadExpectedFiles(tCase.TaskName)

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

func TestPGImport(t *testing.T) {
	skipIfRunningShortTests(t)

	env := initPgTestEnvironment()

	cases := []struct {
		Title       string
		InitQueries []string
		DownQueries []string

		ConfigPath   string
		TaskName     string
		ExpectedRows map[string][]map[string]interface{}
	}{
		{
			Title: "test pg import with yaml-fixtures",
			InitQueries: []string{
				`CREATE TABLE yaml_fixtures_import_users
(
    id   integer NOT NULL,
    name character varying NOT NULL
);`,
			},
			DownQueries: []string{
				"DROP TABLE yaml_fixtures_import_users",
			},
			ConfigPath: "config.yml",
			TaskName:   "pg_yaml-fixtures_import",
			ExpectedRows: map[string][]map[string]interface{}{
				"yaml_fixtures_import_users": {
					{"id": int64(1), "name": "a"},
					{"id": int64(2), "name": "b"},
				},
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Title, func(t *testing.T) {
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

			dl := db.NewDataLoader()

			got := map[string][]map[string]interface{}{}

			conn, err := db.NewOpenedConnection(env.db)
			require.NoError(t, err)

			for table := range tCase.ExpectedRows {
				gotTableRows, err := dl.Load(context.Background(), conn, table)
				if err != nil {
					panic(err)
				}

				got[table] = gotTableRows
			}

			assert.Equal(t, tCase.ExpectedRows, got)
		})
	}
}
