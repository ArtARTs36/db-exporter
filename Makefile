help:
	go run ./cmd/main.go --help

test:
	go test ./...

lint:
	golangci-lint run
