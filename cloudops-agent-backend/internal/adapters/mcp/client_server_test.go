package mcpadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
)

func newTestMCPClient(t *testing.T) *MCPClientAdapter {
	t.Helper()

	repository := cloud.NewFakeCloudAdapter()
	server := NewCloudOpsServer(
		application.NewGetCloudServiceUseCase(repository),
		application.NewGetAllCloudServicesUseCase(repository),
		application.NewSearchCloudServicesUseCase(repository),
		application.NewRestartCloudServiceUseCase(repository),
		application.NewShutdownCloudServiceUseCase(repository),
		application.NewStartCloudServiceUseCase(repository),
	)

	client, err := NewMCPClientAdapter(context.Background(), server)
	if err != nil {
		t.Fatalf("connexion MCP impossible : %v", err)
	}
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Errorf("fermeture MCP impossible : %v", err)
		}
	})

	return client
}

func TestMCPClientAdapter_ListTools(t *testing.T) {
	client := newTestMCPClient(t)

	tools, err := client.ListTools(context.Background())
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}

	expected := map[string]bool{
		GetCloudServiceToolName:      false,
		GetAllCloudServicesToolName:  false,
		FindCloudServicesToolName:    false,
		RestartCloudServiceToolName:  false,
		ShutdownCloudServiceToolName: false,
		StartCloudServiceToolName:    false,
	}
	for _, tool := range tools {
		if _, exists := expected[tool.Name]; exists {
			expected[tool.Name] = true
		}
	}

	if len(tools) != len(expected) {
		t.Errorf("nombre d'outils attendu %d, nombre reçu %d", len(expected), len(tools))
	}
	for name, found := range expected {
		if !found {
			t.Errorf("l'outil %q devait être exposé", name)
		}
	}
}

func TestMCPClientAdapter_CallTool(t *testing.T) {
	t.Run("récupère un service via le cas d'usage existant", func(t *testing.T) {
		client := newTestMCPClient(t)

		output, err := client.CallTool(context.Background(), ports.ToolCall{
			Name:      GetCloudServiceToolName,
			Arguments: `{"id":"OVH-SERVICE-003"}`,
		})
		if err != nil {
			t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
		}

		var service domain.CloudService
		if err := json.Unmarshal([]byte(output), &service); err != nil {
			t.Fatalf("résultat JSON invalide : %v", err)
		}
		if service.ID != "OVH-SERVICE-003" {
			t.Errorf("ID attendu %q, ID reçu %q", "OVH-SERVICE-003", service.ID)
		}
	})

	t.Run("récupère tous les services via le cas d'usage existant", func(t *testing.T) {
		client := newTestMCPClient(t)

		output, err := client.CallTool(context.Background(), ports.ToolCall{
			Name:      GetAllCloudServicesToolName,
			Arguments: `{}`,
		})
		if err != nil {
			t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
		}

		var services []domain.CloudService
		if err := json.Unmarshal([]byte(output), &services); err != nil {
			t.Fatalf("résultat JSON invalide : %v", err)
		}
		if len(services) != 5 {
			t.Errorf("nombre de services attendu %d, nombre reçu %d", 5, len(services))
		}
	})

	t.Run("recherche un service par son nom via le cas d'usage existant", func(t *testing.T) {
		client := newTestMCPClient(t)

		output, err := client.CallTool(context.Background(), ports.ToolCall{
			Name:      FindCloudServicesToolName,
			Arguments: `{"query":"backend"}`,
		})
		if err != nil {
			t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
		}

		var services []domain.CloudService
		if err := json.Unmarshal([]byte(output), &services); err != nil {
			t.Fatalf("résultat JSON invalide : %v", err)
		}
		if len(services) != 1 || services[0].ID != "OVH-SERVICE-001" {
			t.Fatalf("le service backend avec son ID canonique était attendu")
		}
	})

	t.Run("recherche un service avec une faute ou une abréviation", func(t *testing.T) {
		client := newTestMCPClient(t)

		tests := []struct {
			query      string
			expectedID string
		}{
			{query: "authentifaciton", expectedID: "OVH-SERVICE-004"},
			{query: "db service", expectedID: "OVH-SERVICE-003"},
		}

		for _, test := range tests {
			output, err := client.CallTool(context.Background(), ports.ToolCall{
				Name:      FindCloudServicesToolName,
				Arguments: fmt.Sprintf(`{"query":%q}`, test.query),
			})
			if err != nil {
				t.Fatalf("recherche %q impossible : %v", test.query, err)
			}

			var services []domain.CloudService
			if err := json.Unmarshal([]byte(output), &services); err != nil {
				t.Fatalf("résultat JSON invalide : %v", err)
			}
			if len(services) != 1 || services[0].ID != test.expectedID {
				t.Fatalf("service attendu %q pour la recherche %q", test.expectedID, test.query)
			}
		}
	})

	t.Run("retourne une erreur métier comme résultat d'outil", func(t *testing.T) {
		client := newTestMCPClient(t)

		output, err := client.CallTool(context.Background(), ports.ToolCall{
			Name:      GetCloudServiceToolName,
			Arguments: `{"id":"SERVICE-INCONNU"}`,
		})
		if err != nil {
			t.Fatalf("l'erreur métier devait être transmise au LLM, erreur reçue : %v", err)
		}

		var result map[string]string
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			t.Fatalf("résultat d'erreur JSON invalide : %v", err)
		}
		if result["error"] == "" {
			t.Fatal("le résultat devait décrire l'erreur métier")
		}
	})

	t.Run("redémarre un service via le cas d'usage existant", func(t *testing.T) {
		client := newTestMCPClient(t)

		output, err := client.CallTool(context.Background(), ports.ToolCall{
			Name:      RestartCloudServiceToolName,
			Arguments: `{"id":"OVH-SERVICE-003"}`,
		})
		if err != nil {
			t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
		}

		var service domain.CloudService
		if err := json.Unmarshal([]byte(output), &service); err != nil {
			t.Fatalf("résultat JSON invalide : %v", err)
		}
		if service.Status != domain.CloudServiceStatusRunning {
			t.Errorf(
				"statut attendu %q, statut reçu %q",
				domain.CloudServiceStatusRunning,
				service.Status,
			)
		}
		if len(service.Logs) != 4 {
			t.Errorf("nombre de logs attendu %d, nombre reçu %d", 4, len(service.Logs))
		}
	})

	t.Run("arrête puis démarre un service via les cas d'usage existants", func(t *testing.T) {
		client := newTestMCPClient(t)

		shutdownOutput, err := client.CallTool(context.Background(), ports.ToolCall{
			Name:      ShutdownCloudServiceToolName,
			Arguments: `{"id":"OVH-SERVICE-003"}`,
		})
		if err != nil {
			t.Fatalf("arrêt impossible : %v", err)
		}

		var stoppedService domain.CloudService
		if err := json.Unmarshal([]byte(shutdownOutput), &stoppedService); err != nil {
			t.Fatalf("résultat JSON d'arrêt invalide : %v", err)
		}
		if stoppedService.Status != domain.CloudServiceStatusDown || stoppedService.CPUUsage != 0 {
			t.Fatal("le service devait être arrêté avec un CPU nul")
		}

		startOutput, err := client.CallTool(context.Background(), ports.ToolCall{
			Name:      StartCloudServiceToolName,
			Arguments: `{"id":"OVH-SERVICE-003"}`,
		})
		if err != nil {
			t.Fatalf("démarrage impossible : %v", err)
		}

		var startedService domain.CloudService
		if err := json.Unmarshal([]byte(startOutput), &startedService); err != nil {
			t.Fatalf("résultat JSON de démarrage invalide : %v", err)
		}
		if startedService.Status != domain.CloudServiceStatusRunning || startedService.CPUUsage != 0.65 {
			t.Fatal("le service devait être démarré avec son CPU nominal")
		}
	})
}
