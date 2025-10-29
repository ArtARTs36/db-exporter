package graphql

import (
	"github.com/artarts36/db-exporter/internal/shared/iox"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile_Build(t *testing.T) {
	tests := []struct {
		Title    string
		File     File
		Expected string
	}{
		{
			Title:    "empty file",
			File:     File{},
			Expected: "",
		},
		{
			Title: "one object",
			File: File{
				Types: []*Object{
					NewType("User"),
				},
			},
			Expected: `type User {
}`,
		},
		{
			Title: "two objects",
			File: File{
				Types: []*Object{
					NewType("User"),
					NewType("Phone"),
				},
			},
			Expected: `type User {
}

type Phone {
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Title, func(t *testing.T) {
			w := iox.NewWriter()
			test.File.Build(w)

			assert.Equal(t, test.Expected, w.String())
		})
	}
}
