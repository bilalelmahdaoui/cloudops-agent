package mcpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPClientAdapter struct {
	clientSession *mcp.ClientSession
	serverSession *mcp.ServerSession
}

func NewMCPClientAdapter(
	ctx context.Context,
	server *mcp.Server,
) (*MCPClientAdapter, error) {
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		return nil, fmt.Errorf("connexion au serveur MCP : %w", err)
	}

	client := mcp.NewClient(
		&mcp.Implementation{Name: "cloudops-mcp-client", Version: "1.0.0"},
		nil,
	)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		_ = serverSession.Close()
		return nil, fmt.Errorf("connexion du client MCP : %w", err)
	}

	return &MCPClientAdapter{
		clientSession: clientSession,
		serverSession: serverSession,
	}, nil
}

func (a *MCPClientAdapter) ListTools(ctx context.Context) ([]ports.ToolDefinition, error) {
	var definitions []ports.ToolDefinition

	for tool, err := range a.clientSession.Tools(ctx, nil) {
		if err != nil {
			return nil, fmt.Errorf("découverte des outils MCP : %w", err)
		}

		parameters, err := schemaAsMap(tool.InputSchema)
		if err != nil {
			return nil, fmt.Errorf("schéma de l'outil MCP %s : %w", tool.Name, err)
		}

		definitions = append(definitions, ports.ToolDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  parameters,
		})
	}

	return definitions, nil
}

func (a *MCPClientAdapter) CallTool(
	ctx context.Context,
	call ports.ToolCall,
) (string, error) {
	var arguments map[string]any
	if err := json.Unmarshal([]byte(call.Arguments), &arguments); err != nil {
		return "", fmt.Errorf("arguments invalides pour l'outil MCP %s : %w", call.Name, err)
	}

	result, err := a.clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      call.Name,
		Arguments: arguments,
	})
	if err != nil {
		return "", fmt.Errorf("appel de l'outil MCP %s : %w", call.Name, err)
	}

	if result.IsError {
		output, err := json.Marshal(map[string]string{
			"error": textContent(result.Content),
		})
		if err != nil {
			return "", fmt.Errorf("sérialisation de l'erreur MCP %s : %w", call.Name, err)
		}
		return string(output), nil
	}

	if result.StructuredContent != nil {
		output, err := json.Marshal(result.StructuredContent)
		if err != nil {
			return "", fmt.Errorf("sérialisation du résultat MCP %s : %w", call.Name, err)
		}
		return string(output), nil
	}

	return textContent(result.Content), nil
}

func (a *MCPClientAdapter) Close() error {
	clientErr := a.clientSession.Close()
	serverErr := a.serverSession.Close()
	if clientErr != nil {
		return clientErr
	}
	return serverErr
}

func schemaAsMap(schema any) (map[string]any, error) {
	encoded, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	var parameters map[string]any
	if err := json.Unmarshal(encoded, &parameters); err != nil {
		return nil, err
	}
	return parameters, nil
}

func textContent(contents []mcp.Content) string {
	var text strings.Builder
	for _, content := range contents {
		if value, ok := content.(*mcp.TextContent); ok {
			text.WriteString(value.Text)
		}
	}
	return text.String()
}

var _ ports.ToolProvider = (*MCPClientAdapter)(nil)
