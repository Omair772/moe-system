# �-flow_of_action — Ministry of Education System

## 1. Purpose & Scope

This document defines the **Flow of Action** for the Ministry of Education integrated system. While use case scenarios describe the business logic and use case actions describe individual operations, the Flow of Action documents the complete request lifecycle from entry to exit, including all middleware, decision points, and state transitions.

## 2. Request Lifecycle Overview

```
┌──────────────────────────────────────────────────────────────┐
│                    CLIENT REQUEST                            │
│  (React Frontend → Axios → HTTP)                             │
└────────────▲─────────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────────────────────┐
│                    GIN MIDDLEWARE CHAIN                      │
│  1. Correlation ID Generator                                 │
│  2. Authentication (JWT Validation)                          │
│  3. Authorization (RBAC Check)                               │
│  4. Request ID & Timing                                      │
│  5. Input Validation (Zod)                                   │
│  6. Cache Lookup (Redis TTL)                                 │
│  7. Business Use Case Execution                              │
│  8. Response Transformation                                  │
│  9. Cache Store (if applicable)                              │
│  10. Audit Logging                                           │
│  11. HTTP Response                                           │
└────────────▲─────────────────────────────────────────────────┘
              │
              ▼
┌──────────────────────────────────────────────────────────────┐
│                    SERVER RESPONSE                           │
│  (JSON → React State Update → UI)                            │
└──────────────────────────────────────────────────────────────┘
```

## 3. Detailed Middleware Flow

### 3.1 Request Processing Pipeline

```
HTTP Request Entering Gin Framework
    │
    ▼
┌─────────────────────────────────────┐
│ 1. Correlation ID Middleware        │
│    • Generate UUID v4              │
│    • Set X-Request-ID header       │
│    • Store in context for later use│
│    • Example: req-abc123def456     │
└────────────▲────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ 2. Authentication Middleware        │
│    • Extract Authorization header  │
│    • Validate JWT token            │
│    • Check token expiry (2hr)      │
│    • Verify signature with JWT key  │
│    • Set user context:             │
│      - c.Set("userID", "...")      │
│      - c.Set("userRole", "...")    │
│    • Success → c.Next()            │
│    • Fail → c.AbortWithStatus(401) │
└────────────▲────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ 3. Authorization Middleware         │
│    • Get user role from context     │
│    • Check role permissions         │
│    • Validate action access         │
│    • RBAC: SUPER_ADMIN > ADMIN > TEACHER > STUDENT > PARENT │
│    • Success → c.Next()            │
│    • Fail → c.AbortWithStatus(403) │
└────────────▲────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ 4. Request Timing Middleware        │
│    • Start timer: start := time.Now()│
│    • Continue chain                │
│    • After c.Next():               │
│    • Calculate duration            │
│    • Log latency metrics           │
│    • Set X-Response-Time header    │
└────────────▲────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ 5. Cache Lookup Middleware          │
│    • Construct cache key:           │
│      cache:api:{method}{path}        │
│    • Try Redis GET                   │
│    • Cache HIT → Return cached data │
│    • Cache MISS → c.Next()           │
│    • After response:                 │
│      • Store in cache with TTL       │
│      • Example: cache:api:GET/students  │
│        TTL: 1 minute                 │
└────────────▲────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ 6. Input Validation Middleware      │
│    • Extract request body           │
│    • Validate with Zod schema       │
│    • Validation SUCCESS → c.Next()  │
│    • Validation FAIL → c.Abort(400) │
│    • Error format: Zod error issues │
└────────────▲────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ 7. Business Use Case Execution      │
│    • Extract action from context     │
│    • Route to appropriate handler    │
│    • Execute: Prisma → PostgreSQL    │
│    • Apply business rules            │
│    • Return result                   │
│    • May trigger: notifications,     │
│      cache invalidation, audit       │
└────────────▲────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ 8. Response Transformation          │
│    • Format JSON API response        │
│    • success: true/false             │
│    • data: { ... }                     │
│    • meta: { ... }                     │
│    • timestamp: ISO 8601              │
│    • correlation_id: from step 1     │
│    • error handling if applicable    │
└────────────▲────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ 9. Cache Store Middleware (if HIT)  │
│    • If cache HIT earlier:           │
│      • Set content-type              │
│      • Write cached data to response │
│      • c.Abort()                     │
│      • Skip remaining middleware     │
│    • If cache MISS:                  │
│      • Continue to step 7            │
│      • After success:                │
│        • Store result in Redis         │
│        • Set TTL (see TTL spec)        │
└────────────▲────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│ 10. Audit Logging Middleware        │
│     • Write audit entry to           │
│       PostgreSQL via Prisma          │
│     • Fields:                        │
│       - correlation_id               │
│       - user_id (from auth)          │
│       - user_role                    │
│       - action (e.g., "ENROLL_CREATE")│
│       - entity_type (e.g., "ENROLLMENT")│
│       - entity_id                    │
│       - changes (before/after)       │
│       - ip_address (c.ClientIP())    │
│       - user_agent (c.Request.User-Agent)│
│       - success (boolean)            │
│       - duration (time.Since(start)) │
│     • async (non-blocking)           │
└────────────▲────────────────────────┘
             │
             ▼
│ 11. HTTP Response → Client           │
└────────────▲────────────────────────┘
              │
              ▼
```

### 3.2 Middleware Execution Order

| # | Middleware | Purpose | Abort On Failure |
|---|-----------|---------|-----------------|
| 1 | Correlation ID | Generate request tracking ID | No |
| 2 | Authentication | Validate JWT token | Yes (401) |
| 3 | Authorization | Check RBAC permissions | Yes (403) |
| 4 | Request Timing | Measure latency | No |
| 5 | Cache Lookup | Check Redis cache | Yes (if HIT, skip rest) |
| 6 | Input Validation | Zod schema validation | Yes (400) |
| 7 | Business Logic | Execute use case | No (handle errors internally) |
| 8 | Response Transform | Format API response | No |
| 9 | Cache Store | Store in Redis if MISS | No |
| 10 | Audit Log | Write audit entry | No |
| 11 | HTTP Response | Send to client | N/A |

## 4. State Transition Diagram

```
┌─────────────────────┐
│   [UNPROCESSED]     │
│   (HTTP request)    │
└────────▲────────────┘
         │
         ▼
┌─────────────────────┐
│   [CORRELATION]     │
│   (Generate ID)     │
└────────▲────────────┘
         │
         ▼
┌─────────────────────┐
│   [AUTHENTICATED]   │
│   (JWT validated)   │
└────────▲────────────┘
         │
         ▼
┌─────────────────────┐
│   [AUTHORIZED]      │
│   (RBAC check)      │
└────────▲────────────┘
         │
         ▼
┌─────────────────────┐
│   [VALIDATED]       │
│   (Zod schema)      │
└────────▲────────────┘
         │
         ▼
┌─────────────────────┐
│   [CACHE_HIT]       │
│   (or MISS → process)│
└────────▲────────────┘
         │
         ▼
┌─────────────────────┐
│   [EXECUTE]         │
│   (Use case logic)  │
└────────▲────────────┘
         │
         ▼
┌─────────────────────┐
│   [RESPONDED]       │
│   (HTTP response)   │
└────────▲────────────┘
         │
         ▼
┌─────────────────────┐
│   [COMPLETED]       │
│   (Request finished)│
└─────────────────────┘
```

### 4.1 Alternative Paths in State Transition

```
[VALIDATED] ──[Cache HIT]──► [RESPONDED]──► [COMPLETED]
     │
     └─[Cache MISS]──► [EXECUTE]──► [TRANSFORM]──► [STORE_CACHE]──► [RESPONDED]──► [COMPLETED]
     │
     └─[Validation FAIL]──► [RESPONDED_ERROR]──► [COMPLETED]
     │
     └─[Auth FAIL]────────► [RESPONDED_ERROR]──► [COMPLETED]
     │
     └─[Authz FAIL]────────► [RESPONDED_ERROR]──► [COMPLETED]
```

## 5. Cache Integration in Flow

### 5.1 Cache-First Strategy

```
┌─────────────────┐     ┌─────────────────────┐
│  Client Request │     │  Gin Middleware      │
└───────▲─────────┘     └───────▲───────────────┘
        │                     │
        │  1. Construct Key  │
        │    cache:api:GET   │
        │    /students       │
        │                     │
        ├────► [Redis GET]  │
        │         │          │
        │         ├─ HIT ────┤
        │         │          │  Return cached JSON
        │         │          │  c.Abort()
        │         │          │  Skip steps 7-10
        │         │          │
        │         └─ MISS ─┤
        │                     │
        │                     │  Continue to step 7
        │                     │
        └────► [Proceed]      │
                       │       │
                       ▼       │
                ┌─────────────┐
                │  Use Case   │
                │  Execution  │
                └────▲─────────┘
                     │
                     │  On success:
                     │  1. Serialize result
                     │  2. Store in Redis
                     │  3. Set TTL (see spec)
                     │  4. Continue to respond
                     │
                     ▼
                ┌─────────────┐
                │ Response    │
                └─────────────┘
```

### 5.2 Cache Store Pattern (After Use Case Execution)

```go
// After successful use case execution
func storeInCache(key string, value interface{}, ttl time.Duration) {
    ctx := context.Background()
    data, _ := json.Marshal(value)
    err := rdb.Set(ctx, key, data, ttl).Err()
    if err != nil {
        logger.Error("Failed to store in cache", "error", err)
        // Not critical - continue without caching
    }
}

// Example usage after student lookup
student, err := getStudentUseCase(id)
if err == nil {
    storeInCache(
        fmt.Sprintf("student:%s", id),
        student,
        config.StudentProfileTTL, // 30 minutes
    )
}
```
```

## 6. Error Flow Diagrams

### 6.1 Authentication Failure

```
HTTP Request
    │
    ▼
[Auth Middleware]
    │
    ├─ JWT Missing
    │   ▼
    │  c.AbortWithStatusJSON(401, {"error": "Missing token"})
    │
    └─ JWT Invalid/Expired
        ▼
        c.AbortWithStatusJSON(401, {"error": "Invalid or expired token"})
```

### 6.2 Authorization Failure

```
HTTP Request
    │
    ▼
[Auth Middleware] → [Auth OK]
    │
    ▼
[Authz Middleware]
    │
    ├─ Student tries admin action
    │   ▼
    │  c.AbortWithStatusJSON(403, {"error": "Insufficient permissions"})
    │
    └─ Teacher tries unassigned course action
        ▼
        c.AbortWithStatusJSON(403, {"error": "You don't have permission for this course"})
```

### 6.3 Validation Failure

```
HTTP Request
    │
    ▼
[Auth Middleware] → [Auth OK]
    │
    ▼
[Authz Middleware] → [Authz OK]
    │
    ▼
[Validation Middleware]
    │
    ├─ Zod error: email invalid
    │   ▼
    │   c.AbortWithStatusJSON(400, {
    │       "success": false,
    │       "error": {
              "code": "VALIDATION_ERROR",
              "message": "Invalid email format",
           "details": [{"field": "email", "issue": "Must be valid email"}]
       }
   })
   │
   └─ Zod error: missing required field
       ▼
       c.AbortWithStatusJSON(400, similar format)
```

### 6.4 Business Rule Failure

```
HTTP Request
    │
    ▼
[All Middleware OK]
    │
    ▼
[Use Case Execution]
    │
    ├─ Prisma constraint violation
    │   ▼
    │   c.AbortWithStatusJSON(409, {
    │       "success": false,
    │       "error": {
              "code": "CONFLICT",
           "message": "Cannot drop course: grades already entered"
       }
   })
   │
   └─ Capacity full
       ▼
       c.AbortWithStatusJSON(409, {
           "success": false,
           "error": {
               "code": "CAPACITY_FULL",
               "message": "Course is at maximum capacity"
           }
       })
```

## 7. Response Format Standards

### 7.1 Success Response

```json
{
  "success": true,
  "data": {
    // Action-specific data
    "id": "uuid",
    "name": "Student Name",
    "status": "active"
  },
  "meta": {
    "cache_invalidated": ["student:123"],
    "new_ttls": {
      "student:123": "30min"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "correlation_id": "req-abc123def456"
}
```

### 7.2 Error Response

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR|AUTHENTICATION_ERROR|AUTHORIZATION_ERROR|NOT_FOUND|CONFLICT|VALIDATION_ERROR|INTERNAL_ERROR",
    "message": "Human-readable error message",
    "details": [ // optional - field-specific errors
      {"field": "email", "issue": "Must be valid email"}
    ]
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "correlation_id": "req-abc123def456"
}
```

## 8. Performance Standards

| Metric | Target | Measurement Location |
|--------|--------|---------------------|
| Total request latency p95 | < 200ms | Gin middleware timer |
| Authentication latency | < 5ms | JWT validation library |
| Cache check latency | < 2ms | Redis GET command |
| Use case execution | < 100ms | Go function timing |
| Response serialization | < 10ms | JSON.Marshal |
| Cache store latency | < 5ms | Redis SET command |

## 9. Logging Correlation

All logs include correlation ID for request tracing:

```
Log Format (JSON):
{
  "correlation_id": "req-abc123def456",
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "message": "HTTP request processed",
  "fields": {
    "method": "GET",
    "path": "/api/v1/students",
    "user_id": "11111111-1111-1111-1111-111111111111",
    "user_role": "STUDENT",
    "status": 200,
    "latency_ms": 145,
    "ip": "192.168.1.100",
    "user_agent": "Mozilla/5.0..."
  }
}
```

---