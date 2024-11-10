package sql

import (
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/gds"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/artarts36/db-exporter/internal/schema"
)

func TestDDLBuilder_BuildDDL(t *testing.T) {
	cases := []struct {
		Name            string
		Table           *schema.Table
		ExpectedQueries []string
		Opts            BuildDDLParams
	}{
		{
			Name: "empty table",
			Table: &schema.Table{
				Name:    gds.String{Value: "cars"},
				Columns: []*schema.Column{},
			},
			ExpectedQueries: []string{
				"CREATE TABLE cars()",
			},
		},
		{
			Name: "table with 1 column",
			Table: &schema.Table{
				Name: gds.String{Value: "cars"},
				Columns: []*schema.Column{
					{
						Name: *gds.NewString("id"),
						Type: sqltype.PGInteger,
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
				Name: gds.String{Value: "cars"},
				Columns: []*schema.Column{
					{
						Name: *gds.NewString("id"),
						Type: sqltype.PGInteger,
					},
				},
				PrimaryKey: &schema.PrimaryKey{
					Name:         *gds.NewString("cars_pk"),
					ColumnsNames: gds.NewStrings("id"),
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
				Name: gds.String{Value: "users"},
				Columns: []*schema.Column{
					{
						Name: *gds.NewString("id"),
						Type: sqltype.PGInteger,
					},
					{
						Name: *gds.NewString("car_id"),
						Type: sqltype.PGInteger,
					},
					{
						Name: *gds.NewString("mobile_id"),
						Type: sqltype.PGInteger,
					},
				},
				PrimaryKey: &schema.PrimaryKey{
					Name:         *gds.NewString("users_pk"),
					ColumnsNames: gds.NewStrings("id"),
				},
				ForeignKeys: map[string]*schema.ForeignKey{
					"users_car_id_fk": {
						Name:          *gds.NewString("users_car_id_fk"),
						ColumnsNames:  gds.NewStrings("car_id"),
						ForeignTable:  *gds.NewString("cars"),
						ForeignColumn: *gds.NewString("id"),
						IsDeferrable:  true,
					},
					"users_mobile_id_fk": {
						Name:                *gds.NewString("users_mobile_id_fk"),
						ColumnsNames:        gds.NewStrings("mobile_id"),
						ForeignTable:        *gds.NewString("mobiles"),
						ForeignColumn:       *gds.NewString("id"),
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

	builder := NewPostgresDDLBuilder()

	for _, tCase := range cases {
		t.Run(tCase.Name, func(t *testing.T) {
			queries, err := builder.BuildDDL(tCase.Table, tCase.Opts)
			require.NoError(t, err)

			assert.Equal(t, tCase.ExpectedQueries, queries)
		})
	}
}
