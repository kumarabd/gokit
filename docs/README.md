# GoKit Documentation

GoKit is a comprehensive utility library for building Go microservices. It provides standardized abstractions and utilities for common microservice patterns including logging, configuration management, error handling, caching, monitoring, tracing, and HTTP client operations.

## Overview

GoKit is designed to be as durable as vibranium and up to 40 times faster than traditional approaches. It provides a complete toolkit for building production-ready microservices with consistent patterns and best practices.

## Features

- **Logging**: Structured JSON logging with zerolog integration
- **Configuration**: Flexible configuration management with YAML files, environment variables, and command-line flags
- **Error Handling**: Standardized error types with severity levels and error codes
- **Caching**: In-memory caching with TTL support
- **Monitoring**: Prometheus metrics integration
- **Tracing**: OpenTelemetry tracing support
- **HTTP Client**: Simple HTTP client with request/response handling
- **Server Abstractions**: HTTP and gRPC server interfaces

## Quick Start

```go
package main

import (
    "github.com/kumarabd/gokit/config"
    "github.com/kumarabd/gokit/logger"
    "github.com/kumarabd/gokit/errors"
)

type AppConfig struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    LogLevel string `yaml:"log_level"`
}

func main() {
    // Initialize configuration
    var cfg AppConfig
    configObj, err := config.New(&cfg)
    if err != nil {
        panic(err)
    }
    
    // Initialize logger
    logOpts := logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    }
    logger, err := logger.New("myapp", logOpts)
    if err != nil {
        panic(err)
    }
    
    // Use structured logging
    logger.Info().Str("component", "main").Msg("Application started")
    
    // Handle errors with severity
    if err := someOperation(); err != nil {
        appErr := errors.New("OP_001", errors.Critical, "Operation failed:", err)
        logger.Error().Err(appErr).Msg("Critical error occurred")
    }
}
```

## Documentation Sections

- [Configuration Management](./configuration.md) - Learn how to manage application configuration
- [Logging](./logging.md) - Structured logging with zerolog integration
- [Error Handling](./error-handling.md) - Standardized error types and handling
- [Caching](./caching.md) - In-memory caching utilities
- [Monitoring](./monitoring.md) - Prometheus metrics integration
- [Tracing](./tracing.md) - OpenTelemetry tracing support
- [HTTP Client](./http-client.md) - HTTP client utilities
- [Server Abstractions](./server.md) - HTTP and gRPC server interfaces
- [Examples](./examples.md) - Complete examples and use cases
- [API Reference](./api-reference.md) - Complete API documentation

## Installation

```bash
go get github.com/kumarabd/gokit
```

## Requirements

- Go 1.13 or higher
- Dependencies are automatically managed via go.mod

## Contributing

Please read our contributing guidelines and ensure all tests pass before submitting pull requests.

## License

This project is licensed under the terms specified in the LICENSE file.
