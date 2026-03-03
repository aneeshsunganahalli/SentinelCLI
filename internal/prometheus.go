package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type PromClient struct {
	api v1.API
}

func NewPromClient(address string) (*PromClient, error) {
	client, err := api.NewClient(api.Config{Address: address})
	if err != nil {
		return nil, err
	}
	return &PromClient{api: v1.NewAPI(client)}, nil
}

func (p *PromClient) GetNodeStats(lookback string) ([]ClusterStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define our PromQL queries
	queries := map[string]string{
		"cpu":  fmt.Sprintf("100 - (avg by (instance) (irate(node_cpu_seconds_total{mode='idle'}[%s])) * 100)", lookback),
		"mem":  "100 * (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes))",
		"disk": "100 * (1 - (node_filesystem_avail_bytes{mountpoint='/'} / node_filesystem_size_bytes{mountpoint='/'}))",
	}

	// This map will store our combined data, keyed by NodeName
	registry := make(map[string]*ClusterStats)

	for metricType, query := range queries {
		result, _, err := p.api.Query(ctx, query, time.Now())
		if err != nil {
			return nil, err
		}

		vector, ok := result.(model.Vector)
		if !ok {
			continue
		}

		for _, sample := range vector {
			nodeName := string(sample.Metric["instance"])

			// If node isn't in our registry yet, create it
			if _, exists := registry[nodeName]; !exists {
				registry[nodeName] = &ClusterStats{NodeName: nodeName}
			}

			// Assign the value to the correct field
			val := float64(sample.Value)
			switch metricType {
			case "cpu":
				registry[nodeName].CPUUsage = val
			case "mem":
				registry[nodeName].MemUsage = val
			case "disk":
				registry[nodeName].DiskUsage = val
			}
		}
	}

	// Convert map to slice for the UI
	var finalStats []ClusterStats
	for _, v := range registry {
		finalStats = append(finalStats, *v)
	}
	return finalStats, nil
}

func (p *PromClient) GetServiceStats(service string) ([]ServiceStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string
	var metricLabel string

	switch service {
	case "kafka":
		query = "sum(kafka_consumergroup_lag) by (consumergroup)"
		metricLabel = "consumergroup"
	case "opensearch":
		query = "opensearch_cluster_status"
		metricLabel = "cluster"
	case "logstash":
		query = "sum(logstash_stats_events_out) by (instance)"
		metricLabel = "instance"
	default:
		return nil, fmt.Errorf("unsupported service: %s", service)
	}

	result, _, err := p.api.Query(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}

	vector, ok := result.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var stats []ServiceStats
	for _, sample := range vector {
		val := float64(sample.Value)
		metricName := string(sample.Metric[model.LabelName(metricLabel)])

		// Default values
		status := "Active"
		unit := "count"
		description := "Generic Metric"

		// --- Verbosity Enrichment Layer ---
		switch service {
		case "kafka":
			unit = "messages"
			description = "Consumer group lag (backlog)"
			if val > 5000 {
				status = "Warning"
			} else {
				status = "Healthy"
			}

		case "opensearch":
			unit = "code"
			// OpenSearch Status: 0=Green, 1=Yellow, 2=Red
			statusMap := map[float64]string{0: "OK", 1: "WARN", 2: "CRIT"}
			descMap := map[float64]string{0: "Green (Healthy)", 1: "Yellow (Missing Replicas)", 2: "Red (Data Loss)"}
			status = statusMap[val]
			description = descMap[val]

		case "logstash":
			unit = "events"
			description = "Total events processed (Cumulative)"
		}

		stats = append(stats, ServiceStats{
			ServiceName: service,
			MetricName:  metricName,
			Value:       val,
			Status:      status,
			Unit:        unit,
			Description: description,
		})
	}
	return stats, nil
}
