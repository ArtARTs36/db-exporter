package mcp

import (
	"github.com/artarts36/db-exporter/internal/cli/config"
	"github.com/artarts36/db-exporter/internal/cli/mcp/server"
	"github.com/artarts36/db-exporter/internal/cli/mcp/tools"
	"github.com/artarts36/db-exporter/internal/cli/mcp/transport"
	"os"
)

func Create(cfg *config.Config, tp transport.Transport) *server.Server {
	router := tools.NewRouter()

	router.RegisterTool(tools.NewGetDBSchemaTool(cfg))

	return server.NewServer(router, os.Stderr, tp)
}
