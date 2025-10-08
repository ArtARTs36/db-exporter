package webcolor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFix(t *testing.T) {
	cases := []struct {
		Title    string
		Input    string
		Expected string
	}{
		{
			Title:    "color not contains #",
			Input:    "green",
			Expected: "green",
		},
		{
			Title:    "color length < 4",
			Input:    "#ee",
			Expected: "#ee",
		},
		{
			Title:    "color length > 4",
			Input:    "#eeee",
			Expected: "#eeee",
		},
		{
			Title:    "Fixed",
			Input:    "#eee",
			Expected: "#eeeeee",
		},
	}

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			got := Fix(c.Input)
			assert.Equal(t, c.Expected, got)
		})
	}
}
