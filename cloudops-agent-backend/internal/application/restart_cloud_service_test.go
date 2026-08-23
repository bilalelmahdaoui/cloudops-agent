package application

import (
	"context"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
)

func TestRestartCloudServiceUseCase_Execute(t *testing.T) {
	repository := cloud.NewFakeCloudAdapter()
	useCase := NewRestartCloudServiceUseCase(repository)

	service, err := useCase.Execute(context.Background(), "OVH-SERVICE-003")
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}

	if service.Status != domain.CloudServiceStatusRunning {
		t.Errorf(
			"statut attendu %q, statut reçu %q",
			domain.CloudServiceStatusRunning,
			service.Status,
		)
	}
}
