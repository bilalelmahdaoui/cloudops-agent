package mcpadapter

import (
	"context"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	GetCloudServiceToolName     = "get_cloud_service"
	GetAllCloudServicesToolName = "get_all_cloud_services"
	FindCloudServicesToolName   = "find_cloud_services"
	RestartCloudServiceToolName = "restart_cloud_service"
)

type cloudServiceInput struct {
	ID string `json:"id" jsonschema:"identifiant du service cloud"`
}

type findCloudServicesInput struct {
	Query string `json:"query" jsonschema:"nom ou fragment de nom du service cloud"`
}

type cloudServiceReference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type emptyInput struct{}

func NewCloudOpsServer(
	getCloudServiceUseCase *application.GetCloudServiceUseCase,
	getAllCloudServicesUseCase *application.GetAllCloudServicesUseCase,
	searchCloudServicesUseCase *application.SearchCloudServicesUseCase,
	restartCloudServiceUseCase *application.RestartCloudServiceUseCase,
) *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "cloudops-mcp-server", Version: "1.0.0"},
		nil,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name: GetCloudServiceToolName,
			Description: "Retrieves the authoritative current state, CPU usage and logs for one cloud service. " +
				"Use this tool whenever a user asks about a specific service. Requires its canonical id, formatted " +
				"OVH-SERVICE-001; resolve a suffix such as 003 to OVH-SERVICE-003 before calling it.",
		},
		func(
			ctx context.Context,
			_ *mcp.CallToolRequest,
			input cloudServiceInput,
		) (*mcp.CallToolResult, domain.CloudService, error) {
			service, err := getCloudServiceUseCase.Execute(ctx, input.ID)
			return nil, service, err
		},
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name: GetAllCloudServicesToolName,
			Description: "Retrieves the authoritative current state of all cloud services. " +
				"Use this tool for comparisons, global status, availability, and highest or lowest CPU questions.",
		},
		func(
			ctx context.Context,
			_ *mcp.CallToolRequest,
			_ emptyInput,
		) (*mcp.CallToolResult, []domain.CloudService, error) {
			services, err := getAllCloudServicesUseCase.Execute(ctx)
			return nil, services, err
		},
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name: FindCloudServicesToolName,
			Description: "Finds cloud services by a full or partial human-readable name, case-insensitively, and returns " +
				"their canonical ids. Use this before get_cloud_service or restart_cloud_service when the user identifies " +
				"a service by name instead of a canonical OVH-SERVICE-xxx id. Only act automatically when exactly one " +
				"service matches; ask the user to clarify when several services match.",
		},
		func(
			ctx context.Context,
			_ *mcp.CallToolRequest,
			input findCloudServicesInput,
		) (*mcp.CallToolResult, []cloudServiceReference, error) {
			services, err := searchCloudServicesUseCase.Execute(ctx, input.Query)
			if err != nil {
				return nil, nil, err
			}

			references := make([]cloudServiceReference, len(services))
			for index, service := range services {
				references[index] = cloudServiceReference{
					ID:   service.ID,
					Name: service.Name,
				}
			}
			return nil, references, nil
		},
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name: RestartCloudServiceToolName,
			Description: "Restarts the specified cloud service and returns its authoritative state after the restart. " +
				"This tool has a side effect; use it when the user explicitly requests a restart. Requires the canonical " +
				"service id. If the user provides a service name, call find_cloud_services first, then restart the matching id.",
		},
		func(
			ctx context.Context,
			_ *mcp.CallToolRequest,
			input cloudServiceInput,
		) (*mcp.CallToolResult, domain.CloudService, error) {
			service, err := restartCloudServiceUseCase.Execute(ctx, input.ID)
			return nil, service, err
		},
	)

	return server
}
