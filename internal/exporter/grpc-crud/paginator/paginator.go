package paginator

import (
	"github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"
)

type Paginator interface {
	AddPaginationToRequest(req *presentation.Message)
	AddPaginationToResponse(resp *presentation.Message)
}
