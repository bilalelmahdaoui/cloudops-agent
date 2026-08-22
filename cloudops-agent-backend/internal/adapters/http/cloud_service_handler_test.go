package httpadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
)

func TestCloudServiceHandler_GetByID(t *testing.T) {
	repository := cloud.NewFakeCloudAdapter()
	useCase := application.NewGetCloudServiceUseCase(repository)
	handler := NewCloudServiceHandler(useCase)

	t.Run("retourne 200 si le service existe", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/cloud-services/OVH-SERVICE-003",
			nil,
		)

		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("code HTTP attendu %d, code reçu %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("retourne 404 si le service n'existe pas", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/cloud-services/UNKNOWN",
			nil,
		)

		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf(
				"code HTTP attendu %d, code reçu %d",
				http.StatusNotFound,
				rec.Code,
			)
		}
	})
}
