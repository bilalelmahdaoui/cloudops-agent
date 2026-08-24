package ports

import (
	"context"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
)

type CloudServiceRepository interface {
	GetAll(ctx context.Context) ([]domain.CloudService, error)
	GetByID(ctx context.Context, id string) (domain.CloudService, error)
	Search(ctx context.Context, query string, limit int) ([]domain.CloudService, error)
	Restart(ctx context.Context, id string) (domain.CloudService, error)
	Shutdown(ctx context.Context, id string) (domain.CloudService, error)
	Start(ctx context.Context, id string) (domain.CloudService, error)
}
