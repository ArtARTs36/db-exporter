package sql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/sql"
)

func TestDDLBuilder_BuildDDL(t *testing.T) {
	cases := []struct {
		Name            string
		Table           *schema.Table
		ExpectedQueries []string
	}{
		{
			Name: "empty table",
			Table: &schema.Table{
				Name:    ds.String{Value: "cars"},
				Columns: []*schema.Column{},
			},
			ExpectedQueries: []string{
				"CREATE TABLE cars()",
			},
		},
		{
			Name: "table with 1 column",
			Table: &schema.Table{
				Name: ds.String{Value: "cars"},
				Columns: []*schema.Column{
					{
						Name: *ds.NewString("id"),
						Type: *ds.NewString("integer"),
					},
				},
			},
			ExpectedQueries: []string{
				`CREATE TABLE cars
(
    id integer NOT NULL
);`,
			},
		},
		{
			Name: "table with 1 column and primary key",
			Table: &schema.Table{
				Name: ds.String{Value: "cars"},
				Columns: []*schema.Column{
					{
						Name: *ds.NewString("id"),
						Type: *ds.NewString("integer"),
					},
				},
				PrimaryKey: &schema.PrimaryKey{
					Name:         *ds.NewString("cars_pk"),
					ColumnsNames: ds.NewStrings("id"),
				},
			},
			ExpectedQueries: []string{
				`CREATE TABLE cars
(
    id integer NOT NULL,

    CONSTRAINT cars_pk PRIMARY KEY (id)
);`,
			},
		},
	}

	builder := sql.NewDDLBuilder()

	for _, tCase := range cases {
		t.Run(tCase.Name, func(t *testing.T) {
			queries := builder.BuildDDL(tCase.Table)

			assert.Equal(t, tCase.ExpectedQueries, queries)
		})
	}
}
