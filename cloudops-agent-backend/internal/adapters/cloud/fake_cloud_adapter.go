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
	nominalCPU    [5]float64
	restarts      map[int]chan struct{}
	restartDelay  time.Duration
	shutdownDelay time.Duration
	startDelay    time.Duration
}

func NewFakeCloudAdapter() *FakeCloudAdapter {
	adapter := &FakeCloudAdapter{
		restarts:      make(map[int]chan struct{}),
		restartDelay:  4 * time.Second,
		shutdownDelay: 2 * time.Second,
		startDelay:    2 * time.Second,
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

	for index, service := range adapter.cloudServices {
		adapter.nominalCPU[index] = service.CPUUsage
	}

	return adapter
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
	if e.cloudServices[cloudServiceIdx].Status == domain.CloudServiceStatusDown {
		e.mu.Unlock()
		return domain.CloudService{}, fmt.Errorf("cannot restart stopped service %q", id)
	}

	done, alreadyRestarting := e.restarts[cloudServiceIdx]
	if !alreadyRestarting {
		done = make(chan struct{})
		e.restarts[cloudServiceIdx] = done
		e.nominalCPU[cloudServiceIdx] = e.cloudServices[cloudServiceIdx].CPUUsage
		e.cloudServices[cloudServiceIdx].Status = domain.CloudServiceStatusRestarting
		e.cloudServices[cloudServiceIdx].CPUUsage = 0
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
	e.cloudServices[index].CPUUsage = e.nominalCPU[index]
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

func (e *FakeCloudAdapter) Shutdown(ctx context.Context, id string) (domain.CloudService, error) {
	if err := ctx.Err(); err != nil {
		return domain.CloudService{}, err
	}
	if id == "" {
		return domain.CloudService{}, fmt.Errorf("ID is required")
	}

	e.mu.RLock()

	index := e.cloudServiceIndex(id)
	if index == -1 {
		e.mu.RUnlock()
		return domain.CloudService{}, fmt.Errorf("cloud service with ID %q not found", id)
	}
	if _, restarting := e.restarts[index]; restarting {
		e.mu.RUnlock()
		return domain.CloudService{}, fmt.Errorf("cannot stop service %q while it is restarting", id)
	}
	if e.cloudServices[index].Status == domain.CloudServiceStatusDown {
		service := cloneCloudService(e.cloudServices[index])
		e.mu.RUnlock()
		return service, nil
	}
	e.mu.RUnlock()

	if err := waitForOperation(ctx, e.shutdownDelay); err != nil {
		return domain.CloudService{}, err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.nominalCPU[index] = e.cloudServices[index].CPUUsage
	e.cloudServices[index].Status = domain.CloudServiceStatusDown
	e.cloudServices[index].CPUUsage = 0
	e.cloudServices[index].Logs = append(
		e.cloudServices[index].Logs,
		domain.CloudServiceLog{
			DateTime: time.Now(),
			Event:    "Server shut down successfully",
		},
	)

	return cloneCloudService(e.cloudServices[index]), nil
}

func (e *FakeCloudAdapter) Start(ctx context.Context, id string) (domain.CloudService, error) {
	if err := ctx.Err(); err != nil {
		return domain.CloudService{}, err
	}
	if id == "" {
		return domain.CloudService{}, fmt.Errorf("ID is required")
	}

	e.mu.RLock()

	index := e.cloudServiceIndex(id)
	if index == -1 {
		e.mu.RUnlock()
		return domain.CloudService{}, fmt.Errorf("cloud service with ID %q not found", id)
	}
	if _, restarting := e.restarts[index]; restarting {
		e.mu.RUnlock()
		return domain.CloudService{}, fmt.Errorf("cannot start service %q while it is restarting", id)
	}
	if e.cloudServices[index].Status == domain.CloudServiceStatusRunning {
		service := cloneCloudService(e.cloudServices[index])
		e.mu.RUnlock()
		return service, nil
	}
	e.mu.RUnlock()

	if err := waitForOperation(ctx, e.startDelay); err != nil {
		return domain.CloudService{}, err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.cloudServices[index].Status = domain.CloudServiceStatusRunning
	e.cloudServices[index].CPUUsage = e.nominalCPU[index]
	e.cloudServices[index].Logs = append(
		e.cloudServices[index].Logs,
		domain.CloudServiceLog{
			DateTime: time.Now(),
			Event:    "Server started successfully",
		},
	)

	return cloneCloudService(e.cloudServices[index]), nil
}

func waitForOperation(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (e *FakeCloudAdapter) cloudServiceIndex(id string) int {
	for index, service := range e.cloudServices {
		if service.ID == id {
			return index
		}
	}
	return -1
}

func cloneCloudService(service domain.CloudService) domain.CloudService {
	service.Logs = append([]domain.CloudServiceLog(nil), service.Logs...)
	return service
}
