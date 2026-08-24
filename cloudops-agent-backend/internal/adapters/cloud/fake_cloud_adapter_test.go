package cloud

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestFakeCloudAdapter_Search(t *testing.T) {
	adapter := NewFakeCloudAdapter()

	t.Run("recherche un service par une partie de son nom sans tenir compte de la casse", func(t *testing.T) {
		services, err := adapter.Search(context.Background(), "BACKEND", 10)
		if err != nil {
			t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
		}
		if len(services) != 1 {
			t.Fatalf("un service attendu, %d reçus", len(services))
		}
		if services[0].ID != "OVH-SERVICE-001" {
			t.Errorf("ID attendu %q, ID reçu %q", "OVH-SERVICE-001", services[0].ID)
		}
	})

	t.Run("retourne une liste vide sans erreur si aucun service ne correspond", func(t *testing.T) {
		services, err := adapter.Search(context.Background(), "inconnu", 10)
		if err != nil {
			t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
		}
		if len(services) != 0 {
			t.Errorf("aucun service attendu, %d reçus", len(services))
		}
	})

	t.Run("retourne une erreur si la recherche est vide", func(t *testing.T) {
		_, err := adapter.Search(context.Background(), "   ", 10)
		if err == nil {
			t.Fatal("une erreur était attendue")
		}
	})

	t.Run("respecte un contexte annulé", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := adapter.Search(ctx, "backend", 10)
		if !errors.Is(err, context.Canceled) {
			t.Errorf("context.Canceled attendu, erreur reçue : %v", err)
		}
	})

	t.Run("borne le nombre de résultats", func(t *testing.T) {
		services, err := adapter.Search(context.Background(), "service", 2)
		if err != nil {
			t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
		}
		if len(services) != 2 {
			t.Errorf("deux services attendus, %d reçus", len(services))
		}
	})
}

func TestFakeCloudAdapter_Restart(t *testing.T) {
	adapter := NewFakeCloudAdapter()
	adapter.restartDelay = 10 * time.Millisecond

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

func TestFakeCloudAdapter_Restart_ContinueApresAnnulation(t *testing.T) {
	adapter := NewFakeCloudAdapter()
	adapter.restartDelay = 20 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		_, err := adapter.Restart(ctx, "OVH-SERVICE-003")
		done <- err
	}()

	deadline := time.Now().Add(100 * time.Millisecond)
	for {
		service, err := adapter.GetByID(context.Background(), "OVH-SERVICE-003")
		if err != nil {
			t.Fatalf("lecture du service impossible : %v", err)
		}
		if service.Status == domain.CloudServiceStatusRestarting {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("le service devait passer en redémarrage")
		}
		time.Sleep(time.Millisecond)
	}

	cancel()
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Fatalf("context.Canceled attendu, erreur reçue : %v", err)
	}

	time.Sleep(30 * time.Millisecond)
	service, err := adapter.GetByID(context.Background(), "OVH-SERVICE-003")
	if err != nil {
		t.Fatalf("lecture du service impossible : %v", err)
	}
	if service.Status != domain.CloudServiceStatusRunning {
		t.Errorf("le redémarrage devait se terminer malgré l'annulation, statut reçu : %q", service.Status)
	}
}
