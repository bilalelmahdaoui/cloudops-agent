package cloud

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bilalelmahdaoui/cloudops-agent-backend/internal/domain"
)

type FakeCloudAdapter struct {
	mu            sync.RWMutex
	cloudServices [5]domain.CloudService
	restarts      map[int]chan struct{}
	restartDelay  time.Duration
}

func NewFakeCloudAdapter() *FakeCloudAdapter {
	return &FakeCloudAdapter{
		restarts:     make(map[int]chan struct{}),
		restartDelay: 5 * time.Second,
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

func (e *FakeCloudAdapter) GetAll(ctx context.Context) ([]domain.CloudService, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	cloudServices := make([]domain.CloudService, len(e.cloudServices))
	for index, service := range e.cloudServices {
		cloudServices[index] = cloneCloudService(service)
	}

	return cloudServices, nil
}

func (e *FakeCloudAdapter) GetByID(ctx context.Context, id string) (domain.CloudService, error) {
	if err := ctx.Err(); err != nil {
		return domain.CloudService{}, err
	}
	if id == "" {
		return domain.CloudService{}, fmt.Errorf("ID is required")
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, service := range e.cloudServices {
		if service.ID == id {
			return cloneCloudService(service), nil
		}
	}

	return domain.CloudService{}, fmt.Errorf("cloud service with id %q not found", id)
}

func (e *FakeCloudAdapter) Search(
	ctx context.Context,
	query string,
	limit int,
) ([]domain.CloudService, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	if normalizedQuery == "" {
		return nil, fmt.Errorf("search query is required")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("search limit must be positive")
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	services := make([]domain.CloudService, 0)
	for _, service := range e.cloudServices {
		if strings.Contains(strings.ToLower(service.Name), normalizedQuery) ||
			strings.Contains(strings.ToLower(service.ID), normalizedQuery) {
			services = append(services, cloneCloudService(service))
			if len(services) == limit {
				break
			}
		}
	}

	return services, nil
}

func (e *FakeCloudAdapter) Restart(ctx context.Context, id string) (domain.CloudService, error) {
	if err := ctx.Err(); err != nil {
		return domain.CloudService{}, err
	}
	if id == "" {
		return domain.CloudService{}, fmt.Errorf("ID is required")
	}

	e.mu.Lock()

	cloudServiceIdx := -1
	for idx, service := range e.cloudServices {
		if service.ID == id {
			cloudServiceIdx = idx
			break
		}
	}

	if cloudServiceIdx == -1 {
		e.mu.Unlock()
		return domain.CloudService{}, fmt.Errorf("cloud service with id %q not found", id)
	}

	done, alreadyRestarting := e.restarts[cloudServiceIdx]
	if !alreadyRestarting {
		done = make(chan struct{})
		e.restarts[cloudServiceIdx] = done
		e.cloudServices[cloudServiceIdx].Status = domain.CloudServiceStatusRestarting
		e.cloudServices[cloudServiceIdx].Logs = append(
			e.cloudServices[cloudServiceIdx].Logs,
			domain.CloudServiceLog{
				DateTime: time.Now(),
				Event:    "Server restarting...",
			},
		)

		go e.completeRestart(cloudServiceIdx, done)
	}
	e.mu.Unlock()

	select {
	case <-ctx.Done():
		return domain.CloudService{}, ctx.Err()
	case <-done:
	}

	return e.GetByID(context.Background(), id)
}

func (e *FakeCloudAdapter) completeRestart(index int, done chan struct{}) {
	timer := time.NewTimer(e.restartDelay)
	defer timer.Stop()
	<-timer.C

	e.mu.Lock()
	e.cloudServices[index].Status = domain.CloudServiceStatusRunning
	e.cloudServices[index].Logs = append(
		e.cloudServices[index].Logs,
		domain.CloudServiceLog{
			DateTime: time.Now(),
			Event:    "Server started successfully",
		},
	)
	delete(e.restarts, index)
	close(done)
	e.mu.Unlock()
}

func cloneCloudService(service domain.CloudService) domain.CloudService {
	service.Logs = append([]domain.CloudServiceLog(nil), service.Logs...)
	return service
}
