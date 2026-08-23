package application

import (
	"context"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
)

type GetAllCloudServicesUseCase struct {
	repository ports.CloudServiceRepository
}

func NewGetAllCloudServicesUseCase(repository ports.CloudServiceRepository) *GetAllCloudServicesUseCase {
	return &GetAllCloudServicesUseCase{
		repository: repository,
	}
}

func (e *GetAllCloudServicesUseCase) Execute(ctx context.Context) ([]domain.CloudService, error) {
	return e.repository.GetAll(ctx)
}
