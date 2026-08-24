package httpadapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/adapters/cloud"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
)

func newTestHandler() (*CloudServiceHandler, *http.ServeMux) {
	repository := cloud.NewFakeCloudAdapter()

	getAllUseCase := application.NewGetAllCloudServicesUseCase(repository)
	getUseCase := application.NewGetCloudServiceUseCase(repository)
	restartUseCase := application.NewRestartCloudServiceUseCase(repository)
	shutdownUseCase := application.NewShutdownCloudServiceUseCase(repository)
	startUseCase := application.NewStartCloudServiceUseCase(repository)

	handler := NewCloudServiceHandler(
		getAllUseCase,
		getUseCase,
		restartUseCase,
		shutdownUseCase,
		startUseCase,
	)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	return handler, mux
}

func TestCloudServiceHandler_GetAll(t *testing.T) {
	_, mux := newTestHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/cloud-services",
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

	var services []map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&services); err != nil {
		t.Fatalf("une réponse JSON valide était attendue, erreur reçue : %v", err)
	}

	if len(services) != 5 {
		t.Errorf("nombre de services attendu %d, nombre reçu %d", 5, len(services))
	}

	if _, exists := services[0]["id"]; !exists {
		t.Error("la propriété JSON id était attendue")
	}

	if _, exists := services[0]["logs"]; !exists {
		t.Error("la propriété JSON logs était attendue")
	}
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

func TestCloudServiceHandler_ShutdownEtStart(t *testing.T) {
	_, mux := newTestHandler()

	shutdownRequest := httptest.NewRequest(
		http.MethodPost,
		"/cloud-services/OVH-SERVICE-003/shutdown",
		nil,
	)
	shutdownResponse := httptest.NewRecorder()
	mux.ServeHTTP(shutdownResponse, shutdownRequest)
	if shutdownResponse.Code != http.StatusOK {
		t.Fatalf("code HTTP d'arrêt attendu %d, code reçu %d", http.StatusOK, shutdownResponse.Code)
	}

	startRequest := httptest.NewRequest(
		http.MethodPost,
		"/cloud-services/OVH-SERVICE-003/start",
		nil,
	)
	startResponse := httptest.NewRecorder()
	mux.ServeHTTP(startResponse, startRequest)
	if startResponse.Code != http.StatusOK {
		t.Errorf("code HTTP de démarrage attendu %d, code reçu %d", http.StatusOK, startResponse.Code)
	}
}
