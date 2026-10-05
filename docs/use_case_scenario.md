# 📋 Use Case Scenarios — Ministry of Education System (Updated Architecture)

## 1. Purpose & Scope

This document defines the comprehensive use case scenarios for the Ministry of Education integrated system, covering the full student lifecycle from enrollment to grading and reporting. Each use case describes actor interactions, preconditions, main success flows, and alternative paths. **All database ORM references have been updated to reflect native Go ORM implementations.**

## 2. Actor Definitions

| Actor | Role | Permissions | Constraints |
|-------|------|-------------|-------------|
| **STUDENT** | enrolled student | View own data, enroll in courses, view grades, update profile | Can only access own records |
| **TEACHER** | faculty member | Take attendance, enter grades, view class roster | Assigned courses only |
| **ADMIN** | system administrator | Manage users, courses, reports, system config | Full system access |
| **PARENT** | parent/guardian | View child's attendance and grades | Linked to specific student(s) |
| **SUPER_ADMIN** | top-level administrator | All permissions, user management, licensing | Organization-wide authority |

## 3. Use Case Catalog

### UC-01: Student Enrollment in Course

**Primary Actor:** STUDENT  
**Preconditions:** 
- Student is authenticated with valid JWT token
- Student has not exceeded maximum course load
- Course is active and has available capacity
- Prerequisites (if any) are satisfied

**Main Success Scenario (Basic Flow):**
1. Student logs in and navigates to Course Catalog
2. Student searches/selects desired course
3. System validates prerequisites automatically
4. System checks course capacity/enrollment limits
5. If valid → System creates enrollment record
6. Student receives enrollment confirmation
7. Cache invalidated: `cache:course:{code}` (TTL reset)
8. Audit log: `ENROLLMENT_CREATED`

**Alternative Paths:**
- **Prerequisite Not Met:** System shows required prerequisites not completed → Student must complete those first
- **Capacity Full:** System shows "Course at full capacity" → Student can join waitlist or select alternative
- **Already Enrolled:** System shows "Already enrolled in this course" → No action taken
- **Schedule Conflict:** System detects timetable conflict → Student can resolve or choose different course

**Postconditions:**
- Student's enrollment count incremented
- Course student count incremented
- Cache refreshed with TTL: 30 minutes
- Audit entry created with full trail

**Extensions:**
- Parent can enroll student (with STUDENT authentication)
- Multiple course enrollment in single session
- Waitlist management when at capacity

---

### UC-02: Teacher Takes Attendance

**Primary Actor:** TEACHER  
**Preconditions:**
- Teacher is authenticated
- Teacher is assigned to the course/session
- Session date is within active term
- Attendance has not been finalized for this session

**Main Success Scenario:**
1. Teacher logs into system and navigates to Attendance module
2. Teacher selects class and session date
3. System displays student roster for that class
4. Teacher marks each student as Present/Absent/Excused
5. System validates no duplicate attendance for same date
6. Teacher submits attendance
7. **System records attendance in PostgreSQL via native Go ORM (ent/pgx/sqlx)** — *formerly Prisma*
8. **Cache updated: `cache:attendance:{school_id}:{class_id}:{date}`** (TTL: 5 min) — *formerly `cache:attendance:{date}` to prevent key contention during bulk registration*
9. Attendance count incremented for each student
10. If absent count exceeds threshold → Auto-trigger parent notification via Event Queue
11. Audit log: `ATTENDANCE_MARKED`

**Alternative Paths:**
- **Duplicate Entry:** System prevents marking same student twice → Shows error
- **Session Locked:** System shows "Attendance finalized for this date" → Only admin can reopen
- **Invalid Status:** System validates status is Present/Absent/Excused → Rejects other values

**Postconditions:**
- Attendance data persisted in DB
- Cache refreshed with fresh TTL
- Parent notifications triggered if thresholds met (via Event Queue, not synchronous HTTP)
- Teacher can view attendance summary

**Extensions:**
- Bulk attendance import (CSV)
- Late attendance marking (with admin approval)
- Excuse note attachment
- Substitute teacher attendance handling

**Async Notification Trigger (Postcondition):**
> *Parent absence notifications are NOT sent synchronously within the HTTP request cycle. Instead, an `ABSENCE_THRESHOLD_EXCEEDED` event is published to a Message Queue (Redis Streams or RabbitMQ). Background Worker Pools consume this event and deliver SMS/Email notifications asynchronously, ensuring the API response remains fast and the system remains responsive under load.*

---

### UC-03: Parent Views Child's Academic Progress

**Primary Actor:** PARENT  
**Preconditions:**
- Parent is authenticated with valid JWT
- Parent is linked to student as guardian
- Student records exist in system

**Main Success Scenario:**
1. Parent logs into parent portal
2. Parent selects child's profile from linked students
3. System retrieves student data from cache or DB
4. System displays:
   - Current grades per course
   - Attendance percentage per course
   - Overall GPA
   - Progress trends (if available)
5. System formats data for parent-friendly display
6. Cache hit: `cache:student:{id}` (TTL: 30 min) or DB query if miss
7. Audit log: `PARENT_VIEW_ACCESS`

**Alternative Paths:**
- **No Access:** System shows "You don't have access to this student's records" → Denies access
- **Multiple Children:** System shows selector for all linked students → Parent can switch
- **Data Not Available:** System shows "Grades not yet entered" → Inform parent of delay

**Postconditions:**
- Parent views student data securely
- Access logged for compliance
- No modification to student data

**Extensions:**
- Download progress report as PDF
- Email summary to parent
- Set up attendance alerts/thresholds
- View historical data (previous terms)

---

### UC-04: Admin Generates Institutional Report

**Primary Actor:** ADMIN / SUPER_ADMIN  
**Preconditions:**
- Admin has `report:generate` permission
- System has data for requested time period
- Report type is supported

**Main Success Scenario:**
1. Admin navigates to Reports section
2. Admin selects report type (e.g., "Student Performance", "Attendance Summary", "Course Analysis")
3. Admin specifies date range (start/end date)
4. System validates date range is within active terms
5. System queries data from PostgreSQL via native Go ORM (ent/pgx/sqlx) — *formerly Prisma*
6. System processes data and generates report format
7. Report generated in requested format (PDF, CSV, Excel)
8. System provides download link
9. Audit log: `REPORT_GENERATED` with parameters

**Alternative Paths:**
- **Insufficient Data:** System shows "Insufficient data for selected period" → Suggest alternative ranges
- **Permission Denied:** Admin without proper role → Access denied with explanation
- **Large Dataset:** System processes in batches → May take longer, provide progress indicator

**Postconditions:**
- Report file generated and stored
- Download link provided (time-limited, e.g., 24 hours)
- Access audit trail maintained
- No data modification

**Report Types Available:**
- Student Performance Summary (per grade/section)
- Attendance Analysis (by class, by date range)
- Course Effectiveness (pass rates, grade distribution)
- Teacher Performance (student feedback, attendance correlation)
- Institutional KPIs (promotion rates, dropout indicators)

---

### UC-05: Student Views Personal Profile

**Primary Actor:** STUDENT  
**Preconditions:**
- Student authenticated with valid JWT
- Student profile exists in system

**Main Success Scenario:**
1. Student logs into student dashboard
2. System retrieves student profile from cache (`cache:student:{id}`) with TTL: 30 min
3. If cache miss → Query PostgreSQL via native Go ORM (ent/pgx/sqlx) — *formerly Prisma*
4. System displays:
   - Personal information (name, email, ID)
   - Enrolled courses list
   - Attendance summary
   - Current grades
   - Profile update options
5. Student can edit profile (name, contact info)
6. Changes validated with go-playground/validator schema — *formerly Zod; Zod now used exclusively in Frontend*
7. Updated data written to DB
8. Cache invalidated and refreshed
9. Audit log: `PROFILE_VIEWED` or `PROFILE_UPDATED`

**Alternative Paths:**
- **Profile Not Found:** System shows "Profile not found" → Contact admin
- **Old Cache:** System refreshes from DB if stale → Ensures data freshness
- **Validation Failed:** go-playground/validator rejects invalid input → Shows error messages to student

**Postconditions:**
- Student sees current profile data
- Update changes persisted
- Access logged for audit

**Student-Editable Fields:**
- Contact email (with verification)
- Phone number
- Emergency contact information
- Profile picture (if supported)

**Backend Validation (UC-05 Note):**
> *All student profile update validations are now performed using go-playground/validator schemas within the Go Backend. The Zod schema is reserved exclusively for React Frontend form validation. This separation ensures consistent validation logic across the system and leverages Go's type safety for backend operations.*

---

### UC-06: Authentication & Login

**Primary Actor:** ALL (STUDENT, TEACHER, ADMIN, PARENT)  
**Preconditions:** None (public login endpoint)

**Main Success Scenario:**
1. User navigates to login page
2. User enters credentials (email/username + password)
3. System validates credentials against PostgreSQL
4. If valid → System generates JWT token (2hr expiry for Access Token)
5. **Refresh Token generated and stored in Redis with 24-hour TTL** — *supporting instant Revocation*
6. Token stored in HttpOnly cookie or localStorage
7. System returns auth success response
8. Client stores token and redirects user to appropriate dashboard
9. Redis session validated (if token revocation needed)
10. Audit log: `LOGIN_SUCCESS` with IP, user-agent, correlation ID

**JWT Details (Auth & Tokens - UC-06):**
> *The system implements a Dual-Token Strategy:*
> - **Access Token:** Short-lived (15-60 minutes expiry), used for API authentication per request
> - **Refresh Token:** Long-lived (24 hours), stored in Redis with key format `refresh:{token_id}`, supports immediate Revocation via Redis DEL operation
> - **Access Flow:** When Access Token expires, client uses Refresh Token to obtain new Access Token without re-entering credentials
> - **Revocation:** Any admin can invalidate a Refresh Token by deleting its Redis key: `DELETE redis:refresh:{token_id}`, immediately terminating all active sessions for that token

**Alternative Paths:**
- **Invalid Credentials:** System shows "Invalid email or password" → Do NOT reveal which part was wrong
- **Account Locked:** After 5 failed attempts → Lock account 30 min, alert security
- **2FA Required:** If enabled → Send OTP to registered email/phone
- **Password Expired:** If password age > 90 days → Force password change

**Postconditions:**
- User authenticated and authorized
- Session established
- Login event logged for security audit
- Failed attempts logged for brute-force protection

**Security Measures:**
- Passwords hashed with bcrypt cost factor 12
- Rate limiting: max 5 login attempts per 15 min per IP
- Account lockout after 5 failed attempts (30 min)
- All passwords minimum 8 characters, complexity requirements
- GDPR-compliant data handling
- Password reset flow with token verification

---

### UC-07: Student Course Drop/Withdrawal

**Primary Actor:** STUDENT  
**Preconditions:**
- Student is authenticated
- Student is currently enrolled in the course
- Drop period is active (within academic calendar)
- No pending financial obligations for the course

**Main Success Scenario:**
1. Student navigates to enrolled courses
2. Student selects course to drop/withdraw from
3. System confirms enrollment exists and checks drop deadline
4. If within drop period → System soft-deletes enrollment record
5. System updates student's course count
6. Cache invalidated: `cache:student:{id}`, `cache:course:{code}`
7. Student receives withdrawal confirmation
8. Audit log: `ENROLLMENT_DROPPED`

**Alternative Paths:**
- **Past Drop Deadline:** System shows "Drop period has ended" → Inform of next opportunity (next term)
- **Course In Progress:** System warns "Grades already entered, dropping may affect transcript" → Provide warning
- **Financial Hold:** System shows "Cannot drop: outstanding fees" → Resolve financial matter first
- **Required Course:** System prevents dropping required courses → Show alternative paths

**Postconditions:**
- Enrollment status changed to "dropped/withdrawn"
- Student freed from course capacity
- Audit trail maintained for academic records
- Necessary notifications sent (admin, finance, student)

**Extensions:**
- Formal withdrawal process with documentation
- Grade impact assessment before dropping
- Transfer to alternative course
- International student visa implications (if applicable)

---
## 4. Technical Implementation Notes

### 5. Backend Validation & ORM Migration (UC-05)

**ORM Migration (UC-02, UC-04, UC-05):**
> *All database operations in use cases UC-02, UC-04, and UC-05 have been migrated from Prisma ORM to native Go ORM implementations. The following options are supported:*
> - **ent ORM:** Type-safe Go ORM with compile-time query generation
> - **pgx/sqlx:** Low-level PostgreSQL driver with enhanced features and direct SQL control
> 
> *Migration command example:*
> ```bash
> # Using golang-migrate for schema management
> npx migrate create -seq init_schema
> npx migrate up
> # OR with ent:
> ent generate
> ```

**Backend Validation (UC-05):**
> *Student profile update and all input validations in the Go Backend are now performed using `go-playground/validator` schemas. This library provides compile-time validated struct tags and is the de facto standard for Go validation.*
> 
> *Zod Schema Usage Note:*
> > *Zod is reserved exclusively for React Frontend form validation. The Backend no longer imports or uses Zod schemas, ensuring a clean separation of concerns and leveraging Go's type system for backend validation logic.*
> 
> *go-playground/validator Example:*
> ```go
> type StudentUpdate struct {
>     Name    string `validate:"min=2,max=100"`
>     Email   string `validate:"email,required"`
>     Phone   string `validate:"numeric"`
> }
> 
> validator := validator.New()
> if err := validator.Struct(studentUpdate); err != nil {
>     // Handle validation errors
> }
> ```

### 6. Cache Key Structure (UC-02)

**Attendance Cache Refactoring:**
> *The attendance cache key has been refactored from `cache:attendance:{date}` to `cache:attendance:{school_id}:{class_id}:{date}` to prevent key contention and system overload when registering attendance in bulk for entire classes.*
> 
> *Key Structure Breakdown:*
> - `cache:attendance:` — Fixed prefix
> - `{school_id}` — School identifier (e.g., "school-001")
> - `{class_id}` — Class/section identifier (e.g., "grade-5-a")
> - `{date}` — Attendance date (YYYY-MM-DD format)
> 
> *This restructuring ensures that even when all students in a class of 40+ are marked present/absent simultaneously, each key remains unique and cache operations remain performant.*

### 7. Event-Driven Notifications (UC-02)

**Async Notification Architecture:**
> *The system employs an event-driven architecture for parent notifications. When attendance thresholds are exceeded:*
> 
> 1. *Teacher marks attendance via HTTP POST /api/v1/attendance*
> 2. *System validates and persists attendance data to PostgreSQL*
> 3. *If absent count > threshold, system publishes `ABSENCE_THRESHOLD_EXCEEDED` event to Message Queue*
> 4. *Message Queue options: Redis Streams or RabbitMQ*
> 5. *Background Worker Pools consume events and deliver SMS/Email notifications asynchronously*
> 6. *API response returns immediately without waiting for notification delivery*
> 
> *This architecture ensures:*
> - *Fast API response times (no synchronous notification delays)*
> - *Reliable event delivery (at-least-once through Message Queue)*
> - *Scalable notification processing (worker pools can scale independently)*
> - *Resilient system operation (if notification service fails, events remain in queue)*

### 8. Database ORM Reference

**Supported ORM Implementations:**
- **ent ORM:** `import "ent.io/ent"`
- **pgx/sqlx:** `import "github.com/jackc/pgx/v5/stdlib"`
- **Previous (deprecated):** `import "github.com/prisma/prisma-client-go/v5"`

**Migration Path:**
> *Systems currently using Prisma are encouraged to migrate to either ent ORM or pgx/sqlx for improved performance, type safety, and reduced dependency complexity. Migration scripts and patterns are documented in the project's migration guide.*

---
## 5. Summary of Changes

| Area | Previous | Updated |
|------|----------|---------|
| **ORM** | Prisma ORM | ent ORM / pgx / sqlx (native Go) |
| **Attendance Cache Key** | `cache:attendance:{date}` | `cache:attendance:{school_id}:{class_id}:{date}` |
| **Backend Validation** | Zod Schema | go-playground/validator (Backend); Zod only for Frontend |
| **Auth Strategy** | Standard JWT | Dual-Token: Access (15-60min) + Refresh (24h in Redis) |
| **Notification Pattern** | Synchronous in HTTP Request | Async via Message Queue (Redis Streams/RabbitMQ) |
| **Postconditions** | Synchronous notifications | Async event publishing with Background Workers |

---
## 6. Quick Reference: Key Commands

```bash
# ORM Migration (ent)
ent generate

# ORM Migration (pgx/sqlx with golang-migrate)
npx migrate create -seq init_schema
npx migrate up

# Cache Key Format
cache:attendance:{school_id}:{class_id}:{date}

# Auth Tokens in Redis
access:{token_id}    TTL: 15-60 minutes
refresh:{token_id}   TTL: 24 hours

# Event Queue Setup
# Redis Streams
XADD events ABSENCE_THRESHOLD_EXCEEDED ...
# RabbitMQ
Publish to "absence.notifications" exchange
```