# Monitoring and APM

GoKit provides an Application Performance Monitoring (APM) system with Prometheus integration for metrics collection and monitoring.

## Overview

The monitoring system provides:
- **Prometheus metrics integration** for time-series data collection
- **Custom metrics configuration** for counters, gauges, histograms, and summaries
- **HTTP handler** for Prometheus scraping
- **Automatic metric registration** and collection
- **Flexible metric configuration** via YAML/JSON

## Basic Usage

### 1. Initialize APM

```go
package main

import (
    "github.com/kumarabd/gokit/apm"
)

func main() {
    // Create APM options
    opts := apm.Options{
        Prometheus: apm.MetricOptions{
            Enabled: true,
        },
    }
    
    // Initialize APM (configuration will be handled separately)
    // The actual initialization depends on your specific metrics configuration
}
```

### 2. Configure Prometheus Metrics

```go
import (
    "github.com/kumarabd/gokit/apm/prometheus"
    "github.com/prometheus/client_golang/prometheus"
)

func setupMetrics() {
    config := prometheus.Config{
        Counters: []prometheus.CounterOpts{
            {
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            {
                Name: "database_queries_total",
                Help: "Total number of database queries",
            },
        },
        Gauges: []prometheus.GaugeOpts{
            {
                Name: "active_connections",
                Help: "Number of active connections",
            },
        },
        Histograms: []prometheus.HistogramOpts{
            {
                Name:    "http_request_duration_seconds",
                Help:    "HTTP request duration in seconds",
                Buckets: prometheus.DefBuckets,
            },
        },
    }
    
    // Initialize metrics
    prometheus.InitMetrics(config)
}
```

## Prometheus Integration

### 1. HTTP Handler

The Prometheus package provides an HTTP handler for metrics scraping:

```go
import (
    "net/http"
    "github.com/kumarabd/gokit/apm/prometheus"
)

func setupMetricsEndpoint() {
    // Create HTTP handler for Prometheus metrics
    handler := prometheus.GetHTTPHandler()
    
    // Register the handler with your HTTP server
    http.Handle("/metrics", handler)
    
    // Start the server
    go http.ListenAndServe(":9090", nil)
}
```

### 2. Custom Metrics Configuration

```go
import (
    "time"
    "github.com/kumarabd/gokit/apm/prometheus"
    "github.com/prometheus/client_golang/prometheus"
)

func setupCustomMetrics() {
    config := prometheus.Config{
        Counters: []prometheus.CounterOpts{
            {
                Name: "user_registrations_total",
                Help: "Total number of user registrations",
            },
            {
                Name: "api_errors_total",
                Help: "Total number of API errors",
            },
        },
        Gauges: []prometheus.GaugeOpts{
            {
                Name: "cache_hit_ratio",
                Help: "Cache hit ratio percentage",
            },
            {
                Name: "memory_usage_bytes",
                Help: "Current memory usage in bytes",
            },
        },
        Histograms: []prometheus.HistogramOpts{
            {
                Name:    "database_query_duration_seconds",
                Help:    "Database query duration in seconds",
                Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5},
            },
        },
        Summaries: []prometheus.SummaryOpts{
            {
                Name:       "response_size_bytes",
                Help:       "Response size in bytes",
                Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
            },
        },
    }
    
    prometheus.InitMetrics(config)
}
```

## Metric Types

### 1. Counters

Counters are monotonically increasing metrics:

```go
config := prometheus.Config{
    Counters: []prometheus.CounterOpts{
        {
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        {
            Name: "errors_total",
            Help: "Total number of errors",
        },
    },
}
```

### 2. Gauges

Gauges represent a single numerical value that can arbitrarily go up and down:

```go
config := prometheus.Config{
    Gauges: []prometheus.GaugeOpts{
        {
            Name: "active_users",
            Help: "Number of currently active users",
        },
        {
            Name: "queue_size",
            Help: "Current size of the processing queue",
        },
    },
}
```

### 3. Histograms

Histograms track the size and number of events in buckets:

```go
config := prometheus.Config{
    Histograms: []prometheus.HistogramOpts{
        {
            Name:    "request_duration_seconds",
            Help:    "Request duration in seconds",
            Buckets: []float64{0.1, 0.25, 0.5, 1, 2.5, 5, 10},
        },
    },
}
```

### 4. Summaries

Summaries track the size and number of events with quantiles:

```go
config := prometheus.Config{
    Summaries: []prometheus.SummaryOpts{
        {
            Name:       "response_time_seconds",
            Help:       "Response time in seconds",
            Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
        },
    },
}
```

## Complete Example

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    "math/rand"
    
    "github.com/kumarabd/gokit/apm"
    "github.com/kumarabd/gokit/apm/prometheus"
    "github.com/kumarabd/gokit/logger"
    "github.com/prometheus/client_golang/prometheus"
)

type MetricsService struct {
    log *logger.Handler
}

func NewMetricsService() (*MetricsService, error) {
    log, err := logger.New("metrics-service", logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    })
    if err != nil {
        return nil, err
    }
    
    return &MetricsService{log: log}, nil
}

func (s *MetricsService) setupMetrics() {
    config := prometheus.Config{
        Counters: []prometheus.CounterOpts{
            {
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            {
                Name: "user_actions_total",
                Help: "Total number of user actions",
            },
            {
                Name: "errors_total",
                Help: "Total number of errors",
            },
        },
        Gauges: []prometheus.GaugeOpts{
            {
                Name: "active_connections",
                Help: "Number of active connections",
            },
            {
                Name: "cache_size",
                Help: "Current cache size",
            },
        },
        Histograms: []prometheus.HistogramOpts{
            {
                Name:    "request_duration_seconds",
                Help:    "Request duration in seconds",
                Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
            },
        },
        Summaries: []prometheus.SummaryOpts{
            {
                Name:       "response_size_bytes",
                Help:       "Response size in bytes",
                Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
            },
        },
    }
    
    // Initialize metrics
    prometheus.InitMetrics(config)
    
    s.log.Info().Msg("Metrics initialized")
}

func (s *MetricsService) startMetricsServer() {
    // Get Prometheus HTTP handler
    handler := prometheus.GetHTTPHandler()
    
    // Register metrics endpoint
    http.Handle("/metrics", handler)
    
    // Start metrics server
    go func() {
        s.log.Info().Str("port", "9090").Msg("Starting metrics server")
        if err := http.ListenAndServe(":9090", nil); err != nil {
            s.log.Error().Err(err).Msg("Metrics server failed")
        }
    }()
}

func (s *MetricsService) simulateMetrics() {
    // Simulate some metrics collection
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            // Simulate random metrics
            s.log.Info().
                Int("active_connections", rand.Intn(100)).
                Int("cache_size", rand.Intn(1000)).
                Msg("Metrics collected")
        }
    }()
}

func main() {
    service, err := NewMetricsService()
    if err != nil {
        panic(err)
    }
    
    // Setup APM options
    apmOpts := apm.Options{
        Prometheus: apm.MetricOptions{
            Enabled: true,
        },
    }
    
    fmt.Printf("APM Options: %+v\n", apmOpts)
    
    // Setup metrics
    service.setupMetrics()
    
    // Start metrics server
    service.startMetricsServer()
    
    // Simulate metrics collection
    service.simulateMetrics()
    
    // Keep the application running
    select {}
}
```

## Configuration via YAML

You can configure metrics using YAML configuration files:

```yaml
apm:
  prometheus:
    enabled: true
    counters:
      - name: "http_requests_total"
        help: "Total number of HTTP requests"
      - name: "errors_total"
        help: "Total number of errors"
    gauges:
      - name: "active_connections"
        help: "Number of active connections"
    histograms:
      - name: "request_duration_seconds"
        help: "Request duration in seconds"
        buckets: [0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5]
    summaries:
      - name: "response_size_bytes"
        help: "Response size in bytes"
        objectives:
          0.5: 0.05
          0.9: 0.01
          0.99: 0.001
```

## Integration with HTTP Server

```go
package main

import (
    "net/http"
    "time"
    
    "github.com/kumarabd/gokit/apm/prometheus"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    httpRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests",
    })
    
    httpRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
        Name:    "http_request_duration_seconds",
        Help:    "HTTP request duration in seconds",
        Buckets: prometheus.DefBuckets,
    })
)

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Increment request counter
        httpRequestsTotal.Inc()
        
        // Call the next handler
        next.ServeHTTP(w, r)
        
        // Record request duration
        duration := time.Since(start).Seconds()
        httpRequestDuration.Observe(duration)
    })
}

func main() {
    // Setup metrics
    setupMetrics()
    
    // Create HTTP server with metrics middleware
    mux := http.NewServeMux()
    
    // Add your application routes
    mux.HandleFunc("/api/users", handleUsers)
    mux.HandleFunc("/api/health", handleHealth)
    
    // Add metrics endpoint
    mux.Handle("/metrics", prometheus.GetHTTPHandler())
    
    // Apply metrics middleware
    handler := metricsMiddleware(mux)
    
    // Start server
    http.ListenAndServe(":8080", handler)
}

func setupMetrics() {
    config := prometheus.Config{
        Counters: []prometheus.CounterOpts{
            {
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
        },
        Histograms: []prometheus.HistogramOpts{
            {
                Name:    "http_request_duration_seconds",
                Help:    "HTTP request duration in seconds",
                Buckets: prometheus.DefBuckets,
            },
        },
    }
    
    prometheus.InitMetrics(config)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"users": []}`))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"status": "healthy"}`))
}
```

## Prometheus Configuration

To scrape metrics from your application, add this to your `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'gokit-app'
    static_configs:
      - targets: ['localhost:9090']
    metrics_path: /metrics
    scrape_interval: 5s
```

## Best Practices

### 1. Metric Naming

```go
// Good - descriptive and consistent
"http_requests_total"
"database_query_duration_seconds"
"cache_hit_ratio"
"active_user_sessions"

// Avoid - vague or inconsistent
"requests"
"duration"
"ratio"
"sessions"
```

### 2. Help Text

```go
// Good - descriptive help text
{
    Name: "http_requests_total",
    Help: "Total number of HTTP requests by endpoint and status code",
}

// Avoid - generic help text
{
    Name: "http_requests_total",
    Help: "HTTP requests",
}
```

### 3. Bucket Configuration

```go
// Good - appropriate buckets for your use case
Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5, 10}

// Avoid - too many or inappropriate buckets
Buckets: []float64{0.001, 0.002, 0.003, 0.004, 0.005, ...}
```

### 4. Label Usage

```go
// Use labels for dimensions that have a limited set of values
{
    Name: "http_requests_total",
    Help: "Total number of HTTP requests",
    ConstLabels: prometheus.Labels{
        "service": "user-api",
        "version": "v1.0.0",
    },
}
```

## Monitoring and Alerting

With Prometheus metrics, you can create powerful monitoring and alerting:

```yaml
# prometheus.yml
rule_files:
  - "alerts.yml"

# alerts.yml
groups:
  - name: gokit-app
    rules:
      - alert: HighErrorRate
        expr: rate(errors_total[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} errors per second"
      
      - alert: HighResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High response time detected"
          description: "95th percentile response time is {{ $value }} seconds"
```

This provides comprehensive monitoring capabilities for your GoKit applications.
