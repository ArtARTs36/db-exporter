package sql_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/artarts36/db-exporter/internal/sql"
)

func TestInsertBuilder_Build(t *testing.T) {
	cases := []struct {
		Table       *schema.Table
		Rows        []map[string]interface{}
		Expected    string
		ExpectedErr error
	}{
		{
			Table: &schema.Table{
				Name: *ds.NewString("users"),
				Columns: []*schema.Column{
					{
						Name: *ds.NewString("id"),
					},
				},
			},
			ExpectedErr: fmt.Errorf("rows is empty"),
		},
		{
			Table: &schema.Table{
				Name: *ds.NewString("users"),
				Columns: []*schema.Column{
					{
						Name: *ds.NewString("id"),
					},
					{
						Name: *ds.NewString("name"),
					},
				},
			},
			Rows: []map[string]interface{}{
				{
					"id":   1,
					"name": "dev",
				},
				{
					"id": 2,
				},
			},
			Expected: `INSERT INTO users (id, name)
VALUES (
    (1, 'dev'),
    (2, null)
);`,
			ExpectedErr: fmt.Errorf("rows is empty"),
		},
	}

	builder := &sql.QueryBuilder{}

	for i, tCase := range cases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			got, err := builder.BuildInsertQuery(tCase.Table, tCase.Rows)
			if tCase.ExpectedErr == nil {
				require.NoError(t, err)
			}

			assert.Equal(t, tCase.Expected, got)
		})
	}
}
