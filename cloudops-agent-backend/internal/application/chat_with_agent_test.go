package application

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
)

type scriptedAgentLLM struct {
	responses []ports.AgentResponse
	requests  []ports.AgentRequest
}

func (f *scriptedAgentLLM) GenerateWithTools(
	_ context.Context,
	request ports.AgentRequest,
) (ports.AgentResponse, error) {
	index := len(f.requests)
	f.requests = append(f.requests, request)
	if index >= len(f.responses) {
		return ports.AgentResponse{}, fmt.Errorf("aucune réponse préparée")
	}
	return f.responses[index], nil
}

type fakeToolProvider struct {
	tools     []ports.ToolDefinition
	outputs   map[string]string
	listCalls int
	calls     []ports.ToolCall
}

func (f *fakeToolProvider) ListTools(_ context.Context) ([]ports.ToolDefinition, error) {
	f.listCalls++
	return f.tools, nil
}

func (f *fakeToolProvider) CallTool(
	_ context.Context,
	call ports.ToolCall,
) (string, error) {
	f.calls = append(f.calls, call)
	output, exists := f.outputs[call.Name]
	if !exists {
		return "", fmt.Errorf("outil inconnu : %s", call.Name)
	}
	return output, nil
}

func testToolDefinitions() []ports.ToolDefinition {
	return []ports.ToolDefinition{
		{Name: "get_cloud_service"},
		{Name: "get_all_cloud_services"},
		{Name: "find_cloud_services"},
		{Name: "restart_cloud_service"},
		{Name: "shutdown_cloud_service"},
		{Name: "start_cloud_service"},
	}
}

func TestChatWithAgentUseCase_Execute(t *testing.T) {
	llm := &scriptedAgentLLM{responses: []ports.AgentResponse{{
		Message: "Je peux surveiller et piloter vos services cloud.",
	}}}
	provider := &fakeToolProvider{tools: testToolDefinitions()}
	useCase := NewChatWithAgentUseCase(llm, provider)

	response, err := useCase.Execute(context.Background(), "Que peux-tu faire ?", nil)
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}

	if response != "Je peux surveiller et piloter vos services cloud." {
		t.Errorf("réponse inattendue : %q", response)
	}
	if provider.listCalls != 1 {
		t.Errorf("un appel de découverte attendu, %d reçu", provider.listCalls)
	}
	if !strings.Contains(llm.requests[0].Instructions, "CloudOps Agent") {
		t.Error("les instructions CloudOps Agent devaient être transmises")
	}
	if !strings.Contains(llm.requests[0].Instructions, "exactly one service") {
		t.Error("les instructions devaient interdire une action sur une recherche ambiguë")
	}
	if !strings.Contains(llm.requests[0].Instructions, "0.3 is 30%") {
		t.Error("les instructions devaient préciser la conversion du ratio CPU en pourcentage")
	}
	if len(llm.requests[0].Tools) != 6 {
		t.Errorf("six outils MCP attendus, %d reçus", len(llm.requests[0].Tools))
	}
	if len(llm.requests[0].Messages) != 1 ||
		llm.requests[0].Messages[0].Content != "Que peux-tu faire ?" {
		t.Error("le message courant devait être transmis au LLM")
	}
}

func TestChatWithAgentUseCase_Execute_TransmetHistoriqueConversation(t *testing.T) {
	llm := &scriptedAgentLLM{responses: []ports.AgentResponse{{
		Message: "Le service d'authentification va être redémarré.",
	}}}
	provider := &fakeToolProvider{tools: testToolDefinitions()}
	useCase := NewChatWithAgentUseCase(llm, provider)
	history := []ports.AgentMessage{
		{
			Role:    ports.AgentMessageRoleUser,
			Content: "Affiche les logs du service d'authentification.",
		},
		{
			Role:    ports.AgentMessageRoleAssistant,
			Content: "Voici les logs de OVH-SERVICE-004.",
		},
	}

	_, err := useCase.Execute(context.Background(), "Redémarre-le encore.", history)
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}

	messages := llm.requests[0].Messages
	if len(messages) != 3 {
		t.Fatalf("trois messages attendus, %d reçus", len(messages))
	}
	if messages[1].Role != ports.AgentMessageRoleAssistant ||
		messages[1].Content != "Voici les logs de OVH-SERVICE-004." {
		t.Error("la réponse précédente devait être conservée dans le contexte")
	}
	if messages[2].Role != ports.AgentMessageRoleUser ||
		messages[2].Content != "Redémarre-le encore." {
		t.Error("le nouveau message devait terminer la conversation")
	}
}

func TestChatWithAgentUseCase_Execute_BorneHistoriqueConversation(t *testing.T) {
	llm := &scriptedAgentLLM{responses: []ports.AgentResponse{{Message: "Réponse."}}}
	useCase := NewChatWithAgentUseCase(llm, &fakeToolProvider{tools: testToolDefinitions()})
	history := make([]ports.AgentMessage, maxConversationMessages+2)
	for index := range history {
		history[index] = ports.AgentMessage{
			Role:    ports.AgentMessageRoleUser,
			Content: fmt.Sprintf("message-%d", index),
		}
	}

	_, err := useCase.Execute(context.Background(), "nouveau message", history)
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}

	messages := llm.requests[0].Messages
	if len(messages) != maxConversationMessages+1 {
		t.Fatalf("historique borné attendu à %d messages, %d reçus", maxConversationMessages+1, len(messages))
	}
	if messages[0].Content != "message-2" {
		t.Errorf("le plus ancien message conservé est inattendu : %q", messages[0].Content)
	}
}

func TestChatWithAgentUseCase_Execute_MessageVide(t *testing.T) {
	provider := &fakeToolProvider{}
	useCase := NewChatWithAgentUseCase(&scriptedAgentLLM{}, provider)

	_, err := useCase.Execute(context.Background(), "   ", nil)
	if err == nil {
		t.Fatal("une erreur était attendue pour un message vide")
	}
	if provider.listCalls != 0 {
		t.Error("aucun outil ne devait être découvert pour un message vide")
	}
}

func TestChatWithAgentUseCase_Execute_TransmetLeResultatOutil(t *testing.T) {
	const authoritativeOutput = `{"id":"OVH-SERVICE-003","status":"running","cpuUsage":0.65}`
	llm := &scriptedAgentLLM{responses: []ports.AgentResponse{
		{
			ID: "response-1",
			ToolCalls: []ports.ToolCall{{
				ID:        "call-1",
				Name:      "get_cloud_service",
				Arguments: `{"id":"OVH-SERVICE-003"}`,
			}},
		},
		{Message: "OVH-SERVICE-003 est en ligne avec 65 % de CPU."},
	}}
	provider := &fakeToolProvider{
		tools:   testToolDefinitions(),
		outputs: map[string]string{"get_cloud_service": authoritativeOutput},
	}
	useCase := NewChatWithAgentUseCase(llm, provider)

	response, err := useCase.Execute(
		context.Background(),
		"Quel est le statut de OVH-SERVICE-003 ?",
		nil,
	)
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}

	if response != "OVH-SERVICE-003 est en ligne avec 65 % de CPU." {
		t.Errorf("réponse finale inattendue : %q", response)
	}
	if len(provider.calls) != 1 || provider.calls[0].Name != "get_cloud_service" {
		t.Fatalf("l'outil de lecture du service devait être appelé")
	}
	if len(llm.requests) != 2 || len(llm.requests[1].ToolOutputs) != 1 {
		t.Fatalf("le résultat d'outil devait être renvoyé au LLM")
	}
	if llm.requests[1].ToolOutputs[0].Output != authoritativeOutput {
		t.Errorf("résultat transmis inattendu : %s", llm.requests[1].ToolOutputs[0].Output)
	}
}

func TestChatWithAgentUseCase_Execute_RedemarrageExplicite(t *testing.T) {
	llm := &scriptedAgentLLM{responses: []ports.AgentResponse{
		{
			ID: "response-1",
			ToolCalls: []ports.ToolCall{{
				ID:        "call-1",
				Name:      "restart_cloud_service",
				Arguments: `{"id":"OVH-SERVICE-003"}`,
			}},
		},
		{Message: "OVH-SERVICE-003 a été redémarré avec succès."},
	}}
	provider := &fakeToolProvider{
		tools: testToolDefinitions(),
		outputs: map[string]string{
			"restart_cloud_service": `{"id":"OVH-SERVICE-003","status":"running"}`,
		},
	}
	useCase := NewChatWithAgentUseCase(llm, provider)

	_, err := useCase.Execute(context.Background(), "Redémarre OVH-SERVICE-003.", nil)
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}
	if len(provider.calls) != 1 || provider.calls[0].Name != "restart_cloud_service" {
		t.Fatal("l'outil de redémarrage devait être exécuté")
	}
}

func TestChatWithAgentUseCase_Execute_RecherchePuisRedemarreParNom(t *testing.T) {
	llm := &scriptedAgentLLM{responses: []ports.AgentResponse{
		{
			ID: "response-1",
			ToolCalls: []ports.ToolCall{{
				ID:        "call-find",
				Name:      "find_cloud_services",
				Arguments: `{"query":"backend"}`,
			}},
		},
		{
			ID: "response-2",
			ToolCalls: []ports.ToolCall{{
				ID:        "call-restart",
				Name:      "restart_cloud_service",
				Arguments: `{"id":"OVH-SERVICE-001"}`,
			}},
		},
		{Message: "Le service backend a été redémarré avec succès."},
	}}
	provider := &fakeToolProvider{
		tools: testToolDefinitions(),
		outputs: map[string]string{
			"find_cloud_services":   `[{"id":"OVH-SERVICE-001","name":"Backend service"}]`,
			"restart_cloud_service": `{"id":"OVH-SERVICE-001","status":"running"}`,
		},
	}
	useCase := NewChatWithAgentUseCase(llm, provider)

	response, err := useCase.Execute(context.Background(), "Redémarre le service backend.", nil)
	if err != nil {
		t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
	}
	if response != "Le service backend a été redémarré avec succès." {
		t.Errorf("réponse inattendue : %q", response)
	}
	if len(provider.calls) != 2 {
		t.Fatalf("deux appels d'outils attendus, %d reçus", len(provider.calls))
	}
	if provider.calls[0].Name != "find_cloud_services" {
		t.Errorf("la recherche devait être appelée en premier")
	}
	if provider.calls[1].Name != "restart_cloud_service" {
		t.Errorf("le redémarrage devait être appelé après la recherche")
	}
	if llm.requests[2].ToolOutputs[0].Output != provider.outputs["restart_cloud_service"] {
		t.Fatal("le résultat du redémarrage devait être transmis au LLM")
	}
}

func TestChatWithAgentUseCase_Execute_ArretEtDemarrageExplicites(t *testing.T) {
	tests := []struct {
		nom      string
		message  string
		outil    string
		resultat string
		reponse  string
	}{
		{
			nom:      "arrête un service",
			message:  "Arrête OVH-SERVICE-002.",
			outil:    "shutdown_cloud_service",
			resultat: `{"id":"OVH-SERVICE-002","status":"down","cpuUsage":0}`,
			reponse:  "OVH-SERVICE-002 a été arrêté.",
		},
		{
			nom:      "démarre un service",
			message:  "Démarre OVH-SERVICE-002.",
			outil:    "start_cloud_service",
			resultat: `{"id":"OVH-SERVICE-002","status":"running","cpuUsage":0.15}`,
			reponse:  "OVH-SERVICE-002 a été démarré.",
		},
	}

	for _, test := range tests {
		t.Run(test.nom, func(t *testing.T) {
			llm := &scriptedAgentLLM{responses: []ports.AgentResponse{
				{
					ID: "response-1",
					ToolCalls: []ports.ToolCall{{
						ID:        "call-1",
						Name:      test.outil,
						Arguments: `{"id":"OVH-SERVICE-002"}`,
					}},
				},
				{Message: test.reponse},
			}}
			provider := &fakeToolProvider{
				tools:   testToolDefinitions(),
				outputs: map[string]string{test.outil: test.resultat},
			}

			response, err := NewChatWithAgentUseCase(llm, provider).Execute(
				context.Background(),
				test.message,
				nil,
			)
			if err != nil {
				t.Fatalf("aucune erreur attendue, erreur reçue : %v", err)
			}
			if response != test.reponse {
				t.Errorf("réponse attendue %q, réponse reçue %q", test.reponse, response)
			}
			if len(provider.calls) != 1 || provider.calls[0].Name != test.outil {
				t.Errorf("l'outil %q devait être appelé", test.outil)
			}
		})
	}
}

func TestChatWithAgentUseCase_Execute_ExpliqueErreurMetierOutil(t *testing.T) {
	const toolError = `{"error":"cloud service with id \\"003\\" not found"}`
	llm := &scriptedAgentLLM{responses: []ports.AgentResponse{
		{
			ID: "response-1",
			ToolCalls: []ports.ToolCall{{
				ID:        "call-1",
				Name:      "get_cloud_service",
				Arguments: `{"id":"003"}`,
			}},
		},
		{Message: "Je ne trouve aucun service avec cet identifiant."},
	}}
	provider := &fakeToolProvider{
		tools: testToolDefinitions(),
		outputs: map[string]string{
			"get_cloud_service": toolError,
		},
	}
	useCase := NewChatWithAgentUseCase(llm, provider)

	response, err := useCase.Execute(context.Background(), "Quel est le statut du service 003 ?", nil)
	if err != nil {
		t.Fatalf("l'erreur métier ne devait pas interrompre le chat : %v", err)
	}
	if response != "Je ne trouve aucun service avec cet identifiant." {
		t.Errorf("réponse inattendue : %q", response)
	}
	if llm.requests[1].ToolOutputs[0].Output != toolError {
		t.Fatal("l'erreur métier devait être transmise au LLM comme résultat d'outil")
	}
}

func TestChatWithAgentUseCase_Execute_BoucleBornee(t *testing.T) {
	responses := make([]ports.AgentResponse, maxAgentRounds)
	for index := range responses {
		responses[index] = ports.AgentResponse{
			ID: fmt.Sprintf("response-%d", index),
			ToolCalls: []ports.ToolCall{{
				ID:        fmt.Sprintf("call-%d", index),
				Name:      "get_all_cloud_services",
				Arguments: `{}`,
			}},
		}
	}

	llm := &scriptedAgentLLM{responses: responses}
	provider := &fakeToolProvider{
		tools:   testToolDefinitions(),
		outputs: map[string]string{"get_all_cloud_services": `[]`},
	}
	useCase := NewChatWithAgentUseCase(llm, provider)

	_, err := useCase.Execute(context.Background(), "Quels services sont disponibles ?", nil)
	if err == nil || !strings.Contains(err.Error(), "maximum number") {
		t.Fatalf("une erreur de boucle bornée était attendue, erreur reçue : %v", err)
	}
	if len(llm.requests) != maxAgentRounds {
		t.Errorf("nombre de tours attendu %d, nombre reçu %d", maxAgentRounds, len(llm.requests))
	}
}

func TestChatWithAgentUseCase_Execute_OutilInconnu(t *testing.T) {
	llm := &scriptedAgentLLM{responses: []ports.AgentResponse{{
		ID: "response-1",
		ToolCalls: []ports.ToolCall{{
			ID:        "call-1",
			Name:      "outil_inconnu",
			Arguments: `{}`,
		}},
	}}}
	provider := &fakeToolProvider{
		tools:   testToolDefinitions(),
		outputs: map[string]string{},
	}
	useCase := NewChatWithAgentUseCase(llm, provider)

	_, err := useCase.Execute(context.Background(), "Test", nil)
	if err == nil || !strings.Contains(err.Error(), "outil inconnu") {
		t.Fatalf("une erreur sûre était attendue, erreur reçue : %v", err)
	}
}
