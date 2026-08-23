package application

import (
	"context"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
)

func TestGetAllCloudServicesUseCase_Execute(t *testing.T) {
	repository := cloud.NewFakeCloudAdapter()
	useCase := NewGetAllCloudServicesUseCase(repository)

	services, err := useCase.Execute(context.Background())
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}

	if len(services) != 5 {
		t.Errorf("nombre de services attendu %d, nombre reçu %d", 5, len(services))
	}
}
