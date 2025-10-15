package googleapihttp

import "github.com/artarts36/db-exporter/internal/shared/proto"

func Get(url string) *proto.ServiceProcedureOption {
	return opt("get", url)
}

func Post(url string) *proto.ServiceProcedureOption {
	return opt("post", url)
}

func Put(url string) *proto.ServiceProcedureOption {
	return opt("put", url)
}

func Delete(url string) *proto.ServiceProcedureOption {
	return opt("delete", url)
}

func opt(method string, url string) *proto.ServiceProcedureOption {
	return &proto.ServiceProcedureOption{
		Name: "google.api.http",
		Params: map[string]interface{}{
			method: url,
		},
	}
}
