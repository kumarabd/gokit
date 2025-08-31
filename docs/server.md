# Server Abstractions

GoKit provides server abstractions for HTTP and gRPC servers with a consistent interface for managing server lifecycle and configuration.

## Overview

The server abstractions provide:
- **Unified server interface** for HTTP and gRPC servers
- **Server lifecycle management** with start/stop capabilities
- **Configuration management** for server settings
- **Error handling** for server operations
- **Address management** for multiple server endpoints

## Basic Usage

### 1. Server Interface

```go
package main

import (
    "github.com/kumarabd/gokit/server"
)

func main() {
    // Create server options
    opts := server.Options{
        Kind: server.HTTP,
    }
    
    // Create server instance
    // Note: The actual server implementation would need to be provided
    // This is just the interface definition
}
```

## Server Types

### 1. Server Kinds

```go
const (
    GRPC ServerKind = "grpc"
    HTTP ServerKind = "http"
)
```

### 2. Server Interface

```go
type Server interface {
    Run(chan struct{}, chan error)
}
```

## Configuration

### 1. Server Options

```go
type Options struct {
    Kind ServerKind
}
```

### 2. Address Configuration

```go
type HostPort struct {
    Host string `json:"host,omitempty" yaml:"host,omitempty"`
    Port string `json:"port,omitempty" yaml:"port,omitempty"`
}

type Addresses map[string]HostPort
```

## HTTP Server Implementation

### 1. Basic HTTP Server

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/kumarabd/gokit/server"
    "github.com/kumarabd/gokit/logger"
)

type HTTPServer struct {
    server *http.Server
    log    *logger.Handler
}

func NewHTTPServer(addr string) (*HTTPServer, error) {
    log, err := logger.New("http-server", logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    })
    if err != nil {
        return nil, err
    }
    
    mux := http.NewServeMux()
    
    // Add routes
    mux.HandleFunc("/health", handleHealth)
    mux.HandleFunc("/api/users", handleUsers)
    
    httpServer := &http.Server{
        Addr:    addr,
        Handler: mux,
    }
    
    return &HTTPServer{
        server: httpServer,
        log:    log,
    }, nil
}

func (s *HTTPServer) Run(stopCh chan struct{}, errCh chan error) {
    s.log.Info().Str("address", s.server.Addr).Msg("Starting HTTP server")
    
    // Start server in goroutine
    go func() {
        if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            s.log.Error().Err(err).Msg("HTTP server error")
            errCh <- err
        }
    }()
    
    // Wait for stop signal
    <-stopCh
    
    s.log.Info().Msg("Shutting down HTTP server")
    
    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := s.server.Shutdown(ctx); err != nil {
        s.log.Error().Err(err).Msg("Error during server shutdown")
        errCh <- err
    }
    
    s.log.Info().Msg("HTTP server stopped")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"status": "healthy"}`))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"users": []}`))
}

func main() {
    // Create HTTP server
    httpServer, err := NewHTTPServer(":8080")
    if err != nil {
        panic(err)
    }
    
    // Create channels for server management
    stopCh := make(chan struct{})
    errCh := make(chan error, 1)
    
    // Start server
    go httpServer.Run(stopCh, errCh)
    
    // Wait for interrupt signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    select {
    case sig := <-sigCh:
        fmt.Printf("Received signal: %v\n", sig)
        close(stopCh)
    case err := <-errCh:
        fmt.Printf("Server error: %v\n", err)
    }
}
```

### 2. HTTP Server with Middleware

```go
package main

import (
    "net/http"
    "time"
    
    "github.com/kumarabd/gokit/logger"
)

type HTTPServer struct {
    server *http.Server
    log    *logger.Handler
}

func NewHTTPServerWithMiddleware(addr string) (*HTTPServer, error) {
    log, err := logger.New("http-server", logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    })
    if err != nil {
        return nil, err
    }
    
    mux := http.NewServeMux()
    
    // Add routes
    mux.HandleFunc("/health", handleHealth)
    mux.HandleFunc("/api/users", handleUsers)
    
    // Apply middleware
    handler := applyMiddleware(mux, log)
    
    httpServer := &http.Server{
        Addr:    addr,
        Handler: handler,
    }
    
    return &HTTPServer{
        server: httpServer,
        log:    log,
    }, nil
}

func applyMiddleware(next http.Handler, log *logger.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Log request
        log.Info().
            Str("method", r.Method).
            Str("path", r.URL.Path).
            Str("remote_addr", r.RemoteAddr).
            Msg("HTTP request started")
        
        // Call next handler
        next.ServeHTTP(w, r)
        
        // Log response
        duration := time.Since(start)
        log.Info().
            Str("method", r.Method).
            Str("path", r.URL.Path).
            Dur("duration", duration).
            Msg("HTTP request completed")
    })
}

func (s *HTTPServer) Run(stopCh chan struct{}, errCh chan error) {
    s.log.Info().Str("address", s.server.Addr).Msg("Starting HTTP server")
    
    go func() {
        if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            s.log.Error().Err(err).Msg("HTTP server error")
            errCh <- err
        }
    }()
    
    <-stopCh
    
    s.log.Info().Msg("Shutting down HTTP server")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := s.server.Shutdown(ctx); err != nil {
        s.log.Error().Err(err).Msg("Error during server shutdown")
        errCh <- err
    }
    
    s.log.Info().Msg("HTTP server stopped")
}
```

## gRPC Server Implementation

### 1. Basic gRPC Server

```go
package main

import (
    "context"
    "fmt"
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    
    "github.com/kumarabd/gokit/server"
    "github.com/kumarabd/gokit/logger"
)

type GRPCServer struct {
    server *grpc.Server
    log    *logger.Handler
    addr   string
}

func NewGRPCServer(addr string) (*GRPCServer, error) {
    log, err := logger.New("grpc-server", logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    })
    if err != nil {
        return nil, err
    }
    
    // Create gRPC server
    grpcServer := grpc.NewServer()
    
    // Register services
    // pb.RegisterUserServiceServer(grpcServer, &UserService{})
    
    // Enable reflection
    reflection.Register(grpcServer)
    
    return &GRPCServer{
        server: grpcServer,
        log:    log,
        addr:   addr,
    }, nil
}

func (s *GRPCServer) Run(stopCh chan struct{}, errCh chan error) {
    s.log.Info().Str("address", s.addr).Msg("Starting gRPC server")
    
    // Create listener
    lis, err := net.Listen("tcp", s.addr)
    if err != nil {
        s.log.Error().Err(err).Msg("Failed to create listener")
        errCh <- err
        return
    }
    
    // Start server in goroutine
    go func() {
        if err := s.server.Serve(lis); err != nil {
            s.log.Error().Err(err).Msg("gRPC server error")
            errCh <- err
        }
    }()
    
    // Wait for stop signal
    <-stopCh
    
    s.log.Info().Msg("Shutting down gRPC server")
    
    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    done := make(chan struct{})
    go func() {
        s.server.GracefulStop()
        close(done)
    }()
    
    select {
    case <-ctx.Done():
        s.log.Warn().Msg("Forcing gRPC server shutdown")
        s.server.Stop()
    case <-done:
        s.log.Info().Msg("gRPC server stopped gracefully")
    }
}

func main() {
    // Create gRPC server
    grpcServer, err := NewGRPCServer(":9090")
    if err != nil {
        panic(err)
    }
    
    // Create channels for server management
    stopCh := make(chan struct{})
    errCh := make(chan error, 1)
    
    // Start server
    go grpcServer.Run(stopCh, errCh)
    
    // Wait for interrupt signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    select {
    case sig := <-sigCh:
        fmt.Printf("Received signal: %v\n", sig)
        close(stopCh)
    case err := <-errCh:
        fmt.Printf("Server error: %v\n", err)
    }
}
```

## Server Factory

### 1. Server Factory Pattern

```go
package main

import (
    "fmt"
    "github.com/kumarabd/gokit/server"
)

type ServerFactory struct{}

func NewServerFactory() *ServerFactory {
    return &ServerFactory{}
}

func (f *ServerFactory) CreateServer(opts server.Options) (server.Server, error) {
    switch opts.Kind {
    case server.HTTP:
        return NewHTTPServer(":8080")
    case server.GRPC:
        return NewGRPCServer(":9090")
    default:
        return nil, server.ErrInvalidKind
    }
}

func main() {
    factory := NewServerFactory()
    
    // Create HTTP server
    httpOpts := server.Options{Kind: server.HTTP}
    httpServer, err := factory.CreateServer(httpOpts)
    if err != nil {
        panic(err)
    }
    
    // Create gRPC server
    grpcOpts := server.Options{Kind: server.GRPC}
    grpcServer, err := factory.CreateServer(grpcOpts)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("HTTP Server: %T\n", httpServer)
    fmt.Printf("gRPC Server: %T\n", grpcServer)
}
```

## Multiple Server Management

### 1. Server Manager

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
    
    "github.com/kumarabd/gokit/server"
    "github.com/kumarabd/gokit/logger"
)

type ServerManager struct {
    servers map[string]server.Server
    log     *logger.Handler
    wg      sync.WaitGroup
}

func NewServerManager() (*ServerManager, error) {
    log, err := logger.New("server-manager", logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    })
    if err != nil {
        return nil, err
    }
    
    return &ServerManager{
        servers: make(map[string]server.Server),
        log:     log,
    }, nil
}

func (m *ServerManager) AddServer(name string, srv server.Server) {
    m.servers[name] = srv
    m.log.Info().Str("server_name", name).Msg("Server added to manager")
}

func (m *ServerManager) StartAll() {
    m.log.Info().Msg("Starting all servers")
    
    for name, srv := range m.servers {
        m.wg.Add(1)
        go func(name string, srv server.Server) {
            defer m.wg.Done()
            
            stopCh := make(chan struct{})
            errCh := make(chan error, 1)
            
            m.log.Info().Str("server_name", name).Msg("Starting server")
            
            go srv.Run(stopCh, errCh)
            
            select {
            case err := <-errCh:
                m.log.Error().Err(err).Str("server_name", name).Msg("Server error")
            case <-stopCh:
                m.log.Info().Str("server_name", name).Msg("Server stopped")
            }
        }(name, srv)
    }
}

func (m *ServerManager) StopAll() {
    m.log.Info().Msg("Stopping all servers")
    
    // Signal all servers to stop
    for name := range m.servers {
        m.log.Info().Str("server_name", name).Msg("Stopping server")
        // In a real implementation, you would signal each server to stop
    }
    
    // Wait for all servers to stop
    m.wg.Wait()
    m.log.Info().Msg("All servers stopped")
}

func main() {
    // Create server manager
    manager, err := NewServerManager()
    if err != nil {
        panic(err)
    }
    
    // Create servers
    httpServer, _ := NewHTTPServer(":8080")
    grpcServer, _ := NewGRPCServer(":9090")
    
    // Add servers to manager
    manager.AddServer("http", httpServer)
    manager.AddServer("grpc", grpcServer)
    
    // Start all servers
    manager.StartAll()
    
    // Wait for interrupt signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigCh
    fmt.Println("Shutting down...")
    
    // Stop all servers
    manager.StopAll()
}
```

## Configuration Management

### 1. Server Configuration

```go
package main

import (
    "github.com/kumarabd/gokit/config"
    "github.com/kumarabd/gokit/server"
)

type ServerConfig struct {
    HTTP struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"http"`
    
    GRPC struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"grpc"`
    
    Servers []struct {
        Name string `yaml:"name"`
        Type string `yaml:"type"`
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"servers"`
}

func main() {
    var cfg ServerConfig
    
    // Load configuration
    configObj, err := config.New(&cfg)
    if err != nil {
        panic(err)
    }
    
    // Create servers based on configuration
    for _, serverConfig := range cfg.Servers {
        addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
        
        switch serverConfig.Type {
        case "http":
            srv, _ := NewHTTPServer(addr)
            fmt.Printf("Created HTTP server: %s at %s\n", serverConfig.Name, addr)
        case "grpc":
            srv, _ := NewGRPCServer(addr)
            fmt.Printf("Created gRPC server: %s at %s\n", serverConfig.Name, addr)
        }
    }
}
```

## Error Handling

### 1. Server Errors

```go
import "github.com/kumarabd/gokit/server"

var (
    ErrInvalidKind    = errors.New("", errors.Alert, "Unknown server kind")
    ErrInvalidName    = errors.New("", errors.Alert, "Unknown server name")
    ErrInvalidVersion = errors.New("", errors.Alert, "Unknown server version")
)
```

### 2. Error Handling in Servers

```go
func (s *HTTPServer) Run(stopCh chan struct{}, errCh chan error) {
    s.log.Info().Str("address", s.server.Addr).Msg("Starting HTTP server")
    
    go func() {
        if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            s.log.Error().Err(err).Msg("HTTP server error")
            
            // Create structured error
            serverErr := errors.New("SERVER_ERROR", errors.Critical, 
                "HTTP server failed:", err)
            errCh <- serverErr
        }
    }()
    
    <-stopCh
    
    s.log.Info().Msg("Shutting down HTTP server")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := s.server.Shutdown(ctx); err != nil {
        s.log.Error().Err(err).Msg("Error during server shutdown")
        
        shutdownErr := errors.New("SHUTDOWN_ERROR", errors.Warn, 
            "Server shutdown failed:", err)
        errCh <- shutdownErr
    }
    
    s.log.Info().Msg("HTTP server stopped")
}
```

## Best Practices

### 1. Graceful Shutdown

```go
// Always implement graceful shutdown
func (s *Server) gracefulShutdown(ctx context.Context) error {
    s.log.Info().Msg("Starting graceful shutdown")
    
    // Stop accepting new connections
    if err := s.server.Shutdown(ctx); err != nil {
        return err
    }
    
    s.log.Info().Msg("Graceful shutdown completed")
    return nil
}
```

### 2. Health Checks

```go
// Implement health check endpoints
func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}
```

### 3. Logging

```go
// Log server lifecycle events
func (s *Server) logServerEvent(event string, fields ...interface{}) {
    s.log.Info().
        Str("event", event).
        Str("server_type", "http").
        Str("address", s.server.Addr).
        Msg("Server event")
}
```

### 4. Configuration Validation

```go
// Validate server configuration
func validateServerConfig(cfg *ServerConfig) error {
    if cfg.HTTP.Port <= 0 || cfg.HTTP.Port > 65535 {
        return errors.New("CONFIG_ERROR", errors.Alert, "Invalid HTTP port")
    }
    
    if cfg.GRPC.Port <= 0 || cfg.GRPC.Port > 65535 {
        return errors.New("CONFIG_ERROR", errors.Alert, "Invalid gRPC port")
    }
    
    return nil
}
```

This provides a comprehensive foundation for server management in your GoKit applications.
