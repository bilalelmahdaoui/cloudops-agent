package domain

import "time"

type CloudServiceStatus string

const (
	CloudServiceStatusRunning CloudServiceStatus = "running"
	CloudServiceStatusRestarting CloudServiceStatus = "restarting"
	CloudServiceStatusDown    CloudServiceStatus = "down"
)

type CloudService struct {
	ID       string
	Name     string
	Status   CloudServiceStatus
	CPUUsage float64
	Logs     []CloudServiceLog
}

type CloudServiceLog struct {
	DateTime time.Time
	Event    string
}
