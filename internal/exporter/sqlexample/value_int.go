package sqlexample

import (
	"math/rand"

	"github.com/artarts36/db-exporter/internal/schema"
)

type intValueResolver struct{}

func (r *intValueResolver) supports(column *schema.Column) bool {
	return column.Type.IsInteger
}

func (r *intValueResolver) resolve(column *schema.Column) (interface{}, error) {
	return rand.Int(), nil
}
