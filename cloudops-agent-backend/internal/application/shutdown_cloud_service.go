package application

import (
	"context"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
)

type ShutdownCloudServiceUseCase struct {
	repository ports.CloudServiceRepository
}

func NewShutdownCloudServiceUseCase(repository ports.CloudServiceRepository) *ShutdownCloudServiceUseCase {
	return &ShutdownCloudServiceUseCase{repository: repository}
}

func (u *ShutdownCloudServiceUseCase) Execute(
	ctx context.Context,
	id string,
) (domain.CloudService, error) {
	return u.repository.Shutdown(ctx, id)
}
