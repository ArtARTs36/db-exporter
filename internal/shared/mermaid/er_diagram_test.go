package mermaid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErDiagram_Build(t *testing.T) {
	diagram := NewErDiagram()

	diagram.AddEntity(&Entity{
		Name: "Customer",
		Fields: []*EntityField{
			{
				Name: "id",
				Type: "string",
			},
			{
				Name: "name",
				Type: "string",
			},
		},
	})

	diagram.AddEntity(&Entity{
		Name: "Order",
		Fields: []*EntityField{
			{
				Name: "id",
				Type: "string",
			},
			{
				Name: "date",
				Type: "date",
			},
			{
				Name: "customer_id",
				Type: "string",
			},
		},
	})

	diagram.AddRelation(&Relation{
		Owner:   "Order",
		Related: "Customer",
		Action:  "includes",
	})

	expected := `erDiagram
  Order ||--|{ Customer : includes
  Customer {
    string id
    string name
  }
  Order {
    string id
    date date
    string customer_id
  }
`

	assert.Equal(t, expected, diagram.Build())
}
