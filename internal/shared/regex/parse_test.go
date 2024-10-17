package regex

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestParseSingleValue(t *testing.T) {
	cases := []struct {
		Regex    *regexp.Regexp
		Str      string
		Expected string
	}{
		{
			Regex:    regexp.MustCompile("11(.*)22"),
			Str:      "113322",
			Expected: "33",
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Str, func(t *testing.T) {
			got := ParseSingleValue(tCase.Regex, tCase.Str)

			assert.Equal(t, tCase.Expected, got)
		})
	}
}
