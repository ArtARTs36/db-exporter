package sql

import (
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"github.com/artarts36/gds"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/artarts36/db-exporter/internal/schema"
)

func TestDDLBuilder_Build(t *testing.T) {
	cases := []struct {
		Name        string
		Schema      *schema.Schema
		Table       *schema.Table
		ExpectedDDL *DDL
		Opts        BuildDDLOpts
	}{
		{
			Name: "empty table",
			Schema: &schema.Schema{
				Tables: schema.NewTableMap(&schema.Table{
					Name:    gds.String{Value: "cars"},
					Columns: []*schema.Column{},
				}),
				Driver: schema.DatabaseDriverPostgres,
			},
			ExpectedDDL: &DDL{
				Name:        *gds.NewString("init"),
				UpQueries:   []string{"CREATE TABLE cars()"},
				DownQueries: []string{"DROP TABLE cars;"},
			},
		},
		{
			Name: "table with 1 column",
			Schema: &schema.Schema{
				Tables: schema.NewTableMap(&schema.Table{
					Name: gds.String{Value: "cars"},
					Columns: []*schema.Column{
						{
							Name: *gds.NewString("id"),
							Type: sqltype.PGInteger,
						},
					},
				}),
				Driver: schema.DatabaseDriverPostgres,
			},
			ExpectedDDL: &DDL{
				Name: *gds.NewString("init"),
				UpQueries: []string{
					`CREATE TABLE cars
(
    id integer NOT NULL
);`,
				},
				DownQueries: []string{"DROP TABLE cars;"},
			},
		},
		{
			Name: "table with 1 column and primary key",
			Schema: &schema.Schema{
				Tables: schema.NewTableMap(&schema.Table{
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
				}),
				Driver: schema.DatabaseDriverPostgres,
			},
			ExpectedDDL: &DDL{
				Name: *gds.NewString("init"),
				UpQueries: []string{
					`CREATE TABLE cars
(
    id integer NOT NULL,

    CONSTRAINT cars_pk PRIMARY KEY (id)
);`,
				},
				DownQueries: []string{"DROP TABLE cars;"},
			},
		},
		{
			Name: "table with deferrable foreign keys",
			Schema: &schema.Schema{
				Tables: schema.NewTableMap(&schema.Table{
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
				}),
				Driver: schema.DatabaseDriverPostgres,
			},
			ExpectedDDL: &DDL{
				Name: *gds.NewString("init"),
				UpQueries: []string{
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
				DownQueries: []string{"DROP TABLE users;"},
			},
		},
		{
			Name: "column with enum",
			Schema: &schema.Schema{
				Tables: schema.NewTableMap(&schema.Table{
					Name: gds.String{Value: "users"},
					Columns: []*schema.Column{
						{
							Name: *gds.NewString("id"),
							Type: sqltype.PGInteger,
						},
						{
							Name: *gds.NewString("status"),
							Type: schema.Type{Name: "status"},
							Enum: &schema.Enum{
								Name:   gds.NewString("status"),
								Values: []string{"a", "b", "c", "d"},
								Used:   1,
							},
						},
					},
					ForeignKeys: map[string]*schema.ForeignKey{},
					UsingEnums: map[string]*schema.Enum{"status": &schema.Enum{
						Name:   gds.NewString("status"),
						Values: []string{"a", "b", "c", "d"},
						Used:   1,
					}},
				}),
				Enums: map[string]*schema.Enum{"status": &schema.Enum{
					Name:   gds.NewString("status"),
					Values: []string{"a", "b", "c", "d"},
					Used:   1,
				}},
				Driver: schema.DatabaseDriverPostgres,
			},
			ExpectedDDL: &DDL{
				Name: *gds.NewString("init"),
				UpQueries: []string{
					`CREATE TYPE status AS ENUM ('a', 'b', 'c', 'd');`,
					`CREATE TABLE users
(
    id     integer NOT NULL,
    status status NOT NULL
);`,
				},
				DownQueries: []string{"DROP TABLE users;", "DROP TYPE status;"},
			},
		},
		{
			Name: "column with sequence",
			Schema: &schema.Schema{
				Tables: schema.NewTableMap(&schema.Table{
					Name: gds.String{Value: "users"},
					Columns: []*schema.Column{
						{
							Name: *gds.NewString("id"),
							Type: sqltype.PGInteger,
						},
						{
							Name: *gds.NewString("status"),
							Type: schema.Type{Name: "status"},
							UsingSequences: map[string]*schema.Sequence{
								"users_id_seq": &schema.Sequence{
									Name: "users_id_seq",
								},
							},
						},
					},
					ForeignKeys: map[string]*schema.ForeignKey{},
					UsingEnums: map[string]*schema.Enum{"status": &schema.Enum{
						Name:   gds.NewString("status"),
						Values: []string{"a", "b", "c", "d"},
						Used:   1,
					}},
					UsingSequences: map[string]*schema.Sequence{
						"users_id_seq": &schema.Sequence{
							Name: "users_id_seq",
						},
					},
				}),
				Enums: map[string]*schema.Enum{"status": &schema.Enum{
					Name:   gds.NewString("status"),
					Values: []string{"a", "b", "c", "d"},
					Used:   1,
				}},
				Sequences: map[string]*schema.Sequence{
					"users_id_seq": &schema.Sequence{
						Name:     "users_id_seq",
						DataType: sqltype.PGInteger,
					},
				},
				Driver: schema.DatabaseDriverPostgres,
			},
			ExpectedDDL: &DDL{
				Name: *gds.NewString("init"),
				UpQueries: []string{
					`CREATE TYPE status AS ENUM ('a', 'b', 'c', 'd');`,
					`CREATE SEQUENCE users_id_seq as integer;`,
					`CREATE TABLE users
(
    id     integer NOT NULL,
    status status NOT NULL
);`,
				},
				DownQueries: []string{
					"DROP TABLE users;",
					"DROP TYPE status;",
					"DROP SEQUENCE users_id_seq;",
				},
			},
		},
	}

	builder := NewPostgresDDLBuilder()

	for _, tCase := range cases {
		t.Run(tCase.Name, func(t *testing.T) {
			queries, err := builder.Build(tCase.Schema, tCase.Opts)
			require.NoError(t, err)

			assert.Equal(t, tCase.ExpectedDDL, queries)
		})
	}
}
