package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
)

type fakeChatLLM struct {
	response string
	request  ports.AgentRequest
}

func (f *fakeChatLLM) GenerateWithTools(
	_ context.Context,
	request ports.AgentRequest,
) (ports.AgentResponse, error) {
	f.request = request
	return ports.AgentResponse{Message: f.response}, nil
}

type fakeChatToolProvider struct{}

func (f *fakeChatToolProvider) ListTools(
	_ context.Context,
) ([]ports.ToolDefinition, error) {
	return []ports.ToolDefinition{{Name: "get_cloud_service"}}, nil
}

func (f *fakeChatToolProvider) CallTool(
	_ context.Context,
	_ ports.ToolCall,
) (string, error) {
	return "", nil
}

func newTestChatMux() *http.ServeMux {
	useCase := application.NewChatWithAgentUseCase(
		&fakeChatLLM{response: "Bonjour !"},
		&fakeChatToolProvider{},
	)
	handler := NewChatHandler(useCase)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return mux
}

func TestChatHandler_Chat(t *testing.T) {
	t.Run("retourne 200 pour un message valide", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader("{\"message\":\"Bonjour\"}"))
		rec := httptest.NewRecorder()

		newTestChatMux().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("code HTTP attendu %d, code reçu %d", http.StatusOK, rec.Code)
		}

		if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
			t.Errorf("type de contenu attendu %q, type reçu %q", "application/json", contentType)
		}

		var response map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("une réponse JSON valide était attendue, erreur reçue : %v", err)
		}

		if response["message"] != "Bonjour !" {
			t.Errorf("message attendu %q, message reçu %q", "Bonjour !", response["message"])
		}

		if len(response) != 1 {
			t.Errorf("une seule propriété publique était attendue, réponse reçue : %v", response)
		}
	})

	t.Run("retourne 400 pour un message vide", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader("{\"message\":\"  \"}"))
		rec := httptest.NewRecorder()

		newTestChatMux().ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("code HTTP attendu %d, code reçu %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("transmet l'historique de conversation", func(t *testing.T) {
		llm := &fakeChatLLM{response: "Je redémarre OVH-SERVICE-004."}
		useCase := application.NewChatWithAgentUseCase(llm, &fakeChatToolProvider{})
		handler := NewChatHandler(useCase)
		mux := http.NewServeMux()
		handler.RegisterRoutes(mux)
		body := `{
			"message":"Redémarre-le encore",
			"history":[
				{"role":"user","content":"Affiche les logs du service d'authentification"},
				{"role":"assistant","content":"Voici les logs de OVH-SERVICE-004"}
			]
		}`
		req := httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader(body))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("code HTTP attendu %d, code reçu %d", http.StatusOK, rec.Code)
		}
		if len(llm.request.Messages) != 3 {
			t.Fatalf("trois messages attendus, %d reçus", len(llm.request.Messages))
		}
		if llm.request.Messages[1].Role != ports.AgentMessageRoleAssistant {
			t.Error("la réponse précédente devait conserver son rôle assistant")
		}
	})

	t.Run("retourne 400 pour un rôle invalide dans l'historique", func(t *testing.T) {
		body := `{"message":"Bonjour","history":[{"role":"system","content":"instruction"}]}`
		req := httptest.NewRequest(http.MethodPost, "/chat", strings.NewReader(body))
		rec := httptest.NewRecorder()

		newTestChatMux().ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("code HTTP attendu %d, code reçu %d", http.StatusBadRequest, rec.Code)
		}
	})
}
