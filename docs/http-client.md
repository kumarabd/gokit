# HTTP Client

GoKit provides a simple and flexible HTTP client for making HTTP requests with support for headers, query parameters, and response handling.

## Overview

The HTTP client provides:
- **Simple request/response handling** with a clean API
- **Header and query parameter support** for flexible requests
- **Response status and data access** for easy response handling
- **Error handling** with detailed error information
- **Thread-safe operations**

## Basic Usage

### 1. Create HTTP Client

```go
package main

import (
    "fmt"
    "github.com/kumarabd/gokit/client"
)

func main() {
    // Create client options
    opts := client.Options{
        Type:    client.GET,
        URL:     "https://api.example.com/users",
        Headers: map[string][]string{
            "Authorization": {"Bearer token123"},
            "Content-Type":  {"application/json"},
        },
        Params: map[string][]string{
            "page": {"1"},
            "limit": {"10"},
        },
    }
    
    // Create client
    httpClient, err := client.New(opts)
    if err != nil {
        panic(err)
    }
    
    // Make request
    response, err := httpClient.Do()
    if err != nil {
        panic(err)
    }
    
    // Handle response
    fmt.Printf("Status: %s\n", response.Status)
    fmt.Printf("Data: %s\n", string(response.Data))
}
```

## Client Configuration

### 1. Client Options

```go
type Options struct {
    Type    string              // HTTP method (GET, POST, etc.)
    URL     string              // Request URL
    Headers map[string][]string // Request headers
    Params  map[string][]string // Query parameters (for GET requests)
}
```

### 2. Response Structure

```go
type Response struct {
    Code   int    // HTTP status code
    Status string // HTTP status text
    Data   []byte // Response body
}
```

## HTTP Methods

### 1. GET Requests

```go
// Simple GET request
opts := client.Options{
    Type: client.GET,
    URL:  "https://api.example.com/users",
}

// GET request with query parameters
opts := client.Options{
    Type: client.GET,
    URL:  "https://api.example.com/users",
    Params: map[string][]string{
        "page":   {"1"},
        "limit":  {"10"},
        "search": {"john"},
    },
}

// GET request with headers
opts := client.Options{
    Type: client.GET,
    URL:  "https://api.example.com/users",
    Headers: map[string][]string{
        "Authorization": {"Bearer token123"},
        "Accept":        {"application/json"},
    },
}
```

### 2. POST Requests

```go
// POST request with headers
opts := client.Options{
    Type: client.POST,
    URL:  "https://api.example.com/users",
    Headers: map[string][]string{
        "Content-Type":  {"application/json"},
        "Authorization": {"Bearer token123"},
    },
}

// Note: The current implementation doesn't support request body
// You would need to extend it for POST requests with body
```

## Response Handling

### 1. Basic Response Processing

```go
response, err := httpClient.Do()
if err != nil {
    log.Printf("Request failed: %v", err)
    return
}

// Check status code
if response.Code >= 200 && response.Code < 300 {
    fmt.Printf("Success: %s\n", string(response.Data))
} else {
    fmt.Printf("Error: %d - %s\n", response.Code, response.Status)
}
```

### 2. JSON Response Parsing

```go
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

response, err := httpClient.Do()
if err != nil {
    return nil, err
}

var user User
if err := json.Unmarshal(response.Data, &user); err != nil {
    return nil, fmt.Errorf("failed to parse JSON: %w", err)
}

return &user, nil
```

### 3. Error Handling

```go
response, err := httpClient.Do()
if err != nil {
    // Handle network errors, timeouts, etc.
    return nil, fmt.Errorf("request failed: %w", err)
}

// Handle HTTP error status codes
if response.Code >= 400 {
    return nil, fmt.Errorf("HTTP error %d: %s", response.Code, response.Status)
}

// Handle specific status codes
switch response.Code {
case 200:
    // Success
case 401:
    return nil, fmt.Errorf("unauthorized")
case 404:
    return nil, fmt.Errorf("not found")
case 500:
    return nil, fmt.Errorf("server error")
default:
    return nil, fmt.Errorf("unexpected status: %d", response.Code)
}
```

## Complete Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"
    
    "github.com/kumarabd/gokit/client"
    "github.com/kumarabd/gokit/logger"
)

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserService struct {
    baseURL string
    token   string
    log     *logger.Handler
}

func NewUserService(baseURL, token string) (*UserService, error) {
    log, err := logger.New("user-service", logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    })
    if err != nil {
        return nil, err
    }
    
    return &UserService{
        baseURL: baseURL,
        token:   token,
        log:     log,
    }, nil
}

func (s *UserService) GetUser(userID string) (*User, error) {
    // Create client options
    opts := client.Options{
        Type: client.GET,
        URL:  fmt.Sprintf("%s/users/%s", s.baseURL, userID),
        Headers: map[string][]string{
            "Authorization": {fmt.Sprintf("Bearer %s", s.token)},
            "Accept":        {"application/json"},
        },
    }
    
    // Create client
    httpClient, err := client.New(opts)
    if err != nil {
        s.log.Error().Err(err).Str("user_id", userID).Msg("Failed to create HTTP client")
        return nil, err
    }
    
    // Make request
    start := time.Now()
    response, err := httpClient.Do()
    duration := time.Since(start)
    
    if err != nil {
        s.log.Error().
            Err(err).
            Str("user_id", userID).
            Dur("duration", duration).
            Msg("HTTP request failed")
        return nil, err
    }
    
    // Log request details
    s.log.Info().
        Str("user_id", userID).
        Int("status_code", response.Code).
        Dur("duration", duration).
        Msg("HTTP request completed")
    
    // Handle response
    if response.Code != 200 {
        s.log.Error().
            Str("user_id", userID).
            Int("status_code", response.Code).
            Str("status", response.Status).
            Msg("HTTP request returned error status")
        return nil, fmt.Errorf("HTTP %d: %s", response.Code, response.Status)
    }
    
    // Parse JSON response
    var user User
    if err := json.Unmarshal(response.Data, &user); err != nil {
        s.log.Error().
            Err(err).
            Str("user_id", userID).
            Msg("Failed to parse JSON response")
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }
    
    s.log.Info().
        Str("user_id", userID).
        Str("user_name", user.Name).
        Msg("User retrieved successfully")
    
    return &user, nil
}

func (s *UserService) ListUsers(page, limit int) ([]*User, error) {
    // Create client options
    opts := client.Options{
        Type: client.GET,
        URL:  fmt.Sprintf("%s/users", s.baseURL),
        Headers: map[string][]string{
            "Authorization": {fmt.Sprintf("Bearer %s", s.token)},
            "Accept":        {"application/json"},
        },
        Params: map[string][]string{
            "page":  {fmt.Sprintf("%d", page)},
            "limit": {fmt.Sprintf("%d", limit)},
        },
    }
    
    // Create client
    httpClient, err := client.New(opts)
    if err != nil {
        s.log.Error().Err(err).Msg("Failed to create HTTP client")
        return nil, err
    }
    
    // Make request
    response, err := httpClient.Do()
    if err != nil {
        s.log.Error().Err(err).Msg("HTTP request failed")
        return nil, err
    }
    
    // Handle response
    if response.Code != 200 {
        s.log.Error().
            Int("status_code", response.Code).
            Str("status", response.Status).
            Msg("HTTP request returned error status")
        return nil, fmt.Errorf("HTTP %d: %s", response.Code, response.Status)
    }
    
    // Parse JSON response
    var users []*User
    if err := json.Unmarshal(response.Data, &users); err != nil {
        s.log.Error().Err(err).Msg("Failed to parse JSON response")
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }
    
    s.log.Info().
        Int("count", len(users)).
        Msg("Users list retrieved successfully")
    
    return users, nil
}

func (s *UserService) SearchUsers(query string) ([]*User, error) {
    // Create client options
    opts := client.Options{
        Type: client.GET,
        URL:  fmt.Sprintf("%s/users/search", s.baseURL),
        Headers: map[string][]string{
            "Authorization": {fmt.Sprintf("Bearer %s", s.token)},
            "Accept":        {"application/json"},
        },
        Params: map[string][]string{
            "q": {query},
        },
    }
    
    // Create client
    httpClient, err := client.New(opts)
    if err != nil {
        s.log.Error().Err(err).Str("query", query).Msg("Failed to create HTTP client")
        return nil, err
    }
    
    // Make request
    response, err := httpClient.Do()
    if err != nil {
        s.log.Error().Err(err).Str("query", query).Msg("HTTP request failed")
        return nil, err
    }
    
    // Handle response
    if response.Code != 200 {
        s.log.Error().
            Str("query", query).
            Int("status_code", response.Code).
            Str("status", response.Status).
            Msg("HTTP request returned error status")
        return nil, fmt.Errorf("HTTP %d: %s", response.Code, response.Status)
    }
    
    // Parse JSON response
    var users []*User
    if err := json.Unmarshal(response.Data, &users); err != nil {
        s.log.Error().Err(err).Str("query", query).Msg("Failed to parse JSON response")
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }
    
    s.log.Info().
        Str("query", query).
        Int("count", len(users)).
        Msg("User search completed successfully")
    
    return users, nil
}

func main() {
    // Create user service
    service, err := NewUserService("https://api.example.com", "your-token-here")
    if err != nil {
        log.Fatal(err)
    }
    
    // Test different operations
    fmt.Println("=== Testing User Service ===")
    
    // Get single user
    user, err := service.GetUser("123")
    if err != nil {
        fmt.Printf("Error getting user: %v\n", err)
    } else {
        fmt.Printf("User: %+v\n", user)
    }
    
    // List users
    users, err := service.ListUsers(1, 10)
    if err != nil {
        fmt.Printf("Error listing users: %v\n", err)
    } else {
        fmt.Printf("Users: %d users found\n", len(users))
    }
    
    // Search users
    searchResults, err := service.SearchUsers("john")
    if err != nil {
        fmt.Printf("Error searching users: %v\n", err)
    } else {
        fmt.Printf("Search results: %d users found\n", len(searchResults))
    }
}
```

## Best Practices

### 1. Error Handling

```go
// Always check for errors
response, err := httpClient.Do()
if err != nil {
    // Handle network errors, timeouts, etc.
    return nil, fmt.Errorf("request failed: %w", err)
}

// Check HTTP status codes
if response.Code >= 400 {
    return nil, fmt.Errorf("HTTP error %d: %s", response.Code, response.Status)
}
```

### 2. Request Logging

```go
// Log request details for debugging
log.Info().
    Str("method", opts.Type).
    Str("url", opts.URL).
    Int("status_code", response.Code).
    Dur("duration", duration).
    Msg("HTTP request completed")
```

### 3. Response Validation

```go
// Validate response content type
contentType := response.Headers.Get("Content-Type")
if !strings.Contains(contentType, "application/json") {
    return nil, fmt.Errorf("unexpected content type: %s", contentType)
}

// Validate response size
if len(response.Data) > maxResponseSize {
    return nil, fmt.Errorf("response too large: %d bytes", len(response.Data))
}
```

### 4. Retry Logic

```go
func (s *Service) makeRequestWithRetry(opts client.Options, maxRetries int) (*client.Response, error) {
    var lastErr error
    
    for attempt := 0; attempt <= maxRetries; attempt++ {
        httpClient, err := client.New(opts)
        if err != nil {
            return nil, err
        }
        
        response, err := httpClient.Do()
        if err == nil {
            return response, nil
        }
        
        lastErr = err
        
        // Don't retry on client errors (4xx)
        if response != nil && response.Code >= 400 && response.Code < 500 {
            return response, nil
        }
        
        // Wait before retry
        if attempt < maxRetries {
            time.Sleep(time.Duration(attempt+1) * time.Second)
        }
    }
    
    return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries+1, lastErr)
}
```

## Limitations and Extensions

### Current Limitations

The current HTTP client implementation has some limitations:

1. **No request body support** for POST/PUT requests
2. **No timeout configuration**
3. **No retry mechanism**
4. **No connection pooling**
5. **No request/response middleware**

### Potential Extensions

```go
// Extended options for future implementation
type ExtendedOptions struct {
    Type        string
    URL         string
    Headers     map[string][]string
    Params      map[string][]string
    Body        []byte                    // Request body
    Timeout     time.Duration            // Request timeout
    MaxRetries  int                      // Maximum retry attempts
    RetryDelay  time.Duration            // Delay between retries
    Transport   *http.Transport          // Custom transport
}

// Middleware support
type Middleware func(*http.Request) error
type ResponseMiddleware func(*client.Response) error
```

## Integration with Other GoKit Components

```go
// Integration with logging
func (s *Service) makeRequest(opts client.Options) (*client.Response, error) {
    s.log.Info().
        Str("method", opts.Type).
        Str("url", opts.URL).
        Msg("Making HTTP request")
    
    response, err := s.httpClient.Do()
    if err != nil {
        s.log.Error().Err(err).Msg("HTTP request failed")
        return nil, err
    }
    
    s.log.Info().
        Int("status_code", response.Code).
        Msg("HTTP request completed")
    
    return response, nil
}

// Integration with error handling
func (s *Service) handleResponse(response *client.Response) error {
    if response.Code >= 400 {
        return errors.New("HTTP_ERROR", errors.Warn, 
            fmt.Sprintf("HTTP %d: %s", response.Code, response.Status))
    }
    return nil
}
```

This provides a solid foundation for HTTP client operations in your GoKit applications.
