package transport

import (
	"encoding/json"
	"github.com/artarts36/db-exporter/internal/cli/mcp/protocol"
	"io"
	"log/slog"
	"net/http"
)

type HTTP struct {
}

func NewHTTP() *HTTP {
	return &HTTP{}
}

func (h *HTTP) Listen(handler Handler) error {
	srv := &http.Server{Addr: ":8080"}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var protoReq protocol.Request

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			slog.
				With(slog.Any("err", err)).
				ErrorContext(r.Context(), "failed to read request body")
			return
		}

		if err = json.Unmarshal(reqBody, &protoReq); err != nil {
			slog.With(slog.Any("err", err)).ErrorContext(r.Context(), "failed to unmarshal request body")
			return
		}

		resp, err := handler(r.Context(), &protoReq)
		if err != nil {
			slog.With(slog.Any("err", err)).ErrorContext(r.Context(), "failed to handle request")
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			slog.With(slog.Any("err", err)).ErrorContext(r.Context(), "failed to marshal response")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(respBody); err != nil {
			slog.With(slog.Any("err", err)).ErrorContext(r.Context(), "failed to write response")
			return
		}
	})

	srv.Handler = mux

	return srv.ListenAndServe()
}
