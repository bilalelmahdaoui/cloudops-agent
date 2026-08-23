package httpadapter

import (
	"encoding/json"
	"net/http"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
)

type CloudServiceHandler struct {
	getAllCloudServicesUseCase *application.GetAllCloudServicesUseCase
	getCloudServiceUseCase     *application.GetCloudServiceUseCase
	restartCloudServiceUseCase *application.RestartCloudServiceUseCase
}

func NewCloudServiceHandler(
	getAllCloudServicesUseCase *application.GetAllCloudServicesUseCase,
	getCloudServiceUseCase *application.GetCloudServiceUseCase,
	restartCloudServiceUseCase *application.RestartCloudServiceUseCase,
) *CloudServiceHandler {
	return &CloudServiceHandler{
		getAllCloudServicesUseCase: getAllCloudServicesUseCase,
		getCloudServiceUseCase:     getCloudServiceUseCase,
		restartCloudServiceUseCase: restartCloudServiceUseCase,
	}
}

func (h *CloudServiceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /cloud-services", h.GetAll)
	mux.HandleFunc("GET /cloud-services/{id}", h.GetByID)
	mux.HandleFunc("POST /cloud-services/{id}/restart", h.Restart)
}

func (h *CloudServiceHandler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {
	cloudServices, err := h.getAllCloudServicesUseCase.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cloudServices)
}

func (h *CloudServiceHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")

	cloudService, err := h.getCloudServiceUseCase.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cloudService)
}

func (h *CloudServiceHandler) Restart(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := r.PathValue("id")

	cloudService, err := h.restartCloudServiceUseCase.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cloudService)
}
