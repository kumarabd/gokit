# Caching

GoKit provides a flexible caching system with an interface-based design that supports multiple cache backends. Currently, it includes an in-memory cache implementation with TTL (Time To Live) support.

## Overview

The caching system provides:
- **Interface-based design** for easy extension and testing
- **In-memory cache** with automatic expiration
- **TTL support** for cache entries
- **Thread-safe operations**
- **Simple and consistent API**

## Basic Usage

### 1. Initialize Cache

```go
package main

import (
    "time"
    "github.com/kumarabd/gokit/cache"
    "github.com/kumarabd/gokit/cache/inmem"
)

func main() {
    // Create cache options
    opts := inmem.Options{
        Expiration:      5 * time.Minute,  // Default TTL for entries
        CleanupInterval: 10 * time.Minute, // How often to clean expired entries
    }
    
    // Initialize cache
    cache, err := inmem.New(opts)
    if err != nil {
        panic(err)
    }
    
    // Use the cache
    err = cache.Set("user:123", "John Doe", 1*time.Hour)
    if err != nil {
        panic(err)
    }
    
    value, err := cache.Get("user:123")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("User: %v\n", value)
}
```

### 2. Cache Interface

```go
type Handler interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, exp ...time.Duration) error
}
```

## Cache Operations

### Setting Values

```go
// Set with default expiration (from cache options)
err := cache.Set("key1", "value1")

// Set with custom TTL
err = cache.Set("key2", "value2", 30*time.Minute)

// Set with no expiration
err = cache.Set("key3", "value3", 0)

// Set complex data types
user := User{ID: "123", Name: "John"}
err = cache.Set("user:123", user, 1*time.Hour)

// Set with very short TTL for testing
err = cache.Set("temp:data", "temporary", 5*time.Second)
```

### Getting Values

```go
// Get value
value, err := cache.Get("key1")
if err != nil {
    if err == inmem.ErrKeyNotExist {
        fmt.Println("Key not found or expired")
    } else {
        fmt.Println("Cache error:", err)
    }
    return
}

// Type assertion for specific types
if str, ok := value.(string); ok {
    fmt.Printf("String value: %s\n", str)
}

// Get user object
if userValue, err := cache.Get("user:123"); err == nil {
    if user, ok := userValue.(User); ok {
        fmt.Printf("User: %+v\n", user)
    }
}
```

## Cache Configuration

### In-Memory Cache Options

```go
type Options struct {
    Expiration      time.Duration // Default TTL for cache entries
    CleanupInterval time.Duration // Interval for cleaning expired entries
}
```

### Common Configuration Patterns

```go
// Short-lived cache (for API responses)
shortCache, _ := inmem.New(inmem.Options{
    Expiration:      1 * time.Minute,
    CleanupInterval: 5 * time.Minute,
})

// Long-lived cache (for user sessions)
sessionCache, _ := inmem.New(inmem.Options{
    Expiration:      24 * time.Hour,
    CleanupInterval: 1 * time.Hour,
})

// No expiration cache (for configuration)
configCache, _ := inmem.New(inmem.Options{
    Expiration:      0, // No expiration
    CleanupInterval: 0, // No cleanup needed
})
```

## Error Handling

### Cache Errors

```go
import "github.com/kumarabd/gokit/cache/inmem"

// Check for specific cache errors
value, err := cache.Get("nonexistent")
if err != nil {
    switch err {
    case inmem.ErrKeyNotExist:
        fmt.Println("Key does not exist or has expired")
    default:
        fmt.Println("Unexpected cache error:", err)
    }
    return
}
```

### Error Types

```go
// In-memory cache specific errors
var (
    ErrKeyNotExist = errors.New("", errors.Alert, "Key does not exist")
)
```

## Best Practices

### 1. Use Descriptive Key Names

```go
// Good - descriptive and namespaced
cache.Set("user:profile:123", userProfile)
cache.Set("api:response:users:list", userList)
cache.Set("config:database:connection", dbConfig)

// Avoid - generic keys
cache.Set("key1", value1)
cache.Set("data", data)
```

### 2. Appropriate TTL Values

```go
// Static data (configuration, user profiles)
cache.Set("config:app", config, 24*time.Hour)

// Semi-static data (user sessions)
cache.Set("session:123", session, 1*time.Hour)

// Dynamic data (API responses)
cache.Set("api:users:list", users, 5*time.Minute)

// Temporary data (rate limiting)
cache.Set("rate:limit:user:123", count, 1*time.Minute)
```

### 3. Handle Cache Misses Gracefully

```go
func getUserProfile(userID string) (*UserProfile, error) {
    // Try to get from cache first
    if cached, err := cache.Get("user:profile:" + userID); err == nil {
        if profile, ok := cached.(*UserProfile); ok {
            return profile, nil
        }
    }
    
    // Cache miss - fetch from database
    profile, err := fetchUserProfileFromDB(userID)
    if err != nil {
        return nil, err
    }
    
    // Store in cache for next time
    cache.Set("user:profile:"+userID, profile, 1*time.Hour)
    
    return profile, nil
}
```

### 4. Cache Warming

```go
func warmCache() {
    // Pre-load frequently accessed data
    users, _ := fetchAllActiveUsers()
    for _, user := range users {
        cache.Set("user:profile:"+user.ID, user, 1*time.Hour)
    }
    
    // Pre-load configuration
    config, _ := loadConfiguration()
    cache.Set("config:app", config, 24*time.Hour)
}
```

## Complete Example

```go
package main

import (
    "fmt"
    "time"
    "encoding/json"
    
    "github.com/kumarabd/gokit/cache"
    "github.com/kumarabd/gokit/cache/inmem"
    "github.com/kumarabd/gokit/logger"
)

type User struct {
    ID       string    `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Created  time.Time `json:"created"`
}

type UserService struct {
    cache cache.Handler
    log   *logger.Handler
}

func NewUserService() (*UserService, error) {
    // Initialize cache
    cacheOpts := inmem.Options{
        Expiration:      30 * time.Minute,
        CleanupInterval: 5 * time.Minute,
    }
    
    cache, err := inmem.New(cacheOpts)
    if err != nil {
        return nil, err
    }
    
    // Initialize logger
    log, err := logger.New("user-service", logger.Options{
        Format:     logger.JSONLogFormat,
        DebugLevel: true,
    })
    if err != nil {
        return nil, err
    }
    
    return &UserService{
        cache: cache,
        log:   log,
    }, nil
}

func (s *UserService) GetUser(userID string) (*User, error) {
    cacheKey := "user:" + userID
    
    // Try cache first
    if cached, err := s.cache.Get(cacheKey); err == nil {
        if user, ok := cached.(*User); ok {
            s.log.Info().
                Str("user_id", userID).
                Str("source", "cache").
                Msg("User retrieved from cache")
            return user, nil
        }
    }
    
    // Cache miss - fetch from database
    s.log.Info().
        Str("user_id", userID).
        Str("source", "database").
        Msg("User not in cache, fetching from database")
    
    user, err := s.fetchUserFromDB(userID)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    err = s.cache.Set(cacheKey, user, 30*time.Minute)
    if err != nil {
        s.log.Warn().
            Str("user_id", userID).
            Err(err).
            Msg("Failed to cache user")
    }
    
    return user, nil
}

func (s *UserService) UpdateUser(userID string, updates map[string]interface{}) error {
    // Update in database
    err := s.updateUserInDB(userID, updates)
    if err != nil {
        return err
    }
    
    // Invalidate cache
    cacheKey := "user:" + userID
    // Note: The current implementation doesn't have a Delete method,
    // but we can set the value to nil or use a very short TTL
    s.cache.Set(cacheKey, nil, 1*time.Second)
    
    s.log.Info().
        Str("user_id", userID).
        Msg("User updated and cache invalidated")
    
    return nil
}

func (s *UserService) GetUserList() ([]*User, error) {
    cacheKey := "users:list"
    
    // Try cache first
    if cached, err := s.cache.Get(cacheKey); err == nil {
        if users, ok := cached.([]*User); ok {
            s.log.Info().
                Str("source", "cache").
                Int("count", len(users)).
                Msg("User list retrieved from cache")
            return users, nil
        }
    }
    
    // Cache miss - fetch from database
    users, err := s.fetchUsersFromDB()
    if err != nil {
        return nil, err
    }
    
    // Store in cache with shorter TTL for lists
    s.cache.Set(cacheKey, users, 5*time.Minute)
    
    s.log.Info().
        Str("source", "database").
        Int("count", len(users)).
        Msg("User list fetched from database and cached")
    
    return users, nil
}

// Simulate database operations
func (s *UserService) fetchUserFromDB(userID string) (*User, error) {
    // Simulate database delay
    time.Sleep(100 * time.Millisecond)
    
    return &User{
        ID:      userID,
        Name:    "John Doe",
        Email:   "john@example.com",
        Created: time.Now(),
    }, nil
}

func (s *UserService) updateUserInDB(userID string, updates map[string]interface{}) error {
    // Simulate database update
    time.Sleep(50 * time.Millisecond)
    return nil
}

func (s *UserService) fetchUsersFromDB() ([]*User, error) {
    // Simulate database delay
    time.Sleep(200 * time.Millisecond)
    
    return []*User{
        {ID: "1", Name: "John Doe", Email: "john@example.com", Created: time.Now()},
        {ID: "2", Name: "Jane Smith", Email: "jane@example.com", Created: time.Now()},
    }, nil
}

func main() {
    service, err := NewUserService()
    if err != nil {
        panic(err)
    }
    
    // Test cache functionality
    fmt.Println("=== Testing User Cache ===")
    
    // First call - should hit database
    user1, err := service.GetUser("123")
    if err != nil {
        panic(err)
    }
    fmt.Printf("User 1: %+v\n", user1)
    
    // Second call - should hit cache
    user2, err := service.GetUser("123")
    if err != nil {
        panic(err)
    }
    fmt.Printf("User 2: %+v\n", user2)
    
    // Test user list caching
    fmt.Println("\n=== Testing User List Cache ===")
    
    users1, err := service.GetUserList()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Users 1: %d users\n", len(users1))
    
    users2, err := service.GetUserList()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Users 2: %d users\n", len(users2))
    
    // Test cache invalidation
    fmt.Println("\n=== Testing Cache Invalidation ===")
    
    err = service.UpdateUser("123", map[string]interface{}{
        "name": "John Updated",
    })
    if err != nil {
        panic(err)
    }
    
    // Next call should hit database again
    user3, err := service.GetUser("123")
    if err != nil {
        panic(err)
    }
    fmt.Printf("User 3 (after update): %+v\n", user3)
}
```

## Performance Considerations

### Memory Usage

The in-memory cache stores all data in RAM, so consider:

- **Cache size limits** for large datasets
- **Memory monitoring** in production
- **Appropriate TTL values** to prevent memory bloat

### Thread Safety

The in-memory cache implementation is thread-safe and can be used concurrently.

### Cache Hit Ratios

Monitor cache performance:

```go
type CacheStats struct {
    Hits   int64
    Misses int64
}

func (s *UserService) getCacheStats() CacheStats {
    // Implement cache statistics tracking
    return CacheStats{
        Hits:   atomic.LoadInt64(&s.cacheHits),
        Misses: atomic.LoadInt64(&s.cacheMisses),
    }
}
```

## Future Extensions

The interface-based design allows for easy extension with additional cache backends:

- **Redis cache** for distributed caching
- **File-based cache** for persistence
- **Database cache** for large datasets
- **CDN cache** for static content

This provides flexibility to choose the right caching strategy for your specific use case.
