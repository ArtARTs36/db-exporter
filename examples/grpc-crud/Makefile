gen:
	@echo "> generate .proto"
	db-exporter

	@echo "> generate Go structures for proto messages"
	protoc -I ./contracts \
	  --go_out=./pkg/api --go_opt=paths=source_relative \
	  --go-grpc_out=./pkg/api --go-grpc_opt=paths=source_relative \
	  --proto_path=./contracts/vendor \
	  ./contracts/*.proto
