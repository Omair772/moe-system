# 🏗️ Architecture — Ministry of Education System

## 1. Clean Architecture Overview

This system follows **Clean Architecture** (Hexagonal Architecture) with dependencies pointing inward. The core business logic has no dependencies on external frameworks, databases, or UI technologies.

```
┌──────────────────────────────────────────────────────────────┐
│                    DOMAIN LAYER (Core)                       │
│  Entity Interfaces → Repository Ports → Use Cases            │
│  No framework dependencies — pure Go interfaces              │
└──────────────────────────────────────────────────────────────┘

                    ↑ Domain Layer (INNERMOST)

┌──────────────────────────────────────────────────────────────┐
│                APPLICATION LAYER (Use Cases)                 │
│  Orchestrates business rules, DTO transformation, validation │
│  Depends on Domain layer only                                │
└──────────────────────────────────────────────────────────────┘

                    ↑ Application Layer

┌──────────────────────────────────────────────────────────────┐
│              ADAPTER LAYER (Ports & Adapters)                │
│  External interfaces: REST controllers, DB adapters, Cache     │
│  Depends on Application layer                                │
└──────────────────────────────────────────────────────────────┘

                    ↑ Adapter Layer

┌──────────────────────────────────────────────────────────────┐
│                 EXTERNAL INTERFACES (Driving)                │
│  REST API Controllers, External Services, React Frontend       │
│  Depends on Adapter layer                                    │
└──────────────────────────────────────────────────────────────┘
```

## 2. Technology Mapping by Layer

| Layer | Technology | Responsibility |
|-------|-----------|----------------|
| **Domain** | Go interfaces | Student, Course, Enrollment entities; Repository ports |
| **Application** | Go services | EnrollStudent, UpdateGrade, GetAttendance use cases |
| **Adapter** | Go HTTP handlers | Gin/Chi routes; React component bindings |
| **External** | PostgreSQL + Prisma | Relational data; Redis with TTL cache |

## 3. REST API Design (Gin Web Framework)

### 3.1 Base Configuration
```go
// main.go
router := gin.New()
router.Use(gin.Recovery())
router.Use(correlationIDMiddleware())
router.Use(auditMiddleware())
router.Use(limitingMiddleware(100)) // 100 requests/minute
```

### 3.2 Core Endpoints

| Method | Endpoint | Handler | Auth |
|--------|----------|---------|------|
| GET | `/api/v1/health` | `health.Check` | Public |
| GET | `/api/v1/students` | `student.List` | Bearer Token |
| GET | `/api/v1/students/:id` | `student.Get` | Bearer Token |
| POST | `/api/v1/students` | `student.Create` | Bearer Token |
| PUT | `/api/v1/students/:id` | `student.Update` | Bearer Token |
| DELETE | `/api/v1/students/:id` | `student.SoftDelete` | Bearer Token |
| GET | `/api/v1/courses` | `course.List` | Bearer Token |
| POST | `/api/v1/enrollments` | `enrollment.Create` | Bearer Token |
| GET | `/api/v1/attendance/:date` | `attendance.List` | Bearer Token |

### 3.3 Request/Response Middleware Chain

```
Client Request
    │
    ▼
[Correlation ID Generator] → HTTP Header `X-Request-ID`
    │
    ▼
[Authentication Middleware] → JWT Validation (2hr expiry)
    │
    ▼
[Authorization Middleware] → Role Check (RBAC)
    │
    ▼
[Audit Logging] → Write audit entry
    │
    ▼
[Request Validation] → Zod schema validation
    │
    ▼
[Business Use Case] → Execute application logic
    │
    ▼
[Response Transformation] → JSON API format
    │
    ▼
[HTTP Response] → Client
```

## 4. Prisma ORM Configuration

### 4.1 schema.prisma
```prisma
generator client {
  provider = "postgresql"
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
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
  enrollments Enrollment[]
}

model Course {
  id        String   @id @default(uuid())
  title     String
  code      String   @unique
  credits   Int
  students  Enrollment[]
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
}

model Enrollment {
  id         String  @id @default(uuid())
  student    Student @relation(fields: [studentId], references: [id])
  studentId  String
  course     Course  @relation(fields: [courseId], references: [id])
  courseId   String
  grade      String?
  attended   Int     // percentage
  enrolledAt DateTime @default(now())
  UNIQUE(studentId, courseId)
}

model Teacher {
  id        String   @id @default(uuid())
  name      String
  email     String   @unique
  password  String
  courses   Course[]
  createdAt DateTime @default(now__)
  updatedAt DateTime @updatedAt
}
```

### 4.2 Prisma Client Usage

```go
// Creating a student
func CreateStudent(name, email, password string) (*models.Student, error) {
    hashedPassword, _ := bcrypt.HashPassword(password, 12)
    return await prisma.student.Create{
        Name:     name,
        Email:    email,
        Password: hashedPassword,
    }.Save()
}

// Getting student with enrollments
func GetStudentWithEnrollments(id string) (*models.Student, error) {
    return await prisma.Student.FindUnique(
        where: student.ID(id),
    ).Enrollments().Exec()
}

// Updating enrollment grade
func UpdateEnrollmentGrade(enrollmentID, grade string) error {
    return await prisma.enrollment.Update(
        where: enrollment.ID(enrollmentID),
        data:  enrollment.Grade.Set(grade),
    ).Exec()
}
```

## 5. Redis Cache with TTL Strategy

### 5.1 Cache Configuration
```go
// config/redis.go
package config

import ("github.com/redis/go-redis/v9")

func NewRedisClient() *redis.Client {
    ctx := context.Background()
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,
    })
    return client
}
```

### 5.2 TTL Keys Structure

| Key | TTL | Purpose |
|-----|-----|---------|
| `student:{id}` | 30 min | Student profile cache |
| `course:{code}` | 2 hr | Course catalog (read-heavy) |
| `attendance:{date}` | 5 min | Daily attendance (fresh data) |
| `session:{token}` | 24 hr | Auth session validation |
| `api:response:{path}` | 1 min | API response caching |

### 5.3 Cache Integration Pattern

```go
// Middleware for cache-first strategy
func CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        cacheKey := fmt.Sprintf("api:%s%s", c.Request.Method, c.FullPath())
        
        // Try cache first
        cached, err := rdb.Get(ctx, cacheKey).Result()
        if err == nil {
            c.Data(http.StatusOK, "application/json", []byte(cached))
            c.Abort()
            return
        }
        
        // Cache miss - proceed to handler
        c.Next()
        
        // Store in cache after response
        go func() {
            rdb.Set(ctx, cacheKey, c.Writer.String(), ttl)
        }()
    }
}
```

## 6. React Frontend Architecture

### 6.1 Project Structure
```
src/
├─ components/        │
│  ├─ common/        │
│  │  ├─ Button.tsx
│  │  ├─ Input.tsx
│  │  └─ Table.tsx
│  ├─ students/      │
│  │  ├─ List.tsx
│  │  ├─ Form.tsx
│  │  └─ Table.tsx
│  ├─ courses/       │
│  │  └─ similar structure
│  └─ shared/        │
│     ├─ auth/       │
│     └─ api/        │
├─ hooks/             │
│  useAuth.ts
│  useStudents.ts
│  useToast.ts
├─ pages/             │
│  ├─ Login.tsx
│  ├─ Dashboard.tsx
│  ├─ StudentsPage.tsx
│  └─ CoursesPage.tsx
├─ stores/            │
│  useStudentStore.ts
│  useCourseStore.ts
├─ styles/            │
│  └─ tailwind.config.js
└─ types/             │
   └─ api.ts
```

### 6.2 API Client (Axios with Interceptors)
```typescript
// src/api/client.ts
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - handle auth errors
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      router.push('/login');
    }
    return Promise.reject(error);
  }
);

export default api;
```

### 6.3 Component Example - Student List
```tsx
// src/components/students/List.tsx
import React, { useEffect, useState } from 'react';
import api from '@/api/client';
import { Table, TableRow, TableCell, Button } from '@/components/common';

interface Student {
  id: string;
  name: string;
  email: string;
  enrolledAt: string;
}

export default function StudentList() {
  const [students, setStudents] = useState<Student[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get('/students').then(({ data }) => {
      setStudents(data.data);
      setLoading(false);
    });
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <Table>
      <TableRow>
        <TableCell>Name</TableCell>
        <TableCell>Email</TableCell>
        <TableCell>Enrolled</TableCell>
        <TableCell>Actions</TableCell>
      </TableRow>
      {students.map((student) => (
        <TableRow key={student.id}>
          <TableCell>{student.name}</TableCell>
          <TableCell>{student.email}</TableCell>
          <TableCell>{formatDate(student.enrolledAt)}</TableCell>
          <TableCell>
            <Button variant="secondary" size="sm">
              View Details
            </Button>
          </TableCell>
        </TableRow>
      ))}
    </Table>
  );
}
```

## 7. Authentication & Authorization

### 7.1 JWT Implementation
```go
// middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "Missing token"})
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := jwt.Verify(tokenString, jwtKey)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
            return
        }

        c.Set("userID", claims.UserID)
        c.Set("userRole", claims.Role)
        c.Next()
    }
}
```

### 7.2 Role-Based Access Control
```
SUPER_ADMIN: Full access to all endpoints
ADMIN:    Manage students, teachers, courses, reports
TEACHER:  View classes, take attendance, enter grades
STUDENT:  View own data, enroll in courses, view grades
PARENT:   View child's attendance and grades
```

## 8. Observability Stack

### 8.1 Structured Logging
```go
// logging middleware
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)
        
        logEntry := map[string]interface{}{
            "correlation_id": c.GetString("correlationID"),
            "user_id":        c.Get("userID"),
            "user_role":      c.Get("userRole"),
            "method":         c.Request.Method,
            "path":           c.Request.URL.Path,
            "status":         c.Writer.Status(),
            "latency_ms":     duration.Milliseconds(),
            "ip":             c.ClientIP(),
        }
        logger.Info("HTTP request", logEntry)
    }
}
```

### 8.2 Key Metrics
- **p95 API latency**: < 200ms
- **Error rate**: < 1%
- **Cache hit ratio**: > 70%
- **Database query count**: < 5 per request
- **JWT validation time**: < 5ms

## 9. Security Considerations

### 9.1 Data Protection
- **TLS 1.3** for all external communication
- **bcrypt** with cost factor 12 for password hashing
- **Parameterized queries** via Prisma (prevents SQL injection)
- **Input validation** with Zod on both frontend and backend
- **Rate limiting** per IP and per user role

### 9.2 GDPR Compliance
- Right to access student data
- Right to data rectification
- Right to erasure (with legal retention for 7 years for audit logs)
- Data minimization - only collect necessary fields
- Explicit consent for data processing

## 10. Deployment Architecture

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