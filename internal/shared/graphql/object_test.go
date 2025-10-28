package graphql

import (
	"github.com/artarts36/db-exporter/internal/shared/iox"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType_Build(t *testing.T) {
	assertBuild := func(t *testing.T, typ *Object, expected string) {
		w := iox.NewWriter()
		typ.Build(w)
		assert.Equal(t, expected, w.String())
	}

	t.Run("build type without fields", func(t *testing.T) {
		typ := NewType("Character")

		assertBuild(t, typ, `type Character {
}`)
	})

	t.Run("build type with fields", func(t *testing.T) {
		typ := NewType("Character")
		typ.AddField("name").Of(TypeString).Require()
		typ.AddField("appearsIn").ListOf(TypeOfName("Episode")).Require()

		assertBuild(t, typ, `type Character {
  name: String!
  appearsIn: [Episode!]!
}`)
	})

	t.Run("build type with fields and comments", func(t *testing.T) {
		typ := NewType("Character")
		typ.AddField("name").Of(TypeString).Require().Comment("name of character")
		typ.AddField("appearsIn").ListOf(TypeOfName("Episode")).Require()

		assertBuild(t, typ, `type Character {
  # name of character
  name: String!
  appearsIn: [Episode!]!
}`)
	})
}
