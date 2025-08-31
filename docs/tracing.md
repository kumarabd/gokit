# Tracing

GoKit provides a tracing system with OpenTelemetry (OTel) integration for distributed tracing across microservices.

## Overview

The tracing system provides:
- **OpenTelemetry integration** for distributed tracing
- **Standardized tracing interface** for consistent tracing across services
- **Configurable tracing options** for different environments
- **Integration with existing observability tools**

## Basic Usage

### 1. Initialize Tracing

```go
package main

import (
    "github.com/kumarabd/gokit/tracing"
)

func main() {
    // Create tracing options
    opts := tracing.Options{
        OTel: tracing.TracingOptions{
            Enabled: true,
        },
    }
    
    // Initialize tracing (configuration will be handled separately)
    // The actual implementation depends on your specific OpenTelemetry setup
}
```

## OpenTelemetry Integration

### 1. Basic Setup

The tracing system is designed to work with OpenTelemetry. Here's how to set it up:

```go
package main

import (
    "context"
    "log"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
    "go.opentelemetry.io/otel/trace"
    
    "github.com/kumarabd/gokit/tracing"
)

func setupTracing() (*sdktrace.TracerProvider, error) {
    // Create Jaeger exporter
    exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
    if err != nil {
        return nil, err
    }
    
    // Create resource with service information
    res, err := resource.New(context.Background(),
        resource.WithAttributes(
            semconv.ServiceNameKey.String("my-service"),
            semconv.ServiceVersionKey.String("1.0.0"),
        ),
    )
    if err != nil {
        return nil, err
    }
    
    // Create tracer provider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(res),
    )
    
    // Set global tracer provider
    otel.SetTracerProvider(tp)
    
    return tp, nil
}

func main() {
    // Setup tracing
    tp, err := setupTracing()
    if err != nil {
        log.Fatal(err)
    }
    defer tp.Shutdown(context.Background())
    
    // Create tracing options
    opts := tracing.Options{
        OTel: tracing.TracingOptions{
            Enabled: true,
        },
    }
    
    // Use tracing in your application
    tracer := otel.Tracer("my-service")
    
    // Create a span
    ctx, span := tracer.Start(context.Background(), "main-operation")
    defer span.End()
    
    // Your application logic here
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    tracer := otel.Tracer("my-service")
    
    ctx, span := tracer.Start(ctx, "process-request")
    defer span.End()
    
    // Add attributes to the span
    span.SetAttributes(
        semconv.HTTPMethodKey.String("GET"),
        semconv.HTTPRouteKey.String("/api/users"),
    )
    
    // Simulate some work
    time.Sleep(100 * time.Millisecond)
    
    // Create child span
    ctx, childSpan := tracer.Start(ctx, "database-query")
    defer childSpan.End()
    
    // Simulate database query
    time.Sleep(50 * time.Millisecond)
    
    childSpan.SetAttributes(
        semconv.DBSystemKey.String("postgresql"),
        semconv.DBStatementKey.String("SELECT * FROM users"),
    )
}
```

### 2. HTTP Middleware Integration

```go
package main

import (
    "net/http"
    "time"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/trace"
)

func tracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract trace context from request headers
        ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
        
        // Start span
        tracer := otel.Tracer("http-server")
        ctx, span := tracer.Start(ctx, r.URL.Path,
            trace.WithSpanKind(trace.SpanKindServer),
        )
        defer span.End()
        
        // Add request attributes
        span.SetAttributes(
            semconv.HTTPMethodKey.String(r.Method),
            semconv.HTTPRouteKey.String(r.URL.Path),
            semconv.HTTPUserAgentKey.String(r.UserAgent()),
        )
        
        // Create response writer wrapper to capture status code
        wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        // Call next handler
        next.ServeHTTP(wrappedWriter, r.WithContext(ctx))
        
        // Add response attributes
        span.SetAttributes(semconv.HTTPStatusCodeKey.Int(wrappedWriter.statusCode))
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func main() {
    // Setup tracing
    tp, _ := setupTracing()
    defer tp.Shutdown(context.Background())
    
    // Create HTTP server with tracing middleware
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", handleUsers)
    mux.HandleFunc("/api/health", handleHealth)
    
    // Apply tracing middleware
    handler := tracingMiddleware(mux)
    
    // Start server
    http.ListenAndServe(":8080", handler)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    tracer := otel.Tracer("user-handler")
    
    ctx, span := tracer.Start(ctx, "get-users")
    defer span.End()
    
    // Simulate processing
    time.Sleep(100 * time.Millisecond)
    
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"users": []}`))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"status": "healthy"}`))
}
```

### 3. Database Tracing

```go
package main

import (
    "context"
    "database/sql"
    "time"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

type UserService struct {
    db *sql.DB
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    tracer := otel.Tracer("user-service")
    
    ctx, span := tracer.Start(ctx, "get-user",
        trace.WithAttributes(
            attribute.String("user.id", userID),
        ),
    )
    defer span.End()
    
    // Database query with tracing
    ctx, dbSpan := tracer.Start(ctx, "database.query",
        trace.WithAttributes(
            attribute.String("db.system", "postgresql"),
            attribute.String("db.statement", "SELECT * FROM users WHERE id = $1"),
        ),
    )
    defer dbSpan.End()
    
    var user User
    err := s.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = $1", userID).
        Scan(&user.ID, &user.Name, &user.Email)
    
    if err != nil {
        dbSpan.RecordError(err)
        return nil, err
    }
    
    return &user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    tracer := otel.Tracer("user-service")
    
    ctx, span := tracer.Start(ctx, "create-user",
        trace.WithAttributes(
            attribute.String("user.email", user.Email),
        ),
    )
    defer span.End()
    
    // Database insert with tracing
    ctx, dbSpan := tracer.Start(ctx, "database.insert",
        trace.WithAttributes(
            attribute.String("db.system", "postgresql"),
            attribute.String("db.statement", "INSERT INTO users (name, email) VALUES ($1, $2)"),
        ),
    )
    defer dbSpan.End()
    
    _, err := s.db.ExecContext(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
    
    if err != nil {
        dbSpan.RecordError(err)
        return err
    }
    
    return nil
}
```

## Configuration

### 1. Tracing Options

```go
type TracingOptions struct {
    Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

type Options struct {
    OTel TracingOptions `json:"otel,omitempty" yaml:"otel,omitempty"`
}
```

### 2. YAML Configuration

```yaml
tracing:
  otel:
    enabled: true
```

## Complete Example

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
    "go.opentelemetry.io/otel/trace"
    
    "github.com/kumarabd/gokit/tracing"
    "github.com/kumarabd/gokit/logger"
)

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

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

func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    tracer := otel.Tracer("user-service")
    
    ctx, span := tracer.Start(ctx, "get-user",
        trace.WithAttributes(
            semconv.HTTPMethodKey.String("GET"),
            semconv.HTTPRouteKey.String("/api/users/{id}"),
        ),
    )
    defer span.End()
    
    s.log.Info().
        Str("user_id", userID).
        Str("trace_id", span.SpanContext().TraceID().String()).
        Msg("Getting user")
    
    // Simulate database query
    ctx, dbSpan := tracer.Start(ctx, "database.query",
        trace.WithAttributes(
            semconv.DBSystemKey.String("postgresql"),
            semconv.DBStatementKey.String("SELECT * FROM users WHERE id = $1"),
        ),
    )
    defer dbSpan.End()
    
    // Simulate processing time
    time.Sleep(100 * time.Millisecond)
    
    // Simulate user data
    user := &User{
        ID:    userID,
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    s.log.Info().
        Str("user_id", userID).
        Str("trace_id", span.SpanContext().TraceID().String()).
        Msg("User retrieved successfully")
    
    return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    tracer := otel.Tracer("user-service")
    
    ctx, span := tracer.Start(ctx, "create-user",
        trace.WithAttributes(
            semconv.HTTPMethodKey.String("POST"),
            semconv.HTTPRouteKey.String("/api/users"),
        ),
    )
    defer span.End()
    
    s.log.Info().
        Str("user_email", user.Email).
        Str("trace_id", span.SpanContext().TraceID().String()).
        Msg("Creating user")
    
    // Simulate database insert
    ctx, dbSpan := tracer.Start(ctx, "database.insert",
        trace.WithAttributes(
            semconv.DBSystemKey.String("postgresql"),
            semconv.DBStatementKey.String("INSERT INTO users (name, email) VALUES ($1, $2)"),
        ),
    )
    defer dbSpan.End()
    
    // Simulate processing time
    time.Sleep(50 * time.Millisecond)
    
    // Simulate success
    user.ID = "generated-id"
    
    s.log.Info().
        Str("user_id", user.ID).
        Str("trace_id", span.SpanContext().TraceID().String()).
        Msg("User created successfully")
    
    return nil
}

func setupTracing() (*sdktrace.TracerProvider, error) {
    // Create Jaeger exporter
    exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
    if err != nil {
        return nil, err
    }
    
    // Create resource
    res, err := resource.New(context.Background(),
        resource.WithAttributes(
            semconv.ServiceNameKey.String("user-service"),
            semconv.ServiceVersionKey.String("1.0.0"),
        ),
    )
    if err != nil {
        return nil, err
    }
    
    // Create tracer provider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exp),
        sdktrace.WithResource(res),
    )
    
    // Set global tracer provider and propagator
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))
    
    return tp, nil
}

func tracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract trace context
        ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
        
        // Start span
        tracer := otel.Tracer("http-server")
        ctx, span := tracer.Start(ctx, r.URL.Path,
            trace.WithSpanKind(trace.SpanKindServer),
        )
        defer span.End()
        
        // Add request attributes
        span.SetAttributes(
            semconv.HTTPMethodKey.String(r.Method),
            semconv.HTTPRouteKey.String(r.URL.Path),
            semconv.HTTPUserAgentKey.String(r.UserAgent()),
        )
        
        // Create response writer wrapper
        wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        // Call next handler
        next.ServeHTTP(wrappedWriter, r.WithContext(ctx))
        
        // Add response attributes
        span.SetAttributes(semconv.HTTPStatusCodeKey.Int(wrappedWriter.statusCode))
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    service, _ := NewUserService()
    
    userID := r.URL.Query().Get("id")
    if userID == "" {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }
    
    user, err := service.GetUser(ctx, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    service, _ := NewUserService()
    
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if err := service.CreateUser(ctx, &user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func main() {
    // Setup tracing
    tp, err := setupTracing()
    if err != nil {
        log.Fatal(err)
    }
    defer tp.Shutdown(context.Background())
    
    // Create tracing options
    opts := tracing.Options{
        OTel: tracing.TracingOptions{
            Enabled: true,
        },
    }
    
    fmt.Printf("Tracing Options: %+v\n", opts)
    
    // Create HTTP server
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetUser(w, r)
        case "POST":
            handleCreateUser(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
    
    // Apply tracing middleware
    handler := tracingMiddleware(mux)
    
    // Start server
    fmt.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

## Best Practices

### 1. Span Naming

```go
// Good - descriptive span names
tracer.Start(ctx, "user-service.get-user")
tracer.Start(ctx, "database.query.users")
tracer.Start(ctx, "http-client.call-external-api")

// Avoid - generic span names
tracer.Start(ctx, "operation")
tracer.Start(ctx, "query")
tracer.Start(ctx, "call")
```

### 2. Attribute Usage

```go
// Add relevant attributes to spans
span.SetAttributes(
    semconv.HTTPMethodKey.String("GET"),
    semconv.HTTPRouteKey.String("/api/users"),
    attribute.String("user.id", userID),
    attribute.String("db.system", "postgresql"),
)
```

### 3. Error Handling

```go
// Record errors in spans
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
    return err
}
```

### 4. Context Propagation

```go
// Always pass context through your application
func (s *Service) Process(ctx context.Context, data interface{}) error {
    // Use the context for tracing
    ctx, span := tracer.Start(ctx, "service.process")
    defer span.End()
    
    // Pass context to downstream calls
    return s.database.Save(ctx, data)
}
```

## Integration with Observability Tools

The tracing system integrates with various observability tools:

- **Jaeger**: Distributed tracing backend
- **Zipkin**: Distributed tracing system
- **Prometheus**: Metrics collection
- **Grafana**: Visualization and monitoring
- **Elastic APM**: Application performance monitoring

This provides comprehensive observability for your GoKit applications.
