# API Reference

This document provides a comprehensive API reference for all GoKit components.

## Table of Contents

- [Configuration](#configuration)
- [Logging](#logging)
- [Error Handling](#error-handling)
- [Caching](#caching)
- [Monitoring](#monitoring)
- [Tracing](#tracing)
- [HTTP Client](#http-client)
- [Server](#server)

## Configuration

### Package: `github.com/kumarabd/gokit/config`

#### Functions

##### `New(configObject interface{}) (interface{}, error)`

Creates a new configuration instance and loads configuration from multiple sources.

**Parameters:**
- `configObject interface{}` - Pointer to configuration struct

**Returns:**
- `interface{}` - The populated configuration object
- `error` - Error if configuration loading fails

**Example:**
```go
var cfg AppConfig
configObj, err := config.New(&cfg)
if err != nil {
    log.Fatal(err)
}
```

#### Types

##### `Format`

```go
type Format int
```

Logging format constants:
- `JSONLogFormat` - JSON structured logging
- `SyslogLogFormat` - Syslog format (not fully implemented)

##### `Options`

```go
type Options struct {
    Format     Format
    DebugLevel bool
}
```

Logger configuration options.

## Logging

### Package: `github.com/kumarabd/gokit/logger`

#### Functions

##### `New(appname string, opts Options) (*Handler, error)`

Creates a new logger instance.

**Parameters:**
- `appname string` - Application name for logging context
- `opts Options` - Logger configuration options

**Returns:**
- `*Handler` - Logger handler
- `error` - Error if logger creation fails

**Example:**
```go
log, err := logger.New("myapp", logger.Options{
    Format:     logger.JSONLogFormat,
    DebugLevel: true,
})
```

#### Methods

##### `Handler.AsLogrLogger() logr.Logger`

Converts GoKit logger to logr.Logger for compatibility.

**Returns:**
- `logr.Logger` - logr-compatible logger

**Example:**
```go
logrLogger := log.AsLogrLogger()
logrLogger.Info("Application started", "version", "1.0.0")
```

##### `Handler.Info() *zerolog.Event`

Creates an info level log event.

**Returns:**
- `*zerolog.Event` - Log event for chaining

**Example:**
```go
log.Info().Str("user_id", "123").Msg("User logged in")
```

##### `Handler.Error() *zerolog.Event`

Creates an error level log event.

**Returns:**
- `*zerolog.Event` - Log event for chaining

**Example:**
```go
log.Error().Err(err).Msg("Operation failed")
```

##### `Handler.Debug() *zerolog.Event`

Creates a debug level log event.

**Returns:**
- `*zerolog.Event` - Log event for chaining

**Example:**
```go
log.Debug().Str("component", "auth").Msg("Processing request")
```

##### `Handler.Warn() *zerolog.Event`

Creates a warning level log event.

**Returns:**
- `*zerolog.Event` - Log event for chaining

**Example:**
```go
log.Warn().Str("endpoint", "/api/users").Msg("Rate limit approaching")
```

##### `Handler.Fatal() *zerolog.Event`

Creates a fatal level log event and calls `os.Exit(1)`.

**Returns:**
- `*zerolog.Event` - Log event for chaining

**Example:**
```go
log.Fatal().Msg("Critical system error")
```

#### Types

##### `Handler`

```go
type Handler struct {
    zerolog.Logger
}
```

Logger handler that wraps zerolog.Logger.

##### `Options`

```go
type Options struct {
    Format     Format
    DebugLevel bool
}
```

Logger configuration options.

##### `Format`

```go
type Format int
```

Logging format constants:
- `JSONLogFormat` - JSON structured logging
- `SyslogLogFormat` - Syslog format

## Error Handling

### Package: `github.com/kumarabd/gokit/errors`

#### Functions

##### `New(code string, severity Severity, description ...interface{}) *Error`

Creates a new error with code, severity, and description.

**Parameters:**
- `code string` - Error code for categorization
- `severity Severity` - Error severity level
- `description ...interface{}` - Error description parts

**Returns:**
- `*Error` - New error instance

**Example:**
```go
err := errors.New("USER_NOT_FOUND", errors.Warn, "User not found:", userID)
```

##### `GetCode(err error) string`

Extracts error code from GoKit error.

**Parameters:**
- `err error` - Error to extract code from

**Returns:**
- `string` - Error code or empty string

**Example:**
```go
code := errors.GetCode(err)
```

##### `GetSeverity(err error) Severity`

Extracts severity level from GoKit error.

**Parameters:**
- `err error` - Error to extract severity from

**Returns:**
- `Severity` - Error severity level

**Example:**
```go
severity := errors.GetSeverity(err)
```

##### `Is(err error) bool`

Checks if error is a GoKit error.

**Parameters:**
- `err error` - Error to check

**Returns:**
- `bool` - True if GoKit error

**Example:**
```go
if errors.Is(err) {
    // Handle GoKit error
}
```

#### Methods

##### `Error.Error() string`

Returns error description as string.

**Returns:**
- `string` - Error description

**Example:**
```go
message := err.Error()
```

#### Types

##### `Error`

```go
type Error struct {
    Code        string
    Severity    Severity
    Description []interface{}
}
```

GoKit error structure.

##### `Severity`

```go
type Severity string
```

Error severity levels:
- `Emergency` - System unusable
- `Alert` - Action must be taken immediately
- `Critical` - Critical conditions
- `Warn` - Warning conditions
- `NoneSeverity` - No severity (default)

## Caching

### Package: `github.com/kumarabd/gokit/cache`

#### Interface

##### `Handler`

```go
type Handler interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, exp ...time.Duration) error
}
```

Cache handler interface.

**Methods:**
- `Get(key string) (interface{}, error)` - Retrieve value from cache
- `Set(key string, value interface{}, exp ...time.Duration) error` - Store value in cache

**Example:**
```go
value, err := cache.Get("user:123")
if err != nil {
    // Handle cache miss or error
}

err = cache.Set("user:123", user, 1*time.Hour)
if err != nil {
    // Handle cache error
}
```

### Package: `github.com/kumarabd/gokit/cache/inmem`

#### Functions

##### `New(opts Options) (cache.Handler, error)`

Creates a new in-memory cache instance.

**Parameters:**
- `opts Options` - Cache configuration options

**Returns:**
- `cache.Handler` - Cache handler
- `error` - Error if cache creation fails

**Example:**
```go
cache, err := inmem.New(inmem.Options{
    Expiration:      5 * time.Minute,
    CleanupInterval: 10 * time.Minute,
})
```

#### Types

##### `Options`

```go
type Options struct {
    Expiration      time.Duration
    CleanupInterval time.Duration
}
```

In-memory cache configuration options.

##### `inmem`

```go
type inmem struct {
    handler *gocache.Cache
}
```

In-memory cache implementation.

#### Variables

##### `ErrKeyNotExist`

```go
var ErrKeyNotExist = errors.New("", errors.Alert, "Key does not exist")
```

Error returned when cache key does not exist.

## Monitoring

### Package: `github.com/kumarabd/gokit/apm`

#### Types

##### `Options`

```go
type Options struct {
    Prometheus MetricOptions `json:"prometheus,omitempty" yaml:"prometheus,omitempty"`
}
```

APM configuration options.

##### `MetricOptions`

```go
type MetricOptions struct {
    Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}
```

Metric configuration options.

##### `MetricsType`

```go
type MetricsType string
```

Metrics type identifier.

#### Variables

##### `Prometheus`

```go
var Prometheus MetricsType = "prometheus"
```

Prometheus metrics type.

### Package: `github.com/kumarabd/gokit/apm/prometheus`

#### Functions

##### `GetHTTPHandler() http.Handler`

Returns HTTP handler for Prometheus metrics endpoint.

**Returns:**
- `http.Handler` - Prometheus metrics handler

**Example:**
```go
handler := prometheus.GetHTTPHandler()
http.Handle("/metrics", handler)
```

##### `InitMetrics(config Config)`

Initializes Prometheus metrics with configuration.

**Parameters:**
- `config Config` - Metrics configuration

**Example:**
```go
config := prometheus.Config{
    Counters: []prometheus.CounterOpts{
        {Name: "http_requests_total", Help: "Total HTTP requests"},
    },
}
prometheus.InitMetrics(config)
```

#### Types

##### `Config`

```go
type Config struct {
    Counters   []prometheus.CounterOpts   `json:"counters,omitempty" yaml:"counters,omitempty"`
    Gauges     []prometheus.GaugeOpts     `json:"gauges,omitempty" yaml:"gauges,omitempty"`
    Histograms []prometheus.HistogramOpts `json:"histograms,omitempty" yaml:"histograms,omitempty"`
    Summaries  []prometheus.SummaryOpts   `json:"summaries,omitempty" yaml:"summaries,omitempty"`
}
```

Prometheus metrics configuration.

## Tracing

### Package: `github.com/kumarabd/gokit/tracing`

#### Types

##### `Options`

```go
type Options struct {
    OTel TracingOptions `json:"otel,omitempty" yaml:"otel,omitempty"`
}
```

Tracing configuration options.

##### `TracingOptions`

```go
type TracingOptions struct {
    Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}
```

Tracing configuration options.

##### `TracingType`

```go
type TracingType string
```

Tracing type identifier.

#### Variables

##### `OTel`

```go
var OTel TracingType = "otel"
```

OpenTelemetry tracing type.

## HTTP Client

### Package: `github.com/kumarabd/gokit/client`

#### Functions

##### `New(opts Options) (*Handler, error)`

Creates a new HTTP client instance.

**Parameters:**
- `opts Options` - Client configuration options

**Returns:**
- `*Handler` - HTTP client handler
- `error` - Error if client creation fails

**Example:**
```go
client, err := client.New(client.Options{
    Type: client.GET,
    URL:  "https://api.example.com/users",
    Headers: map[string][]string{
        "Authorization": {"Bearer token123"},
    },
})
```

#### Methods

##### `Handler.Do() (*Response, error)`

Executes HTTP request.

**Returns:**
- `*Response` - HTTP response
- `error` - Error if request fails

**Example:**
```go
response, err := client.Do()
if err != nil {
    // Handle request error
}
fmt.Printf("Status: %s\n", response.Status)
```

#### Types

##### `Options`

```go
type Options struct {
    Type    string
    URL     string
    Headers map[string][]string
    Params  map[string][]string
}
```

HTTP client configuration options.

##### `Response`

```go
type Response struct {
    Code   int
    Status string
    Data   []byte
}
```

HTTP response structure.

##### `Handler`

```go
type Handler struct {
    client *http.Client
    req    *http.Request
}
```

HTTP client handler.

#### Constants

##### HTTP Methods

```go
const (
    GET  string = "GET"
    POST        = "POST"
)
```

Supported HTTP methods.

## Server

### Package: `github.com/kumarabd/gokit/server`

#### Interface

##### `Server`

```go
type Server interface {
    Run(chan struct{}, chan error)
}
```

Server interface for HTTP and gRPC servers.

**Methods:**
- `Run(stopCh chan struct{}, errCh chan error)` - Start server and handle lifecycle

**Example:**
```go
server := NewHTTPServer(":8080")
stopCh := make(chan struct{})
errCh := make(chan error, 1)

go server.Run(stopCh, errCh)

// Stop server
close(stopCh)
```

#### Types

##### `Options`

```go
type Options struct {
    Kind ServerKind
}
```

Server configuration options.

##### `ServerKind`

```go
type ServerKind string
```

Server type identifier.

##### `HostPort`

```go
type HostPort struct {
    Host string `json:"host,omitempty" yaml:"host,omitempty"`
    Port string `json:"port,omitempty" yaml:"port,omitempty"`
}
```

Host and port configuration.

##### `Addresses`

```go
type Addresses map[string]HostPort
```

Map of server addresses.

#### Constants

##### Server Kinds

```go
const (
    GRPC ServerKind = "grpc"
    HTTP ServerKind = "http"
)
```

Supported server types.

#### Variables

##### Server Errors

```go
var (
    ErrInvalidKind    = errors.New("", errors.Alert, "Unknown server kind")
    ErrInvalidName    = errors.New("", errors.Alert, "Unknown server name")
    ErrInvalidVersion = errors.New("", errors.Alert, "Unknown server version")
)
```

Server-specific error definitions.

## Common Patterns

### Error Handling

```go
// Create structured error
err := errors.New("OPERATION_FAILED", errors.Critical, "Operation failed:", cause)

// Check error type
if errors.Is(err) {
    code := errors.GetCode(err)
    severity := errors.GetSeverity(err)
    // Handle GoKit error
}

// Log error with context
log.Error().
    Err(err).
    Str("operation", "database_query").
    Str("error_code", errors.GetCode(err)).
    Msg("Operation failed")
```

### Configuration Loading

```go
// Define configuration struct
type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
}

// Load configuration
var cfg Config
configObj, err := config.New(&cfg)
if err != nil {
    log.Fatal("Failed to load configuration:", err)
}

// Use configuration
fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
```

### Logging with Context

```go
// Create logger
log, err := logger.New("myapp", logger.Options{
    Format:     logger.JSONLogFormat,
    DebugLevel: true,
})

// Log with structured fields
log.Info().
    Str("user_id", "123").
    Str("operation", "login").
    Int("response_time", 150).
    Msg("User logged in successfully")

// Log errors
log.Error().
    Err(err).
    Str("operation", "database_query").
    Msg("Database operation failed")
```

### Caching

```go
// Create cache
cache, err := inmem.New(inmem.Options{
    Expiration:      5 * time.Minute,
    CleanupInterval: 10 * time.Minute,
})

// Store value
err = cache.Set("user:123", user, 1*time.Hour)

// Retrieve value
value, err := cache.Get("user:123")
if err != nil {
    if err == inmem.ErrKeyNotExist {
        // Cache miss
    } else {
        // Cache error
    }
}
```

### HTTP Client Usage

```go
// Create client
client, err := client.New(client.Options{
    Type: client.GET,
    URL:  "https://api.example.com/users",
    Headers: map[string][]string{
        "Authorization": {"Bearer token123"},
    },
    Params: map[string][]string{
        "page": {"1"},
        "limit": {"10"},
    },
})

// Make request
response, err := client.Do()
if err != nil {
    // Handle request error
}

// Process response
if response.Code == 200 {
    var users []User
    json.Unmarshal(response.Data, &users)
}
```

### Server Management

```go
// Create server
server := NewHTTPServer(":8080")

// Start server
stopCh := make(chan struct{})
errCh := make(chan error, 1)

go server.Run(stopCh, errCh)

// Wait for error or stop signal
select {
case err := <-errCh:
    log.Error().Err(err).Msg("Server error")
case <-stopCh:
    log.Info().Msg("Server stopped")
}
```

## Best Practices

### Error Codes

Use consistent error code naming:
- Format: `[MODULE]_[ACTION]_[REASON]`
- Examples: `USER_NOT_FOUND`, `DB_CONNECTION_FAILED`, `AUTH_INVALID_TOKEN`

### Logging

- Use structured logging with relevant fields
- Include operation context in log messages
- Use appropriate log levels
- Log errors with full context

### Configuration

- Use YAML tags for configuration fields
- Provide sensible defaults
- Validate configuration on startup
- Use environment variables for secrets

### Caching

- Use descriptive cache keys
- Set appropriate TTL values
- Handle cache misses gracefully
- Monitor cache performance

### HTTP Client

- Always check response status codes
- Handle network errors appropriately
- Use timeouts for requests
- Log request/response details for debugging

### Server Management

- Implement graceful shutdown
- Handle server errors properly
- Use health check endpoints
- Monitor server metrics

This API reference provides comprehensive documentation for all GoKit components and their usage patterns.
