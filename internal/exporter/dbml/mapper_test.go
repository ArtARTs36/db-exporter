package dbml

import (
	"github.com/artarts36/gds"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/dbml"
)

func TestMapper_mapDefault(t *testing.T) {
	tests := []struct {
		Title    string
		Column   *schema.Column
		Expected dbml.ColumnDefault
	}{
		{
			Title: "nil",
			Column: &schema.Column{
				Default: nil,
			},
			Expected: dbml.ColumnDefault{},
		},
		{
			Title: "bool: true",
			Column: &schema.Column{
				Default: &schema.ColumnDefault{
					Type:  schema.ColumnDefaultTypeValue,
					Value: true,
				},
			},
			Expected: dbml.ColumnDefault{
				Type:  dbml.ColumnDefaultTypeBoolean,
				Value: "true",
			},
		},
		{
			Title: "bool: false",
			Column: &schema.Column{
				Default: &schema.ColumnDefault{
					Type:  schema.ColumnDefaultTypeValue,
					Value: false,
				},
			},
			Expected: dbml.ColumnDefault{
				Type:  dbml.ColumnDefaultTypeBoolean,
				Value: "false",
			},
		},
		{
			Title: "string",
			Column: &schema.Column{
				Default: &schema.ColumnDefault{
					Type:  schema.ColumnDefaultTypeValue,
					Value: "test-string",
				},
			},
			Expected: dbml.ColumnDefault{
				Type:  dbml.ColumnDefaultTypeString,
				Value: "test-string",
			},
		},
		{
			Title: "int > number",
			Column: &schema.Column{
				Default: &schema.ColumnDefault{
					Type:  schema.ColumnDefaultTypeValue,
					Value: 12,
				},
			},
			Expected: dbml.ColumnDefault{
				Type:  dbml.ColumnDefaultTypeNumber,
				Value: "12",
			},
		},
		{
			Title: "func",
			Column: &schema.Column{
				Default: &schema.ColumnDefault{
					Type:  schema.ColumnDefaultTypeFunc,
					Value: "NOW",
				},
			},
			Expected: dbml.ColumnDefault{
				Type:  dbml.ColumnDefaultTypeExpression,
				Value: "NOW",
			},
		},
	}

	mper := &mapper{}

	for _, test := range tests {
		got, err := mper.mapDefault(test.Column)
		require.NoError(t, err)
		assert.Equal(t, test.Expected, got)
	}
}

func TestMapper_mapEnums(t *testing.T) {
	tests := []struct {
		Title    string
		Schema   *schema.Schema
		Expected []*dbml.Enum
	}{
		{
			Title:    "empty",
			Schema:   &schema.Schema{},
			Expected: []*dbml.Enum{},
		},
		{
			Title: "filled",
			Schema: &schema.Schema{
				Enums: map[string]*schema.Enum{
					"my_enum": {
						Name: gds.NewString("my_enum"),
						Values: []string{
							"ok", "good",
						},
					},
				},
			},
			Expected: []*dbml.Enum{
				{
					Name: "my_enum",
					Values: []dbml.EnumValue{
						{
							Name: "ok",
						},
						{
							Name: "good",
						},
					},
				},
			},
		},
	}

	mper := &mapper{}

	for _, test := range tests {
		t.Run(test.Title, func(t *testing.T) {
			got := mper.mapEnums(test.Schema)

			assert.Equal(t, test.Expected, got)
		})
	}
}
