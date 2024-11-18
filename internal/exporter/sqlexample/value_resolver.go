package sqlexample

import "github.com/artarts36/db-exporter/internal/schema"

type valueResolver interface {
	supports(column *schema.Column) bool
	resolve(column *schema.Column) (interface{}, error)
}
