package httpadapter

import (
	"encoding/json"
	"net/http"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
)

type CloudServiceHandler struct {
	getAllCloudServicesUseCase  *application.GetAllCloudServicesUseCase
	getCloudServiceUseCase      *application.GetCloudServiceUseCase
	restartCloudServiceUseCase  *application.RestartCloudServiceUseCase
	shutdownCloudServiceUseCase *application.ShutdownCloudServiceUseCase
	startCloudServiceUseCase    *application.StartCloudServiceUseCase
}

func NewCloudServiceHandler(
	getAllCloudServicesUseCase *application.GetAllCloudServicesUseCase,
	getCloudServiceUseCase *application.GetCloudServiceUseCase,
	restartCloudServiceUseCase *application.RestartCloudServiceUseCase,
	shutdownCloudServiceUseCase *application.ShutdownCloudServiceUseCase,
	startCloudServiceUseCase *application.StartCloudServiceUseCase,
) *CloudServiceHandler {
	return &CloudServiceHandler{
		getAllCloudServicesUseCase:  getAllCloudServicesUseCase,
		getCloudServiceUseCase:      getCloudServiceUseCase,
		restartCloudServiceUseCase:  restartCloudServiceUseCase,
		shutdownCloudServiceUseCase: shutdownCloudServiceUseCase,
		startCloudServiceUseCase:    startCloudServiceUseCase,
	}
}

func (h *CloudServiceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /cloud-services", h.GetAll)
	mux.HandleFunc("GET /cloud-services/{id}", h.GetByID)
	mux.HandleFunc("POST /cloud-services/{id}/restart", h.Restart)
	mux.HandleFunc("POST /cloud-services/{id}/shutdown", h.Shutdown)
	mux.HandleFunc("POST /cloud-services/{id}/start", h.Start)
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

func (h *CloudServiceHandler) Shutdown(
	w http.ResponseWriter,
	r *http.Request,
) {
	cloudService, err := h.shutdownCloudServiceUseCase.Execute(r.Context(), r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cloudService)
}

func (h *CloudServiceHandler) Start(
	w http.ResponseWriter,
	r *http.Request,
) {
	cloudService, err := h.startCloudServiceUseCase.Execute(r.Context(), r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cloudService)
}
