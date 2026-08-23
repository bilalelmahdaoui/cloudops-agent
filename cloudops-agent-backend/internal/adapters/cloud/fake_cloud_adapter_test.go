package cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
)

func TestFakeCloudAdapter_GetAll(t *testing.T) {
	adapter := NewFakeCloudAdapter()

	t.Run("retourne les cinq services cloud", func(t *testing.T) {
		services, err := adapter.GetAll(context.Background())

		if err != nil {
			t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
		}

		if len(services) != 5 {
			t.Errorf("nombre de services attendu %d, nombre reçu %d", 5, len(services))
		}
	})

	t.Run("retourne une erreur si le contexte est annulé", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := adapter.GetAll(ctx)

		if !errors.Is(err, context.Canceled) {
			t.Errorf("context.Canceled attendu, erreur reçue : %v", err)
		}
	})
}

func TestFakeCloudAdapter_GetByID(t *testing.T) {
	adapter := NewFakeCloudAdapter()

	t.Run("retourne le service cloud si l'ID existe", func(t *testing.T) {
		service, err := adapter.GetByID(context.Background(), "OVH-SERVICE-003")

		if err != nil {
			t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
		}

		if service.ID != "OVH-SERVICE-003" {
			t.Errorf("ID attendu %q, ID reçu %q", "OVH-SERVICE-003", service.ID)
		}

		if service.Name != "Database service" {
			t.Errorf("nom attendu %q, nom reçu %q", "Database service", service.Name)
		}
	})

	t.Run("retourne une erreur si l'ID n'existe pas", func(t *testing.T) {
		_, err := adapter.GetByID(context.Background(), "UNKNOWN")

		if err == nil {
			t.Fatal("une erreur était attendue, nil reçu")
		}
	})

	t.Run("retourne une erreur si le contexte est annulé", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := adapter.GetByID(ctx, "OVH-SERVICE-003")

		if !errors.Is(err, context.Canceled) {
			t.Errorf("context.Canceled attendu, erreur reçue : %v", err)
		}
	})
}

func TestFakeCloudAdapter_Restart(t *testing.T) {
	adapter := NewFakeCloudAdapter()

	service, err := adapter.Restart(context.Background(), "OVH-SERVICE-003")
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

	if len(service.Logs) < 2 {
		t.Fatal("des logs de redémarrage étaient attendus")
	}

	lastLog := service.Logs[len(service.Logs)-1]

	if lastLog.Event != "Server started successfully" {
		t.Errorf(
			"dernier log attendu %q, log reçu %q",
			"Server started successfully",
			lastLog.Event,
		)
	}
}
