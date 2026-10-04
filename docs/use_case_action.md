# 🎯 Use Case Actions — Ministry of Education System

## 1. Purpose & Scope

This document defines the specific **actions** that can be performed within each use case scenario for the Ministry of Education integrated system. While use case scenarios describe the "what" (the overall flow), use case actions describe the "how" — the specific operations, commands, and API endpoints that implement each scenario.

## 2. Action Categorization

Actions are categorized by:
- **Domain**: Student, Course, Attendance, Grade, Report, Auth, Profile
- **Operation**: Create, Read, Update, Delete (CRUD), plus special operations
- **Authorization**: Which roles can perform each action
- **TTL Impact**: Which caches are affected and with what TTL reset

## 3. Action Reference Table

| Action ID | Use Case | Operation | Endpoint | Roles Allowed | Cache Invalidation | TTL |
|-----------|----------|-----------|----------|--------------|-------------------|-----|
| AC-01 | Enrollment | Create | POST /api/v1/enrollments | STUDENT, PARENT, ADMIN | cache:student:{id}, cache:course:{code} | 30 min / 2 hr |
| AC-02 | Enrollment | Read | GET /api/v1/enrollments/:id | STUDENT (own), TEACHER (assigned), ADMIN | None | - |
| AC-03 | Enrollment | Delete/Withdraw | DELETE /api/v1/enrollments/:id | STUDENT (own), ADMIN | cache:student:{id} | 30 min |
| AC-04 | Attendance | Create/Mark | POST /api/v1/attendance | TEACHER (assigned), ADMIN | cache:attendance:{date} | 5 min |
| AC-05 | Attendance | Read | GET /api/v1/attendance/:date | STUDENT (own), PARENT (child), ADMIN | cache:attendance:{date} | 5 min |
| AC-06 | Attendance | Update/Edit | PUT /api/v1/attendance/:id | TEACHER (with admin auth), ADMIN | cache:attendance:{date} | 5 min |
| AC-07 | Grade | Create/Enter | POST /api/v1/grades | TEACHER (assigned), ADMIN | cache:student:{id} | 30 min |
| AC-08 | Grade | Read | GET /api/v1/grades/:id | STUDENT (own), PARENT (child), ADMIN | None | - |
| AC-09 | Grade | Update/Edit | PUT /api/v1/grades/:id | TEACHER (assigned), ADMIN | cache:student:{id} | 30 min |
| AC-10 | Student Profile | View | GET /api/v1/students/:id | STUDENT (own), PARENT (linked), ADMIN | cache:student:{id} | 30 min |
| AC-11 | Student Profile | Update | PUT /api/v1/students/:id | STUDENT (own), ADMIN | cache:student:{id} | 30 min |
| AC-12 | Course | List | GET /api/v1/courses | All roles (public with filters) | cache:course:{code} | 2 hr |
| AC-13 | Course | Create | POST /api/v1/courses | ADMIN, SUPER_ADMIN | cache:course:{code} (invalidate) | 2 hr |
| AC-14 | Course | Update | PUT /api/v1/courses/:id | ADMIN, SUPER_ADMIN | cache:course:{code} | 2 hr |
| AC-15 | Course | Delete | DELETE /api/v1/courses/:id | SUPER_ADMIN | cache:course:{code} | 2 hr |
| AC-16 | Authentication | Login | POST /api/v1/auth/login | All (public endpoint) | session:{token} | 24 hr |
| AC-17 | Authentication | Logout | POST /api/v1/auth/logout | All authenticated | session:{token} (blacklist) | Until expiry |
| AC-18 | Authentication | Refresh | POST /api/v1/auth/refresh | All authenticated | session:{token} (new) | 24 hr |
| AC-19 | Report | Generate | POST /api/v1/reports/generate | ADMIN, SUPER_ADMIN | Various (depends on report) | Varies |
| AC-20 | Report | Download | GET /api/v1/reports/:id | ADMIN, SUPER_ADMIN | None | - |
| AC-21 | Parent | View Child | GET /api/v1/parents/children | PARENT (linked) | cache:student:{id} | 30 min |
| AC-22 | Parent | Update Preferences | PUT /api/v1/parents/preferences | PARENT (own) | None | - |

## 4. Detailed Action Specifications

### AC-01: Enrollment Create Action

**Action Name:** `enrollment_create`  
**Endpoint:** `POST /api/v1/enrollments`  
**Request Body:**
```json
{
  "student_id": "string (UUID)",
  "course_code": "string (course code)",
  "semester": "string (e.g., 'Fall 2024')",
  "section": "string (optional)"
}
```

**Authorization:** 
- STUDENT: Can enroll themselves
- PARENT: Can enroll linked student
- ADMIN: Can enroll any student

**Validation Rules:**
1. Student must be authenticated
2. Course must exist and be active
3. Student must not already be enrolled in this course
4. Student must not exceed maximum course load (typically 6 courses)
5. Prerequisites must be satisfied (if course has prerequisites)
6. Course capacity must not be full

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "enrollment_id": "uuid",
    "student_id": "uuid",
    "course_code": "CS101",
    "enrolled_at": "2024-01-15T10:30:00Z",
    "status": "active"
  },
  "meta": {
    "cache_invalidated": ["student:123", "course:CS101"],
    "new_ttls": {
      "student:123": "30min",
      "course:CS101": "2hr"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Error Responses:**
- 400: Prerequisite not met, capacity full, max course load exceeded
- 401: Invalid/expired token
- 403: Student already enrolled in this course
- 422: Validation failed

**Cache Impact:**
- `DEL cache:student:{student_id}` — TTL reset to 30 min
- `DEL cache:course:{course_code}` — TTL reset to 2 hr
- `SET cache:api:response:enrollments:{student_id}` — TTL 1 min (temporary)

---

### AC-04: Attendance Mark Action

**Action Name:** `attendance_mark`  
**Endpoint:** `POST /api/v1/attendance`  
**Request Body:**
```json
{
  "class_id": "uuid",
  "session_date": "2024-01-15",
  "marks": [
    {
      "student_id": "uuid",
      "status": "Present|Absent|Excused"
    },
    ...
  ]
}
```

**Authorization:**
- TEACHER: Can mark attendance for assigned classes
- ADMIN: Can mark attendance for any class

**Status Values:**
- `Present` — Student was present
- `Absent` — Student was absent (unexcused)
- `Excused` — Student was absent with valid excuse

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "marked_count": 30, // number of students marked
    "present_count": 28,
    "absent_count": 2,
    "excused_count": 0
  },
  "meta": {
    "cache_invalidated": ["attendance:2024-01-15"],
    "new_ttl": "5min"
  },
  "timestamp": "2024-01-15T10:35:00Z"
}
```

**Alternative Paths:**
- **Duplicate Mark:** If student already marked for this date → Error "Attendance already marked for this student"
- **Session Locked:** System shows "Attendance finalized for this date" → Only admin can reopen
- **Invalid Status:** If status not in {Present, Absent, Excused} → Error "Invalid attendance status"

**Cache Impact:**
- `DEL cache:attendance:{date}` — TTL reset to 5 min (fresh data required)
- Individual student caches may also be invalidated if GPA calculation affected

**Post-Action Triggers:**
- If absent count > threshold (e.g., 3+ absences in a week) → Auto-trigger parent notification
- Update student attendance percentage
- Generate attendance reports for administration

---

### AC-07: Grade Enter Action

**Action Name:** `grade_enter`  
**Endpoint:** `POST /api/v1/grades`  
**Request Body:**
```json
{
  "enrollment_id": "uuid",
  "student_id": "uuid",
  "course_code": "string",
  "grade": "string (A, B, C, D, F or numeric 0-100)",
  "assignment_name": "string (optional)",
  "max_points": "number (optional, if numeric grade)",
  "exam_date": "date (optional)"
}
```

**Authorization:**
- TEACHER: Can enter grades for assigned courses
- ADMIN: Can enter grades for any course

**Grade Scale:**
- Letter grades: A, B, C, D, F
- Numeric: 0-100 scale
- Pass/Fail options available

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "grade_id": "uuid",
    "enrollment_id": "uuid",
    "grade": "A",
    "letter_equivalent": "4.0 GPA",
    "updated_at": "2024-01-15T10:40:00Z"
  },
  "meta": {
    "cache_invalidated": ["student:123", "gpa:123"],
    "new_ttls": {
      "student:123": "30min",
      "gpa:123": "1hr"
    }
  },
  "timestamp": "2024-01-15T10:40:00Z"
}
```

**GPA Calculation (Automatic):**
- Each grade maps to grade points (A=4, B=3, C=2, D=1, F=0)
- GPA = Σ(grade_points × credit_hours) / Σ(credit_hours)
- Cached under `gpa:{student_id}` with TTL: 1 hour
- Cache invalidated on any grade change for that student

**Error Responses:**
- 400: Invalid grade value, enrollment not found, student not in course
- 401: Invalid/expired token
- 403: Teacher not assigned to this course
- 422: Validation failed (grade out of range)

**Cache Impact:**
- `DEL cache:student:{student_id}` — TTL reset to 30 min
- `DEL cache:gpa:{student_id}` — TTL reset to 1 hour (recalculate on next access)
- `SET cache:api:response:grades:{enrollment_id}` — TTL 1 min

---

### AC-11: Student Profile Update Action

**Action Name:** `profile_update`  
**Endpoint:** `PUT /api/v1/students/:id`  
**Request Body (Zod validation):**
```json
{
  "name": "string (optional, 2-100 chars)",
  "email": "string (optional, valid email, max 100 chars)",
  "phone": "string (optional, 10-20 chars)",
  "emergency_contact": "object (optional)",
  {
    "name": "string",
    "phone": "string",
    "relationship": "string"
  }
}
```

**Authorization:**
- STUDENT: Can update own profile (limited fields)
- ADMIN: Can update any student profile

**Validation (Zod Schema):**
- Email must be valid format if provided
- Phone must match format if provided
- Emergency contact relationship must be "parent", "guardian", or "spouse" if provided
- Name cannot be empty if provided

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "student_id": "uuid",
    "updated_fields": ["name", "email"], // which fields were actually updated
    "updated_at": "2024-01-15T10:45:00Z"
  },
  "meta": {
    "cache_invalidated": ["student:123"],
    "new_ttl": "30min"
  },
  "timestamp": "2024-01-15T10:45:00Z"
}
```

**Cache Impact:**
- `DEL cache:student:{student_id}` — Entire cache entry invalidated
- New data fetched from DB on next access
- TTL reset to 30 minutes

**Restricted Fields (Students Cannot Change):**
- Student ID (immutable)
- Enrollment status (admin-controlled)
- Password (handled via separate auth flow)
- Role/permissions (admin-controlled)

**Error Responses:**
- 400: Validation failed (Zod errors)
- 401: Invalid/expired token
- 403: Student trying to update another student's profile
- 404: Student not found
- 409: Email already taken (if email being changed)

---

### AC-19: Report Generate Action

**Action Name:** `report_generate`  
**Endpoint:** `POST /api/v1/reports/generate`  
**Request Body:**
```json
{
  "report_type": "string (enum: student_performance, attendance_analysis, course_effectiveness, teacher_performance, institutional_kpis)",
  "filters": {
    "start_date": "string (ISO date, optional)",
    "end_date": "string (ISO date, optional)",
    "student_id": "UUID (optional, for per-student reports)",
    "class_id": "UUID (optional, for class-specific reports)",
    "teacher_id": "UUID (optional)",
    "grade_level": "string (optional, e.g., '10', '11', '12')"
  },
  "format": "string (enum: pdf, csv, excel, html)",
  "include_charts": "boolean (default true)"
}
```

**Authorization:**
- ADMIN: Can generate any report
- SUPER_ADMIN: Can generate any report + manage report templates
- Other roles: Denied access

**Report Types:**
1. **student_performance** — Per-student grades, attendance, GPA
2. **attendance_analysis** — Absence patterns, trends, alerts
3. **course_effectiveness** — Pass rates, grade distribution, feedback
4. **teacher_performance** — Student feedback, attendance correlation
5. **institutional_kpis** — Promotion rates, dropout indicators, metrics

**Success Response (202 Accepted):**
```json
{
  "success": true,
  "data": {
    "report_id": "uuid",
    "status": "processing", // or "completed"
    "estimated_completion": "2024-01-15T11:00:00Z"
  },
  "meta": {
    "cache_invalidated": depending on filters,
    "estimated_size_mb": 2.5
  },
  "timestamp": "2024-01-15T10:50:00Z"
}
```

**Polling Status:**
```
GET /api/v1/reports/:id/status
{
  "status": "processing|completed|failed",
  "progress": "0-100%",
  "download_url": "if completed"
}
```

**Error Responses:**
- 400: Invalid report type, invalid date range, conflicting filters
- 401: Invalid/expired token
- 403: Insufficient permissions
- 422: Validation failed

**Cache Impact (depends on filters):**
- If `student_id` provided: `DEL cache:student:{student_id}`
- If class-specific: `DEL cache:course:{code}` for relevant courses
- If institutional: May invalidate multiple caches
- General: `SET cache:api:response:reports:{report_type}` — TTL 5 min (stale while processing)

**Post-Processing:**
- Report stored in PostgreSQL + optionally S3 for PDF/Excel files
- Download link expires after 24 hours
- Audit log: `REPORT_GENERATED` with parameters
- Notification sent to requester when complete

## 5. Action Validation Flow

All actions follow a consistent validation flow:

```
Request Received
    │
    ▼
[Authentication Middleware] → JWT Validation
    │
    ▼
[Authorization Middleware] → Role Check
    │
    ▼
[Input Validation] → Zod schema validation
    │
    ▼
[Business Rule Check] → Domain-specific validation
    │
    ▼
[Database Operation] → Prisma / PostgreSQL
    │
    ▼
[Cache Management] → Invalidate/Update TTL
    │
    ▼
[Audit Logging] → Write audit entry
    │
    ▼
[Response] → JSON API format
```

## 6. Error Handling for Actions

| Error Code | Scenario | Response Format |
|------------|----------|----------------|
| 400 | Bad request — validation failed | `{"success": false, "error": { "code": "VALIDATION_ERROR", "message": "Details", "details": [...] }}` |
| 401 | Unauthenticated — missing/invalid token | `{"success": false, "error": { "code": "AUTHENTICATION_ERROR", "message": "Missing or invalid token" }}` |
| 403 | Forbidden — insufficient permissions | `{"success": false, "error": { "code": "AUTHORIZATION_ERROR", "message": "You don't have permission" }}` |
| 404 | Not found — resource doesn't exist | `{"success": false, "error": { "code": "NOT_FOUND", "message": "Resource not found" }}` |
| 409 | Conflict — duplicate enrollment, etc. | `{"success": false, "error": { "code": "CONFLICT", "message": "Already enrolled in this course" }}` |
| 422 | Unprocessable entity — business rule violation | `{"success": false, "error": { "code": "VALIDATION_ERROR", "message": "Details" }}` |
| 500 | Internal server error | `{"success": false, "error": { "code": "INTERNAL_ERROR", "message": "Unexpected error" }}` |

## 7. Success Response Format (All Actions)

All successful actions return consistent format:

```json
{
  "success": true,
  "data": { ... }, // action-specific data
  "meta": {
    "cache_invalidated": [...], // cache keys invalidated
    "new_ttls": { // new TTLs set for invalidated keys
      "cache:key1": "30min",
      "cache:key2": "2hr"
    }
  },
  "timestamp": "2024-01-15T10:50:00Z",
  "correlation_id": "req-abc123def456" // for tracing
}
```

## 8. Action Performance Standards

| Metric | Target | Measurement |
|--------|--------|-------------|
| Action response time (p95) | < 200ms | Gin middleware timing |
| Cache invalidation latency | < 50ms | Redis DEL operation |
| GPA recalculation (on grade enter) | < 2s | Background job or sync |
| Report generation (small) | < 30s | Queue-based processing |
| Report generation (large) | < 2min | Batch processing with progress |

---