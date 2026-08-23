package ports

import (
	"context"
	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
)

type CloudServiceRepository interface {
	GetByID(ctx context.Context, id string) (domain.CloudService, error)
	Restart(ctx context.Context, id string) (domain.CloudService, error)
}
