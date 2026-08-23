package domain

import "time"

type CloudServiceStatus string

const (
	CloudServiceStatusRunning    CloudServiceStatus = "running"
	CloudServiceStatusRestarting CloudServiceStatus = "restarting"
	CloudServiceStatusDown       CloudServiceStatus = "down"
)

type CloudService struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Status   CloudServiceStatus `json:"status"`
	CPUUsage float64            `json:"cpuUsage"`
	Logs     []CloudServiceLog  `json:"logs"`
}

type CloudServiceLog struct {
	DateTime time.Time `json:"dateTime"`
	Event    string    `json:"event"`
}
