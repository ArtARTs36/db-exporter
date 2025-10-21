package dbml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnum_Render(t *testing.T) {
	tests := []struct {
		Title    string
		Enum     Enum
		Expected string
	}{
		{
			Title: "Simple Enum",
			Enum: Enum{
				Name: "Status",
				Values: []EnumValue{
					{
						Name: "running",
					},
					{
						Name: "stopped",
					},
				},
			},
			Expected: `Enum Status {
  "running"
  "stopped"
}`,
		},
		{
			Title: "Enum with note",
			Enum: Enum{
				Name: "Status",
				Values: []EnumValue{
					{
						Name: "running",
						Settings: EnumValueSettings{
							Note: "job running",
						},
					},
					{
						Name: "stopped",
						Settings: EnumValueSettings{
							Note: "job stopped",
						},
					},
				},
			},
			Expected: `Enum Status {
  "running" [note: 'job running']
  "stopped" [note: 'job stopped']
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Title, func(t *testing.T) {
			w := &strings.Builder{}

			test.Enum.Render(w)

			assert.Equal(t, test.Expected, w.String())
		})
	}
}
