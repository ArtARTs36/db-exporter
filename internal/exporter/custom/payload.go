package custom

import "github.com/artarts36/db-exporter/internal/schema"

type exportingSchema struct {
	Tables []*schema.Table
}
