package httpadapter

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/application"
)

type CloudServiceHandler struct {
	getCloudServiceUseCase *application.GetCloudServiceUseCase
}

func NewCloudServiceHandler(getCloudServiceUseCase *application.GetCloudServiceUseCase) *CloudServiceHandler {
	return &CloudServiceHandler{
		getCloudServiceUseCase: getCloudServiceUseCase,
	}
}

func (h *CloudServiceHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := strings.TrimPrefix(r.URL.Path, "/cloud-services/")
	cloudService, err := h.getCloudServiceUseCase.Execute(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cloudService)
}
