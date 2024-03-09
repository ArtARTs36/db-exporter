CURRENT_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')

BUILD_FLAGS := -ldflags="-X 'main.Version=v0.1.0' -X 'main.BuildDate=${CURRENT_DATE}'"

build:
	go build ${BUILD_FLAGS} -o db-exporter cmd/main.go

help:
	go run ./cmd/main.go --help

test:
	go test ./...

lint:
	golangci-lint run
