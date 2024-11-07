package dbml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRender_Build(t *testing.T) {
	cases := []struct {
		Title    string
		Table    Table
		Expected string
	}{
		{
			Title: "empty table",
			Table: Table{
				Name: "posts",
			},
			Expected: `Table posts {
}`,
		},
		{
			Title: "table with one field",
			Table: Table{
				Name: "posts",
				Columns: []*Column{
					{
						Name: "id",
						Type: "varchar",
					},
				},
			},
			Expected: `Table posts {
  id varchar
}`,
		},
		{
			Title: "table with primary key",
			Table: Table{
				Name: "posts",
				Columns: []*Column{
					{
						Name: "id",
						Type: "integer",
						Settings: ColumnSettings{
							PrimaryKey: true,
						},
					},
					{
						Name: "title",
						Type: "varchar",
					},
				},
			},
			Expected: `Table posts {
  id integer [primary key]
  title varchar
}`,
		},
		{
			Title: "table with primary key and note",
			Table: Table{
				Name: "posts",
				Columns: []*Column{
					{
						Name: "id",
						Type: "integer",
						Settings: ColumnSettings{
							PrimaryKey: true,
						},
					},
					{
						Name: "title",
						Type: "varchar",
					},
					{
						Name: "body",
						Type: "varchar",
						Settings: ColumnSettings{
							Note: "Content of the post",
						},
					},
				},
			},
			Expected: `Table posts {
  id integer [primary key]
  title varchar
  body varchar [note: 'Content of the post']
}`,
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Title, func(t *testing.T) {
			got := tCase.Table.Render()

			assert.Equal(t, tCase.Expected, got)
		})
	}
}
