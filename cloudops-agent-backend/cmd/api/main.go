package main

import (
	"context"
	"fmt"
	"net/http"

	cloudadapter "github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
	httpadapter "github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/http"
	llmadapter "github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/llm"
	mcpadapter "github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/mcp"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	fakeCloudAdapter := cloudadapter.NewFakeCloudAdapter()

	openAIAdapter := llmadapter.NewOpenAIAdapter(
		cfg.OpenAIAPIKey,
		cfg.OpenAIModel,
	)

	getAllCloudServicesUseCase := application.NewGetAllCloudServicesUseCase(fakeCloudAdapter)
	getCloudServiceUseCase := application.NewGetCloudServiceUseCase(fakeCloudAdapter)
	searchCloudServicesUseCase := application.NewSearchCloudServicesUseCase(fakeCloudAdapter)
	restartCloudServiceUseCase := application.NewRestartCloudServiceUseCase(fakeCloudAdapter)
	shutdownCloudServiceUseCase := application.NewShutdownCloudServiceUseCase(fakeCloudAdapter)
	startCloudServiceUseCase := application.NewStartCloudServiceUseCase(fakeCloudAdapter)

	mcpServer := mcpadapter.NewCloudOpsServer(
		getCloudServiceUseCase,
		getAllCloudServicesUseCase,
		searchCloudServicesUseCase,
		restartCloudServiceUseCase,
		shutdownCloudServiceUseCase,
		startCloudServiceUseCase,
	)
	mcpClient, err := mcpadapter.NewMCPClientAdapter(context.Background(), mcpServer)
	if err != nil {
		panic(err)
	}
	defer mcpClient.Close()

	chatWithAgentUseCase := application.NewChatWithAgentUseCase(
		openAIAdapter,
		mcpClient,
	)

	cloudServiceHandler := httpadapter.NewCloudServiceHandler(
		getAllCloudServicesUseCase,
		getCloudServiceUseCase,
		restartCloudServiceUseCase,
		shutdownCloudServiceUseCase,
		startCloudServiceUseCase,
	)
	chatHandler := httpadapter.NewChatHandler(chatWithAgentUseCase)

	mux := http.NewServeMux()
	cloudServiceHandler.RegisterRoutes(mux)
	chatHandler.RegisterRoutes(mux)

	handler := httpadapter.CORSMiddleware(mux)

	fmt.Println("API running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
