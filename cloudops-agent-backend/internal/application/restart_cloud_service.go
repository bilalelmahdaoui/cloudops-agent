package application

import (
	"context"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
)

type RestartCloudServiceUseCase struct {
	repository ports.CloudServiceRepository
}

func NewRestartCloudServiceUseCase(repository ports.CloudServiceRepository) *RestartCloudServiceUseCase {
	return &RestartCloudServiceUseCase{
		repository: repository,
	}
}

func (e *RestartCloudServiceUseCase) Execute(ctx context.Context, id string) (domain.CloudService, error) {
	return e.repository.Restart(ctx, id)
}
