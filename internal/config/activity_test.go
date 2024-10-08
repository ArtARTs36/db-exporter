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
			Title: "Parse go-entities spec",
			Content: `
export: go-entities
database: db1
spec:
    package: model
`,
			Expected: Activity{
				Export: ExportActivity{
					Format: "go-entities",
					Spec: &GoStructsExportSpec{
						Package: "model",
					},
				},
				Database: "db1",
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
