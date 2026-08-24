package ports

import "context"

type LLM interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type ToolDefinition struct {
	Name        string
	Description string
	Parameters  map[string]any
}

type ToolCall struct {
	ID        string
	Name      string
	Arguments string
}

type ToolOutput struct {
	CallID string
	Output string
}

type AgentMessageRole string

const (
	AgentMessageRoleUser      AgentMessageRole = "user"
	AgentMessageRoleAssistant AgentMessageRole = "assistant"
)

type AgentMessage struct {
	Role    AgentMessageRole
	Content string
}

type AgentRequest struct {
	Messages           []AgentMessage
	Instructions       string
	MaxOutputTokens    int64
	Tools              []ToolDefinition
	PreviousResponseID string
	ToolOutputs        []ToolOutput
}

type AgentResponse struct {
	ID        string
	Message   string
	ToolCalls []ToolCall
}

type AgentLLM interface {
	GenerateWithTools(ctx context.Context, request AgentRequest) (AgentResponse, error)
}
