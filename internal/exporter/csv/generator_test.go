package csv

import (
	"github.com/artarts36/db-exporter/internal/shared/iox"
	"testing"

	"github.com/artarts36/db-exporter/internal/infrastructure/data"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate(t *testing.T) {
	gen := &generator{}

	cases := []struct {
		Title     string
		Data      *transformingData
		Delimiter string
		Expected  string
	}{
		{
			Title:     "empty",
			Data:      &transformingData{},
			Delimiter: ",",
			Expected:  "",
		},
		{
			Title: "single",
			Data: &transformingData{
				cols: []string{"id", "name", "email_verified", "phone", "weight"},
				rows: data.TableData{
					{
						"id":             "1",
						"name":           "Artem",
						"email_verified": true,
						"phone":          123,
						"weight":         80.5,
					},
				},
			},
			Delimiter: ",",
			Expected: "id,name,email_verified,phone,weight" + "\n" +
				"\"1\",\"Artem\",true,123,80.500000",
		},
		{
			Title: "two",
			Data: &transformingData{
				cols: []string{"id", "name", "email_verified", "phone"},
				rows: data.TableData{
					{
						"id":             "1",
						"name":           "Artem",
						"email_verified": true,
						"phone":          123,
					},
					{
						"id":             "2",
						"name":           "Ivan",
						"email_verified": false,
						"phone":          456,
					},
				},
			},
			Delimiter: ",",
			Expected: "id,name,email_verified,phone" + "\n" +
				"\"1\",\"Artem\",true,123" + "\n" +
				"\"2\",\"Ivan\",false,456",
		},
	}

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			buf := iox.NewWriter()

			gen.generate(c.Data, c.Delimiter, buf)
			assert.Equal(t, c.Expected, buf.String())
		})
	}
}
