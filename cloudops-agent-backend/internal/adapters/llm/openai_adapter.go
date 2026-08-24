package llm

import (
	"context"
	"fmt"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

type OpenAIAdapter struct {
	client openai.Client
	model  string
}

func NewOpenAIAdapter(apiKey string, model string) *OpenAIAdapter {
	return &OpenAIAdapter{
		client: openai.NewClient(
			option.WithAPIKey(apiKey),
		),
		model: model,
	}
}

func (a *OpenAIAdapter) Generate(
	ctx context.Context,
	prompt string,
) (string, error) {
	response, err := a.client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model: a.model,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(prompt),
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate LLM response: %w", err)
	}

	return response.OutputText(), nil
}

func (a *OpenAIAdapter) GenerateWithTools(
	ctx context.Context,
	request ports.AgentRequest,
) (ports.AgentResponse, error) {
	tools := make([]responses.ToolUnionParam, 0, len(request.Tools))
	for _, definition := range request.Tools {
		tool := responses.ToolParamOfFunction(
			definition.Name,
			openAIToolParameters(definition.Parameters),
			true,
		)
		tool.OfFunction.Description = openai.String(definition.Description)
		tools = append(tools, tool)
	}

	params := responses.ResponseNewParams{
		Model:           a.model,
		Tools:           tools,
		Instructions:    openai.String(request.Instructions),
		MaxOutputTokens: openai.Int(request.MaxOutputTokens),
		Text: responses.ResponseTextConfigParam{
			Verbosity: responses.ResponseTextConfigVerbosityLow,
		},
		Reasoning: shared.ReasoningParam{
			Effort: shared.ReasoningEffortLow,
		},
	}

	if request.PreviousResponseID != "" {
		params.PreviousResponseID = openai.String(request.PreviousResponseID)
	}

	if len(request.ToolOutputs) > 0 {
		input := make(responses.ResponseInputParam, 0, len(request.ToolOutputs))
		for _, toolOutput := range request.ToolOutputs {
			input = append(
				input,
				responses.ResponseInputItemParamOfFunctionCallOutput(toolOutput.CallID, toolOutput.Output),
			)
		}
		params.Input.OfInputItemList = input
	} else {
		params.Input.OfInputItemList = openAIConversationInput(request.Messages)
	}

	response, err := a.client.Responses.New(ctx, params)
	if err != nil {
		return ports.AgentResponse{}, fmt.Errorf("failed to generate LLM response: %w", err)
	}

	toolCalls := make([]ports.ToolCall, 0)
	for _, output := range response.Output {
		if output.Type != "function_call" {
			continue
		}

		functionCall := output.AsFunctionCall()
		toolCalls = append(toolCalls, ports.ToolCall{
			ID:        functionCall.CallID,
			Name:      functionCall.Name,
			Arguments: functionCall.Arguments,
		})
	}

	return ports.AgentResponse{
		ID:        response.ID,
		Message:   response.OutputText(),
		ToolCalls: toolCalls,
	}, nil
}

func openAIConversationInput(messages []ports.AgentMessage) responses.ResponseInputParam {
	input := make(responses.ResponseInputParam, 0, len(messages))
	for _, message := range messages {
		role := responses.EasyInputMessageRoleUser
		if message.Role == ports.AgentMessageRoleAssistant {
			role = responses.EasyInputMessageRoleAssistant
		}

		input = append(
			input,
			responses.ResponseInputItemParamOfMessage(message.Content, role),
		)
	}
	return input
}

func openAIToolParameters(parameters map[string]any) map[string]any {
	normalized := make(map[string]any, len(parameters)+2)
	for key, value := range parameters {
		normalized[key] = value
	}

	if normalized["type"] == "object" {
		if _, exists := normalized["properties"]; !exists {
			normalized["properties"] = map[string]any{}
		}
		if _, exists := normalized["required"]; !exists {
			normalized["required"] = []string{}
		}
	}

	return normalized
}

var _ ports.LLM = (*OpenAIAdapter)(nil)
var _ ports.AgentLLM = (*OpenAIAdapter)(nil)
