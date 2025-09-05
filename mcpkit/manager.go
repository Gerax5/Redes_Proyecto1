package mcpkit

import (
	"context"
	"fmt"
	"strings"

	mcpc "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type Manager struct {
	client *mcpc.Client
}

// New crea un cliente MCP apuntando a una URL streamable-http, p.ej. http://127.0.0.1:3001/mcp
func New(ctx context.Context, url string) (*Manager, error) {
	trans, err := transport.NewStreamableHTTP(url)
	if err != nil {
		return nil, fmt.Errorf("crear transporte: %w", err)
	}
	c := mcpc.NewClient(trans)
	if err := c.Start(ctx); err != nil {
		return nil, fmt.Errorf("iniciar cliente: %w", err)
	}
	if _, err := c.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo:      mcp.Implementation{Name: "GoChat", Version: "0.1.0"},
			Capabilities:    mcp.ClientCapabilities{},
		},
	}); err != nil {
		return nil, fmt.Errorf("initialize: %w", err)
	}
	return &Manager{client: c}, nil
}

// ListTools devuelve la lista de tools del servidor MCP.
func (m *Manager) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	resp, err := m.client.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Tools, nil
}

// Call ejecuta una tool del MCP y devuelve un texto concatenado (solo maneja contenido "text").
func (m *Manager) Call(ctx context.Context, name string, args map[string]any) (string, error) {
	cr, err := m.client.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      name,
			Arguments: args,
		},
	})
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, ct := range cr.Content {
		switch v := ct.(type) {
		case mcp.TextContent:
			if sb.Len() > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(v.Text)
		default:
			// ignora otros tipos por ahora
		}
	}
	return sb.String(), nil
}
