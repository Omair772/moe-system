# 📊 Data Flow Diagram — Ministry of Education System

## 1. Purpose & Scope

This document defines the **Data Flow Diagrams (DFD)** for the Ministry of Education integrated system, showing how data flows between the core system and adjacent government systems. The diagrams follow the **Clean Architecture** principles from pro-skills and the multi-agent orchestration patterns from delegate-skills.

## 2. Context Diagram (Level 0)

```
┌──────────────────────────────────────────────────────────────┐
│                    MINISTRY OF EDUCATION SYSTEM               │
│  ┌─────────────┐  ┌─────────────────────────────────────┐ │
│  │  React UI   │  │  External Systems (Integration)       │ │
│  └─────▲───────┘  └───────▲──────────────────────┘ │
│        │                │ │                     │ │
│        │                │ │                     │ │
│        │                │ │                     │ │
│        │                │ │                     │ │
│  REST  │                API│ │  Government API │ │
│  API   │                │ │  Endpoints      │ │
│  Gin   │                │ │                 │ │
│        │                │ │                 │ │
│        │                │ │                 │ │
│        ▼                ▼ ▼                     ▼ ▼ │
│  ┌─────────────────────────────────────────────────┐ │
│  │  Go Backend (Gin + Prisma)                      │ │
│  │  ── PostgreSQL ──►  Database                      │ │
│  │  ── Redis Cache ────►  TTL Cache                   │ │
│  └─────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────┘
```

### 2.1 External Systems Integration

| System | Integration Purpose | Data Exchange | API Contract |
|--------|-------------------|---------------|--------------|
| **Civil Status Department** | Student registration, nationality data | Name, ID, parent info, birth date | POST /api/v1/civil-status/verify, GET /api/v1/civil-status/{id} |
| **Higher Education Ministry** | Graduate tracking, equivalency | Degree info, graduation date, university | POST /api/v1/higher-edu/status, GET /api/v1/higher-edu/{studentId} |
| **Financial Ministry** | Fee management, scholarships | Tuition fees, scholarship amounts, payment status | POST /api/v1/financial/scholarship, GET /api/v1/financial/fees |
| **Labor Ministry** | Graduate employment, work permits | Job placement, work visa status | POST /api/v1/labor/visa, GET /api/v1/labor/{graduateId} |

## 3. Level-1 DFD: Core System Processes

```
─────────────────────────────────────────────────────────────────────
│                          PROCESS 0: CORE SYSTEM                    │
│  (Ministry of Education Integrated System)                         │
├─────────────────────────────────────────────────────────────────────
│  │                                                                 │
│  │  1.0  Student Management                                         │
│  │ ────────────────────────                                        │
│  │  1.1  Enroll Student                         │─────►│  1.1.1 Validate Prerequisites              │
│  │                                     │       │  1.1.2 Check Course Capacity               │
│  │                                     │       │  1.1.3 Create Enrollment (Prisma)         │
│  │                                     │       │  1.1.4 Update Redis Cache (TTL: 30min)   │
│  │                                     │       │  1.1.5 Publish ENROLLMENT_CREATED Event   │
│  │                                                                 │
│  │  1.2  Update Student Profile           │─────►│  1.2.1 Validate Zod Schema                 │
│  │                                     │       │  1.2.2 Update PostgreSQL (Prisma)         │
│  │                                     │       │  1.2.3 DELETE cache:student:{id}          │
│  │                                     │       │  1.2.4 SET cache:student:{id} (TTL: 30min)│
│  │                                                                 │
│  │  2.0  Course Management                                              │
│  │ ────────────────────────                                            │
│  │  2.1  Create/Update Course             │─────►│  2.1.1 Zod Validation                      │
│  │                                     │       │  2.1.2 Prisma Upsert                      │
│  │                                     │       │  2.1.3 DELETE cache:course:{code}         │
│  │                                     │       │  2.1.4 SET cache:course:{code} (TTL: 2hr) │
│  │                                                                 │
│  │  3.0  Attendance & Grading                                             │
│  │ ────────────────────────                                            │
│  │  3.1  Mark Attendance              │─────►│  3.1.1 Prisma Attendance Record             │
│  │                                     │       │  3.1.2 DELETE cache:attendance:{date}     │
│  │                                     │       │  3.1.3 TTL: 5min                         │
│  │                                     │       │  3.1.4 If absence > threshold → Notify   │
│  │                                                                 │
│  │  3.2  Enter Grade                   │─────►│  3.2.1 Prisma Grade Update                  │
│  │                                     │       │  3.2.2 DELETE cache:student:{id}          │
│  │                                     │       │  3.2.3 DELETE cache:gpa:{studentId}       │
│  │                                     │       │  3.2.4 TTL: 30min / 1hr                   │
│  │                                     │       │  3.2.5 Publish GRADE_ENTERED Event        │
│  │                                                                 │
│  │  4.0  Reporting & Analytics                                               │
│  │ ────────────────────────                                            │
│  │  4.1  Generate Report              │─────►│  4.1.1 Query PostgreSQL via Prisma          │
│  │                                     │       │  4.1.2 Format (PDF/CSV/Excel)             │
│  │                                     │       │  4.1.3 Store in DB + S3                   │
│  │                                     │       │  4.1.4 SET cache:api:response:reports (TTL: 5min) │
│  │                                                                 │
├─────────────────────────────────────────────────────────────────────
```

## 4. Level-2 DFD: Student Enrollment Flow

```
─────────────────────────────────────────────────────────────────────
│                          PROCESS 1.1: ENROLL STUDENT               │
│  (Student enrolls in a course)                                     │
├────────────▲───────────────────────────────────────────────────────┤
│              │                                                   │
│  Input:      │  Output:                                          │
│  • student_id│  • enrollment_id                                    │
│  • course_code│  • student updated cache                           │
│  • semester  │  • course updated cache                           │
│              │  • audit log                                       │
├────────────▼───────────────────────────────────────────────────────┤
│  │ 1. Receive REST API POST /api/v1/enrollments                     │
│  │  │                                                             │
│  │  │  1.1 Validate JWT token (2hr expiry)                        │
│  │  │                                                             │
│  │  │  1.2 Zod validation (student_id format, course existence)   │
│  │  │                                                             │
│  │  │  1.3 Check Prisma: course exists & active                     │
│  │  │                                                             │
│  │  │  1.4 Check Redis cache:course:{code} (TTL: 2hr, HIT→refresh) │
│  │  │                                                             │
│  │  │  1.5 If cache MISS → Prisma GET course                       │
│  │  │                                                             │
│  │  │  1.6 Validate prerequisites (Prisma query)                  │
│  │  │                                                             │
│  │  │  1.7 Check capacity (Prisma: course.students count)          │
│  │  │                                                             │
│  │  │  1.8 If valid → Prisma ENROLLMENT CREATE                      │
│  │  │                                                             │
│  │  │  1.9 DELETE cache:student:{studentId} (TTL reset to 30min)   │
│  │  │                                                             │
│  │  │  1.10 DELETE cache:course:{code} (TTL reset to 2hr)          │
│  │  │                                                             │
│  │  │  1.11 Publish ENROLLMENT_CREATED event (Redis Pub/Sub)       │
│  │  │                                                             │
│  │  │  1.12 Return 201 with enrollment confirmation               │
│  │  │                                                             │
│  │  └─[Validation FAIL] → 400 Error                                │
│  │                                                             │
│  └─[Auth FAIL] → 401 Unauthorized                                  │
├─────────────────────────────────────────────────────────────────────
```

## 5. Level-2 DFD: Attendance Marking Flow

```
────────────▲───────────────────────────────────────────────────────
│              PROCESS 3.1: MARK ATTENDANCE                         │
│  (Teacher marks class attendance)                                 │
├────────────▼───────────────────────────────────────────────────────
│  │                                                                 │
│  Input:      │  Output:                                                  │
│  • class_id  │  • marked_count                                         │
│  • date      │  • updated attendance cache                            │
│  • marks[]   │  • absence alerts (if threshold exceeded)             │
│              │  • audit log                                              │
│  │                                                                 │
│  ├─ 1. Receive POST /api/v1/attendance                           │
│  │                                                                 │
│  │  1.1 JWT Auth (Teacher role check)                            │
│  │                                                                 │
│  │  1.2 Zod validation (marks array, valid statuses)             │
│  │                                                                 │
│  │  1.3 Prisma: Find class, validate teacher assignment          │
│  │                                                                 │
│  │  1.4 Check Redis: cache:attendance:{date} (TTL: 5min)         │
│  │                                                                 │
│  │  1.5 If HIT → use cached data, else → Prisma query            │
│  │                                                                 │
│  │  1.6 For each student mark:                                    │
│  │  │   │  Prisma: Attendance UPSERT                                │
│  │  │   │  │  DELETE cache:student:{id} (if GPA affected)        │
│  │  │   │  │                                                    │
│  │  │   │  └─ Absent count > threshold → Publish ABSENCE_EVENT    │
│  │  │   │                                                    │
│  │  │   └─ Present/Absent/Excused → Record                       │
│  │  │                                                   │
│  │  1.7 After all marks:                                          │
│  │  │   │  DELETE cache:attendance:{date} (TTL reset to 5min)   │
│  │  │   │                                                    │
│  │  │   │  Publish ATTENDANCE_MARKED event                       │
│  │  │                                                   │
│  │  │   │  Return success response                               │
│  │  │                                                   │
│  │  └─ [Validation FAIL] → 400 Error                             │
│  │                                                                 │
│  └─ [Auth FAIL] → 403 Forbidden                                 │
├─────────────────────────────────────────────────────────────────────
```

## 6. Data Flow Between Adjacent Systems

### 6.1 Ministry of Education ↔ Civil Status Department

```
┌─────────────────┐      ┌─────────────────────┐
│   Education Sys │      │   Civil Status Dept │
│   (Go Backend)  │      │   (Government API)  │
├─────────────────┤      ├─────────────────────┤
│  GET /students   ├────► │  verifyStudent(id)  │
│  POST /enrollment│     │  return: {name, dad, │
│  │               │     │   mom, birthDate}   │
│  └───────────────┘      └─────────────────────┘
│                       │                           │
│                       │ 1. Request verification     │
│                       │ 2. Return student data        │
│                       │ 3. Education system caches    │
│                       │     the verified data (TTL: 2hr)│
│                       ▼                           ▼
│                Cache: student:{id}         Status: verified
└─────────────────────────────────────────────────────────────────
```

**Data Exchange Fields:**
```
Student Verification Response:
{
  "success": true,
  "data": {
    "full_name": "Mohammed Ali",
    "father_name": "Ali Hassan",
    "mother_name": "Fatima Hassan",
    "national_id": "1234567890",
    "birth_date": "1995-03-15",
    "nationality": "Saudi",
    "verified_at": "2024-01-15T10:30:00Z"
  }
}
```

### 6.2 Ministry of Education ↔ Higher Education Ministry

```
┌─────────────────┐      ┌─────────────────────┐
│   Education Sys │      │   Higher Edu Dept   │
│   (Go Backend)  │      │   (Government API)  │
├─────────────────┤      ├─────────────────────┤
│  POST /graduation │    │  checkGraduation(   │
│  │               │    │   studentId, year)  │
│  └───────────────┘    │  return: {degree,     │
│                       │           university}   │
│                       ▼                           ▼
│                Cache: gpa:{studentId}   Status: graduated
└─────────────────┘      └─────────────────────┘
```

**Graduation Status Response:**
```json
{
  "success": true,
  "data": {
    "student_id": "uuid",
    "graduated": true,
    "degree": "Bachelor of Science in Computer Science",
    "university": "King Saud University",
    "graduation_date": "2023-06-30",
    "gpa": "3.84/4.0",
    "verified_at": "2024-01-15T10:30:00Z"
  }
}
```

### 6.3 Data Flow with TTL Caching

```
┌─────────────────────────────────────────────────────────────────┐
│                    TTL CACHING STRATEGY                         │
├────────────▲────────────────────────────▲───────────────────────┤
│            │                     │                      │
│  Student   │     Course        │   Attendance         │
│  Cache:    │     Cache:        │   Cache:             │
│  student:{id}│   course:{code}   │   attendance:{date}  │
│            │                     │                      │
│  TTL: 30min│   TTL: 2hr        │   TTL: 5min          │
│            │                     │                      │
│  Flow:     │   Flow:           │   Flow:              │
│  │           │                   │                      │
│  │  GET →    │  GET →            │  GET →               │
│  │   HIT?    │   HIT?            │  HIT?                │
│  │   │       │   │                   │   │                    │
│  │   ├─ YES  │   ├─ YES      │   ├─ YES             │
│  │   │   →   │   →           │   → Return cached    │
│  │   │   │     │   │        data (fresh)              │
│  │   ├─ NO   │   ├─ NO       │   ├─ NO              │
│  │   │   │     │   │        │   │  → Prisma Query    │
│  │   │   │     │   │        │   │  │                 │
│  │   │   │     │   │        │   │  │  → DB Query      │
│  │   │   │     │   │        │   │  │                 │
│  │   │   │     │   │        │   │  │  → Update Cache  │
│  │   │   │     │   │        │   │  │  │  (TTL Set)     │
│  │   │   │     │   │        │   │  │  │              │
│  │   │   │     │   │        │   │  │  │              │
│  │   │   │     │   │        │   │  │  │              │
│  │   │   │     │   │        │   │  │  │              │
│  │   │   │     │   │        │   │  │  │              │
│  │   │   │     │   │        │   │  │  │              │
│  │   │   │     │   │        │   │  │  │              │
│  └──┬────────┘   └─────────────┘   └────────────────┘
     │                   │                   │
   Refresh           Refresh           Refresh
   (after write)     (after write)     (after write)
```

## 7. API Response Flow (With Cache)

```
Client Request
    │
    ▼
[React Frontend] → axios.get('/api/v1/students/123')
    │
    ▼
[Gin Middleware Chain]
    │
    ├─ [Correlation ID] → req-abc123
    │
    ├─ [Auth] → JWT validated (user=STUDENT)
    │
    ├─ [Rate Limit] → 100/min per IP
    │
    ├─ [Cache Lookup] → Redis HGET student:123
    │       │
    │       ├─ HIT → Return cached JSON (TTL: 30min)
    │       │   └─ c.Abort() → Skip remaining middleware
    │       │
    │       └─ MISS → Continue to step 2
    │
    ├─ [Validation] → Zod schema check
    │
    ├─ [Business Logic] → Prisma Student FindUnique
    │       │
    │       ├─ Success →
    │       │   ├─ Serialize to JSON
    │       │   ├─ Redis SET student:123 (TTL: 30min)
    │       │   └─ Return 200 OK
    │       │
    │       └─ Failure →
    │           └─ 404 Not Found
    │
    └─ [Response] → JSON → React Component → UI Update
```

## 8. Event Flow Between Systems

```
┌──────────────────────────────────────────────────────────────┐
│                    EVENT PROPAGATION                           │
├────────────▲────────────────────────────▲───────────────────────┤
│            │                     │                      │
│  Education │     Civil Status    │   Higher Edu         │
│  System    │     Department      │   Ministry           │
│            │                     │                      │
│  │         │                     │                      │
│  │  1. Enroll Student          │                      │
│  │         │                     │                      │
│  │  2. Publish ENROLLMENT_CREATED        │                      │
│  │         │                     │                      │
│  │  3. Redis Pub/Sub → events.ENROLLMENT_CREATED │              │
│  │         │                     │                      │
│  │  4. Civil Status Subscriber      │                      │
│  │         │                     │                      │
│  │  5. Verify student data        │                      │
│  │         │                     │                      │
│  │  6. Update civil records        │                      │
│  │         │                     │                      │
│  └───────────────────────────────────────────────────────────────┘
```