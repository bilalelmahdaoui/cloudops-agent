package cloud

import (
	"context"
	"fmt"
	"time"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
)

type FakeCloudAdapter struct {
	cloudServices [5]domain.CloudService
}

func NewFakeCloudAdapter() *FakeCloudAdapter {
	return &FakeCloudAdapter{
		cloudServices: [5]domain.CloudService{
			{
				ID:       "OVH-SERVICE-001",
				Name:     "Backend service",
				Status:   domain.CloudServiceStatusRunning,
				CPUUsage: 0.30,
				Logs: []domain.CloudServiceLog{
					{
						DateTime: time.Date(2026, 8, 22, 10, 0, 0, 0, time.UTC),
						Event:    "Server started successfully",
					},
					{
						DateTime: time.Date(2026, 8, 22, 10, 15, 0, 0, time.UTC),
						Event:    "Database connection established",
					},
				},
			},
			{
				ID:       "OVH-SERVICE-002",
				Name:     "Frontend service",
				Status:   domain.CloudServiceStatusRunning,
				CPUUsage: 0.15,
				Logs: []domain.CloudServiceLog{
					{
						DateTime: time.Date(2026, 8, 22, 9, 30, 0, 0, time.UTC),
						Event:    "Frontend service started",
					},
					{
						DateTime: time.Date(2026, 8, 22, 9, 45, 0, 0, time.UTC),
						Event:    "Static assets loaded",
					},
				},
			},
			{
				ID:       "OVH-SERVICE-003",
				Name:     "Database service",
				Status:   domain.CloudServiceStatusRunning,
				CPUUsage: 0.65,
				Logs: []domain.CloudServiceLog{
					{
						DateTime: time.Date(2026, 8, 22, 8, 0, 0, 0, time.UTC),
						Event:    "Database initialized",
					},
					{
						DateTime: time.Date(2026, 8, 22, 11, 20, 0, 0, time.UTC),
						Event:    "Backup completed successfully",
					},
				},
			},
			{
				ID:       "OVH-SERVICE-004",
				Name:     "Authentication service",
				Status:   domain.CloudServiceStatusRunning,
				CPUUsage: 0.42,
				Logs: []domain.CloudServiceLog{
					{
						DateTime: time.Date(2026, 8, 22, 7, 45, 0, 0, time.UTC),
						Event:    "Authentication service started",
					},
					{
						DateTime: time.Date(2026, 8, 22, 12, 5, 0, 0, time.UTC),
						Event:    "Token refresh completed",
					},
				},
			},
			{
				ID:       "OVH-SERVICE-005",
				Name:     "Worker service",
				Status:   domain.CloudServiceStatusRunning,
				CPUUsage: 0.82,
				Logs: []domain.CloudServiceLog{
					{
						DateTime: time.Date(2026, 8, 22, 6, 30, 0, 0, time.UTC),
						Event:    "Worker service started",
					},
					{
						DateTime: time.Date(2026, 8, 22, 13, 10, 0, 0, time.UTC),
						Event:    "Processed 250 queued jobs",
					},
				},
			},
		}}
}

func (e *FakeCloudAdapter) GetByID(ctx context.Context, id string) (domain.CloudService, error) {
	if err := ctx.Err(); err != nil {
		return domain.CloudService{}, err
	}
	if id == "" {
		return domain.CloudService{}, fmt.Errorf("ID is required")
	}

	for _, service := range e.cloudServices {
		if service.ID == id {
			return service, nil
		}
	}

	return domain.CloudService{}, fmt.Errorf("cloud service with id %q not found", id)
}

func (e *FakeCloudAdapter) Restart(ctx context.Context, id string) (domain.CloudService, error) {
	if err := ctx.Err(); err != nil {
		return domain.CloudService{}, err
	}
	if id == "" {
		return domain.CloudService{}, fmt.Errorf("ID is required")
	}

	cloudServiceIdx := -1

	for idx, service := range e.cloudServices {
		if service.ID == id {
			cloudServiceIdx = idx
			break
		}
	}

	if cloudServiceIdx == -1 {
		return domain.CloudService{}, fmt.Errorf("cloud service with id %q not found", id)
	}

	e.cloudServices[cloudServiceIdx].Status = domain.CloudServiceStatusRestarting
	e.cloudServices[cloudServiceIdx].Logs = append(
		e.cloudServices[cloudServiceIdx].Logs,
		domain.CloudServiceLog{
			DateTime: time.Now(),
			Event:    "Server restarting...",
		},
	)

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return domain.CloudService{}, ctx.Err()

	case <-timer.C:
		e.cloudServices[cloudServiceIdx].Status = domain.CloudServiceStatusRunning
		e.cloudServices[cloudServiceIdx].Logs = append(e.cloudServices[cloudServiceIdx].Logs, domain.CloudServiceLog{
			DateTime: time.Now(),
			Event:    "Server started successfully",
		})
	}

	return e.cloudServices[cloudServiceIdx], nil
}
