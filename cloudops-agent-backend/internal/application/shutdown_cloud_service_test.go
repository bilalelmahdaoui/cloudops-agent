package application

import (
	"context"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
)

func TestShutdownCloudServiceUseCase_Execute(t *testing.T) {
	useCase := NewShutdownCloudServiceUseCase(cloud.NewFakeCloudAdapter())

	service, err := useCase.Execute(context.Background(), "OVH-SERVICE-001")
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}
	if service.Status != domain.CloudServiceStatusDown {
		t.Errorf("statut attendu %q, statut reçu %q", domain.CloudServiceStatusDown, service.Status)
	}
}
