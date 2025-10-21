package paginator

import "github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"

type None struct{}

func (None) AddPaginationToRequest(_ *presentation.Message) {}

func (None) AddPaginationToResponse(_ *presentation.Message) {}
