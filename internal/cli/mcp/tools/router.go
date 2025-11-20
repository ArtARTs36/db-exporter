package tools

import (
	"context"
	"github.com/artarts36/gds"

	"github.com/artarts36/db-exporter/internal/cli/mcp/protocol"
)

// Tool defines the interface for MCP tools
type Tool interface {
	Info() protocol.ToolInfo

	Execute(ctx context.Context, args map[string]interface{}) (any, error)
}

type Router struct {
	tools *gds.Map[string, Tool]
}

func NewRouter() *Router {
	return &Router{
		tools: gds.NewMap[string, Tool](),
	}
}

func (r *Router) RegisterTool(t Tool) {
	r.tools.Set(t.Info().Name, t)
}

func (r *Router) Find(name string) (Tool, bool) {
	return r.tools.Get(name)
}

func (r *Router) List() []Tool {
	return r.tools.List()
}
