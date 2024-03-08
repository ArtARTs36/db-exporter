package schema_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/artarts36/db-exporter/internal/shared/ds"
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
						Name: *ds.NewString("a"),
						ForeignKeys: map[string]*schema.ForeignKey{
							"a_b": {
								ForeignTable: *ds.NewString("b"),
							},
						},
					},
					&schema.Table{
						Name: *ds.NewString("b"),
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
						Name: *ds.NewString("cars"),
					},
					&schema.Table{
						Name: *ds.NewString("users"),
					},
					&schema.Table{
						Name: *ds.NewString("a"),
						ForeignKeys: map[string]*schema.ForeignKey{
							"a_b": {
								ForeignTable: *ds.NewString("b"),
							},
							"a_cars": {
								ForeignTable: *ds.NewString("cars"),
							},
						},
					},
					&schema.Table{
						Name: *ds.NewString("b"),
						ForeignKeys: map[string]*schema.ForeignKey{
							"b_a": {
								ForeignTable: *ds.NewString("a"),
							},
							"b_users": {
								ForeignTable: *ds.NewString("users"),
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
