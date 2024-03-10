package functest //nolint: testpackage // not need

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for pg driver
	"github.com/stretchr/testify/assert"
)

type pgTestEnvironment struct {
	DSN string
	db  *sqlx.DB
}

func initPgTestEnvironment() *pgTestEnvironment {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		panic("PG_DSN not found")
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed connect to db: %s", err))
	}

	return &pgTestEnvironment{
		DSN: dsn,
		db:  db,
	}
}

func TestPG(t *testing.T) {
	skipIfRunningShortTests(t)

	env := initPgTestEnvironment()

	cases := []struct {
		Name        string
		InitQueries []string
		DownQueries []string

		BinArgs []string
	}{
		{
			Name: "test pg with go-structs",
			InitQueries: []string{
				`CREATE TABLE users
(
    id   integer NOT NULL,
    name character varying,

    CONSTRAINT users_pk PRIMARY KEY (id)
);`,
			},
			DownQueries: []string{
				"DROP TABLE users",
			},
			BinArgs: []string{
				"pg",
				env.DSN,
				"go-structs",
				"out",
			},
		},
	}

	for i, tCase := range cases {
		t.Run(tCase.Name, func(t *testing.T) {
			expectedFiles := loadExpectedFiles("pg_test", i)

			for _, query := range tCase.InitQueries {
				_, err := env.db.Exec(query)
				if err != nil {
					panic(err)
				}
			}

			cmdErr := exec.Command("./../db-exporter", tCase.BinArgs...).Run()

			assert.NoError(t, cmdErr)

			outFiles := loadFiles("out")

			for expFileName, expFileContent := range expectedFiles {
				outFileContent, outFileExists := outFiles[expFileName]

				assert.True(t, outFileExists)
				assert.Equal(t, expFileContent, outFileContent)
			}

			removeDir("out")

			for _, query := range tCase.DownQueries {
				_, err := env.db.Exec(query)
				if err != nil {
					panic(err)
				}
			}
		})
	}
}
