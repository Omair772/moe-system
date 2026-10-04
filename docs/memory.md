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
    
    // Verify connection
    pong := client.Ping(ctx)
    _ = pong
    return client
}
```

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

## 4. Cache-Aside Pattern (Go Implementation)

### 4.1 Student Cache Get/Set
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

### 4.2 Course Cache Lookup
```go
// cache/courseCache.go
func GetCourse(code string) (*models.Course, error) {
    cacheKey := fmt.Sprintf("course:%s", code)
    
    val, err := rdb.Get(ctx, cacheKey).Result()
    if err == nil {
        var course models.Course
        json.Unmarshal([]byte(val), &course)
        return &course, nil
    }
    
    // Cache miss
    course, err := prisma.Course.FindUnique(
        prisma.Course.Code(code),
    ).Exec(ctx)
    if err != nil {
        return nil, err
    }
    
    // Store in cache with TTL
    go setCourseCache(course)
    
    return &course, nil
}

func setCourseCache(course *models.Course) error {
    cacheKey := fmt.Sprintf("course:%s", course.Code)
    data, _ := json.Marshal(course)
    
    return rdb.Set(ctx, cacheKey, data, config.CourseCatalogTTL).Err()
}
```

## 5. Session Management

### 5.1 JWT Token Structure
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
    JWTID        string `json:"jti"`       // Unique token ID
    IssuedAt     int64  `json:"iat"`       // Issued at timestamp
    Expiry       int64  `json:"exp"`       // Expiration timestamp
    NotBefore    int64  `json:"nbf"`       // Not before timestamp
    Issuer       string `json:"iss"`       // Issuer
    Subject      string `json:"sub"`       // Subject (user ID)
    Audience     string `json:"aud"`       // Audience
    SessionIndex int64  `json:"sidx"`      // Session index for revocation
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

### 5.2 Session Blacklist (Revocation)
```go
// auth/revocation.go
func RevokeToken(tokenString string) error {
    // Extract JTI from token
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
    
    return rdb.Set(ctx, blacklistKey, "revoked", time.Until(expiry)).Err()
}

func IsTokenRevoked(tokenString string) bool {
    ctx := context.Background()
    
    // Check blacklist
    isBlacklisted, err := rdb.Exists(ctx, fmt.Sprintf("token_blacklist:%s", extractJTI(tokenString))).Result()
    if err != nil || isBlacklisted == 0 {
        return false // Not blacklisted or error checking
    }
    
    return true
}
```

## 6. Memory Optimization Strategies

### 6.1 Go Memory Management
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

### 6.2 Connection Pool Configuration
```go
// db/postgres.go
package db

import (
    "database/sql"
    "fmt"
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

## 7. Cache Invalidation Strategies

### 7.1 Write-Through Cache
```
Scenario: Student profile update

1. API request receives PUT /students/:id
2. Validate input with Zod
3. Update student in PostgreSQL via Prisma
4. If DB update successful:
   a. DELETE cache:student:{id}
   b. Return updated student to client
5. Next READ request → Cache miss → DB query → Fresh cache populated

Benefits: Cache always consistent with DB
Trade-off: Slightly slower writes (extra DELETE operation)
```

### 7.2 Write-Behind Cache
```
Scenario: Attendance marking

1. Teacher marks attendance via React UI
2. Immediately write to Redis cache (write-behind)
3. ACK response to client instantly
4. Background worker batch-writes to PostgreSQL every 30 seconds
5. If worker fails, retry with exponential backoff

Benefits: Fast write response
Trade-off: Potential data loss if process crashes before background write
Recommended for: Non-critical data, frequent small updates
```

### 7.3 Selective Invalidation
```go
// cache/invalidation.go
func InvalidateStudentCache(studentID string) error {
    ctx := context.Background()
    keysToInvalidate := []string{
        fmt.Sprintf("student:%s", studentID),
        fmt.Sprintf("gpa:%s", studentID),
        fmt.Sprintf("report:%s", studentID), // if exists
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
    
    // Invalidate course cache and related
    keys := []string{
        fmt.Sprintf("course:%s", courseCode),
        // Could also invalidate enrollment caches
    }
    
    pipe := rdb.Pipeline(ctx)
    for _, key := range keys {
        pipe.Del(ctx, key)
    }
    
    _, err := pipe.Exec(ctx)
    return err
}
```

## 8. Memory Monitoring & Alerting

### 8.1 Key Metrics
| Metric | Normal Range | Alert Threshold | Monitoring Tool |
|--------|-------------|-----------------|-----------------|
| Redis Used Memory | < 70% of limit | > 85% | Redis MONITOR + CloudWatch |
| Go Heap Alloc | < 200 MB | > 500 MB | Prometheus + Grafana |
| Active Goroutines | < 200 | > 500 | Custom exporter |
| Cache Hit Ratio | > 70% | < 50% | Redis INFO command |
| DB Connection Pool Usage | < 80% | > 95% | pg_stat_activity |
| Request Latency p95 | < 200ms | > 500ms | OpenTelemetry |

### 8.2 Alert Rules
```
# Redis memory alert
alert "redis_memory_high" {
    expr = redis_memory_percent > 85
    for = 2m
    labels = { severity = "warning" }
    annotations = {
        summary = "Redis memory usage high"
        description = "Redis is using {{ $value }}% of available memory"
    }
}

# Go memory leak alert
alert "go_memory_leak" {
    expr = process_resident_memory_bytes > 1000000000
    for = 5m
    labels = { severity = "critical" }
    annotations = {
        summary = "Potential Go memory leak"
        description = "Process memory exceeds 1GB: {{ $value }} bytes"
    }
}
```

## 9. Data Retention & Cleanup

### 9.1 Automated Cleanup Jobs
```go
// jobs/cleanup.go
package jobs

import (
    "context"
    "time"
    "ministry-education/config"
    "github.com/redis/go-redis/v9"
)

func CleanupExpiredSessions() {
    ctx := context.Background()
    cutoff := time.Now().Add(-24 * time.Hour) // Remove sessions older than 24h + buffer
    
    // Use SCAN to find keys matching pattern
    var cursor uint64
    for {
        var keys []string
        cursor, keys = rdb.Scan(ctx, cursor, "session:*", 100).Result()
        
        if len(keys) > 0 {
            // Check TTL and delete expired
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
    
    // Clean up stale API response caches
    var cursor uint64
    for {
        var keys []string
        cursor, keys = rdb.Scan(ctx, cursor, "api:response:*", 100).Result()
        
        if len(keys) > 0 {
            // Delete all API response caches (they have short TTLs anyway)
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

### 9.2 Retention Policies
| Data Type | Retention Period | Storage | Cleanup Method |
|-----------|-----------------|---------|----------------|
| Session data | 24-48 hours | Redis | TTL automatic expiration |
| API response cache | 1-5 minutes | Redis | TTL automatic expiration |
| Audit logs | 10 years | PostgreSQL + S3 | Scheduled pg_dump + S3 transfer |
| Error logs | 1 year | Elasticsearch | Index lifecycle management |
| Session blacklist | Until token expiry | Redis | TTL = token expiry time |
| GPA calculations | 1 hour | Redis | TTL automatic expiration |
| Report exports | 30 days | PostgreSQL | Weekly cleanup job |

## 10. Disaster Recovery

### 10.1 Backup Strategy
| Component | Backup Frequency | Retention | Recovery Point Objective (RPO) | Recovery Time Objective (RTO) |
|-----------|-----------------|-----------|-------------------------------|-------------------------------|
| PostgreSQL data | Daily full + hourly incremental | 10 years | < 1 hour | < 4 hours |
| Redis data | Snapshotting every 5 min | 7 days | < 5 minutes | < 1 hour |
| Static assets (React build) | On deploy | Indefinite | N/A | N/A |
| Configuration files | On change | 90 days | N/A | N/A |

### 10.2 Redis Persistence
```go
// Ensure Redis persistence is enabled
// redis.conf settings (should be set in production)
/
   save 900 1    # Save at least 1 key after 900 seconds (15 min)
   save 300 100  # Save after 300 seconds if at least 100 keys changed
   save 60 10000 # Save after 60 seconds if at least 10000 keys changed
/

appendonly yes   # AOF (Append Only File) for durability
appendfsync everysec  # fsync every second (balance of durability/performance)
no-appendfsync-on-rewrite no  # Allow AOF rewrite during fsync
```
---