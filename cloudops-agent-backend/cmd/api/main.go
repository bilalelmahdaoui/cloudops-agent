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
	getCloudServiceUseCase := application.NewGetCloudServiceUseCase(fakeCloudAdapter)

	cloudServiceHandler := httpadapter.NewCloudServiceHandler(getCloudServiceUseCase)

	http.HandleFunc("/cloud-services/", cloudServiceHandler.GetByID)

	fmt.Println("API running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
