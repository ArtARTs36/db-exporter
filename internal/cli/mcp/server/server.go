package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/artarts36/db-exporter/internal/cli/mcp/protocol"
	"github.com/artarts36/db-exporter/internal/cli/mcp/tools"
	"io"
	"os"
	"time"
)

// Server represents an MCP server
type Server struct {
	tools  *tools.Router
	stdout io.Writer
	stderr io.Writer
}

// NewServer creates a new MCP server
func NewServer(
	router *tools.Router,
	stdout,
	stderr io.Writer,
) *Server {
	return &Server{
		tools:  router,
		stdout: stdout,
		stderr: stderr,
	}
}

const (
	reqTimeout = time.Minute
)

// Run starts the MCP server
func (s *Server) Run() error {
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(s.stdout)

	for {
		var request protocol.Request
		if err := decoder.Decode(&request); err != nil {
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

		ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)

		// Process the request
		response := s.handleRequest(ctx, &request)
		cancel()

		// Only encode and send a response if there is one
		// Notifications don't require a response
		if response != nil {
			if err := encoder.Encode(response); err != nil {
				return fmt.Errorf("encode response: %w", err)
			}
		}
	}
}

// handleRequest handles a single JSON-RPC request
func (s *Server) handleRequest(ctx context.Context, req *protocol.Request) *protocol.Response {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "notifications/initialized":
		_, _ = fmt.Fprintf(s.stderr, "Received initialized notification, client is ready\n")
		return nil
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(ctx, req)
	default:
		return &protocol.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   protocol.MethodNotFoundError(req.Method),
		}
	}
}

// handleInitialize handles initialization requests according to the MCP protocol
func (s *Server) handleInitialize(req *protocol.Request) *protocol.Response {
	// Parse the initialize parameters
	var params protocol.InitializeParams
	if len(req.Params) > 0 {
		if err := json.Unmarshal(req.Params, &params); err != nil {
			_, _ = fmt.Fprintf(s.stderr, "Warning: failed to parse initialize params: %v\n", err)
			return &protocol.Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   protocol.InvalidParamsError(err),
			}
		}
	}

	// Log information about the client
	if params.ClientInfo.Name != "" {
		_, _ = fmt.Fprintf(s.stderr, "Client info: %s %s\n", params.ClientInfo.Name, params.ClientInfo.Version)
	}

	// Log the client's protocol version
	if params.ProtocolVersion != "" {
		_, _ = fmt.Fprintf(s.stderr, "Client protocol version: %s\n", params.ProtocolVersion)
	}

	// Check if we support the client's protocol version
	// We only support 2024-11-05
	supportedVersion := "2024-11-05"
	if params.ProtocolVersion != "" && params.ProtocolVersion != supportedVersion {
		_, _ = fmt.Fprintf(s.stderr, "Warning: Client requested protocol version %s, but we're responding with %s\n",
			params.ProtocolVersion, supportedVersion)
		// We still continue with our supported version - the client will judge compatibility
	}

	// Create initialize result with proper MCP protocol capabilities
	result := protocol.InitializeResult{
		ProtocolVersion: supportedVersion,
		ServerInfo: protocol.ServerInfo{
			Name:    "db-exporter-mcp",
			Version: "0.1.0",
		},
		Capabilities: protocol.ServerCapabilities{
			// We only support tools with listChanged capability
			Tools: map[string]interface{}{
				"listChanged": true,
			},
		},
		Instructions: "db-exporter",
	}

	return &protocol.Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (s *Server) handleToolsList(req *protocol.Request) *protocol.Response {
	toolInfos := make([]protocol.ToolInfo, 0, len(s.tools.List()))
	for _, tool := range s.tools.List() {
		toolInfos = append(toolInfos, tool.Info())
	}

	return &protocol.Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: &protocol.ListToolsResponse{
			Tools: toolInfos,
		},
	}
}

func (s *Server) handleToolsCall(ctx context.Context, req *protocol.Request) *protocol.Response {
	var payload protocol.CallToolPayload
	if err := json.Unmarshal(req.Params, &payload); err != nil {
		return &protocol.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &protocol.Error{
				Code:    -32700, // Parse error
				Message: fmt.Sprintf("Invalid payload: %v", err),
			},
		}
	}

	tool, toolFound := s.tools.Find(payload.Name)
	if !toolFound {
		return &protocol.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   protocol.MethodNotFoundError(payload.Name),
		}
	}

	// Execute the tool
	result, err := tool.Execute(ctx, payload.Arguments)
	if err != nil {
		return &protocol.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &protocol.Error{
				Code:    -32000, // Server error
				Message: fmt.Sprintf("Tool execution failed: %v", err),
			},
		}
	}

	resultText, err := json.Marshal(result)
	if err != nil {
		return &protocol.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &protocol.Error{
				Code:    -32000, // Server error
				Message: fmt.Sprintf("Failed to marshal result: %v", err),
			},
		}
	}

	return &protocol.Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: &protocol.CallToolResponse{
			Content: []protocol.ContentItem{
				{
					Type: "text",
					Text: string(resultText),
				},
			},
			IsError: false,
		},
	}
}
