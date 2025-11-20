package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/artarts36/db-exporter/internal/cli/mcp/protocol"
	"io"
	"time"
)

type Console struct {
	timeout time.Duration
	decoder *json.Decoder
	encoder *json.Encoder
}

func NewConsole(
	timeout time.Duration,
	in io.Reader,
	out io.Writer,
) *Console {
	return &Console{
		timeout: timeout,
		decoder: json.NewDecoder(in),
		encoder: json.NewEncoder(out),
	}
}

func (c *Console) Listen(handler Handler) error {
	for {
		var request protocol.Request
		if err := c.decoder.Decode(&request); err != nil {
			if err == io.EOF {
				// Normal termination
				return nil
			}
			return fmt.Errorf("decode request: %w", err)
		}

		// Ensure JSONRPC version is set
		if request.JSONRPC == "" {
			request.JSONRPC = "2.0"
		}

		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)

		// Process the request
		response, err := handler(ctx, &request)
		cancel()

		if err != nil {
			return fmt.Errorf("handle request: %w", err)
		}

		if err = c.encoder.Encode(response); err != nil {
			return fmt.Errorf("encode response: %w", err)
		}
	}
}
