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
				Name:     "id",
				DataType: "string",
				KeyType:  KeyTypePK,
			},
			{
				Name:     "name",
				DataType: "string",
			},
		},
	})

	diagram.AddEntity(&Entity{
		Name: "Order",
		Fields: []*EntityField{
			{
				Name:     "id",
				DataType: "string",
				KeyType:  KeyTypePK,
			},
			{
				Name:     "date",
				DataType: "date",
			},
			{
				Name:     "customer_id",
				DataType: "string",
				KeyType:  KeyTypeFK,
			},
		},
	})

	diagram.AddRelation(&Relation{
		Owner:   "Order",
		Related: "Customer",
		Action:  "includes",
	})

	expected := `erDiagram
  Customer ||--o{ Order : includes
  Customer {
    string id PK
    string name
  }
  Order {
    string id PK
    date date
    string customer_id FK
  }
`

	assert.Equal(t, expected, diagram.Build())
}
