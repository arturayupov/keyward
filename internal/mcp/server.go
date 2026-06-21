package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arturayupov/keyward/internal/approval"
	"github.com/arturayupov/keyward/internal/broker"
	"github.com/arturayupov/keyward/internal/model"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Handlers holds the dependencies; methods are kept small for unit testing.
type Handlers struct {
	Store  *model.Store
	Broker *broker.Broker
}

// listKeys returns value-free metadata as JSON.
func (h *Handlers) listKeys(_ context.Context, namespace string) (string, error) {
	b, err := json.Marshal(h.Store.Meta(namespace))
	return string(b), err
}

// requestKey runs the broker and returns its value-free Result as JSON.
func (h *Handlers) requestKey(_ context.Context, name, namespace, target, reason string) (string, error) {
	res, err := h.Broker.Request(approval.Request{
		Tool: "mcp", Name: name, Namespace: namespace, Target: target, Reason: reason,
	})
	if err != nil {
		return "", err
	}
	b, _ := json.Marshal(res)
	return string(b), nil
}

// NewServer wires the two tools onto an MCP server.
func NewServer(h *Handlers) *server.MCPServer {
	s := server.NewMCPServer("keyward", "0.1.0")

	list := mcpgo.NewTool("list_keys",
		mcpgo.WithDescription("List available secret names and namespaces. Never returns values."),
		mcpgo.WithString("namespace", mcpgo.Description("optional namespace filter")),
	)
	s.AddTool(list, func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		out, err := h.listKeys(ctx, req.GetString("namespace", ""))
		if err != nil {
			return mcpgo.NewToolResultError(err.Error()), nil
		}
		return mcpgo.NewToolResultText(out), nil
	})

	reqTool := mcpgo.NewTool("request_key",
		mcpgo.WithDescription("Request a secret be injected into a target .env file. Requires human approval. Returns only a confirmation, never the value."),
		mcpgo.WithString("name", mcpgo.Required(), mcpgo.Description("secret name")),
		mcpgo.WithString("namespace", mcpgo.Required(), mcpgo.Description("secret namespace/project")),
		mcpgo.WithString("target", mcpgo.Required(), mcpgo.Description("target .env path")),
		mcpgo.WithString("reason", mcpgo.Description("why the key is needed")),
	)
	s.AddTool(reqTool, func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
		out, err := h.requestKey(ctx,
			req.GetString("name", ""), req.GetString("namespace", ""),
			req.GetString("target", ""), req.GetString("reason", ""))
		if err != nil {
			return mcpgo.NewToolResultError(fmt.Sprintf("request failed: %v", err)), nil
		}
		return mcpgo.NewToolResultText(out), nil
	})

	return s
}
