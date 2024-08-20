package proto

import "github.com/artarts36/db-exporter/internal/shared/ds"

type File struct {
	Package  string
	Services []*Service
	Messages []*Message
	Imports  *ds.Set
	Options  map[string]string
}
