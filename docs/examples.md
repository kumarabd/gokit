# Examples

This section provides complete, working examples of GoKit components in action. Each example demonstrates real-world usage patterns and best practices.

## Table of Contents

- [Basic Microservice](#basic-microservice)
- [Configuration Management](#configuration-management)
- [Logging and Error Handling](#logging-and-error-handling)
- [Caching with HTTP Server](#caching-with-http-server)
- [Monitoring and Metrics](#monitoring-and-metrics)
- [Complete API Service](#complete-api-service)

## Basic Microservice

A simple microservice demonstrating basic GoKit usage.

### Project Structure

```
basic-service/
├── main.go
├── config.yaml
├── go.mod
└── README.md
```

### Configuration

```yaml
# config.yaml
app:
  name: "basic-service"
  version: "1.0.0"

server:
  host: "0.0.0.0"
  port: 8080

database:
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  ssl: false

logging:
  level: "info"
  format: "json"
```

### Main Application

```go
// main.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/kumarabd/gokit/config"
    "github.com/kumarabd/gokit/logger"
    "github.com/kumarabd/gokit/errors"
)

type AppConfig struct {
    App struct {
        Name    string `yaml:"name"`
        Version string `yaml:"version"`
    } `yaml:"app"`
    
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        SSL      bool   `yaml:"ssl"`
    } `yaml:"database"`
    
    Logging struct {
        Level  string `yaml:"level"`
        Format string `yaml:"format"`
    } `yaml:"logging"`
}

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserService struct {
    log *logger.Handler
    cfg *AppConfig
}

func NewUserService(cfg *AppConfig) (*UserService, error) {
    log, err := logger.New(cfg.App.Name, logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: cfg.Logging.Level == "debug",
    })
    if err != nil {
        return nil, err
    }
    
    return &UserService{
        log: log,
        cfg: cfg,
    }, nil
}

func (s *UserService) GetUser(userID string) (*User, error) {
    s.log.Info().
        Str("user_id", userID).
        Str("operation", "get_user").
        Msg("Getting user")
    
    // Simulate database query
    if userID == "notfound" {
        return nil, errors.New("USER_NOT_FOUND", errors.Warn, "User not found:", userID)
    }
    
    user := &User{
        ID:    userID,
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    s.log.Info().
        Str("user_id", userID).
        Str("user_name", user.Name).
        Msg("User retrieved successfully")
    
    return user, nil
}

func (s *UserService) CreateUser(user *User) error {
    s.log.Info().
        Str("user_email", user.Email).
        Str("operation", "create_user").
        Msg("Creating user")
    
    // Simulate database insert
    if user.Email == "duplicate@example.com" {
        return errors.New("USER_EXISTS", errors.Alert, "User already exists:", user.Email)
    }
    
    user.ID = "generated-id"
    
    s.log.Info().
        Str("user_id", user.ID).
        Str("user_email", user.Email).
        Msg("User created successfully")
    
    return nil
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}

func handleGetUser(service *UserService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.URL.Query().Get("id")
        if userID == "" {
            http.Error(w, "User ID is required", http.StatusBadRequest)
            return
        }
        
        user, err := service.GetUser(userID)
        if err != nil {
            if errors.Is(err) {
                switch errors.GetCode(err) {
                case "USER_NOT_FOUND":
                    http.Error(w, err.Error(), http.StatusNotFound)
                default:
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                }
            } else {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}

func handleCreateUser(service *UserService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
        if err := service.CreateUser(&user); err != nil {
            if errors.Is(err) {
                switch errors.GetCode(err) {
                case "USER_EXISTS":
                    http.Error(w, err.Error(), http.StatusConflict)
                default:
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                }
            } else {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
    }
}

func main() {
    // Load configuration
    var cfg AppConfig
    configObj, err := config.New(&cfg)
    if err != nil {
        log.Fatal("Failed to load configuration:", err)
    }
    
    // Set defaults
    if cfg.Server.Host == "" {
        cfg.Server.Host = "localhost"
    }
    if cfg.Server.Port == 0 {
        cfg.Server.Port = 8080
    }
    
    // Create user service
    service, err := NewUserService(&cfg)
    if err != nil {
        log.Fatal("Failed to create user service:", err)
    }
    
    // Create HTTP server
    mux := http.NewServeMux()
    mux.HandleFunc("/health", handleHealth)
    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            handleGetUser(service)(w, r)
        case "POST":
            handleCreateUser(service)(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
    
    server := &http.Server{
        Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
        Handler: mux,
    }
    
    // Start server
    go func() {
        service.log.Info().
            Str("address", server.Addr).
            Msg("Starting HTTP server")
        
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            service.log.Error().Err(err).Msg("HTTP server error")
            log.Fatal(err)
        }
    }()
    
    // Wait for interrupt signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigCh
    service.log.Info().Msg("Shutting down server")
    
    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        service.log.Error().Err(err).Msg("Error during server shutdown")
    }
    
    service.log.Info().Msg("Server stopped")
}
```

### Running the Example

```bash
# Run with default configuration
go run main.go

# Run with custom config file
go run main.go --config config.yaml

# Run with environment variables
export DB_HOST=production-db.example.com
export DB_PASSWORD=secret123
go run main.go --from-env database.host::DB_HOST --from-env database.password::DB_PASSWORD

# Test the API
curl http://localhost:8080/health
curl http://localhost:8080/api/users?id=123
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'
```

## Configuration Management

A comprehensive example of configuration management with multiple sources.

### Configuration Structure

```go
// config.go
package main

import (
    "time"
)

type Config struct {
    App struct {
        Name        string        `yaml:"name"`
        Version     string        `yaml:"version"`
        Environment string        `yaml:"environment"`
        ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
    } `yaml:"app"`
    
    Server struct {
        HTTP struct {
            Host string `yaml:"host"`
            Port int    `yaml:"port"`
        } `yaml:"http"`
        
        GRPC struct {
            Host string `yaml:"host"`
            Port int    `yaml:"port"`
        } `yaml:"grpc"`
    } `yaml:"server"`
    
    Database struct {
        Host     string        `yaml:"host"`
        Port     int           `yaml:"port"`
        Username string        `yaml:"username"`
        Password string        `yaml:"password"`
        Name     string        `yaml:"name"`
        SSL      bool          `yaml:"ssl"`
        Timeout  time.Duration `yaml:"timeout"`
        MaxConnections int     `yaml:"max_connections"`
    } `yaml:"database"`
    
    Cache struct {
        Redis struct {
            Host     string        `yaml:"host"`
            Port     int           `yaml:"port"`
            Password string        `yaml:"password"`
            DB       int           `yaml:"db"`
            Timeout  time.Duration `yaml:"timeout"`
        } `yaml:"redis"`
        
        InMemory struct {
            Expiration      time.Duration `yaml:"expiration"`
            CleanupInterval time.Duration `yaml:"cleanup_interval"`
        } `yaml:"in_memory"`
    } `yaml:"cache"`
    
    Monitoring struct {
        Prometheus struct {
            Enabled bool   `yaml:"enabled"`
            Port    int    `yaml:"port"`
            Path    string `yaml:"path"`
        } `yaml:"prometheus"`
        
        Tracing struct {
            Enabled bool   `yaml:"enabled"`
            Jaeger  string `yaml:"jaeger"`
        } `yaml:"tracing"`
    } `yaml:"monitoring"`
    
    Logging struct {
        Level     string `yaml:"level"`
        Format    string `yaml:"format"`
        Output    string `yaml:"output"`
        Timestamp bool   `yaml:"timestamp"`
    } `yaml:"logging"`
}
```

### Configuration Files

```yaml
# config/default.yaml
app:
  name: "config-service"
  version: "1.0.0"
  environment: "development"
  shutdown_timeout: 30s

server:
  http:
    host: "0.0.0.0"
    port: 8080
  grpc:
    host: "0.0.0.0"
    port: 9090

database:
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  name: "myapp"
  ssl: false
  timeout: 10s
  max_connections: 10

cache:
  redis:
    host: "localhost"
    port: 6379
    password: ""
    db: 0
    timeout: 5s
  in_memory:
    expiration: 5m
    cleanup_interval: 10m

monitoring:
  prometheus:
    enabled: true
    port: 9090
    path: "/metrics"
  tracing:
    enabled: false
    jaeger: "http://localhost:14268/api/traces"

logging:
  level: "info"
  format: "json"
  output: "stdout"
  timestamp: true
```

```yaml
# config/production.yaml
app:
  environment: "production"
  shutdown_timeout: 60s

server:
  http:
    port: 80
  grpc:
    port: 443

database:
  host: "${DB_HOST}"
  port: "${DB_PORT}"
  username: "${DB_USERNAME}"
  password: "${DB_PASSWORD}"
  ssl: true
  max_connections: 100

cache:
  redis:
    host: "${REDIS_HOST}"
    port: "${REDIS_PORT}"
    password: "${REDIS_PASSWORD}"

monitoring:
  tracing:
    enabled: true
    jaeger: "${JAEGER_ENDPOINT}"

logging:
  level: "warn"
  output: "file"
```

### Configuration Manager

```go
// config_manager.go
package main

import (
    "fmt"
    "os"
    "strings"
    
    "github.com/kumarabd/gokit/config"
    "github.com/kumarabd/gokit/errors"
)

type ConfigManager struct {
    config *Config
}

func NewConfigManager() (*ConfigManager, error) {
    var cfg Config
    
    // Load configuration
    configObj, err := config.New(&cfg)
    if err != nil {
        return nil, fmt.Errorf("failed to load configuration: %w", err)
    }
    
    // Validate configuration
    if err := validateConfig(&cfg); err != nil {
        return nil, fmt.Errorf("invalid configuration: %w", err)
    }
    
    // Set defaults
    setDefaults(&cfg)
    
    return &ConfigManager{
        config: &cfg,
    }, nil
}

func (cm *ConfigManager) GetConfig() *Config {
    return cm.config
}

func (cm *ConfigManager) GetDatabaseDSN() string {
    cfg := cm.config.Database
    
    sslMode := "disable"
    if cfg.SSL {
        sslMode = "require"
    }
    
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name, sslMode)
}

func (cm *ConfigManager) GetHTTPAddr() string {
    return fmt.Sprintf("%s:%d", cm.config.Server.HTTP.Host, cm.config.Server.HTTP.Port)
}

func (cm *ConfigManager) GetGRPCAddr() string {
    return fmt.Sprintf("%s:%d", cm.config.Server.GRPC.Host, cm.config.Server.GRPC.Port)
}

func validateConfig(cfg *Config) error {
    if cfg.App.Name == "" {
        return errors.New("CONFIG_ERROR", errors.Alert, "App name is required")
    }
    
    if cfg.Server.HTTP.Port <= 0 || cfg.Server.HTTP.Port > 65535 {
        return errors.New("CONFIG_ERROR", errors.Alert, "Invalid HTTP port")
    }
    
    if cfg.Server.GRPC.Port <= 0 || cfg.Server.GRPC.Port > 65535 {
        return errors.New("CONFIG_ERROR", errors.Alert, "Invalid gRPC port")
    }
    
    if cfg.Database.Host == "" {
        return errors.New("CONFIG_ERROR", errors.Alert, "Database host is required")
    }
    
    if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
        return errors.New("CONFIG_ERROR", errors.Alert, "Invalid database port")
    }
    
    return nil
}

func setDefaults(cfg *Config) {
    if cfg.App.Environment == "" {
        cfg.App.Environment = "development"
    }
    
    if cfg.App.ShutdownTimeout == 0 {
        cfg.App.ShutdownTimeout = 30 * time.Second
    }
    
    if cfg.Server.HTTP.Host == "" {
        cfg.Server.HTTP.Host = "localhost"
    }
    
    if cfg.Server.GRPC.Host == "" {
        cfg.Server.GRPC.Host = "localhost"
    }
    
    if cfg.Database.Timeout == 0 {
        cfg.Database.Timeout = 10 * time.Second
    }
    
    if cfg.Database.MaxConnections == 0 {
        cfg.Database.MaxConnections = 10
    }
    
    if cfg.Cache.InMemory.Expiration == 0 {
        cfg.Cache.InMemory.Expiration = 5 * time.Minute
    }
    
    if cfg.Cache.InMemory.CleanupInterval == 0 {
        cfg.Cache.InMemory.CleanupInterval = 10 * time.Minute
    }
    
    if cfg.Logging.Level == "" {
        cfg.Logging.Level = "info"
    }
    
    if cfg.Logging.Format == "" {
        cfg.Logging.Format = "json"
    }
    
    if cfg.Logging.Output == "" {
        cfg.Logging.Output = "stdout"
    }
}

func (cm *ConfigManager) PrintConfig() {
    fmt.Printf("Configuration:\n")
    fmt.Printf("  App: %s v%s (%s)\n", 
        cm.config.App.Name, 
        cm.config.App.Version, 
        cm.config.App.Environment)
    fmt.Printf("  HTTP Server: %s\n", cm.GetHTTPAddr())
    fmt.Printf("  gRPC Server: %s\n", cm.GetGRPCAddr())
    fmt.Printf("  Database: %s:%d\n", 
        cm.config.Database.Host, 
        cm.config.Database.Port)
    fmt.Printf("  Logging: %s (%s)\n", 
        cm.config.Logging.Level, 
        cm.config.Logging.Format)
}
```

### Usage Example

```go
// main.go
package main

import (
    "log"
)

func main() {
    // Load configuration
    configManager, err := NewConfigManager()
    if err != nil {
        log.Fatal("Failed to load configuration:", err)
    }
    
    // Print configuration
    configManager.PrintConfig()
    
    // Use configuration
    cfg := configManager.GetConfig()
    
    fmt.Printf("Database DSN: %s\n", configManager.GetDatabaseDSN())
    fmt.Printf("HTTP Address: %s\n", configManager.GetHTTPAddr())
    fmt.Printf("gRPC Address: %s\n", configManager.GetGRPCAddr())
}
```

### Running with Different Configurations

```bash
# Development
go run main.go --config config/default.yaml

# Production with environment variables
export DB_HOST=prod-db.example.com
export DB_PORT=5432
export DB_USERNAME=prod_user
export DB_PASSWORD=secret123
export REDIS_HOST=prod-redis.example.com
export JAEGER_ENDPOINT=http://jaeger:14268/api/traces

go run main.go --config config/production.yaml \
  --from-env database.host::DB_HOST \
  --from-env database.port::DB_PORT \
  --from-env database.username::DB_USERNAME \
  --from-env database.password::DB_PASSWORD \
  --from-env cache.redis.host::REDIS_HOST \
  --from-env monitoring.tracing.jaeger::JAEGER_ENDPOINT

# Override specific values
go run main.go --config config/default.yaml \
  --server.http.port 9000 \
  --logging.level debug \
  --monitoring.prometheus.enabled false
```

## Logging and Error Handling

A comprehensive example demonstrating structured logging and error handling.

### Error Definitions

```go
// errors.go
package main

import (
    "github.com/kumarabd/gokit/errors"
)

// Error codes
const (
    ErrCodeUserNotFound     = "USER_NOT_FOUND"
    ErrCodeUserExists       = "USER_EXISTS"
    ErrCodeInvalidInput     = "INVALID_INPUT"
    ErrCodeDatabaseError    = "DATABASE_ERROR"
    ErrCodeValidationError  = "VALIDATION_ERROR"
    ErrCodeAuthentication   = "AUTHENTICATION_ERROR"
    ErrCodeAuthorization    = "AUTHORIZATION_ERROR"
    ErrCodeRateLimit        = "RATE_LIMIT_ERROR"
    ErrCodeInternalError    = "INTERNAL_ERROR"
)

// Error constructors
func NewUserNotFoundError(userID string) error {
    return errors.New(ErrCodeUserNotFound, errors.Warn, 
        "User not found:", userID)
}

func NewUserExistsError(email string) error {
    return errors.New(ErrCodeUserExists, errors.Alert, 
        "User already exists:", email)
}

func NewInvalidInputError(field, value string) error {
    return errors.New(ErrCodeInvalidInput, errors.Warn, 
        "Invalid input for field:", field, "value:", value)
}

func NewDatabaseError(operation string, cause error) error {
    return errors.New(ErrCodeDatabaseError, errors.Critical, 
        "Database error during:", operation, "cause:", cause)
}

func NewValidationError(field, reason string) error {
    return errors.New(ErrCodeValidationError, errors.Warn, 
        "Validation failed for field:", field, "reason:", reason)
}

func NewAuthenticationError(reason string) error {
    return errors.New(ErrCodeAuthentication, errors.Alert, 
        "Authentication failed:", reason)
}

func NewAuthorizationError(resource, action string) error {
    return errors.New(ErrCodeAuthorization, errors.Alert, 
        "Authorization failed for resource:", resource, "action:", action)
}

func NewRateLimitError(limit, window string) error {
    return errors.New(ErrCodeRateLimit, errors.Warn, 
        "Rate limit exceeded:", limit, "per", window)
}

func NewInternalError(operation string, cause error) error {
    return errors.New(ErrCodeInternalError, errors.Critical, 
        "Internal error during:", operation, "cause:", cause)
}
```

### Service with Error Handling

```go
// user_service.go
package main

import (
    "context"
    "time"
    "regexp"
    
    "github.com/kumarabd/gokit/logger"
    "github.com/kumarabd/gokit/errors"
)

type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
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

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    // Start operation span
    s.log.Info().
        Str("operation", "create_user").
        Str("user_email", user.Email).
        Msg("Creating user")
    
    // Validate input
    if err := s.validateUser(user); err != nil {
        s.log.Error().
            Err(err).
            Str("operation", "create_user").
            Str("user_email", user.Email).
            Msg("User validation failed")
        return err
    }
    
    // Check if user exists
    if s.userExists(user.Email) {
        err := NewUserExistsError(user.Email)
        s.log.Warn().
            Err(err).
            Str("operation", "create_user").
            Str("user_email", user.Email).
            Msg("User already exists")
        return err
    }
    
    // Simulate database operation
    if err := s.saveUserToDatabase(user); err != nil {
        s.log.Error().
            Err(err).
            Str("operation", "create_user").
            Str("user_email", user.Email).
            Msg("Failed to save user to database")
        return err
    }
    
    s.log.Info().
        Str("operation", "create_user").
        Str("user_id", user.ID).
        Str("user_email", user.Email).
        Msg("User created successfully")
    
    return nil
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    s.log.Info().
        Str("operation", "get_user").
        Str("user_id", userID).
        Msg("Getting user")
    
    // Validate user ID
    if userID == "" {
        err := NewInvalidInputError("user_id", userID)
        s.log.Error().
            Err(err).
            Str("operation", "get_user").
            Msg("Invalid user ID")
        return nil, err
    }
    
    // Simulate database query
    user, err := s.getUserFromDatabase(userID)
    if err != nil {
        if errors.Is(err) && errors.GetCode(err) == ErrCodeUserNotFound {
            s.log.Warn().
                Err(err).
                Str("operation", "get_user").
                Str("user_id", userID).
                Msg("User not found")
        } else {
            s.log.Error().
                Err(err).
                Str("operation", "get_user").
                Str("user_id", userID).
                Msg("Database error")
        }
        return nil, err
    }
    
    s.log.Info().
        Str("operation", "get_user").
        Str("user_id", userID).
        Str("user_name", user.Name).
        Msg("User retrieved successfully")
    
    return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
    s.log.Info().
        Str("operation", "update_user").
        Str("user_id", userID).
        Msg("Updating user")
    
    // Validate user ID
    if userID == "" {
        err := NewInvalidInputError("user_id", userID)
        s.log.Error().
            Err(err).
            Str("operation", "update_user").
            Msg("Invalid user ID")
        return err
    }
    
    // Check if user exists
    if _, err := s.getUserFromDatabase(userID); err != nil {
        if errors.Is(err) && errors.GetCode(err) == ErrCodeUserNotFound {
            s.log.Warn().
                Err(err).
                Str("operation", "update_user").
                Str("user_id", userID).
                Msg("User not found")
        }
        return err
    }
    
    // Validate updates
    if err := s.validateUpdates(updates); err != nil {
        s.log.Error().
            Err(err).
            Str("operation", "update_user").
            Str("user_id", userID).
            Msg("Update validation failed")
        return err
    }
    
    // Simulate database update
    if err := s.updateUserInDatabase(userID, updates); err != nil {
        s.log.Error().
            Err(err).
            Str("operation", "update_user").
            Str("user_id", userID).
            Msg("Failed to update user in database")
        return err
    }
    
    s.log.Info().
        Str("operation", "update_user").
        Str("user_id", userID).
        Msg("User updated successfully")
    
    return nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
    s.log.Info().
        Str("operation", "delete_user").
        Str("user_id", userID).
        Msg("Deleting user")
    
    // Validate user ID
    if userID == "" {
        err := NewInvalidInputError("user_id", userID)
        s.log.Error().
            Err(err).
            Str("operation", "delete_user").
            Msg("Invalid user ID")
        return err
    }
    
    // Check if user exists
    if _, err := s.getUserFromDatabase(userID); err != nil {
        if errors.Is(err) && errors.GetCode(err) == ErrCodeUserNotFound {
            s.log.Warn().
                Err(err).
                Str("operation", "delete_user").
                Str("user_id", userID).
                Msg("User not found")
        }
        return err
    }
    
    // Simulate database delete
    if err := s.deleteUserFromDatabase(userID); err != nil {
        s.log.Error().
            Err(err).
            Str("operation", "delete_user").
            Str("user_id", userID).
            Msg("Failed to delete user from database")
        return err
    }
    
    s.log.Info().
        Str("operation", "delete_user").
        Str("user_id", userID).
        Msg("User deleted successfully")
    
    return nil
}

func (s *UserService) validateUser(user *User) error {
    if user.Name == "" {
        return NewValidationError("name", "Name is required")
    }
    
    if user.Email == "" {
        return NewValidationError("email", "Email is required")
    }
    
    // Validate email format
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(user.Email) {
        return NewValidationError("email", "Invalid email format")
    }
    
    return nil
}

func (s *UserService) validateUpdates(updates map[string]interface{}) error {
    for field, value := range updates {
        switch field {
        case "name":
            if name, ok := value.(string); ok && name == "" {
                return NewValidationError("name", "Name cannot be empty")
            }
        case "email":
            if email, ok := value.(string); ok {
                if email == "" {
                    return NewValidationError("email", "Email cannot be empty")
                }
                emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
                if !emailRegex.MatchString(email) {
                    return NewValidationError("email", "Invalid email format")
                }
            }
        }
    }
    
    return nil
}

func (s *UserService) userExists(email string) bool {
    // Simulate database check
    return email == "existing@example.com"
}

func (s *UserService) saveUserToDatabase(user *User) error {
    // Simulate database operation
    if user.Email == "error@example.com" {
        return NewDatabaseError("save_user", fmt.Errorf("connection timeout"))
    }
    
    user.ID = "generated-id"
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    return nil
}

func (s *UserService) getUserFromDatabase(userID string) (*User, error) {
    // Simulate database operation
    if userID == "notfound" {
        return nil, NewUserNotFoundError(userID)
    }
    
    if userID == "error" {
        return nil, NewDatabaseError("get_user", fmt.Errorf("connection timeout"))
    }
    
    return &User{
        ID:        userID,
        Name:      "John Doe",
        Email:     "john@example.com",
        CreatedAt: time.Now().Add(-24 * time.Hour),
        UpdatedAt: time.Now(),
    }, nil
}

func (s *UserService) updateUserInDatabase(userID string, updates map[string]interface{}) error {
    // Simulate database operation
    if userID == "error" {
        return NewDatabaseError("update_user", fmt.Errorf("connection timeout"))
    }
    
    return nil
}

func (s *UserService) deleteUserFromDatabase(userID string) error {
    // Simulate database operation
    if userID == "error" {
        return NewDatabaseError("delete_user", fmt.Errorf("connection timeout"))
    }
    
    return nil
}
```

### HTTP Handler with Error Handling

```go
// handler.go
package main

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "github.com/kumarabd/gokit/errors"
)

type ErrorResponse struct {
    Error struct {
        Code    string `json:"code"`
        Message string `json:"message"`
        Severity string `json:"severity"`
    } `json:"error"`
}

type UserService struct {
    service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        h.handleError(w, NewInvalidInputError("body", "Invalid JSON"), http.StatusBadRequest)
        return
    }
    
    if err := h.service.CreateUser(r.Context(), &user); err != nil {
        h.handleServiceError(w, err)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    if userID == "" {
        h.handleError(w, NewInvalidInputError("id", "User ID is required"), http.StatusBadRequest)
        return
    }
    
    user, err := h.service.GetUser(r.Context(), userID)
    if err != nil {
        h.handleServiceError(w, err)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    if userID == "" {
        h.handleError(w, NewInvalidInputError("id", "User ID is required"), http.StatusBadRequest)
        return
    }
    
    var updates map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
        h.handleError(w, NewInvalidInputError("body", "Invalid JSON"), http.StatusBadRequest)
        return
    }
    
    if err := h.service.UpdateUser(r.Context(), userID, updates); err != nil {
        h.handleServiceError(w, err)
        return
    }
    
    w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    if userID == "" {
        h.handleError(w, NewInvalidInputError("id", "User ID is required"), http.StatusBadRequest)
        return
    }
    
    if err := h.service.DeleteUser(r.Context(), userID); err != nil {
        h.handleServiceError(w, err)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) handleServiceError(w http.ResponseWriter, err error) {
    if !errors.Is(err) {
        h.handleError(w, NewInternalError("service_operation", err), http.StatusInternalServerError)
        return
    }
    
    code := errors.GetCode(err)
    severity := errors.GetSeverity(err)
    
    var statusCode int
    switch code {
    case ErrCodeUserNotFound:
        statusCode = http.StatusNotFound
    case ErrCodeUserExists:
        statusCode = http.StatusConflict
    case ErrCodeInvalidInput, ErrCodeValidationError:
        statusCode = http.StatusBadRequest
    case ErrCodeAuthentication:
        statusCode = http.StatusUnauthorized
    case ErrCodeAuthorization:
        statusCode = http.StatusForbidden
    case ErrCodeRateLimit:
        statusCode = http.StatusTooManyRequests
    case ErrCodeDatabaseError, ErrCodeInternalError:
        statusCode = http.StatusInternalServerError
    default:
        statusCode = http.StatusInternalServerError
    }
    
    h.handleError(w, err, statusCode)
}

func (h *UserHandler) handleError(w http.ResponseWriter, err error, statusCode int) {
    var response ErrorResponse
    
    if errors.Is(err) {
        response.Error.Code = errors.GetCode(err)
        response.Error.Message = err.Error()
        response.Error.Severity = string(errors.GetSeverity(err))
    } else {
        response.Error.Code = "UNKNOWN_ERROR"
        response.Error.Message = err.Error()
        response.Error.Severity = "critical"
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(response)
}
```

### Main Application

```go
// main.go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // Create user service
    service, err := NewUserService()
    if err != nil {
        log.Fatal("Failed to create user service:", err)
    }
    
    // Create handler
    handler := NewUserHandler(service)
    
    // Create HTTP server
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "POST":
            handler.CreateUser(w, r)
        case "GET":
            handler.GetUser(w, r)
        case "PUT":
            handler.UpdateUser(w, r)
        case "DELETE":
            handler.DeleteUser(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
    
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    
    // Start server
    go func() {
        service.log.Info().Msg("Starting HTTP server on :8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            service.log.Error().Err(err).Msg("HTTP server error")
            log.Fatal(err)
        }
    }()
    
    // Wait for interrupt signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigCh
    service.log.Info().Msg("Shutting down server")
    
    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        service.log.Error().Err(err).Msg("Error during server shutdown")
    }
    
    service.log.Info().Msg("Server stopped")
}
```

### Testing the Error Handling

```bash
# Start the server
go run main.go

# Test successful operations
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

curl http://localhost:8080/api/users?id=123

# Test error scenarios
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"","email":"invalid-email"}'

curl http://localhost:8080/api/users?id=notfound

curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"existing@example.com"}'

curl http://localhost:8080/api/users?id=error
```

This example demonstrates comprehensive error handling with structured logging, proper error codes, and appropriate HTTP status codes.
