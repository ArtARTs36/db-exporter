package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType_Build(t *testing.T) {
	t.Run("build type without fields", func(t *testing.T) {
		typ := NewType("Character")

		assert.Equal(t, `type Character {
}`, typ.Build())
	})

	t.Run("build type with fields", func(t *testing.T) {
		typ := NewType("Character")
		typ.AddField("name").Of(TypeString).Require()
		typ.AddField("appearsIn").ListOf(TypeOfName("Episode")).Require()

		assert.Equal(t, `type Character {
  name: String!
  appearsIn: [Episode!]!
}`, typ.Build())
	})
}
