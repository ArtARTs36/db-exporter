package config

import (
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportActivity_UnmarshalYAML(t *testing.T) {
	cases := []struct {
		Title    string
		Content  string
		Expected Activity
	}{
		{
			Title: "Parse go-structs spec",
			Content: `
export: go-structs
database: db1
spec:
    package: model
`,
			Expected: Activity{
				Export:   "go-structs",
				Database: "db1",
				Spec: &GoStructsExportSpec{
					Package: "model",
				},
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.Title, func(t *testing.T) {
			var activity Activity

			err := yaml.Unmarshal([]byte(tCase.Content), &activity)
			require.NoError(t, err)

			assert.Equal(t, tCase.Expected, activity)
		})
	}
}
