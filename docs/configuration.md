# Configuration Management

GoKit provides a flexible and powerful configuration management system that supports multiple configuration sources including YAML files, environment variables, and command-line flags.

## Overview

The configuration system automatically generates command-line flags from your configuration struct, supports YAML configuration files, and allows environment variable overrides. It uses reflection to dynamically create flags based on your struct tags and field types.

## Basic Usage

### 1. Define Your Configuration Struct

```go
type AppConfig struct {
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
    
    LogLevel string `yaml:"log_level"`
    Debug    bool   `yaml:"debug"`
}
```

### 2. Initialize Configuration

```go
func main() {
    var cfg AppConfig
    configObj, err := config.New(&cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    // Use the configuration
    fmt.Printf("Server will run on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
}
```

## Configuration Sources

### 1. YAML Configuration Files

Create a configuration file (e.g., `config.yaml`):

```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  host: "localhost"
  port: 5432
  username: "myuser"
  password: "mypassword"
  ssl: true

log_level: "info"
debug: false
```

Run your application with the config file:

```bash
./myapp --config config.yaml
```

### 2. Command-Line Flags

The configuration system automatically generates flags for all exported fields:

```bash
# Set individual values
./myapp --server.host 127.0.0.1 --server.port 9000

# Set nested configuration
./myapp --database.host localhost --database.port 5432 --database.ssl

# Set boolean flags
./myapp --debug

# Set string values
./myapp --log_level debug
```

### 3. Environment Variables

Use the `--from-env` flag to map environment variables to configuration paths:

```bash
# Set environment variables
export DB_HOST=production-db.example.com
export DB_PORT=5432
export API_KEY=secret123

# Map them to configuration
./myapp --from-env database.host::DB_HOST --from-env database.port::DB_PORT --from-env api_key::API_KEY
```

### 4. Environment Variable References in YAML

You can also reference environment variables directly in your YAML file:

```yaml
database:
  host: "${DB_HOST}"
  port: "${DB_PORT}"
  password: "${DB_PASSWORD}"

api:
  key: "${API_KEY}"
```

## Supported Data Types

The configuration system supports the following Go types:

- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `bool`
- `float32`, `float64`
- Structs (nested configuration)
- Pointers to any of the above types

## Advanced Features

### Nested Configuration

```go
type Config struct {
    App struct {
        Name    string `yaml:"name"`
        Version string `yaml:"version"`
    } `yaml:"app"`
    
    Services struct {
        Auth struct {
            URL     string `yaml:"url"`
            Timeout int    `yaml:"timeout"`
        } `yaml:"auth"`
        
        Payment struct {
            URL     string `yaml:"url"`
            Timeout int    `yaml:"timeout"`
        } `yaml:"payment"`
    } `yaml:"services"`
}
```

### Pointer Fields

```go
type Config struct {
    Optional *struct {
        Feature string `yaml:"feature"`
        Enabled bool   `yaml:"enabled"`
    } `yaml:"optional"`
}
```

### Custom YAML Tags

```go
type Config struct {
    // Use custom YAML field names
    ServerHost string `yaml:"server_host"`
    ServerPort int    `yaml:"server_port"`
    
    // Use JSON tags as fallback
    DatabaseURL string `json:"database_url" yaml:"database_url"`
}
```

## Error Handling

The configuration system provides detailed error messages for common issues:

```go
configObj, err := config.New(&cfg)
if err != nil {
    // Handle configuration errors
    switch {
    case strings.Contains(err.Error(), "failed to read config file"):
        log.Fatal("Configuration file not found or unreadable")
    case strings.Contains(err.Error(), "failed to unmarshal config file"):
        log.Fatal("Invalid YAML format in configuration file")
    default:
        log.Fatal("Configuration error:", err)
    }
}
```

## Best Practices

### 1. Use Descriptive Field Names

```go
// Good
type Config struct {
    DatabaseHost string `yaml:"database_host"`
    DatabasePort int    `yaml:"database_port"`
}

// Avoid
type Config struct {
    Host string `yaml:"host"` // Too generic
    Port int    `yaml:"port"`
}
```

### 2. Group Related Configuration

```go
type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    
    Database struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"database"`
}
```

### 3. Use Environment Variables for Secrets

```bash
# Never put secrets in YAML files
# Instead, use environment variables
export DB_PASSWORD=secret123
./myapp --from-env database.password::DB_PASSWORD
```

### 4. Provide Sensible Defaults

```go
type Config struct {
    Server struct {
        Host string `yaml:"host"` // Will default to empty string
        Port int    `yaml:"port"` // Will default to 0
    } `yaml:"server"`
}

// Set defaults after loading
if cfg.Server.Host == "" {
    cfg.Server.Host = "localhost"
}
if cfg.Server.Port == 0 {
    cfg.Server.Port = 8080
}
```

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/kumarabd/gokit/config"
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
    
    LogLevel string `yaml:"log_level"`
    Debug    bool   `yaml:"debug"`
}

func main() {
    var cfg AppConfig
    
    // Load configuration
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
    if cfg.LogLevel == "" {
        cfg.LogLevel = "info"
    }
    
    // Use configuration
    fmt.Printf("Starting %s v%s\n", cfg.App.Name, cfg.App.Version)
    fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
    fmt.Printf("Database: %s:%d (SSL: %t)\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.SSL)
    fmt.Printf("Log Level: %s, Debug: %t\n", cfg.LogLevel, cfg.Debug)
}
```

Run with different configuration sources:

```bash
# With config file
./myapp --config config.yaml

# With command-line flags
./myapp --app.name "My Service" --server.port 9000 --debug

# With environment variables
export DB_HOST=prod-db.example.com
export DB_PASSWORD=secret123
./myapp --from-env database.host::DB_HOST --from-env database.password::DB_PASSWORD
```
