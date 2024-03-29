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
		{
			Name: "table with deferrable foreign keys",
			Table: &schema.Table{
				Name: ds.String{Value: "users"},
				Columns: []*schema.Column{
					{
						Name: *ds.NewString("id"),
						Type: *ds.NewString("integer"),
					},
					{
						Name: *ds.NewString("car_id"),
						Type: *ds.NewString("integer"),
					},
					{
						Name: *ds.NewString("mobile_id"),
						Type: *ds.NewString("integer"),
					},
				},
				PrimaryKey: &schema.PrimaryKey{
					Name:         *ds.NewString("users_pk"),
					ColumnsNames: ds.NewStrings("id"),
				},
				ForeignKeys: map[string]*schema.ForeignKey{
					"users_car_id_fk": {
						Name:          *ds.NewString("users_car_id_fk"),
						ColumnsNames:  ds.NewStrings("car_id"),
						ForeignTable:  *ds.NewString("cars"),
						ForeignColumn: *ds.NewString("id"),
						IsDeferrable:  true,
					},
					"users_mobile_id_fk": {
						Name:                *ds.NewString("users_mobile_id_fk"),
						ColumnsNames:        ds.NewStrings("mobile_id"),
						ForeignTable:        *ds.NewString("mobiles"),
						ForeignColumn:       *ds.NewString("id"),
						IsDeferrable:        true,
						IsInitiallyDeferred: true,
					},
				},
			},
			ExpectedQueries: []string{
				`CREATE TABLE users
(
    id        integer NOT NULL,
    car_id    integer NOT NULL,
    mobile_id integer NOT NULL,

    CONSTRAINT users_pk PRIMARY KEY (id),
    CONSTRAINT users_car_id_fk FOREIGN KEY (car_id) REFERENCES cars (id) DEFERRABLE,
    CONSTRAINT users_mobile_id_fk FOREIGN KEY (mobile_id) REFERENCES mobiles (id) DEFERRABLE INITIALLY DEFERRED
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
