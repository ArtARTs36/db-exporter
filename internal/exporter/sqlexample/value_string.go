package sqlexample

import (
	"github.com/artarts36/db-exporter/internal/schema"
	"github.com/google/uuid"
)

type uuidValueResolver struct{}

type stringValueResolver struct{}

func (v *stringValueResolver) supports(column *schema.Column) bool {
	return column.Type.IsStringable
}

func (v *stringValueResolver) resolve(column *schema.Column) (interface{}, error) {
	if column.Name.Equal("first_name") {
		return "John", nil
	} else if column.Name.Equal("middle_name") {
		return "Bob", nil
	} else if column.Name.Equal("last_name") {
		return "Jane", nil
	}

	return "test-str", nil
}

func (v *uuidValueResolver) supports(column *schema.Column) bool {
	return column.Type.IsUUID
}

func (v *uuidValueResolver) resolve(column *schema.Column) (interface{}, error) {
	return uuid.NewString(), nil
}
