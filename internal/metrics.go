package internal


// NodeCPUStats holds the compute profile for an HPC node
type NodeCPUStats struct {
	NodeName    string
	Utilization float64 // CPU % usage
	Load1m      float64 // 1-minute load average
	Throttles   float64 // Core throttles in the last 5m
}

// ServiceCPUStats holds the compute profile for a specific service
type ServiceCPUStats struct {
	ServiceName string
	EntityName  string
	LoadValue   float64
	Description string
}
