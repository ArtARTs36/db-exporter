package transport

import (
	"context"

	"github.com/artarts36/db-exporter/internal/cli/mcp/protocol"
)

type Transport interface {
	Listen(handler Handler) error
}

type Handler func(ctx context.Context, req *protocol.Request) (*protocol.Response, error)
