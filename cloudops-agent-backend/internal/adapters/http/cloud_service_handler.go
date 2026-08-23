package httpadapter

import (
	"encoding/json"
	"net/http"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
)

type CloudServiceHandler struct {
	getCloudServiceUseCase     *application.GetCloudServiceUseCase
	restartCloudServiceUseCase *application.RestartCloudServiceUseCase
}

func NewCloudServiceHandler(
	getCloudServiceUseCase *application.GetCloudServiceUseCase,
	restartCloudServiceUseCase *application.RestartCloudServiceUseCase,
) *CloudServiceHandler {
	return &CloudServiceHandler{
		getCloudServiceUseCase:     getCloudServiceUseCase,
		restartCloudServiceUseCase: restartCloudServiceUseCase,
	}
}

func (h *CloudServiceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /cloud-services/{id}", h.GetByID)
	mux.HandleFunc("GET /cloud-services/{id}/restart", h.Restart)
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
