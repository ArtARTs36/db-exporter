package paginator

import "github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"

type Offset struct {
}

func (o *Offset) AddPaginationToRequest(req *presentation.Message) {
	req.CreateField("page_size", func(field *presentation.Field) {
		field.SetType("uint64").SetTopComment("Maximum number of results per page.")
	})

	req.CreateField("offset", func(field *presentation.Field) {
		field.SetType("uint64").SetTopComment("The number of records to skip before starting to fetch the records.")
	})
}

func (o *Offset) AddPaginationToResponse(_ *presentation.Message) {}
