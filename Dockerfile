# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder

ARG APP_VERSION="undefined"
ARG BUILD_TIME="undefined"

WORKDIR /go/src/github.com/artarts36/db-exporter

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --update gcc libc-dev
RUN GOOS=linux CGO_ENABLED=1 go build -ldflags="-s -w -extldflags=-static -X 'main.Version=${APP_VERSION}' -X 'main.BuildDate=${BUILD_TIME}'" -o /go/bin/db-exporter /go/src/github.com/artarts36/db-exporter/cmd/main.go

######################################################

FROM alpine

RUN apk add tini git

COPY --from=builder /go/bin/db-exporter /go/bin/db-exporter

# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.title="db-exporter"
LABEL org.opencontainers.image.description="simple app for export db schema to many formats"
LABEL org.opencontainers.image.url="https://github.com/artarts36/db-exporter"
LABEL org.opencontainers.image.source="https://github.com/artarts36/db-exporter"
LABEL org.opencontainers.image.vendor="ArtARTs36"
LABEL org.opencontainers.image.version="$APP_VERSION"
LABEL org.opencontainers.image.created="$BUILD_TIME"
LABEL org.opencontainers.image.licenses="MIT"

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x ./docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
