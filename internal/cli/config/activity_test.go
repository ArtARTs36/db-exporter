package config

import (
	goentity "github.com/artarts36/db-exporter/internal/exporter/go-entity"
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
format: go-entities
database: db1
spec:
    package: model
`,
			Expected: Activity{
				Format: "go-entities",
				Spec: &goentity.EntitySpecification{
					Package: "model",
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
