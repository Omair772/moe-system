# 📋 Implementation Plan — Ministry of Education System

## 1. Project Phases & Timeline

### Phase 1: Foundation & Setup (Weeks 1-2)
| Week | Deliverables | Technology |
|------|-------------|------------|
| 1 | • Repository initialization<br>• Go module setup<br>• Docker configuration<br>• CI/CD pipeline (GitHub Actions) | Go 1.22, Docker, GitHub Actions |
| 2 | • PostgreSQL database setup<br>• Prisma ORM initialization<br>• schema.prisma design<br>• Redis cache configuration | PostgreSQL 15, Prisma v5, Redis |

### Phase 2: Core Domain & API (Weeks 3-5)
| Week | Deliverables | Technology |
|------|-------------|------------|
| 3 | • Define Domain Entities (Student, Course, Enrollment)<br>• Repository Interfaces<br>• Prisma Client Generation<br>• Basic CRUD APIs | Go interfaces, Prisma, REST API |
| 4 | • Authentication middleware (JWT)<br>• Authorization (RBAC)<br>• Input validation (Zod)<br>• Error handling standards | Go, JWT, Zod |
| 5 | • Student management APIs<br>• Course management APIs<br>• Enrollment endpoints | Go Gin, Prisma, PostgreSQL |

### Phase 3: Feature Implementation (Weeks 6-10)
| Week | Deliverables | Technology |
|------|-------------|------------|
| 6 | • Attendance tracking APIs<br>• Grade entry endpoints<br>• Teacher assignment APIs | Go, Prisma |
| 7 | • Parent portal APIs<br>• Report generation endpoints<br>• Search & pagination | Go, SQL queries |
| 8 | • Redis cache integration (TTL)<br>• Cache invalidation strategies<br>• Cache-first middleware | Go, Redis |
| 9 | • API rate limiting<br>• Request logging (audit)<br>• Health check endpoints | Go, Middleware |
| 10 | • Unit test suite (>80% coverage)<br>• Integration tests<br>• API documentation (Swagger) | Go test, Swagger |

### Phase 4: Frontend & Integration (Weeks 11-14)
| Week | Deliverables | Technology |
|------|-------------|------------|
| 11 | • React project setup (Vite + TypeScript)<br>• Component library (common UI)<br>• API client configuration | React 18, TypeScript, Axios |
| 12 | • Student dashboard UI<br>• Course enrollment flow<br>• Attendance view component | React, API client |
| 13 | • Teacher interface (grade entry, attendance)<br>• Parent portal (view grades, attendance)<br>• Notification system | React, WebSockets/SSE |
| 14 | • End-to-end testing<br>• Cross-browser testing<br>• Performance optimization | Cypress, Lighthouse |

### Phase 5: Production Release (Weeks 15-16)
| Week | Deliverables | Technology |
|------|-------------|------------|
| 15 | • Staging deployment<br>• Performance testing (load.io)<br>• Security penetration testing | Kubernetes, Docker |
| 16 | • Production deployment<br>• Monitoring setup (Prometheus+Grafana)<br>• User acceptance testing | AWS/Azure, Observability |

## 2. Technology Stack Implementation Details

### 2.1 go.mod (Updated)

```go
module api.ministry.education

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/prisma/prisma-client-go/v5 v5.14.0
    github.com/redis/go-redis/v9 v9.5.0
    github.com/go-playground/validator/v10 v10.17.0
    github.com/golang-jwt/jwt/v5 v5.2.1
    github.com/overseer/v2 v2.1.0 // rate limiting
)
```

**JWT Library Update:**
- **Old:** `github.com/dgrijalva/jwt-go v5.0.0` (deprecated, unmaintained)
- **New:** `github.com/golang-jwt/jwt/v5 v5.2.1` (modern, actively maintained)

### 2.2 Prisma Schema (schema.prisma) (Updated)

```prisma
generator client {
  provider = "postgresql"
  output   = "./prisma/client"  // Fixed: was output without path
}

datasource db {
  provider = "postgresql"
  url      = env("DATABASE_URL")
}

model Student {
  id        String   @id @default(uuid())
  name      String
  email     String   @unique
  password  String
  createdAt DateTime @default(now())  // Fixed: was now__()
  updatedAt DateTime @updatedAt
  enrollments Enrollment[]
}

model Course {
  id        String   @id @default(uuid())
  title     String
  code      String   @unique
  credits   Int
  teacher   Teacher? @relation(fields: [teacherId], references: [id])
  teacherId String?
  students  Enrollment[]
  createdAt DateTime @default(now())  // Fixed: was now__()
  updatedAt DateTime @updatedAt
}

model Enrollment {
  id         String  @id @default(uuid())
  student    Student @relation(fields: [studentId], references: [id])
  studentId  String
  course     Course  @relation(fields: [courseId], references: [id])
  courseId   String
  grade      String?
  attended   Int     // percentage (0-100)
  enrolledAt DateTime @default(now())  // Fixed: was now__()
  @@unique([studentId, courseId])  // Fixed: was UNIQUE(studentId, courseId)
}

model Teacher {
  id        String   @id @default(uuid())
  name      String
  email     String   @unique
  password  String
  courses   Course[]
  createdAt DateTime @default(now())  // Fixed: was now__()
  updatedAt DateTime @updatedAt
}
```

**Fixed Prisma Commands:**
```bash
npx prisma migrate dev --name init
npx prisma generate
```

### 2.3 REST API Design (Gin Web Framework)

### 2.4 Database Schema & Prisma Client

### 2.5 Redis Cache with TTL Strategy (Updated Go Code)

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

func NewRedisClient() *redis.Client {
    ctx := context.Background()  // Fixed: ensured context initialization
    client := redis.NewClient(&redis.Options{
        Addr:        "localhost:6379",
        Password:    "", // no password set
        DB:          0,
        PoolSize:    10,
        MinIdleConns: 5,
    })
    
    // Verify connection
    pong := client.Ping(ctx)
    _ = pong
    return client
}
```

**Context Initialization in Redis Functions:**
- All Redis handler functions must include `ctx := context.Background()` at the start
- Example pattern:
```go
func GetStudentCache(id string) (*models.Student, error) {
    ctx := context.Background()  // Fixed: context initialization
    cacheKey := fmt.Sprintf("student:%s", id)
    // ... rest of function
}
```

### 2.6 API Route Structure (Gin)

```go
// routes/studentRoutes.go
func StudentRoutes(router *gin.Engine, studentHandler *StudentHandler) {
    students := router.Group("/api/v1/students")
    {
        students.GET("", studentHandler.List)           // List with pagination
        students.GET(":id", studentHandler.Get)         // Get by ID
        students.POST("", studentHandler.Create)        // Create student
        students.PUT(":id", studentHandler.Update)      // Update student
        students.DELETE(":id", studentHandler.SoftDelete) // Soft-delete
    }
}
```

### 2.7 Authentication & Authorization

| Role | Permissions |
|------|-------------|
| SUPER_ADMIN | Full access to all endpoints |
| ADMIN | Manage students, teachers, courses, reports |
| TEACHER | View classes, take attendance, enter grades |
| STUDENT | View own data, enroll in courses, view grades |
| PARENT | View child's attendance and grades |

### 2.8 Observability Stack

- **Structured Logging:** JSON logs with correlation IDs
- **Metrics:** p95 latency < 200ms, error rates, cache hit ratios
- **Health Checks:** `/health/live` and `/health/ready`
- **Tracing:** Correlation IDs for request tracking

### 2.9 Security Considerations

- **TLS 1.3** for all external communication
- **bcrypt** with cost factor 12 for password hashing
- **Parameterized queries** via Prisma (prevents SQL injection)
- **Input validation** with Zod on both frontend and backend
- **Rate limiting** per IP and per user role

### 2.10 Deployment Architecture

```
┌─────────────────────────────────────────────────────┐
│                    PRODUCTION                        │
│  ┌─────────────┐  ┌─────────────────────────┐   │
│  │  Nginx      │  │  Go REST API (8080)      │   │
│  │  Reverse    │  │  - Gin Framework         │   │
│  │  Proxy      │  │  - JWT Auth              │   │
│  └─────────────┘  └─────────────────────────┘   │
│                                      │          │
│  ┌─────────────┐  ┌─────────────────────────┐   │
│  │  Redis      │  │  PostgreSQL + Prisma ORM │   │
│  │  Cache/TTL  │  │  - Schema Migrations     │   │
│  └─────────────┘  └─────────────────────────┘   │
│                                      │          │
│  ┌─────────────┐  ┌─────────────────────────┐   │
│  │  React CDN  │  │  Backups & Monitoring    │   │
│  │  SPA        │  └─────────────────────────┘   │
│  └─────────────┘                             │
└─────────────────────────────────────────────────────┘
```

## 3. Development Workflow

### 3.1 Git Branching Strategy
```
main          ← Production-ready (protected branch)
develop       ← Integration branch (protected)
feature/*     ← New feature branches (from develop)
hotfix/*      ← Production fix branches (from main)
release/*     ← Release candidates (from develop)
```

### 3.2 Branch Naming Conventions
```
feature/student-api-v1
feature/course-enrollment
fix/jwt-auth-timeout
release/v1.0.0-beta
```

### 3.3 Commit Message Convention
```
<type>(<scope>): <subject>

types:
- feat:     New feature (e.g., feat(student): add create endpoint)
- fix:      Bug fix (e.g., fix(attendance): correct grade calc)
- docs:     Documentation changes
- style:    Code formatting (gofmt, prettier)
- refactor: Code refactoring
- test:     Adding/updating tests
- chore:   Build/tooling changes

examples:
feat(student): add create endpoint with validation
fix(course): resolve enrollment duplicate issue
docs(readme): update API authentication guide
refactor(student): improve repository pattern
test(enrollment): add integration tests
```

### 3.4 Pull Request Process
1. Create feature branch from `develop`
2. Implement feature with tests
3. Push branch and create PR to `develop`
4. PR template must include:
   - Description of changes
   - Related issue numbers
   - Testing approach
   - Screenshots (if UI changes)
5. Minimum 2 reviewer approvals
6. All CI checks must pass:
   - Go unit tests
   - Linting (golangci-lint)
   - SQL migration validation
   - Docker build test
7. Squash merge to `develop`

### 3.5 Daily Standup Format
```
Yesterday:
- What I completed (feature/xyz, bug/abc)

Today:
- What I will work on
- Any dependencies needed

Blockers:
- Issue 1: description
- Issue 2: description
```

## 4. Milestones & Deliverables

| Milestone | Target Week | Acceptance Criteria |
|-----------|-------------|---------------------|
| M1: Project Charter | Week 1 | ✅ Repo initialized<br>✅ Docker working<br>✅ Team onboarded |
| M2: Core API Ready | Week 5 | ✅ All CRUD endpoints documented<br>✅ Auth working<br>✅ DB schema migrated |
| M3: Feature Complete | Week 10 | ✅ Student/Course/Enrollment APIs<br>✅ Attendance tracking<br>✅ 80% test coverage |
| M4: Beta Release | Week 14 | ✅ React UI functional<br>✅ User flows complete<br>✅ Bugs < 5 critical |
| M5: Production Ready | Week 16 | ✅ Staging deployed<br>✅ Monitoring active<br>✅ UAT signed off |

## 5. Risk Matrix

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Scope creep | High | Medium | Strict change control, backlog grooming every sprint |
| Database migration breaking | Medium | High | Test migrations on staging first; backup before prod |
| Performance bottlenecks | Medium | Medium | Load testing at week 8; optimize critical queries |
| Security vulnerabilities | Low | High | Regular security reviews; dependency updates; penetration test |
| Team availability | Medium | Medium | Buffer capacity 20%; knowledge sharing sessions |
| Cache inconsistency | Low | Medium | TTL-based invalidation; write-through pattern |

## 6. Success Criteria (Definition of Done)

- ✅ 100% of API endpoints have automated tests
- ✅ API documentation auto-generated (Swagger/OpenAPI)
- ✅ Performance: p95 latency < 200ms for standard endpoints
- ✅ 80%+ test coverage on critical paths (enrollment, grading, attendance)
- ✅ Zero CRITICAL security vulnerabilities (OWASP Top 10)
- ✅ Cache hit ratio > 70% after warm-up period
- ✅ Database query count < 5 per average API call
- ✅ GDPR compliance features implemented (data export/erasure)
- ✅ Audit logging complete for all data-modifying operations
- ✅ Production deployment successful with monitoring