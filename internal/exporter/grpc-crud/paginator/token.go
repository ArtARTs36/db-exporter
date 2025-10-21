package paginator

import "github.com/artarts36/db-exporter/internal/exporter/grpc-crud/presentation"

type Token struct {
}

func (o *Token) AddPaginationToRequest(req *presentation.Message) {
	req.CreateField("page_size", func(field *presentation.Field) {
		field.SetType("uint64").SetTopComment("Maximum number of results per page.")
	})

	req.CreateField("page_token", func(field *presentation.Field) {
		field.SetType("string").SetTopComment("Token of the requested results page.")
	})
}

func (o *Token) AddPaginationToResponse(resp *presentation.Message) {
	resp.CreateField("next_page_token", func(field *presentation.Field) {
		field.SetType("string").SetTopComment("Next page token.")
	})
}
