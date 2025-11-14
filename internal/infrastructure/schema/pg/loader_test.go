package pg

import (
	"database/sql"
	"github.com/artarts36/db-exporter/internal/infrastructure/sqltype"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/artarts36/db-exporter/internal/schema"
)

func TestPGLoader_parseColumnDefault(t *testing.T) {
	cases := []struct {
		Title    string
		Column   *schema.Column
		Expected *schema.ColumnDefault
	}{
		{
			Title: "null",
			Column: &schema.Column{
				DefaultRaw: sql.NullString{
					Valid: false,
				},
			},
			Expected: nil,
		},
		{
			Title: "parse integer value",
			Column: &schema.Column{
				DataType: sqltype.PGInteger,
				DefaultRaw: sql.NullString{
					Valid:  true,
					String: "123",
				},
			},
			Expected: &schema.ColumnDefault{
				Type:  schema.ColumnDefaultTypeValue,
				Value: 123,
			},
		},
		{
			Title: "false value",
			Column: &schema.Column{
				DataType: sqltype.PGBoolean,
				DefaultRaw: sql.NullString{
					Valid:  true,
					String: "false",
				},
			},
			Expected: &schema.ColumnDefault{
				Type:  schema.ColumnDefaultTypeValue,
				Value: false,
			},
		},
		{
			Title: "true value",
			Column: &schema.Column{
				DataType: sqltype.PGBoolean,
				DefaultRaw: sql.NullString{
					Valid:  true,
					String: "true",
				},
			},
			Expected: &schema.ColumnDefault{
				Type:  schema.ColumnDefaultTypeValue,
				Value: true,
			},
		},
		{
			Title: "parse string value",
			Column: &schema.Column{
				DataType: sqltype.PGText,
				DefaultRaw: sql.NullString{
					Valid:  true,
					String: "'str'::character varying",
				},
			},
			Expected: &schema.ColumnDefault{
				Type:  schema.ColumnDefaultTypeValue,
				Value: "str",
			},
		},
		{
			Title: "parse func value",
			Column: &schema.Column{
				DefaultRaw: sql.NullString{
					Valid:  true,
					String: "now()",
				},
			},
			Expected: &schema.ColumnDefault{
				Type:  schema.ColumnDefaultTypeFunc,
				Value: "now",
			},
		},
		{
			Title: "parse autoincrement",
			Column: &schema.Column{
				DefaultRaw: sql.NullString{
					Valid:  true,
					String: "nextval('users_id'::regclass)",
				},
			},
			Expected: &schema.ColumnDefault{
				Type:  schema.ColumnDefaultTypeAutoincrement,
				Value: "users_id",
			},
		},
	}

	loader := &Loader{}

	for _, tCase := range cases {
		t.Run(tCase.Title, func(t *testing.T) {
			got := loader.parseColumnDefault(tCase.Column)

			assert.Equal(t, tCase.Expected, got)
		})
	}
}
