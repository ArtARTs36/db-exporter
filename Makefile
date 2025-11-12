CURRENT_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
BUILD_FLAGS := -ldflags="-X 'main.Version=v0.1.0' -X 'main.BuildDate=${CURRENT_DATE}'"

.PHONY: build
build:
	go build ${BUILD_FLAGS} -o db-exporter cmd/db-exporter/main.go

.PHONY: install
install:
	make build
	cp ./db-exporter /usr/local/bin/db-exporter

.PHONY: install
help:
	go run ./cmd/db-exporter/main.go --help

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run --fix

.PHONY: functest
functest: functest/pg functest/mysql

.PHONY: functest/pg
functest/pg:
	#docker-compose down
	go build -o ./functest/db-exporter cmd/db-exporter/main.go
	docker-compose up postgres -d
	sleep 5
	FUNCTEST=on DB_EXPORTER_BIN=${PWD}/functest/db-exporter PG_DSN="host=localhost port=5419 user=test password=test dbname=users sslmode=disable" go test ./functest
	docker-compose down
	rm ./functest/db-exporter

.PHONY: functest/mysql
functest/mysql:
	docker-compose down
	go build -o ./functest/db-exporter cmd/db-exporter/main.go
	docker-compose up mysql -d
	sleep 20
	FUNCTEST=on DB_EXPORTER_BIN=${PWD}/functest/db-exporter MYSQL_DSN="test:test@tcp(localhost:3306)/users" go test ./functest
	docker-compose down
	rm ./functest/db-exporter

.PHONY: check
check: lint validate-test-cases test functest

.PHONY: validate-test-cases
validate-test-cases:
	protolint lint -config_path=functest/protolint.yaml ./functest
