package pg

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDDLParser_Parse(t *testing.T) {
	tests := []struct {
		Title    string
		Query    string
		Expected *DDL
	}{
		{
			Title: "create users and orders table",
			Query: `CREATE TABLE orders
(
    id integer NOT NULL primary key,
    user_id integer NOT NULL
);

CREATE TABLE users
(
    id   integer NOT NULL PRIMARY KEY,
    name character varying NOT NULL
);

ALTER TABLE orders ADD CONSTRAINT orders_user_id FOREIGN KEY (user_id) REFERENCES users(id);
`,
			Expected: &DDL{
				Queries: []Query{
					&CreateTableQuery{
						Table: "orders",
						Columns: []*Column{
							{
								Name:         "id",
								DataType:     "integer",
								IsPrimaryKey: true,
							},
							{
								Name:     "user_id",
								DataType: "integer",
								Nullable: false,
							},
						},
					},
					&CreateTableQuery{
						Table: "users",
						Columns: []*Column{
							{
								Name:         "id",
								DataType:     "integer",
								Nullable:     false,
								IsPrimaryKey: true,
							},
							{
								Name:     "name",
								DataType: "character varying",
								Nullable: false,
							},
						},
					},
				},
			},
		},
	}

	parser := NewDDLParser()

	for _, test := range tests {
		t.Run(test.Title, func(t *testing.T) {
			got, err := parser.Parse(test.Query)
			require.NoError(t, err)

			assert.Equal(t, test.Expected, got)
		})
	}
}
