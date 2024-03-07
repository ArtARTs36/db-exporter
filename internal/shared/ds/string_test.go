package ds_test

import (
	"fmt"
	"testing"

	"github.com/artarts36/db-exporter/internal/shared/ds"
	"github.com/stretchr/testify/assert"
)

func TestSplitWords(t *testing.T) {
	cases := []struct {
		String        string
		ExpectedWords []string
	}{
		{
			String:        "",
			ExpectedWords: []string{},
		},
		{
			String:        "a",
			ExpectedWords: []string{"a"},
		},
		{
			String:        "a_",
			ExpectedWords: []string{"a"},
		},
		{
			String:        "a-",
			ExpectedWords: []string{"a"},
		},
		{
			String:        "a ",
			ExpectedWords: []string{"a"},
		},
		{
			String:        "ab_cdf",
			ExpectedWords: []string{"ab", "cdf"},
		},
		{
			String:        "AbCdf",
			ExpectedWords: []string{"Ab", "Cdf"},
		},
	}

	for i, tCase := range cases {
		t.Run(fmt.Sprintf("%d: %s", i, tCase.String), func(t *testing.T) {
			str := ds.NewString(tCase.String)

			assert.Equal(t, tCase.ExpectedWords, str.SplitWords())
		})
	}
}
