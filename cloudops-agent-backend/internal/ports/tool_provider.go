package ports

import "context"

type ToolProvider interface {
	ListTools(ctx context.Context) ([]ToolDefinition, error)
	CallTool(ctx context.Context, call ToolCall) (string, error)
}
