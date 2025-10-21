package dbml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_Render(t *testing.T) {
	file := &File{
		Tables: []*Table{
			{
				Name: "users",
				Columns: []*Column{
					{
						Name: "id",
						Type: "integer",
						Settings: ColumnSettings{
							PrimaryKey: true,
						},
					},
					{
						Name: "name",
						Type: "varchar",
					},
				},
			},
			{
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
						Name: "body",
						Type: "text",
						Settings: ColumnSettings{
							Note: "Content of the post",
						},
					},
				},
			},
		},
		Refs: []*Ref{
			{
				From: "posts.user_id",
				Type: ">",
				To:   "users.id",
			},
		},
		Enums: []*Enum{
			{
				Name: "UserStatus",
				Values: []EnumValue{
					{
						Name: "new",
					},
				},
			},
		},
	}

	expected := `Table users {
  id integer [primary key, not null]
  name varchar [not null]
}

Table posts {
  id integer [primary key, not null]
  body text [not null, note: 'Content of the post']
}

Ref: posts.user_id > users.id

Enum UserStatus {
  "new"
}`

	assert.Equal(t, expected, file.Render())
}
