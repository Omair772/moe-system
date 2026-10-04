# 🧠 Mind Map — Ministry of Education System

## 1. High-Level System Overview

```
MINISTRY OF EDUCATION INTEGRATED SYSTEM
├─ Authentication & Authorization
│  ├─ JWT Tokens (2hr expiry)
│  ├─ RBAC (Roles: SUPER_ADMIN, ADMIN, TEACHER, STUDENT, PARENT)
│  └─ Permission Granularity
├─ Student Management
│  ├─ Enrollment & Drop
│  ├─ Personal Information
│  ├─ Academic Records
│  └─ Attendance Tracking
├─ Course Management
│  ├─ Curriculum Catalog
│  ├─ Course Assignments
│  ├─ Prerequisites & Credits
│  └─ Learning Outcomes
├─ Attendance & Assessment
│  ├─ Daily Attendance Records
│  ├─ Grade Entry & Calculation
│  ├─ GPA Computation
│  └─ Progress Reports
├─ Timetable & Scheduling
│  ├─ Class Schedules
│  ├─ Room Assignments
│  ├─ Substitution Management
│  └─ Conflict Resolution
├─ Reporting & Analytics
│  ├─ Student Performance Dashboards
│  ├─ Teacher Efficiency Metrics
│  ├─ Course Effectiveness
│  └─ Institutional Analytics
└─ Notification & Communication
   ├─ In-App Messages
   ├─ Email Notifications
   ├─ SMS Alerts
   └─ Weekly Digests
```

---

## 2. Entity Relationship Diagram (ERD)

```
STUDENT ┬─ 1:n ENROLLMENT ─── n:1 COURSE
          │                   │
          │                   ├─ n:1 STUDENT (repeated for each course)
          │                   └─ attributes: grade, attended%
          │
          ├─ 1:n ATTENDANCE_RECORDS (daily)
          │
          └─ 1:n GRADE_HISTORY

COURSE ┬─ 1:n ENROLLMENT ─── n:1 STUDENT (repeated for each student)
       │
       ├─ 1:n TEACHES (junction: Teacher-Course)
       │    └─ n:1 TEACHER
       │
       ├─ 1:n PREREQUISITES (self-referential)
       │    └─ course_id → prerequisite_course_id
       │
       └─ 1:n SCHEDULE_ENTRIES

TEACHER ┬─ 1:n TEACHES ─── n:1 COURSE
        │
        ├─ 1:m STUDENTS (via courses)
        │
        └─ 1:n GRADE_ENTRIES

ENROLLMENT ┬─ many:1 STUDENT
            └─ many:1 COURSE
               attributes: enrollment_date, status (active/completed/dropped)
```

---

## 3. Data Flow Mind Map

```
                    INPUT LAYER
                            │
                    ┌───────────────▼───────────────┐
                    │  API Requests (React → Go)    │
                    └───────────────▲───────────────┘
                                    │
               ┌────────────────────┼─────────────────────┐
               │                    │                     │
               ▼                    ▼                     ▼
    [Auth Middleware]   [Validation]   [Cache Lookup]   │
               │                    │                     │
               ▼                    ▼                     ▼
    ┌─────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐
    │  Use Case Layer │  │  Prisma ORM Layer   │  │  Redis Cache Layer  │
    └─────────────────┘  └─────────────────────┘  └─────────────────────┘
               │                    │                     │
               ▼                    ▼                     ▼
    ┌───────────────────────────────────────────────────────────────┐
    │              DATABASE LAYER (PostgreSQL)                      │
    └───────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
                    ┌─────────────────────────────────────┐
                    │   OUTPUT LAYER — React Components   │
                    └─────────────────────────────────────┘
```

### 3.1 Request Processing Flow

```
Client Browser
    │
    ▼
React Component → axios/api Client
    │
    ▼
HTTP Request → Gin Middleware Chain
    │
    ├─ [Correlation ID] → X-Request-ID header
    ├─ [Auth] → JWT Validation (2hr expiry)
    ├─ [Rate Limit] → Per-IP/User throttling
    ├─ [Validation] → Zod schema check
    └─ [Cache] → Redis cache-first (TTL-based)
        │
        ├─ Cache HIT → Return cached JSON
        └─ Cache MISS → Continue to handler
            │
            ▼
    Use Case (Business Logic)
        │
        ├─ Query Prisma ORM
        ├─ Transform data (DTOs)
        ├─ Apply business rules
        └─ Return result
            │
            ▼
    → Prisma → PostgreSQL
            │
            ▼
    ← Response Transformation
            │
            ▼
    → Cache Store (TTL set)
            │
            ▼
    → HTTP Response → React Component
            │
            ▼
    → UI Update (state, re-render)
```

### 3.2 Cache Warming & Invalidation Strategy

```
Key Structure:          TTL:                Invalidation Trigger:
cache:student:{id}     30 minutes         Student data UPDATE
cache:course:{code}    2 hours           Course catalog CHANGE
cache:attendance:{date} 5 minutes        New attendance RECORD
cache:session:{token}   24 hours          Session LOGOUT/INVALIDATE
cache:api:{path}        1 minute          Any API response COMPLETE
```

```
┌──────────────────────────────────────────────────────────────┐
│         CACHE INVALIDATION WORKFLOW                            │
├──────────────────────────────────────────────────────────────┤
│  1. API endpoint receives UPDATE request                       │
│  2. Use case updates DB via Prisma                              │
│  3. After successful DB write:                                   │
│     ├─ DEL cache:student:{id}                                   │
│     ├─ DEL cache:course:{code}                                  │
│     ├─ DEL cache:attendance:{date}                              │
│     └─ SET cache:api:response:{path} (new TTL)                 │
│  4. Return updated data to client                               │
│  5. Next READ request → Cache MISS → DB query → Fresh cache    │
│                                                                 │
│  Alternative: Write-Through Cache                              │
│  ───────────────────────────────────────────────────────────── │
│  1. API receives WRITE request                                   │
│  2. Update DB via Prisma                                         │
│  3. Immediately UPDATE cache with new value                      │
│  4. Cache always consistent with DB                              │
└──────────────────────────────────────────────────────────────┘
```

---

## 4. Workflow Diagrams

### 4.1 Student Enrollment Workflow

```
[Student Logs In]
    │
    ▼
[Select Courses from Catalog]
    │
    ▼
[Validate Prerequisites]  ──[Pass]──► [Check Capacity]
    │                              │
    │                              ├─[Pass]──► [Create Enrollment]
    │                              │        │
    │                              │        ▼
    │                              │   [Send Confirmation Email/SMS]
    │                              │        │
    │                              └─[Fail]─► [Show Error: Capacity Full]
    │
    └─[Fail]──────────────────────► [Show Error: Prerequisite Not Met]

[Enrollment Complete]
    │
    ▼
[Invalidate Cache: cache:course:{code}]  (TTL reset)
    │
    ▼
[Notify Student & Parent] (via email/SMS)
    │
    ▼
[Log Audit Entry] (CREATE_ENROLLMENT action)
    │
    ▼
[Refresh Student Dashboard]
```

### 4.2 Attendance Tracking Workflow

```
[Class Session Begins]
    │
    ▼
[Teacher Marks Attendance]
    │
    ├─ Present   │ Absent │ Excused
    │
    ▼
[Validate Attendance Rules]
    │
    ├─ ✓ └─ [Record in PostgreSQL via Prisma]
    │        │
    │        ▼
    │   [Update Redis: cache:attendance:{date}]
    │        │         TTL: 5 minutes
    │        │
    │        ▼
    │   [Trigger Alert if Absent > Threshold]
    │        │
    │        └─ [Notify Parents auto-matically if > 3 absences/week]
    │
    └─ ✗ └─ [Show Error: Invalid status/date]
           │
           ▼
   [Log: attendance_validation_failed]

[End of Day]
    │
    ▼
[Generate Daily Attendance Report]
    │
    ▼
[Export to PDF/CSV for administration]
    │
    ▼
[Parent Notification (if needed)]
```

### 4.3 Grade Entry Workflow

```
[Teacher Enters Grade]
    │
    ▼
[Validate Grade Range] (0-100 or A-F)
    │
    ├─ ✓ └─ [Prisma: enrollment.update grade]
    │        │
    │        ▼
    │   [Invalidate Cache: cache:student:{id}]
    │        TTL: 30 minutes (refresh student profile)
    │
    │        ▼
    │   [Calculate GPA Impact]
    │        │
    │        └─ [Update cumulative GPA in Student model]
    │
    │        ▼
    │   [Generate Grade Report]
    │        │
    │        └─ [Send to Student/Parent]
    │
    └─ ✗ └─ [Show Error: Grade out of range]
           │
           ▼
   [Audit: invalid_grade_attempt]
```

---

## 5. Priority Levels & Critical Paths

```
CRITICAL (Immediate — within 1 hour)
├─ System outage (API down)
├─ Security breach detected
├─ Data corruption in grades
└─ Parent-facing system down

HIGH (Within 4 hours)
├─ Authentication failure (can't login)
├─ Grade calculation error
├─ Attendance tracking failure
┐

MEDIUM (Within 24 hours)
├─ Search not returning expected results
├─ Cache miss rate > 30%
├─ Report generation slow
┐

LOW (Within 1 week)
├─ UI cosmetic issues
├─ Non-critical feature enhancement
├─ Documentation updates
└─ Performance tuning (non-critical paths)
```

### 5.1 Critical Success Paths

```
1. Student Login → Enroll in Course → Attendance Taken → Grade Entered → Report Generated
   │                                                                     │
   └─✓ Each step must complete in < 2s                                  └─✓ End-to-end flow time < 10s

2. Admin Generate Report → Export Data → Download PDF
   │                                                                     │
   └─✓ Must handle 10,000+ records                                     └─✓ Must complete in < 30s

3. Parent View Child's Grades → Authenticate → Display Data
   │                                                                     │
   └─✓ Must show real-time data (cache refresh)                          └─✓ Must not show other students' data
```

---

## 6. Technology Integration Points

```
                          ┌────────────────────┐
                          │   GO BACKEND (Gin) │
                          └───────▲───────▲─────┘
                                │       │
          ┌─────────────────────┘       └─────────────────────┐
          │                                                   │
   ┌────────────▼─────────────┐              ┌────────────▼────────────┐
   │  Prisma ORM → PostgreSQL │              │   Redis (TTL Cache)    │
   └────────────▲─────────────┘              └────────────▲─────────────┘
              │                                                  │
              ▼                                                  ▼
   ┌─────────────────────┐          ┌───────────────────────┐
   │  JWT Auth Middleware│          │  Correlation ID Gen    │
   └────────────▲─────────┘          └────────────▲───────────┘
              │                                  │
              ▼                                  ▼
   ┌─────────────────────┐          ┌───────────────────────┐
   │  Rate Limiting (overseer)│  →   │  Structured Logging   │
   └────────────▲─────────────┘          └────────────▲───────────┘
              │                                  │
              ▼                                  ▼
   ┌─────────────────────┐          ┌───────────────────────┐
   │  Zod Validation      │  →   │  Swagger/OpenAPI       │
   └────────────▲─────────┘          └────────────▲───────────┘
              │                                  │
              ▼                                  ▼
   ┌───────────────────────┐        ┌───────────────────────┐
   │  Go Tests (unit/integration)│  │  React Component Lib  │
   └────────────▲─────────────┘        └────────────▲───────────┘
              │                                  │
              ▼                                  ▼
   ┌───────────────────────┐        ┌───────────────────────┐
   │  Docker & Kubernetes  │  →   │  Vite + React 18      │
   └───────────────────────┘        └────────────▲───────────┘
                                                 │
                                                 ▼
                                          ┌─────────────────┐
                                          │   CDN / Static  │
                                          └─────────────────┘
```

---

## 7. Audit Log Mind Map

```
AUDIT LOG ENTITIES
│
├─ ACTIONS
│  ├─ AUTHENTICATION: LOGIN, LOGOUT, TOKEN_REFRESH
│  ├─ STUDENT: CREATE, UPDATE, SOFT_DELETE, ENROLL, DROP
│  ├─ COURSE: CREATE, UPDATE, ASSIGN_TEACHER, CHANGE_CREDITS
│  ├─ GRADE: ENTER, UPDATE, REMOVE
│  ├─ ATTENDANCE: MARK, UPDATE, EXPORT
│  ├─ ADMIN: ROLE_ASSIGN, PERMISSION_CHANGE, DATA_EXPORT
│  └─ SYSTEM: API_CALL, CACHE_ACCESS, ERROR_LOG
│
├─ ENTITIES AFFECTED
│  ├─ Student (id, name, email, enrollment_status)
│  ├─ Course (id, title, code, credits, active_status)
│  ├─ Enrollment (id, student_id, course_id, grade, attended%)
│  ├─ Teacher (id, name, email, assigned_courses)
│  └─ Attendance record (id, student_id, date, status)
│
├─ TIMESTAMPS
│  ├─ Action timestamp (UTC)
│  ├─ Request start time
│  ├─ DB query duration
│  └─ Response completion time
│
├─ USER CONTEXT
│  ├─ User ID (authenticated)
│  ├─ User role/permissions
│  ├─ IP address
│  ├─ User-Agent (browser/app)
│  └─ Correlation ID (request tracing)
│
├─ CHANGES DETAIL
│  ├─ Before state (JSON)
│  ├─ After state (JSON)
│  ├─ Fields modified
│  └─ Old value → New value
│
└─ METADATA
   ├─ Success/failure flag
   ├─ Error codes (if failed)
   ├─ Rate limit status
   └─ Compliance flags (GDPR, etc.)
```