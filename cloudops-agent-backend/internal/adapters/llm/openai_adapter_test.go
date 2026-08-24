package llm

import (
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
	"github.com/openai/openai-go/v3/responses"
)

func TestOpenAIToolParameters(t *testing.T) {
	t.Run("complète le schéma d'un objet vide", func(t *testing.T) {
		source := map[string]any{
			"type":                 "object",
			"additionalProperties": false,
		}

		parameters := openAIToolParameters(source)

		if _, exists := parameters["properties"]; !exists {
			t.Error("la propriété properties était attendue")
		}
		if _, exists := parameters["required"]; !exists {
			t.Error("la propriété required était attendue")
		}
		if _, mutated := source["properties"]; mutated {
			t.Error("le schéma source ne devait pas être modifié")
		}
	})

	t.Run("conserve un schéma déjà complet", func(t *testing.T) {
		properties := map[string]any{
			"id": map[string]any{"type": "string"},
		}
		required := []string{"id"}
		source := map[string]any{
			"type":       "object",
			"properties": properties,
			"required":   required,
		}

		parameters := openAIToolParameters(source)

		if parameters["properties"] == nil {
			t.Error("les propriétés existantes devaient être conservées")
		}
		if parameters["required"] == nil {
			t.Error("les champs requis existants devaient être conservés")
		}
	})
}

func TestOpenAIConversationInput(t *testing.T) {
	input := openAIConversationInput([]ports.AgentMessage{
		{Role: ports.AgentMessageRoleUser, Content: "Affiche les logs du service d'authentification."},
		{Role: ports.AgentMessageRoleAssistant, Content: "Voici les logs de OVH-SERVICE-004."},
		{Role: ports.AgentMessageRoleUser, Content: "Redémarre-le encore."},
	})

	if len(input) != 3 {
		t.Fatalf("trois messages attendus, %d reçus", len(input))
	}
	if input[0].OfMessage == nil || input[0].OfMessage.Role != responses.EasyInputMessageRoleUser {
		t.Fatal("le premier message devait avoir le rôle user")
	}
	if input[1].OfMessage == nil || input[1].OfMessage.Role != responses.EasyInputMessageRoleAssistant {
		t.Fatal("le deuxième message devait avoir le rôle assistant")
	}
	if input[2].OfMessage == nil || input[2].OfMessage.Content.OfString.Value != "Redémarre-le encore." {
		t.Fatal("le dernier message utilisateur est inattendu")
	}
}
