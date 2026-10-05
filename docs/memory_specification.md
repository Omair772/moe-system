# 💾 Memory Specification — Ministry of Education System

## 1. Purpose & Scope

This document defines the memory management strategy for the Ministry of Education integrated system, covering data persistence, caching strategies, session management, and memory optimization patterns using Go, Redis with TTL, and PostgreSQL with Prisma ORM.

## 2. Memory Layers

### 2.1 Data Persistence Layer (PostgreSQL)

| Data Type | Table | Retention | Backup Strategy | Indexes |
|-----------|-------|-----------|-----------------|---------|
| Students | students | 7 years (GDPR) | Daily full, hourly incremental | Primary key (id), Unique (email), Index (created_at) |
| Courses | courses | 7 years | Daily full, hourly incremental | Primary key (id), Unique (code), Index (title) |
| Enrollments | enrollments | 7 years | Daily full, hourly incremental | Primary key (id), Unique (studentId,courseId), Index (studentId) |
| Teachers | teachers | 7 years | Daily full, hourly incremental | Primary key (id), Unique (email) |
| Attendance | attendance | 3 years | Daily full | Primary key (id), Index (studentId, date) |
| AuditLogs | audit_logs | 10 years (legal) | Weekly full, daily transaction log | Primary key (id), Index (timestamp), Index (user_id) |

### 2.2 Cache Layer (Redis with TTL)

| Cache Key | TTL | Data Type | Eviction Policy | Purpose |
|-----------|-----|-----------|-----------------|---------|
| `student:{id}` | 30 min | Hash | LRU (Least Recently Used) | Student profile caching |
| `course:{code}` | 2 hours | Hash | LRU | Course catalog caching |
| `attendance:{date}` | 5 min | List | LRU | Daily attendance cache |
| `session:{token}` | 24 hours | String | TTL expiration | Auth session validation |
| `api:response:{path}` | 1 min | String | TTL expiration | API response caching |
| `gpa:{studentId}` | 1 hour | String | TTL | Cached GPA calculation |
| `report:{id}` | 30 min | String | TTL | Report generation cache |

### 2.3 Session Memory (In-Memory)

| Property | Value | Description |
|----------|-------|-------------|
| Token Expiry | 2 hours | JWT access token validity |
| Refresh Token TTL | 30 days | Optional refresh token validity |
| Max Session Size | 5 KB | Limit per session data |
| Concurrent Sessions | 3 per user | Maximum active sessions |
| Token Revocation | Immediate | Via Redis blacklist |

## 3. Redis Configuration

### 3.1 Connection Settings

```go
// config/redis.go
package config

import (
    "context"
    "log"
    "github.com/redis/go-redis/v9"
    "github.com/google/uuid"
    "time"
)

type RedisConfig struct {
    Addr     string
    Password string
    DB       int
    PoolSize     int
    MinIdleConns int
    MaxRetries  int
}

func NewRedisClient(cfg *RedisConfig) *redis.Client {
    ctx := context.Background()
    client := redis.NewClient(&redis.Options{
        Addr:        cfg.Addr,
        Password:    cfg.Password,
        DB:          cfg.DB,
        PoolSize:    cfg.PoolSize,
        MinIdleConns: cfg.MinIdleConns,
        MaxRetries:  cfg.MaxRetries,
    })
    
    // Verify connection with proper error handling
    if err := client.Ping(ctx).Err(); err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    return client
}
```

**Key Fix:** The connection verification now properly checks the error from `Ping()`:
- **Old code:** `pong := client.Ping(ctx)` followed by `_ = pong` (ignoring errors)
- **New code:** `if err := client.Ping(ctx).Err(); err != nil { log.Fatalf(...) }` (proper error handling)

### 3.2 TTL Settings Configuration
```go
// config/ttl.go
package config

const (
    // Student profile cache - medium priority, 30 minutes
    StudentProfileTTL = 30 * time.Minute // 1800 seconds
    
    // Course catalog - read-heavy, infrequent changes
    CourseCatalogTTL = 2 * time.Hour     // 7200 seconds
    
    // Attendance data - must be fresh
    AttendanceTTL = 5 * time.Minute      // 300 seconds
    
    // Auth session - standard 24 hours
    SessionTTL = 24 * time.Hour          // 86400 seconds
    
    // API response - very short term
    APIResponseTTL = 1 * time.Minute     // 60 seconds
    
    // GPA calculation - medium term
    GPATTL = 1 * time.Hour               // 3600 seconds
    
    // Report cache - short term
    ReportTTL = 30 * time.Minute         // 1800 seconds
)
```

### 3.3 Cache-Aside Pattern (Go Implementation)

```go
// cache/studentCache.go
package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "ministry-education/config"
    "ministry-education/models"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = config.NewRedisClient()

// Get student from cache or DB
func GetStudent(id string) (*models.Student, error) {
    cacheKey := fmt.Sprintf("student:%s", id)
    
    // Try cache first
    val, err := rdb.HGet(ctx, cacheKey, "data").Result()
    if err == nil {
        var student models.Student
        json.Unmarshal([]byte(val), &student)
        return &student, nil
    }
    
    // Cache miss - fetch from DB
    student, err := prisma.Student.FindUnique(
        prisma.Student.ID(id),
    ).Exec(ctx)
    if err != nil {
        return nil, err
    }
    
    // Store in cache with TTL
    go setStudentCache(student)
    
    return &student, nil
}

// Set student cache asynchronously
func setStudentCache(student *models.Student) error {
    cacheKey := fmt.Sprintf("student:%s", student.ID)
    data, _ := json.Marshal(student)
    
    return rdb.HSet(ctx, cacheKey, "data", data).Err()
    // Also set TTL
    // rdb.Expire(ctx, cacheKey, config.StudentProfileTTL)
}
```

### 3.4 Session Management

```go
// auth/jwt.go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    Name   string `json:"name"`
    JWTStandardClaims
}

type JWTStandardClaims struct {
    JWTID        string `json:"jti"`        // Unique token ID
    IssuedAt     int64  `json:"iat"`        // Issued at timestamp
    Expiry       int64  `json:"exp"`        // Expiration timestamp
    NotBefore    int64  `json:"nbf"`        // Not before timestamp
    Issuer       string `json:"iss"`        // Issuer
    Subject      string `json:"sub"`        // Subject (user ID)
    Audience     string `json:"aud"`        // Audience
    SessionIndex int64  `json:"sidx"`       // Session index for revocation
}

// Generate JWT token
func GenerateToken(userID, role, name string) (string, error) {
    now := time.Now()
    expiry := now.Add(2 * time.Hour) // 2 hour expiry
    
    claim := &Claims{
        UserID: userID,
        Role:   role,
        Name:   name,
        JWTStandardClaims: JWTStandardClaims{
            JWTID:        uuid.New().String(),
            IssuedAt:     now.Unix(),
            Expiry:       expiry.Unix(),
            NotBefore:    now.Unix(),
            Issuer:       "ministry-education-system",
            Subject:      userID,
            Audience:     "ministry-education-api",
            SessionIndex: 0,
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    return token.SignedString(jwtKey)
}
```

### 3.5 Session Blacklist (Revocation)

```go
// auth/revocation.go
package auth

import (
    "context"
    "time"
    "minute"
    "ministry-education/config"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = config.NewRedisClient()

func RevokeToken(tokenString string) error {
    claim, err := jwt.ParseWithClaims(
        tokenString,
        &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        },
    )
    
    if err != nil || !claim.Valid {
        return err
    }
    
    jti := claim.Claims.(Claims).JWTID
    
    // Add to blacklist with TTL until expiry
    ctx := context.Background()
    blacklistKey := fmt.Sprintf("token_blacklist:%s", jti)
    expiry := time.Now().Add(2 * time.Hour) // Until original expiry
    
    if err := rdb.Set(ctx, blacklistKey, "revoked", time.Until(expiry)).Err(); err != nil {
        return err
    }
    
    return nil
}

func IsTokenRevoked(tokenString string) bool {
    ctx := context.Background()
    
    isBlacklisted, err := rdb.Exists(ctx, fmt.Sprintf("token_blacklist:%s", extractJTI(tokenString))).Result()
    if err != nil || isBlacklisted == 0 {
        return false // Not blacklisted or error checking
    }
    
    return true
}
```

### 3.6 Memory Optimization Strategies

```go
// optimize/memory.go
package optimize

import "runtime"

// Enable GC profiling
func SetMemoryProfiler() {
    runtime.GC()
    // Runtime memory profiling
    runtime.MemProfileRate = 512 // Sample every 512 allocations
}

// Check memory stats
func PrintMemoryStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Alloc: %v MiB\n", bToMb(m.Alloc))
    fmt.Printf("TotalAlloc: %v MiB\n", bToMb(m.TotalAlloc))
    fmt.Printf("Sys: %v MiB\n", bToMb(m.Sys))
    fmt.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())
    fmt.Printf("GCNext: %v MiB\n", bToMb(m.GCNext))
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
```

### 3.7 Connection Pool Configuration

```go
// db/postgres.go
package db

import (
    "database/sql"
    "log"
    "os"
    "time"
    
    _ "github.com/jackc/pgx/v5/stdlib"
)

func InitPostgres() *sql.DB {
    dsn := os.Getenv("DATABASE_URL")
    
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        log.Fatal("Failed to open database connection:", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(100)      // Maximum open connections
    db.SetMaxIdleConns(50)       // Maximum idle connections
    db.SetConnMaxLifetime(30 * time.Minute) // Connection lifetime
    db.SetConnMaxIdleTime(10 * time.Minute)   // Idle connection lifetime
    
    // Verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    
    log.Println("Database connection pool configured successfully")
    return db
}
```

### 3.8 Cache Invalidation Strategies

```go
// cache/invalidation.go
package cache

import (
    "context"
    "ministry-education/config"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = config.NewRedisClient()

func InvalidateStudentCache(studentID string) error {
    ctx := context.Background()
    keysToInvalidate := []string{
        fmt.Sprintf("student:%s", studentID),
        fmt.Sprintf("gpa:%s", studentID),
        fmt.Sprintf("report:%s", studentID),
    }
    
    // Pipeline for efficient batch deletion
    pipe := rdb.Pipeline(ctx)
    for _, key := range keysToInvalidate {
        pipe.Del(ctx, key)
    }
    
    _, err := pipe.Exec(ctx)
    return err
}

func InvalidateCourseCache(courseCode string) error {
    ctx := context.Background()
    
    keys := []string{
        fmt.Sprintf("course:%s", courseCode),
    }
    
    pipe := rdb.Pipeline(ctx)
    for _, key := range keys {
        pipe.Del(ctx, key)
    }
    
    _, err := pipe.Exec(ctx)
    return err
}
```

### 3.9 Memory Monitoring & Alerting

```go
// monitoring/alerts.go
package monitoring

import (
    "time"
    "ministry-education/config"
    "github.com/redis/go-redis/v9"
    "github.com/prometheus/client_golang/prometheus"
)

var ctx = context.Background()
var rdb = config.NewRedisClient()

// Prometheus metrics for Redis
var (
    redisMemoryUsage = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "redis_memory_usage_bytes",
            Help: "Current Redis memory usage in bytes",
        },
        ["instance"],
    )
    redisCacheHitRatio = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "redis_cache_hit_ratio",
            Help: "Redis cache hit ratio percentage",
        },
        ["instance"],
    )
)

func init() {
    prometheus.Register(redisMemoryUsage)
    prometheus.Register(redisCacheHitRatio)
}

func UpdateRedisMetrics() {
    // Get Redis memory info
    info, err := rdb.Info(ctx, "memory").Result()
    if err == nil {
        // Parse used_memory from info string
        // This is a simplified example
        redisMemoryUsage.WithLabelValues("primary").Set(1024 * 1024 * 500) // 500MB example
    }
    
    // Get cache hit ratio (simplified)
    redisCacheHitRatio.WithLabelValues("primary").Set(75.0) // 75% example
}
```

### 3.9 Data Retention & Cleanup

```go
// jobs/cleanup.go
package jobs

import (
    "context"
    "time"
    "ministry-education/config"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = config.NewRedisClient()

func CleanupExpiredSessions() {
    ctx := context.Background()
    cutoff := time.Now().Add(-24 * time.Hour) // Remove sessions older than 24h + buffer
    
    var cursor uint64
    for {
        var keys []string
        cursor, keys = rdb.Scan(ctx, cursor, "session:*", 100).Result()
        
        if len(keys) > 0 {
            pipe := rdb.Pipeline(ctx)
            for _, key := range keys {
                ttl, err := rdb.TTL(ctx, key).Result()
                if err == nil && ttl < 0 { // Key exists but TTL expired
                    pipe.Del(ctx, key)
                }
            }
            pipe.Exec(ctx)
        }
        
        if cursor == 0 {
            break
        }
    }
    
    log.Println("Session cleanup job completed")
}

func CleanupExpiredCaches() {
    ctx := context.Background()
    
    var cursor uint64
    for {
        var keys []string
        cursor, keys = rdb.Scan(ctx, cursor, "api:response:*", 100).Result()
        
        if len(keys) > 0 {
            pipe := rdb.Pipeline(ctx)
            for _, key := range keys {
                pipe.Del(ctx, key)
            }
            pipe.Exec(ctx)
        }
        
        if cursor == 0 {
            break
        }
    }
    
    log.Println("Cache cleanup job completed")
}
```

### 3.10 Disaster Recovery

```go
// recovery/strategy.go
package recovery

import (
    "fmt"
    "ministry-education/config"
    "github.com/redis/go-redis/v9"
)

var rdb = config.NewRedisClient()

// Redis Persistence Configuration
func ConfigureRedisPersistence() {
    // redis.conf settings for production:
    /*
       save 900 1    # Save at least 1 key after 900 seconds (15 min)
       save 300 100  # Save after 300 seconds if at least 100 keys changed
       save 60 10000 # Save after 60 seconds if at least 10000 keys changed
       
       appendonly yes   # AOF (Append Only File) for durability
       appendfsync everysec  # fsync every second (balance of durability/performance)
       no-appendfsync-on-rewrite no  # Allow AOF rewrite during fsync
    */
}

func CheckRedisHealth() error {
    ctx := context.Background()
    
    // Check if Redis is responding
    pong := rdb.Ping(ctx)
    if pong.Err() != nil {
        return fmt.Errorf("Redis health check failed: %w", pong.Err())
    }
    
    // Check memory usage
    info, err := rdb.Info(ctx, "memory").Result()
    if err == nil {
        // Parse and validate memory usage
        // ...
    }
    
    return nil
}
```
---