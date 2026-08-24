package httpadapter

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
)

type ChatHandler struct {
	chatWithAgentUseCase *application.ChatWithAgentUseCase
}

func NewChatHandler(chatWithAgentUseCase *application.ChatWithAgentUseCase) *ChatHandler {
	return &ChatHandler{chatWithAgentUseCase: chatWithAgentUseCase}
}

func (h *ChatHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /chat", h.Chat)
}

func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Message string `json:"message"`
		History []struct {
			Role    ports.AgentMessageRole `json:"role"`
			Content string                 `json:"content"`
		} `json:"history"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "corps de requête JSON invalide", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(request.Message) == "" {
		http.Error(w, "le message est requis", http.StatusBadRequest)
		return
	}

	history := make([]ports.AgentMessage, len(request.History))
	for index, historyMessage := range request.History {
		if historyMessage.Role != ports.AgentMessageRoleUser &&
			historyMessage.Role != ports.AgentMessageRoleAssistant {
			http.Error(w, "rôle invalide dans l'historique", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(historyMessage.Content) == "" {
			http.Error(w, "message vide dans l'historique", http.StatusBadRequest)
			return
		}
		history[index] = ports.AgentMessage{
			Role:    historyMessage.Role,
			Content: historyMessage.Content,
		}
	}

	message, err := h.chatWithAgentUseCase.Execute(r.Context(), request.Message, history)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		log.Printf("erreur pendant le chat : %v", err)
		http.Error(w, "impossible de générer une réponse", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{Message: message})
}
