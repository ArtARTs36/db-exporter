package functest

import (
	"context"
	"fmt"
	"github.com/artarts36/db-exporter/internal/infrastructure/conn"
	"github.com/artarts36/db-exporter/internal/infrastructure/data"
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
				`CREATE TABLE countries
		(
		    id   integer NOT NULL,
		    code character varying NOT NULL,
		    name character varying NOT NULL
		);`,
				`INSERT INTO users (id, name) VALUES
						(1, 'a'),
						(2, 'b')
						`,
				`INSERT INTO countries (id, code, name) VALUES
						(1, 'RU', 'Russia')
						`,
			},
			DownQueries: []string{
				"DROP TABLE users",
				"DROP TABLE countries",
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
			Title: "test pg with go-entities",
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
			TaskName:   "pg_go-entities",
		},
		{
			Title: "test pg with go-entity-repository",
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
			TaskName:   "pg_go-entity-repository",
		},
		{
			Title: "test pg with go-entity-repository with external interfaces",
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
			TaskName:   "pg_go-entity-repository_interfaces_external",
		},
		{
			Title: "test pg with go-entity-repository with internal interfaces",
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
			TaskName:   "pg_go-entity-repository_interfaces_internal",
		},
		{
			Title: "test pg with laravel-models",
			InitQueries: []string{
				`CREATE TABLE users
		(
		    id   serial NOT NULL,
		    name character varying NOT NULL,
		
		    CONSTRAINT users_pk PRIMARY KEY (id)
		);`,
				`CREATE TABLE entities
		(
		    entity_type character varying NOT NULL,
		    entity_id character varying NOT NULL,
		    name character varying NOT NULL,
		
		    CONSTRAINT entities_pk PRIMARY KEY (entity_type, entity_id)
		);`,
			},
			DownQueries: []string{
				"DROP TABLE users",
				"DROP TABLE entities",
			},
			ConfigPath: "config.yml",
			TaskName:   "pg_laravel-models_export",
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
				`CREATE TYPE mood AS ENUM ('ok', 'happy');`,
				`CREATE TABLE users
		(
		    id   integer NOT NULL,
		    name character varying NOT NULL,
		    current_mood mood NOT NULL,
		
		    CONSTRAINT users_pk PRIMARY KEY (id)
		);`,
				`INSERT INTO users (id, name, current_mood) VALUES
		(1, 'a', 'ok'),
		(2, 'b', 'ok')
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

			dl := data.NewLoader()

			got := map[string][]map[string]interface{}{}

			cn, err := conn.NewOpenedConnection(env.db)
			require.NoError(t, err)

			for table := range tCase.ExpectedRows {
				gotTableRows, lerr := dl.Load(context.Background(), cn, table)
				if lerr != nil {
					panic(lerr)
				}

				got[table] = gotTableRows
			}

			assert.Equal(t, tCase.ExpectedRows, got)
		})
	}
}
