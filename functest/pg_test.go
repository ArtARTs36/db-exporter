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

	cases := []struct {
		Title       string
		InitQueries []string
		DownQueries []string

		ConfigPath          string
		TaskName            string
		CheckOnlyFileExists bool
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
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/csv_export",
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
			ConfigPath:          "pg_test.yml",
			TaskName:            "pg/diagram",
			CheckOnlyFileExists: true,
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
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/go-entities",
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
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/go-entity-repository",
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
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/go-entity-repository_interfaces_external",
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
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/go-entity-repository_interfaces_internal",
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
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/laravel-models_export",
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
		    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		
		    CONSTRAINT users_pk PRIMARY KEY (id)
		);`,
				`INSERT INTO users (id, name, current_mood) VALUES
		(1, 'a', 'ok'),
		(2, 'b', 'ok')
		`,
			},
			DownQueries: []string{
				"DROP TABLE users",
				"DROP TYPE mood",
			},
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/grpc-crud",
		},
		{
			Title: "test pg with grpc-crud with all options",
			InitQueries: []string{
				`CREATE TYPE mood AS ENUM ('ok', 'happy');`,
				`CREATE TABLE users
		(
		    id   integer NOT NULL,
		    name character varying NOT NULL,
		    current_mood mood NOT NULL,
		    pay_period interval,
		    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		
		    CONSTRAINT users_pk PRIMARY KEY (id)
		);`,
				`INSERT INTO users (id, name, current_mood) VALUES
		(1, 'a', 'ok'),
		(2, 'b', 'ok')
		`,
			},
			DownQueries: []string{
				"DROP TABLE users",
				"DROP TYPE mood",
			},
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/grpc-crud-with-all",
		},
		{
			Title: "test pg with custom",
			InitQueries: []string{
				`CREATE TABLE users
		(
		    id   integer NOT NULL PRIMARY KEY,
		    name character varying NOT NULL
		);`,
				`CREATE TABLE phones
		(
		    user_id   integer NOT NULL,
		    number character varying NOT NULL
		);`,
				`ALTER TABLE phones ADD CONSTRAINT phone_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id);`,
			},
			DownQueries: []string{
				"DROP TABLE phones",
				"DROP TABLE users",
			},
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/custom",
		},
		{
			Title: "test pg with ddl",
			InitQueries: []string{
				`CREATE TABLE users
		(
		    id   integer NOT NULL,
		    name character varying NOT NULL
		);`,
			},
			DownQueries: []string{
				"DROP TABLE users",
			},
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/ddl",
		},
		{
			Title: "test pg with dbml (all in one file)",
			InitQueries: []string{
				`CREATE TYPE mood AS ENUM ('ok', 'happy');`,
				`CREATE TABLE users
		(
		    id   integer NOT NULL PRIMARY KEY,
		    name character varying NOT NULL,
			current_mood mood
		);`,
				`CREATE TABLE phones
		(
		    user_id   integer NOT NULL,
		    number character varying NOT NULL
		);`,
				`ALTER TABLE phones ADD CONSTRAINT phone_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id);`,
			},
			DownQueries: []string{
				"DROP TABLE phones",
				"DROP TABLE users",
				"DROP TYPE mood",
			},
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/dbml/all",
		},
		{
			Title: "test pg with dbml (per table",
			InitQueries: []string{
				`CREATE TYPE mood AS ENUM ('ok', 'happy');`,
				`CREATE TABLE users
		(
		    id   integer NOT NULL PRIMARY KEY,
		    name character varying NOT NULL,
			current_mood mood
		);`,
				`CREATE TABLE phones
		(
		    user_id   integer NOT NULL,
		    number character varying NOT NULL
		);`,
				`ALTER TABLE phones ADD CONSTRAINT phone_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id);`,
			},
			DownQueries: []string{
				"DROP TABLE phones",
				"DROP TABLE users",
				"DROP TYPE mood",
			},
			ConfigPath: "pg_test.yml",
			TaskName:   "pg/dbml/per-table",
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Title, func(t *testing.T) {
			mustExecQueries(env.db, tCase.InitQueries)

			t.Cleanup(func() {
				mustExecQueries(env.db, tCase.DownQueries)
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
