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

// GetNodeCPUStats fetches a pure compute profile for all nodes
func (p *PromClient) GetNodeCPUStats(lookback string) ([]NodeCPUStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queries := map[string]string{
		"utilization": fmt.Sprintf("100 - (avg by (instance) (irate(node_cpu_seconds_total{mode='idle'}[%s])) * 100)", lookback),
		"load1m":      "node_load1",
		"throttles":   fmt.Sprintf("sum by (instance) (increase(node_cpu_core_throttles_total[%s]))", lookback),
	}

	registry := make(map[string]*NodeCPUStats)

	for metricType, query := range queries {
		result, _, err := p.api.Query(ctx, query, time.Now())
		if err != nil {
			continue // Skip if a specific metric fails, but keep processing others
		}

		vector, ok := result.(model.Vector)
		if !ok {
			continue
		}

		for _, sample := range vector {
			nodeName := string(sample.Metric["instance"])
			if _, exists := registry[nodeName]; !exists {
				registry[nodeName] = &NodeCPUStats{NodeName: nodeName}
			}

			val := float64(sample.Value)
			switch metricType {
			case "utilization":
				registry[nodeName].Utilization = val
			case "load1m":
				registry[nodeName].Load1m = val
			case "throttles":
				registry[nodeName].Throttles = val
			}
		}
	}

	var finalStats []NodeCPUStats
	for _, v := range registry {
		finalStats = append(finalStats, *v)
	}
	return finalStats, nil
}

// GetServiceCPUStats fetches compute load for specific services
func (p *PromClient) GetServiceCPUStats(service string) ([]ServiceCPUStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string
	var entityLabel string
	var description string

	switch service {
	case "opensearch":
		query = "opensearch_os_load_average_one_minute"
		entityLabel = "cluster" // or "instance" depending on your exact labels
		description = "1m OS Load Average (JVM Proxy)"
	case "kafka", "logstash":
		// Being candid: Since the exporters lack CPU metrics, we return an explicit error
		return nil, fmt.Errorf("CPU metrics are not currently exported by the %s exporter", service)
	default:
		return nil, fmt.Errorf("unsupported service: %s", service)
	}

	result, _, err := p.api.Query(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}

	vector, ok := result.(model.Vector)
	if !ok || len(vector) == 0 {
		return nil, fmt.Errorf("no CPU metrics found for %s", service)
	}

	var stats []ServiceCPUStats
	for _, sample := range vector {
		labelValue := string(sample.Metric[model.LabelName(entityLabel)])
		if labelValue == "" {
			labelValue = "cluster-wide"
		}

		stats = append(stats, ServiceCPUStats{
			ServiceName: service,
			EntityName:  labelValue,
			LoadValue:   float64(sample.Value),
			Description: description,
		})
	}
	return stats, nil
}
