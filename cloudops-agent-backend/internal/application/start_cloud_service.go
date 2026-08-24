package application

import (
	"context"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
)

type StartCloudServiceUseCase struct {
	repository ports.CloudServiceRepository
}

func NewStartCloudServiceUseCase(repository ports.CloudServiceRepository) *StartCloudServiceUseCase {
	return &StartCloudServiceUseCase{repository: repository}
}

func (u *StartCloudServiceUseCase) Execute(
	ctx context.Context,
	id string,
) (domain.CloudService, error) {
	return u.repository.Start(ctx, id)
}
