package application

import (
	"context"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
)

func TestGetCloudServiceUseCase_Execute(t *testing.T) {
	repository := cloud.NewFakeCloudAdapter()
	useCase := NewGetCloudServiceUseCase(repository)

	t.Run("retourne le service cloud si l'ID existe", func(t *testing.T) {
		service, err := useCase.Execute(context.Background(), "OVH-SERVICE-003")

		if err != nil {
			t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
		}

		if service.ID != "OVH-SERVICE-003" {
			t.Errorf("ID attendu %q, ID reçu %q", "OVH-SERVICE-003", service.ID)
		}
	})

	t.Run("retourne une erreur si l'ID n'existe pas", func(t *testing.T) {
		_, err := useCase.Execute(context.Background(), "UNKNOWN")

		if err == nil {
			t.Fatal("une erreur était attendue, nil reçu")
		}
	})
}
