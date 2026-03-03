package internal

import "time"

// ClusterStats holds the unified metrics for our HPC nodes
type ClusterStats struct {
	NodeName  string
	CPUUsage  float64
	MemUsage  float64 // Percentage used
	DiskUsage float64 // Percentage used
	Timestamp time.Time
}

// ServiceStats represents health for Kafka, Logstash, or OpenSearch
type ServiceStats struct {
	ServiceName string
	MetricName  string // e.g., "logstash-hpc-consumer"
	Value       float64
	Status      string // "HEALTHY", "DEGRADED", "CRITICAL"
	Unit        string // "messages", "events/s", "status-code"
	Description string // Human readable explanation
}

// Provider defines the interface for fetching data.
// This makes it easy to mock for tests later!
type Provider interface {
	GetNodeStats(lookback string) ([]ClusterStats, error)
}
