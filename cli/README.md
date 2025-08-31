# GoKit CLI

A command-line tool for scaffolding and managing Go microservices with GoKit.

## Installation

```bash
# Build from source
go build -o gokit cli/main.go

# Or install globally
go install github.com/kumarabd/gokit/cli@latest
```

## Usage

### Create a New Service

```bash
# Create a new HTTP service
gokit new service --name user-service --template http

# Create a new gRPC service
gokit new service --name payment-service --template grpc

# Create a new event-driven service
gokit new service --name notification-service --template event

# Create a new worker service
gokit new service --name email-worker --template worker

# Specify output directory
gokit new service --name user-service --template http --output ./services

# Force overwrite existing directory
gokit new service --name user-service --template http --force
```

### Add Features to Existing Service

```bash
# Add monitoring (Prometheus metrics)
gokit add monitoring --service ./user-service

# Add tracing (OpenTelemetry)
gokit add tracing --service ./payment-service

# Add caching (In-memory cache)
gokit add caching --service ./notification-service

# Add HTTP client utilities
gokit add client --service ./api-gateway

# Add common middleware
gokit add middleware --service ./user-service
```

### Show Version

```bash
gokit version
```

## Supported Templates

### HTTP Service
- Standard HTTP API service
- RESTful endpoints
- Health checks
- Graceful shutdown
- Structured logging
- Configuration management

### gRPC Service
- gRPC server setup
- Protocol buffer support
- Service registration
- Structured logging
- Configuration management

### Event-Driven Service
- Event processing framework
- Message queue integration
- Background processing
- Structured logging
- Configuration management

### Worker Service
- Background job processing
- Scheduled tasks
- Long-running processes
- Structured logging
- Configuration management

## Supported Features

### Monitoring
- Prometheus metrics
- HTTP request counters
- Request duration histograms
- Custom metrics support
- Metrics endpoint

### Tracing
- OpenTelemetry integration
- Jaeger exporter
- Distributed tracing
- Span attributes
- Trace context propagation

### Caching
- In-memory caching
- TTL support
- Cache cleanup
- Thread-safe operations
- Configurable settings

### HTTP Client
- RESTful client utilities
- JSON serialization
- Custom headers
- Timeout support
- Error handling

### Middleware
- Request logging
- CORS handling
- Panic recovery
- Request timeout
- Response wrapping

## Project Structure

Generated services follow the standard Go project layout:

```
service-name/
├── cmd/
│   └── main.go          # Application entry point
├── internal/
│   ├── handler/         # HTTP handlers
│   ├── service/         # Business logic
│   ├── repository/      # Data access
│   ├── middleware/      # Custom middleware
│   ├── cache/           # Caching layer
│   ├── client/          # HTTP client
│   ├── monitoring/      # Metrics setup
│   └── tracing/         # Tracing setup
├── pkg/
│   └── utils/           # Public utilities
├── api/
│   └── proto/           # Protocol buffers
├── ci/
│   └── docker/          # CI/CD configuration
├── docs/
│   └── config.md        # Configuration docs
├── go.mod               # Go module file
├── go.sum               # Dependency checksums
├── Makefile             # Build commands
├── .gitignore           # Git ignore rules
└── README.md            # Service documentation
```

## Configuration

Services are configured using:
- YAML configuration files
- Environment variables
- Command-line flags

Example configuration:
```yaml
server:
  host: "0.0.0.0"
  port: 8080

log:
  format: "json"
  debug_level: "info"

monitoring:
  enabled: true
  port: 9090
  path: "/metrics"

tracing:
  enabled: true
  endpoint: "http://localhost:14268/api/traces"
  service_name: "user-service"

cache:
  enabled: true
  ttl: 300
  cleanup: 600
```

## Development

### Building

```bash
# Build the CLI
go build -o gokit cli/main.go

# Run tests
go test ./...

# Lint code
golangci-lint run
```

### Adding New Templates

1. Create template files in `cli/templates/`
2. Update the `createMainGo` function in `cli/commands/new.go`
3. Add template validation in `validateTemplate`

### Adding New Features

1. Create feature implementation in `cli/commands/add.go`
2. Add feature validation in `validateFeature`
3. Update documentation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
