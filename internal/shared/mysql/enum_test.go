package mysql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseEnumType(t *testing.T) {
	cases := []struct {
		Title    string
		Input    string
		Expected []string
	}{
		{
			Title:    "empty definition",
			Expected: nil,
		},
		{
			Title:    "one value",
			Input:    "enum('active')",
			Expected: []string{"active"},
		},
		{
			Title:    "two value",
			Input:    "enum('active','banned')",
			Expected: []string{"active", "banned"},
		},
		{
			Title:    "three values",
			Input:    "enum('active','banned', 'three')",
			Expected: []string{"active", "banned", "three"},
		},
	}

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			values := ParseEnumType(c.Input)
			assert.Equal(t, c.Expected, values)
		})
	}
}
