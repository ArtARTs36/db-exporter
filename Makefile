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

.PHONY: functest
functest:
	go build -o ./functest/db-exporter cmd/main.go
	docker-compose up postgres -d
	sleep 5
	FUNCTEST=on DB_EXPORTER_BIN=${PWD}/functest/db-exporter PG_DSN="host=localhost port=5499 user=test password=test dbname=users sslmode=disable" go test ./functest
	docker-compose down
