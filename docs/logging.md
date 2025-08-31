# Logging

GoKit provides a structured logging system built on top of zerolog, offering high-performance JSON logging with flexible configuration options.

## Overview

The logging system provides:
- **Structured JSON logging** for better parsing and analysis
- **High performance** with minimal allocations
- **Flexible configuration** with different log formats
- **Integration with logr** for compatibility with other logging libraries
- **Automatic timestamp and application name tagging**

## Basic Usage

### 1. Initialize Logger

```go
package main

import (
    "github.com/kumarabd/gokit/logger"
)

func main() {
    // Create logger options
    opts := logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    }
    
    // Initialize logger
    log, err := logger.New("myapp", opts)
    if err != nil {
        panic(err)
    }
    
    // Use the logger
    log.Info().Msg("Application started")
}
```

### 2. Log Levels

```go
// Debug level (only when DebugLevel is true)
log.Debug().Str("component", "auth").Msg("Processing authentication request")

// Info level
log.Info().Str("user_id", "123").Msg("User logged in successfully")

// Warning level
log.Warn().Str("endpoint", "/api/users").Msg("Rate limit approaching")

// Error level
log.Error().Err(err).Str("operation", "database_query").Msg("Database connection failed")

// Fatal level (calls os.Exit(1))
log.Fatal().Msg("Critical system error")
```

## Configuration Options

### Logger Options

```go
type Options struct {
    Format     Format // JSONLogFormat or SyslogLogFormat
    DebugLevel bool   // Enable debug level logging
}
```

### Format Types

```go
const (
    JSONLogFormat  = iota // JSON structured logging
    SyslogLogFormat       // Syslog format (not fully implemented)
)
```

## Structured Logging

### Adding Fields

```go
// String fields
log.Info().
    Str("user_id", "123").
    Str("action", "login").
    Str("ip", "192.168.1.1").
    Msg("User authentication")

// Numeric fields
log.Info().
    Int("response_time", 150).
    Float64("cpu_usage", 45.2).
    Msg("Request processed")

// Boolean fields
log.Info().
    Bool("authenticated", true).
    Bool("admin", false).
    Msg("User session created")

// Time fields
log.Info().
    Time("created_at", time.Now()).
    Msg("Record created")

// Duration fields
log.Info().
    Dur("duration", 150*time.Millisecond).
    Msg("Operation completed")
```

### Error Logging

```go
func processUser(userID string) error {
    if err := validateUser(userID); err != nil {
        log.Error().
            Err(err).
            Str("user_id", userID).
            Str("operation", "validate_user").
            Msg("User validation failed")
        return err
    }
    return nil
}
```

### Contextual Logging

```go
// Create a logger with context
userLogger := log.With().Str("user_id", "123").Logger()

// All subsequent logs will include the user_id
userLogger.Info().Msg("User action performed")
userLogger.Error().Err(err).Msg("User operation failed")
```

## Integration with logr

GoKit logger can be converted to a logr.Logger for compatibility with libraries that use logr:

```go
import (
    "github.com/go-logr/logr"
    "github.com/kumarabd/gokit/logger"
)

func main() {
    // Initialize GoKit logger
    gokitLogger, err := logger.New("myapp", logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    })
    if err != nil {
        panic(err)
    }
    
    // Convert to logr.Logger
    logrLogger := gokitLogger.AsLogrLogger()
    
    // Use with logr-compatible libraries
    logrLogger.Info("Application started", "version", "1.0.0")
    logrLogger.Error(err, "Operation failed", "operation", "database_query")
}
```

## Best Practices

### 1. Use Structured Fields

```go
// Good - structured and searchable
log.Info().
    Str("user_id", userID).
    Str("action", "login").
    Str("ip", clientIP).
    Int("response_time", responseTime).
    Msg("User login successful")

// Avoid - unstructured text
log.Info().Msgf("User %s logged in from %s in %dms", userID, clientIP, responseTime)
```

### 2. Consistent Field Names

```go
// Use consistent field names across your application
const (
    FieldUserID     = "user_id"
    FieldAction     = "action"
    FieldIP         = "ip"
    FieldResponseTime = "response_time"
    FieldError      = "error"
    FieldOperation  = "operation"
)

log.Info().
    Str(FieldUserID, userID).
    Str(FieldAction, "login").
    Str(FieldIP, clientIP).
    Int(FieldResponseTime, responseTime).
    Msg("User login successful")
```

### 3. Error Context

```go
// Always include context with errors
if err := database.Query(); err != nil {
    log.Error().
        Err(err).
        Str("operation", "database_query").
        Str("table", "users").
        Str("query", "SELECT * FROM users").
        Msg("Database query failed")
    return err
}
```

### 4. Performance Considerations

```go
// Good - conditional logging
if log.GetLevel() <= zerolog.DebugLevel {
    log.Debug().
        Str("user_id", userID).
        Str("request_data", string(requestBody)).
        Msg("Processing user request")
}

// Avoid - expensive operations in log statements
log.Debug().Str("data", expensiveOperation()).Msg("Debug info")
```

## Complete Example

```go
package main

import (
    "errors"
    "time"
    
    "github.com/kumarabd/gokit/logger"
)

type UserService struct {
    log *logger.Handler
}

func NewUserService() (*UserService, error) {
    log, err := logger.New("user-service", logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    })
    if err != nil {
        return nil, err
    }
    
    return &UserService{log: log}, nil
}

func (s *UserService) CreateUser(userID, email string) error {
    start := time.Now()
    
    s.log.Info().
        Str("user_id", userID).
        Str("email", email).
        Str("operation", "create_user").
        Msg("Creating new user")
    
    // Simulate some work
    time.Sleep(100 * time.Millisecond)
    
    // Simulate error
    if userID == "error" {
        err := errors.New("user already exists")
        s.log.Error().
            Err(err).
            Str("user_id", userID).
            Str("operation", "create_user").
            Msg("Failed to create user")
        return err
    }
    
    duration := time.Since(start)
    s.log.Info().
        Str("user_id", userID).
        Str("operation", "create_user").
        Dur("duration", duration).
        Msg("User created successfully")
    
    return nil
}

func main() {
    service, err := NewUserService()
    if err != nil {
        panic(err)
    }
    
    // Create a user logger with context
    userLogger := service.log.With().Str("component", "main").Logger()
    userLogger.Info().Msg("User service started")
    
    // Test successful user creation
    if err := service.CreateUser("123", "user@example.com"); err != nil {
        userLogger.Error().Err(err).Msg("Failed to create user")
    }
    
    // Test error case
    if err := service.CreateUser("error", "error@example.com"); err != nil {
        userLogger.Error().Err(err).Msg("Failed to create user")
    }
}
```

## Output Format

The logger produces structured JSON output:

```json
{"level":"info","app":"user-service","time":"2024-01-15T10:30:00Z","user_id":"123","email":"user@example.com","operation":"create_user","message":"Creating new user"}
{"level":"info","app":"user-service","time":"2024-01-15T10:30:00Z","user_id":"123","operation":"create_user","duration":100000000,"message":"User created successfully"}
{"level":"error","app":"user-service","time":"2024-01-15T10:30:00Z","error":"user already exists","user_id":"error","operation":"create_user","message":"Failed to create user"}
```

## Performance

The logging system is designed for high performance:
- **Zero allocations** for disabled log levels
- **Minimal overhead** for enabled log levels
- **Efficient JSON encoding** using zerolog
- **Thread-safe** operations

## Integration with Monitoring

The structured JSON logs can be easily parsed by log aggregation systems like:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Fluentd
- Splunk
- CloudWatch Logs
- Google Cloud Logging

This enables powerful log analysis, alerting, and monitoring capabilities.
