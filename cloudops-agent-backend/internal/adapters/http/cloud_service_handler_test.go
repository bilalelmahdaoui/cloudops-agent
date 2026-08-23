package httpadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
)

func newTestHandler() (*CloudServiceHandler, *http.ServeMux) {
	repository := cloud.NewFakeCloudAdapter()

	getUseCase := application.NewGetCloudServiceUseCase(repository)
	restartUseCase := application.NewRestartCloudServiceUseCase(repository)

	handler := NewCloudServiceHandler(
		getUseCase,
		restartUseCase,
	)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return handler, mux
}

func TestCloudServiceHandler_GetByID(t *testing.T) {
	_, mux := newTestHandler()

	t.Run("retourne 200 si le service existe", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/cloud-services/OVH-SERVICE-003",
			nil,
		)

		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf(
				"code HTTP attendu %d, code reçu %d",
				http.StatusOK,
				rec.Code,
			)
		}
	})

	t.Run("retourne 404 si le service n'existe pas", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/cloud-services/UNKNOWN",
			nil,
		)

		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf(
				"code HTTP attendu %d, code reçu %d",
				http.StatusNotFound,
				rec.Code,
			)
		}
	})
}

func TestCloudServiceHandler_Restart(t *testing.T) {
	_, mux := newTestHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/cloud-services/OVH-SERVICE-003/restart",
		nil,
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf(
			"code HTTP attendu %d, code reçu %d",
			http.StatusOK,
			rec.Code,
		)
	}
}
