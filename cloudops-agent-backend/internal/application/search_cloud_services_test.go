package application

import (
	"context"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
)

func TestSearchCloudServicesUseCase_Execute(t *testing.T) {
	useCase := NewSearchCloudServicesUseCase(cloud.NewFakeCloudAdapter())

	services, err := useCase.Execute(context.Background(), "backend")
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}
	if len(services) != 1 {
		t.Fatalf("un service attendu, %d reçus", len(services))
	}
	if services[0].ID != "OVH-SERVICE-001" {
		t.Errorf("ID attendu %q, ID reçu %q", "OVH-SERVICE-001", services[0].ID)
	}
}
