package main

import (
	"fmt"
	"net/http"

	cloudadapter "github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
	httpadapter "github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/http"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
)

func main() {
	fakeCloudAdapter := cloudadapter.NewFakeCloudAdapter()

	getAllCloudServicesUseCase := application.NewGetAllCloudServicesUseCase(fakeCloudAdapter)
	getCloudServiceUseCase := application.NewGetCloudServiceUseCase(fakeCloudAdapter)
	restartCloudServiceUseCase := application.NewRestartCloudServiceUseCase(fakeCloudAdapter)

	cloudServiceHandler := httpadapter.NewCloudServiceHandler(
		getAllCloudServicesUseCase,
		getCloudServiceUseCase,
		restartCloudServiceUseCase,
	)

	mux := http.NewServeMux()
	cloudServiceHandler.RegisterRoutes(mux)

	fmt.Println("API running on http://localhost:8080")

	handler := httpadapter.CORSMiddleware(mux)

	err := http.ListenAndServe(":8080", handler)

	if err != nil {
		panic(err)
	}
}
