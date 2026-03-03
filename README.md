# Sentinel CLI

A command-line monitoring tool for HPC clusters that fetches metrics from Prometheus for nodes and services (Kafka, OpenSearch, Logstash).

## Installation

**Prerequisites:** Go 1.24.4+, Prometheus with node_exporter

```bash
# Clone and build
git clone https://github.com/vinyas-bharadwaj/sentinel.git
cd sentinel
make build

# Or install to $GOPATH/bin
make install
```

The Makefile uses `go build` with version flags embedded via ldflags.

## Features

### Node Monitoring
Monitor CPU, memory, and disk usage across HPC cluster nodes with configurable time ranges.

```bash
sentinel node [-t time_range]
```

- `-t, --history`: Lookback period (default: `5m`). Examples: `30s`, `1h`, `1d`
- `--url`: Prometheus URL (default: `http://localhost:9090`)

**Output:** Formatted table with NODE, CPU (%), MEM (%), DISK (%)

### Service Monitoring  
Monitor Kafka, OpenSearch, or Logstash with health status indicators.

```bash
sentinel service [kafka|opensearch|logstash] [-w]
```

- `-w, --watch`: Auto-refresh every 5 seconds
- `--url`: Prometheus URL (default: `http://localhost:9090`)

**Services:**
- **Kafka**: Consumer group lag (warns if >5000 messages)
- **OpenSearch**: Cluster status (Green=0, Yellow=1, Red=2)
- **Logstash**: Event processing metrics

**Output:** Color-coded table with ENTITY, VALUE, UNIT, STATUS, DESCRIPTION

## Usage Examples

```bash
# Monitor nodes (last 5 minutes)
sentinel node

# Monitor nodes with 1 hour lookback
sentinel node -t 1h

# Check Kafka consumer lag
sentinel service kafka

# Watch OpenSearch health with auto-refresh
sentinel service opensearch -w

# Use custom Prometheus URL
sentinel node --url http://prometheus.example.com:9090
```

## Project Structure

```
├── main.go                 # Entry point
├── makefile                # Build automation (build, install)
├── cmd/                    
│   ├── root.go            # Root command & global flags
│   ├── node.go            # Node monitoring
│   └── service.go         # Service monitoring
└── internal/               
    ├── metrics.go         # Data structures
    └── prometheus.go      # Prometheus client & queries
```

## Development

**Build:** `make build` or `go build -o sentinel main.go`  
**Dependencies:** Cobra (CLI), Prometheus client, json-iterator
