package application

import (
	"context"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
)

func TestStartCloudServiceUseCase_Execute(t *testing.T) {
	repository := cloud.NewFakeCloudAdapter()
	if _, err := repository.Shutdown(context.Background(), "OVH-SERVICE-001"); err != nil {
		t.Fatalf("préparation de l'arrêt impossible : %v", err)
	}
	useCase := NewStartCloudServiceUseCase(repository)

	service, err := useCase.Execute(context.Background(), "OVH-SERVICE-001")
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}
	if service.Status != domain.CloudServiceStatusRunning {
		t.Errorf("statut attendu %q, statut reçu %q", domain.CloudServiceStatusRunning, service.Status)
	}
}
