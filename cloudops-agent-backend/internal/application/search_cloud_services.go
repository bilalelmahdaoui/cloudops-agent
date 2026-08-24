package application

import (
	"context"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/ports"
)

const maxCloudServiceSearchResults = 10

type SearchCloudServicesUseCase struct {
	repository ports.CloudServiceRepository
}

func NewSearchCloudServicesUseCase(repository ports.CloudServiceRepository) *SearchCloudServicesUseCase {
	return &SearchCloudServicesUseCase{repository: repository}
}

func (u *SearchCloudServicesUseCase) Execute(
	ctx context.Context,
	query string,
) ([]domain.CloudService, error) {
	return u.repository.Search(ctx, query, maxCloudServiceSearchResults)
}
