package application

import (
	"context"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
)

type GetCloudServiceUseCase struct {
	repository ports.CloudServiceRepository
}

func NewGetCloudServiceUseCase(repository ports.CloudServiceRepository) *GetCloudServiceUseCase {
	return &GetCloudServiceUseCase{
		repository: repository,
	}
}

func (e *GetCloudServiceUseCase) Execute(ctx context.Context, id string) (domain.CloudService, error) {
	return e.repository.GetByID(ctx, id)
}