# Error Handling

GoKit provides a standardized error handling system with severity levels, error codes, and structured error information for better error management in microservices.

## Overview

The error handling system provides:
- **Structured error types** with severity levels and error codes
- **Consistent error format** across your application
- **Severity-based error handling** for different error types
- **Error code system** for easy identification and categorization
- **Integration with logging** for comprehensive error tracking

## Basic Usage

### 1. Creating Errors

```go
package main

import (
    "github.com/kumarabd/gokit/errors"
)

func main() {
    // Create a simple error
    err := errors.New("AUTH_001", errors.Critical, "Authentication failed")
    
    // Create an error with multiple description parts
    err = errors.New("DB_001", errors.Alert, "Database connection failed:", "connection timeout")
    
    // Create an error with severity
    err = errors.New("API_001", errors.Warn, "Rate limit approaching")
}
```

### 2. Error Severity Levels

```go
const (
    Emergency Severity = "emergency"  // System unusable
    Alert     Severity = "alert"      // Action must be taken immediately
    Critical  Severity = "critical"   // Critical conditions
    Warn      Severity = "warn"       // Warning conditions
    NoneSeverity Severity = "none"    // No severity (default)
)
```

## Error Structure

### Error Type

```go
type Error struct {
    Code        string        // Unique error code
    Severity    Severity      // Error severity level
    Description []interface{} // Error description parts
}
```

### Creating Different Types of Errors

```go
// Critical errors - system failures
dbErr := errors.New("DB_CONN_001", errors.Critical, "Database connection lost")

// Alert errors - immediate action required
authErr := errors.New("AUTH_FAIL_001", errors.Alert, "Invalid API key provided")

// Warning errors - potential issues
rateLimitErr := errors.New("RATE_LIMIT_001", errors.Warn, "Rate limit at 80%")

// Emergency errors - system unusable
systemErr := errors.New("SYS_CRASH_001", errors.Emergency, "System crash detected")
```

## Error Handling Functions

### 1. Error Information

```go
// Get error code
code := errors.GetCode(err)
fmt.Println("Error code:", code) // Output: Error code: DB_001

// Get error severity
severity := errors.GetSeverity(err)
fmt.Println("Severity:", severity) // Output: Severity: critical

// Check if error is a GoKit error
if errors.Is(err) {
    fmt.Println("This is a GoKit error")
}
```

### 2. Error String Representation

```go
err := errors.New("AUTH_001", errors.Critical, "Authentication failed:", "invalid token")
fmt.Println(err.Error()) // Output: Authentication failed: invalid token
```

## Integration with Logging

### Structured Error Logging

```go
package main

import (
    "github.com/kumarabd/gokit/errors"
    "github.com/kumarabd/gokit/logger"
)

func processUser(userID string) error {
    if userID == "" {
        return errors.New("USER_001", errors.Alert, "User ID is required")
    }
    
    // Simulate some operation
    if userID == "invalid" {
        return errors.New("USER_002", errors.Critical, "User not found:", userID)
    }
    
    return nil
}

func main() {
    // Initialize logger
    log, _ := logger.New("myapp", logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    })
    
    // Process user with error handling
    if err := processUser(""); err != nil {
        if errors.Is(err) {
            // Log structured error information
            log.Error().
                Err(err).
                Str("error_code", errors.GetCode(err)).
                Str("severity", string(errors.GetSeverity(err))).
                Msg("User processing failed")
        } else {
            // Handle non-GoKit errors
            log.Error().Err(err).Msg("Unknown error occurred")
        }
    }
}
```

## Error Code System

### Recommended Error Code Format

```go
// Format: [MODULE]_[NUMBER]
// Examples:
const (
    // Authentication errors
    ErrAuthInvalidToken    = "AUTH_001"
    ErrAuthExpiredToken    = "AUTH_002"
    ErrAuthMissingToken    = "AUTH_003"
    
    // Database errors
    ErrDBConnectionFailed  = "DB_001"
    ErrDBQueryFailed       = "DB_002"
    ErrDBTransactionFailed = "DB_003"
    
    // API errors
    ErrAPIRateLimit        = "API_001"
    ErrAPIInvalidRequest   = "API_002"
    ErrAPIServerError      = "API_003"
    
    // User errors
    ErrUserNotFound        = "USER_001"
    ErrUserInvalidInput    = "USER_002"
    ErrUserAlreadyExists   = "USER_003"
)
```

### Error Code Constants

```go
package myapp

import "github.com/kumarabd/gokit/errors"

// Define error codes as constants
const (
    // Authentication errors
    ErrCodeAuthInvalidToken = "AUTH_001"
    ErrCodeAuthExpiredToken = "AUTH_002"
    
    // Database errors
    ErrCodeDBConnectionFailed = "DB_001"
    ErrCodeDBQueryFailed = "DB_002"
    
    // Business logic errors
    ErrCodeUserNotFound = "USER_001"
    ErrCodeUserInvalidInput = "USER_002"
)

// Create error functions for consistency
func NewAuthInvalidTokenError() error {
    return errors.New(ErrCodeAuthInvalidToken, errors.Alert, "Invalid authentication token")
}

func NewUserNotFoundError(userID string) error {
    return errors.New(ErrCodeUserNotFound, errors.Warn, "User not found:", userID)
}

func NewDBConnectionError(cause error) error {
    return errors.New(ErrCodeDBConnectionFailed, errors.Critical, "Database connection failed:", cause)
}
```

## Best Practices

### 1. Use Consistent Error Codes

```go
// Good - consistent naming
const (
    ErrCodeAuthInvalidToken = "AUTH_001"
    ErrCodeAuthExpiredToken = "AUTH_002"
    ErrCodeAuthMissingToken = "AUTH_003"
)

// Avoid - inconsistent naming
const (
    ErrCodeAuthInvalidToken = "AUTH_001"
    ErrCodeExpiredToken = "TOKEN_002"  // Different prefix
    ErrCodeMissingToken = "AUTH_003"
)
```

### 2. Appropriate Severity Levels

```go
// Emergency - System is unusable
errors.New("SYS_CRASH_001", errors.Emergency, "System crash detected")

// Alert - Immediate action required
errors.New("AUTH_FAIL_001", errors.Alert, "Multiple failed login attempts")

// Critical - Critical conditions
errors.New("DB_CONN_001", errors.Critical, "Database connection lost")

// Warn - Warning conditions
errors.New("RATE_LIMIT_001", errors.Warn, "Rate limit approaching")
```

### 3. Descriptive Error Messages

```go
// Good - descriptive and actionable
err := errors.New("USER_001", errors.Warn, "User not found:", userID, "in database:", dbName)

// Avoid - vague messages
err := errors.New("USER_001", errors.Warn, "Error occurred")
```

### 4. Error Wrapping

```go
func processUser(userID string) error {
    if err := validateUser(userID); err != nil {
        // Wrap the original error with context
        return errors.New("USER_001", errors.Warn, "User validation failed:", err)
    }
    
    if err := saveUser(userID); err != nil {
        return errors.New("USER_002", errors.Critical, "Failed to save user:", err)
    }
    
    return nil
}
```

## Complete Example

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/kumarabd/gokit/errors"
    "github.com/kumarabd/gokit/logger"
)

// Error codes
const (
    ErrCodeAuthInvalidToken = "AUTH_001"
    ErrCodeUserNotFound = "USER_001"
    ErrCodeDBConnectionFailed = "DB_001"
    ErrCodeAPIRateLimit = "API_001"
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

func (s *UserService) AuthenticateUser(token, userID string) error {
    // Simulate authentication
    if token == "" {
        return errors.New(ErrCodeAuthInvalidToken, errors.Alert, "Authentication token is required")
    }
    
    if token == "invalid" {
        return errors.New(ErrCodeAuthInvalidToken, errors.Alert, "Invalid authentication token:", token)
    }
    
    // Simulate user lookup
    if userID == "notfound" {
        return errors.New(ErrCodeUserNotFound, errors.Warn, "User not found:", userID)
    }
    
    return nil
}

func (s *UserService) ProcessUserRequest(token, userID string) error {
    // Authenticate user
    if err := s.AuthenticateUser(token, userID); err != nil {
        // Log the error with structured information
        if errors.Is(err) {
            s.log.Error().
                Err(err).
                Str("error_code", errors.GetCode(err)).
                Str("severity", string(errors.GetSeverity(err))).
                Str("user_id", userID).
                Msg("Authentication failed")
        }
        return err
    }
    
    // Simulate rate limiting
    if userID == "rate_limited" {
        rateLimitErr := errors.New(ErrCodeAPIRateLimit, errors.Warn, "Rate limit exceeded for user:", userID)
        s.log.Warn().
            Err(rateLimitErr).
            Str("error_code", errors.GetCode(rateLimitErr)).
            Str("user_id", userID).
            Msg("Rate limit warning")
        return rateLimitErr
    }
    
    s.log.Info().
        Str("user_id", userID).
        Msg("User request processed successfully")
    
    return nil
}

func main() {
    service, err := NewUserService()
    if err != nil {
        panic(err)
    }
    
    // Test different error scenarios
    testCases := []struct {
        token  string
        userID string
        desc   string
    }{
        {"", "123", "Missing token"},
        {"invalid", "123", "Invalid token"},
        {"valid", "notfound", "User not found"},
        {"valid", "rate_limited", "Rate limited"},
        {"valid", "123", "Success"},
    }
    
    for _, tc := range testCases {
        fmt.Printf("\nTesting: %s\n", tc.desc)
        if err := service.ProcessUserRequest(tc.token, tc.userID); err != nil {
            if errors.Is(err) {
                fmt.Printf("Error: %s (Code: %s, Severity: %s)\n", 
                    err.Error(), 
                    errors.GetCode(err), 
                    errors.GetSeverity(err))
            } else {
                fmt.Printf("Unknown error: %v\n", err)
            }
        } else {
            fmt.Println("Success")
        }
    }
}
```

## Error Monitoring and Alerting

The structured error system enables powerful monitoring and alerting:

```go
// Monitor error rates by severity
func monitorErrors(err error) {
    if errors.Is(err) {
        severity := errors.GetSeverity(err)
        code := errors.GetCode(err)
        
        // Send alerts for critical and emergency errors
        if severity == errors.Critical || severity == errors.Emergency {
            sendAlert(severity, code, err.Error())
        }
        
        // Track error metrics
        incrementErrorCounter(severity, code)
    }
}
```

## Integration with External Systems

The error system can be easily integrated with external monitoring systems:

- **Prometheus**: Track error rates and severity distributions
- **AlertManager**: Send alerts based on error severity
- **ELK Stack**: Parse and analyze error logs
- **Sentry**: Error tracking and performance monitoring

This provides comprehensive error visibility and management capabilities for production systems.
