# 📋 Use Case Scenarios — Ministry of Education System

## 1. Purpose & Scope

This document defines the comprehensive use case scenarios for the Ministry of Education integrated system, covering the full student lifecycle from enrollment to grading and reporting. Each use case describes actor interactions, preconditions, main success flows, and alternative paths.

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
7. System records attendance in PostgreSQL via Prisma
8. Redis cache updated: `cache:attendance:{date}` (TTL: 5 min)
9. Attendance count incremented for each student
10. If absent count exceeds threshold → Auto-trigger parent notification
11. Audit log: `ATTENDANCE_MARKED`

**Alternative Paths:**
- **Duplicate Entry:** System prevents marking same student twice → Shows error
- **Session Locked:** System shows "Attendance finalized for this date" → Only admin can reopen
- **Invalid Status:** System validates status is Present/Absent/Excused → Rejects other values

**Postconditions:**
- Attendance data persisted in DB
- Cache refreshed with fresh TTL
- Parent notifications triggered if thresholds met
- Teacher can view attendance summary

**Extensions:**
- Bulk attendance import (CSV)
- Late attendance marking (with admin approval)
- Excuse note attachment
- Substitute teacher attendance handling

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
5. System queries data from PostgreSQL via Prisma
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
3. If cache miss → Query PostgreSQL via Prisma
4. System displays:
   - Personal information (name, email, ID)
   - Enrolled courses list
   - Attendance summary
   - Current grades
   - Profile update options
5. Student can edit profile (name, contact info)
6. Changes validated with Zod schema
7. Updated data written to DB
8. Cache invalidated and refreshed
9. Audit log: `PROFILE_VIEWED` or `PROFILE_UPDATED`

**Alternative Paths:**
- **Profile Not Found:** System shows "Profile not found" → Contact admin
- **Old Cache:** System refreshes from DB if stale → Ensures data freshness
- **Validation Failed:** Zod rejects invalid input → Shows error messages to student

**Postconditions:**
- Student sees current profile data
- Update changes persisted
- Access logged for audit

**Student-Editable Fields:**
- Contact email (with verification)
- Phone number
- Emergency contact information
- Profile picture (if supported)

---

### UC-06: Authentication & Login

**Primary Actor:** ALL (STUDENT, TEACHER, ADMIN, PARENT)  
**Preconditions:** None (public login endpoint)

**Main Success Scenario:**
1. User navigates to login page
2. User enters credentials (email/username + password)
3. System validates credentials against PostgreSQL
4. If valid → System generates JWT token (2hr expiry)
5. Token stored in HttpOnly cookie or localStorage
6. System returns auth success response
7. Client stores token and redirects user to appropriate dashboard
8. Redis session validated (if token revocation needed)
9. Audit log: `LOGIN_SUCCESS` with IP, user-agent, correlation ID

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