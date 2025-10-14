package schema

import (
	"encoding/json"
	"github.com/artarts36/gds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestDetermineRelationships(t *testing.T) {
	customersTable := &Table{
		Name: *gds.NewString("customers"),
		Columns: []*Column{
			{
				Name: *gds.NewString("id"),
			},
			{
				Name: *gds.NewString("name"),
			},
		},
	}
	ordersTable := &Table{
		Name: *gds.NewString("orders"),
		Columns: []*Column{
			{
				Name: *gds.NewString("id"),
			},
			{
				Name: *gds.NewString("customer_id"),
				ForeignKey: &ForeignKey{
					Name: *gds.NewString("orders_customer_id"),

					Table:        *gds.NewString("orders"),
					ColumnsNames: gds.NewStrings("customer_id"),

					ForeignTable:  *gds.NewString("customers"),
					ForeignColumn: *gds.NewString("id"),
				},
			},
		},
		ForeignKeys: map[string]*ForeignKey{
			"orders_customer_id": {
				Name: *gds.NewString("orders_customer_id"),

				Table:        *gds.NewString("orders"),
				ColumnsNames: gds.NewStrings("customer_id"),

				ForeignTable:  *gds.NewString("customers"),
				ForeignColumn: *gds.NewString("id"),
			},
		},
	}
	productsTable := &Table{
		Name: *gds.NewString("products"),
		Columns: []*Column{
			{
				Name: *gds.NewString("id"),
			},
			{
				Name: *gds.NewString("name"),
			},
		},
	}
	orderProductTable := &Table{
		Name: *gds.NewString("order_product"),
		Columns: []*Column{
			{
				Name: *gds.NewString("order_id"),
				ForeignKey: &ForeignKey{
					Name: *gds.NewString("order_product_order_id"),

					Table:        *gds.NewString("order_product"),
					ColumnsNames: gds.NewStrings("order_id"),

					ForeignTable:  *gds.NewString("orders"),
					ForeignColumn: *gds.NewString("id"),
				},
			},
			{
				Name: *gds.NewString("product_id"),
				ForeignKey: &ForeignKey{
					Name: *gds.NewString("order_product_product_id"),

					Table:        *gds.NewString("products"),
					ColumnsNames: gds.NewStrings("product_id"),

					ForeignTable:  *gds.NewString("products"),
					ForeignColumn: *gds.NewString("id"),
				},
			},
		},
		ForeignKeys: map[string]*ForeignKey{
			"order_product_order_id": {
				Name: *gds.NewString("order_product_order_id"),

				Table:        *gds.NewString("order_product"),
				ColumnsNames: gds.NewStrings("order_id"),

				ForeignTable:  *gds.NewString("orders"),
				ForeignColumn: *gds.NewString("id"),
			},
			"order_product_product_id": {
				Name: *gds.NewString("order_product_product_id"),

				Table:        *gds.NewString("products"),
				ColumnsNames: gds.NewStrings("product_id"),

				ForeignTable:  *gds.NewString("products"),
				ForeignColumn: *gds.NewString("id"),
			},
		},
	}

	schema := &Schema{
		Tables: NewTableMap(
			customersTable,
			ordersTable,
			productsTable,
			orderProductTable,
		),
	}

	expected := map[string][]*Relationship{
		"customers": {
			{
				Type:          RelationShipOneToMany,
				OwnerTable:    customersTable,
				OwnerColumn:   customersTable.Columns[0],
				RelatedTable:  ordersTable,
				RelatedColumn: ordersTable.Columns[1],
			},
		},
		"orders": {
			{
				Type:          RelationshipTypeOneToOne,
				OwnerTable:    ordersTable,
				OwnerColumn:   ordersTable.Columns[1],
				RelatedTable:  customersTable,
				RelatedColumn: customersTable.Columns[0],
			},
			{
				Type:          RelationshipTypeManyToMany,
				OwnerTable:    ordersTable,
				OwnerColumn:   ordersTable.Columns[0],
				RelatedTable:  productsTable,
				RelatedColumn: productsTable.Columns[0],
			},
		},
		"products": {
			{
				Type:          RelationshipTypeManyToMany,
				OwnerTable:    productsTable,
				OwnerColumn:   productsTable.Columns[0],
				RelatedTable:  ordersTable,
				RelatedColumn: ordersTable.Columns[0],
			},
		},
	}

	relationships, err := DetermineRelationships(schema)
	require.NoError(t, err)

	assertEqual(t, expected, relationships)
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	e, err := json.Marshal(expected)
	require.NoError(t, err)

	os.WriteFile("expected.json", e, 0755)

	a, err := json.Marshal(actual)
	require.NoError(t, err)

	os.WriteFile("actual.json", a, 0755)

	assert.Equal(t, string(e), string(a))
}
