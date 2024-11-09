package schema_test

import (
	"github.com/artarts36/gds"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/artarts36/db-exporter/internal/schema"
)

func TestSchema_SortByRelations(t *testing.T) {
	cases := []struct {
		Name               string
		Schema             *schema.Schema
		ExpectedTableOrder []string
	}{
		{
			Name: "move a after b",
			Schema: &schema.Schema{
				Tables: schema.NewTableMap(
					&schema.Table{
						Name: *gds.NewString("a"),
						ForeignKeys: map[string]*schema.ForeignKey{
							"a_b": {
								ForeignTable: *gds.NewString("b"),
							},
						},
					},
					&schema.Table{
						Name: *gds.NewString("b"),
					},
				),
			},
			ExpectedTableOrder: []string{
				"b",
				"a",
			},
		},
		{
			Name: "resolve deffered a,b",
			Schema: &schema.Schema{
				Tables: schema.NewTableMap(
					&schema.Table{
						Name: *gds.NewString("cars"),
					},
					&schema.Table{
						Name: *gds.NewString("users"),
					},
					&schema.Table{
						Name: *gds.NewString("a"),
						ForeignKeys: map[string]*schema.ForeignKey{
							"a_b": {
								ForeignTable: *gds.NewString("b"),
							},
							"a_cars": {
								ForeignTable: *gds.NewString("cars"),
							},
						},
					},
					&schema.Table{
						Name: *gds.NewString("b"),
						ForeignKeys: map[string]*schema.ForeignKey{
							"b_a": {
								ForeignTable: *gds.NewString("a"),
							},
							"b_users": {
								ForeignTable: *gds.NewString("users"),
							},
						},
					},
				),
			},
			ExpectedTableOrder: []string{
				"cars",
				"users",
				"a",
				"b",
			},
		},
	}

	for _, tCase := range cases {
		tCase.Schema.SortByRelations()

		assert.Equal(t, tCase.ExpectedTableOrder, tCase.Schema.TablesNames())
	}
}
