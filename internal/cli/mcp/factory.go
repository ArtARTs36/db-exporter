package mcp

import (
	"github.com/artarts36/db-exporter/internal/cli/mcp/transport"
	"os"
	"time"

	"github.com/artarts36/db-exporter/internal/cli/config"
	"github.com/artarts36/db-exporter/internal/cli/mcp/server"
	"github.com/artarts36/db-exporter/internal/cli/mcp/tools"
)

func Create(cfg *config.Config) *server.Server {
	router := tools.NewRouter()

	router.RegisterTool(tools.NewGetDBSchemaTool(cfg))

	return server.NewServer(router, os.Stderr, transport.NewConsole(
		time.Minute,
		os.Stdin,
		os.Stdout,
	))
}
