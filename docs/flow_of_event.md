# Flow of Event — Ministry of Education System

## 1. Purpose & Scope

This document defines the **Flow of Event** for the Ministry of Education integrated system. Events represent significant state changes or occurrences within the system that trigger actions, updates, and notifications. Unlike the Flow of Action which focuses on individual API requests, the Flow of Event covers background processes, async operations, and system-wide event propagation.

## 2. Event Categories

Events are categorized by their source and nature:

| Category | Event Types | Description |
|----------|-------------|-------------|
| **STUDENT** | ENROLLMENT_CREATED, ENROLLMENT_DROPPED, GRADE_ENTERED, ATTENDANCE_MARKED, PROFILE_UPDATED, LOGIN_SUCCESS, LOGOUT | Student lifecycle events |
| **COURSE** | COURSE_CREATED, COURSE_UPDATED, COURSE_DELETED, PREREQUISITE_CHANGED, CAPACITY_CHANGED | Course management events |
| **ATTENDANCE** | ATTENDANCE_MARKED, ATTENDANCE_RESET, ABSENCE_ALERT, THRESHOLD_EXCEEDED | Attendance tracking events |
| **GRADE** | GRADE_ENTERED, GRADE_UPDATED, GPA_CALCULATED, GRADE_RELEASED | Grade and GPA events |
| **REPORT** | REPORT_GENERATED, REPORT_DOWNLOADED, EXPORT_COMPLETE | Report generation events |
| **SYSTEM** | AUTH_LOGIN, AUTH_LOGOUT, CACHE_INVALIDATED, SCHEDULE_UPDATE | System-level events |

## 3. Event Propagation Flow

### 3.1 Event Production (When Events Are Generated)

```
┌──────────────────────────────────────────────────────────────┐
│                    EVENT PRODUCTION                          │
│  (Occurs within use case execution)                          │
└───────▲────────────────────────────▲─────────────────────────┘
        │                          │
        │                          │
        │  1. Event Trigger        │
        │     (within Go function)│
        │                          │
        ├────► [Generate Event]  │
        │      • Event struct     │
        │      • Topic/Channel    │
        │      • Payload          │
        │                          │
        │  2. Publish Event       │
        │    • Redis Pub/Sub      │
        │    • Kafka topic        │
        │    • RabbitMQ exchange  │
        │                          │
        │  3. Acknowledge         │
        │    • Confirm publication│
        │    • Log event ID       │
        ▼                          │
┌──────────────────────────────────────────────────────────────┐
│                    EVENT DISTRIBUTION                        │
│  (Fans out to interested consumers)                         │
└───────▲────────────────────────────▲─────────────────────────┘
        │                          │
        │                          │
        │  1. Subscribers          │
        │    • Auth service        │
        │    • Notification service│
        │    • Reporting service  │
        │    • Monitoring         │
        │                          │
        ├────► [Redis Pub/Sub]   │
        │    • Pattern: events:*  │
        │    • QoS: at least once │
        │                          │
        │  2. Process Events       │
        │    • Deserialize payload│
        │    • Execute handlers   │
        │    • Update state       │
        │                          │
        │  3. Acknowledge         │
        │    • Confirm processing│
        │    • Update offsets     │
        ▼                          │
┌──────────────────────────────────────────────────────────────┐
│                    EVENT CONSUMPTION                         │
│  (Handlers execute side effects)                            │
└──────────────────────────────────────────────────────────────┘
```

## 4. Event Definitions & Schemas

### 4.1 Student Events

```go
// events/student.go
package events

import (
    "time"
    "github.com/google/uuid"
)

// Base event structure
type BaseEvent struct {
    EventID   string    `json:"event_id"`
    OccurredAt time.Time `json:"occurred_at"`
    Source    string    `json:"source"`     // e.g., "api", "worker", "cron"
    Version   string    `json:"version"`    // "v1"
    CorrelationID string `json:"correlation_id"`
}

// ENROLLMENT_CREATED event
type EnrollmentCreated struct {
    BaseEvent
    EventType   string `json:"event_type"`   // "ENROLLMENT_CREATED"
    EnrollmentID string `json:"enrollment_id"`
    StudentID   string `json:"student_id"`
    CourseCode  string `json:"course_code"`
    Semester    string `json:"semester"`
    EnrolledAt  string `json:"enrolled_at"`
}

// ENROLLMENT_DROPPED event
type EnrollmentDropped struct {
    BaseEvent
    EventType    string `json:"event_type"`     // "ENROLLMENT_DROPPED"
    EnrollmentID string `json:"enrollment_id"`
    StudentID    string `json:"student_id"`
    CourseCode   string `json:"course_code"`
    DroppedAt    string `json:"dropped_at"`
    Reason       string `json:"reason"` // optional
}
```

### 4.2 Attendance Events

```go
// events/attendance.go
type AttendanceMarked struct {
    BaseEvent
    EventType      string `json:"event_type"`     // "ATTENDANCE_MARKED"
    AttendanceID   string `json:"attendance_id"`
    StudentID      string `json:"student_id"`
    CourseCode     string `json:"course_code"`
    SessionDate    string `json:"session_date"`   // YYYY-MM-DD
    Status         string `json:"status"`         // "Present", "Absent", "Excused"
    MarkedBy       string `json:"marked_by"`      // "teacher_id" or "admin"
    AbsenceCount   int    `json:"absence_count"`  // cumulative for student
    TriggerAlert   bool   `json:"trigger_alert"`  // if threshold exceeded
}

type AbsenceThresholdExceeded struct {
    BaseEvent
    EventType      string `json:"event_type"`     // "ABSENCE_THRESHOLD_EXCEEDED"
    StudentID      string `json:"student_id"`
    CourseCode     string `json:"course_code"`
    CurrentAbsences int    `json:"current_absences"`
    Threshold      int    `json:"threshold"`      // e.g., 5 absences
    DaysPeriod     string `json:"days_period"`    // "week", "month", "term"
    AlertSent      bool   `json:"alert_sent"`     // whether notification was sent
}
```

### 4.3 Grade Events

```go
// events/grade.go
type GradeEntered struct {
    BaseEvent
    EventType   string `json:"event_type"`   // "GRADE_ENTERED"
    GradeID     string `json:"grade_id"`
    EnrollmentID string `json:"enrollment_id"`
    StudentID   string `json:"student_id"`
    CourseCode  string `json:"course_code"`
    Grade       string `json:"grade"`        // "A", "B", "C", etc.
    Points      int    `json:"points"`       // numeric 0-100 or grade points
    EnteredBy   string `json:"entered_by"`   // "teacher_id" or "admin"
    EnteredAt   string `json:"entered_at"`
    UpdatedGPA  string `json:"updated_gpa"`  // new GPA value
}

type GPACalculated struct {
    BaseEvent
    EventType   string `json:"event_type"`   // "GPA_CALCULATED"
    StudentID   string `json:"student_id"`
    OldGPA      string `json:"old_gpa"`      // previous GPA
    NewGPA      string `json:"new_gpa"`      // recalculated GPA
    CreditHours int    `json:"credit_hours"` // total credit hours considered
    Term        string `json:"term"`         // "Fall 2024", "Spring 2025"
}
```

### 4.4 System Events

```go
// events/system.go
type CacheInvalidated struct {
    BaseEvent
    EventType      string `json:"event_type"`     // "CACHE_INVALIDATED"
    CacheKey       string `json:"cache_key"`      // e.g., "student:123"
    Reason         string `json:"reason"`         // "update", "expire", "manual"
    PreviousTTL    string `json:"previous_ttl"`   // e.g., "30min"
    NewTTL         string `json:"new_ttl"`        // e.g., "30min"
}

type ScheduleUpdated struct {
    BaseEvent
    EventType     string `json:"event_type"`     // "SCHEDULE_UPDATED"
    ScheduleID    string `json:"schedule_id"`
    AffectedDays  string `json:"affected_days"`  // e.g., "Monday-Wednesday"
    UpdatedBy     string `json:"updated_by"`
    UpdatedAt     string `json:"updated_at"`
}
```

## 5. Event Bus Architecture

### 5.1 Redis Pub/Sub Implementation

```go
// events/publisher.go
package events

import (
    "context"
    "encoding/json"
    "time"
    "ministry-education/config"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = config.NewRedisClient()

// Publish event to Redis channel
func Publish(eventType string, payload interface{}) error {
    // Create event struct
    var event BaseEvent
    
    // Set correlation ID from context or generate new
    event.CorrelationID = generateCorrelationID()
    event.EventID = uuid.New().String()
    event.OccurredAt = time.Now().UTC()
    event.Source = "api-gateway"
    event.Version = "v1"
    event.EventType = eventType
    
    // Serialize payload
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    
    // Publish to channel
    channel := fmt.Sprintf("events.%s", eventType)
    err = rdb.Publish(ctx, channel, payloadBytes).Err()
    
    if err != nil {
        logger.Error("Failed to publish event", "error", err, "channel", channel)
    }
    
    return err
}

// Example usage after enrollment creation
func PublishEnrollmentCreated(enrollmentID, studentID, courseCode, semester string) {
    event := EnrollmentCreated{
        EnrollmentID: enrollmentID,
        StudentID:    studentID,
        CourseCode:   courseCode,
        Semester:     semester,
        BaseEvent:    BaseEvent{EventType: "ENROLLMENT_CREATED"},
    }
    
    _ = Publish("ENROLLMENT_CREATED", event)
}
```

### 5.2 Event Subscriber

```go
// events/subscriber.go
package events

import (
    "context"
    "encoding/json"
    "fmt"
    "ministry-education/services/notification"
    "ministry-education/services/audit"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = config.NewRedisClient()

// Subscribe to all events
func Subscribe() {
    pubsub := rdb.Subscribe(ctx, "events.*")
    defer pubsub.Close(ctx)
    
    count := 0
    for msg := range pubsub.Channel() {
        count++
        
        // Parse event
        var base BaseEvent
        if err := json.Unmarshal([]byte(msg.Payload), &base); err != nil {
            logger.Error("Failed to parse event", "error", err, "channel", msg.Channel)
            continue
        }
        
        // Route to appropriate handler
        switch base.EventType {
        case "ENROLLMENT_CREATED":
            handleEnrollmentCreated(msg.Payload)
        case "ATTENDANCE_MARKED":
            handleAttendanceMarked(msg.Payload)
        case "GRADE_ENTERED":
            handleGradeEntered(msg.Payload)
        case "ABSENCE_THRESHOLD_EXCEEDED":
            handleAbsenceThresholdExceeded(msg.Payload)
        case "CACHE_INVALIDATED":
            handleCacheInvalidated(msg.Payload)
        default:
            logger.Warn("Unhandled event type", "event_type", base.EventType, "channel", msg.Channel)
        }
        
        // Acknowledge (Redis Pub/Sub doesn't require explicit ack, but log processing)
        if count%100 == 0 {
            logger.Info("Event processing", "count", count, "channel", msg.Channel)
        }
    }
}

// Handle Enrollment Created event
func handleEnrollmentCreated(payload []byte) {
    var event EnrollmentCreated
    if err := json.Unmarshal(payload, &event); err != nil {
        logger.Error("Failed to unmarshal EnrollmentCreated", "error", err)
        return
    }
    
    // Trigger secondary actions
    // 1. Send welcome email to student
    _ = notification.SendEnrollmentWelcome(
        studentID: event.StudentID,
        courseCode: event.CourseCode,
        semester: event.Semester,
    )
    
    // 2. Update analytics
    _ = analytics.TrackEnrollment(
        studentID: event.StudentID,
        courseCode: event.CourseCode,
    )
    
    // 3. Log audit entry (secondary)
    _ = audit.LogEvent("ENROLLMENT_CREATED_SECONDARY", event)
}
```

## 6. Event Handlers & Side Effects

### 6.1 Attendance Event Handlers

When `ATTENDANCE_MARKED` event is published:

```
┌──────────────────────────────────────────────────────────────┐
│          ATTENDANCE_MARKED EVENT PROPAGATION                 │
├──────────────────────────────────────────────────────────────┤
│  1. Publish:  Redis → events.ATTENDANCE_MARKED               │
│       Payload: {student_id, course_code, status, ...}        │
│──────────────────────────────────────────────────────────────│
│  2. Subscriber 1: Notification Service                       │
│    • If absence count > threshold → Send SMS/Email to parent  │
│    • Format: "Student X was absent from Y class on Z date"   │
│──────────────────────────────────────────────────────────────│
│  3. Subscriber 2: Analytics Service                          │
│    • Increment absence counter for student                   │
│    • Update attendance percentage                            │
│    • Store in Redis: cache:attendance:{date}                 │
│    • Set TTL: 5 minutes                                      │
│──────────────────────────────────────────────────────────────│
│  4. Subscriber 3: Reporting Service                          │
│    • Update daily attendance aggregates                      │
│    • Prepare summary statistics                              │
│──────────────────────────────────────────────────────────────│
│  5. Subscriber 4: Audit Service                              │
│    • Write audit entry: ATTENDANCE_MARKED                    │
│    • Store in PostgreSQL audit_logs table                    │
└──────────────────────────────────────────────────────────────┘
```

### 6.2 Absence Threshold Handler

```go
// handlers/absence.go
package handlers

import (
    "context"
    "encoding/json"
    "ministry-education/config"
    "github.com/redis/go-redis/v9"
    "ministry-education/models"
    "ministry-education/services/notification"
)

var ctx = context.Background()
var rdb = config.NewRedisClient()

// Handle ABSENCE_THRESHOLD_EXCEEDED event
func HandleAbsenceThresholdExceeded(payload []byte) {
    var event models.AbsenceThresholdExceeded
    if err := json.Unmarshal(payload, &event); err != nil {
        logger.Error("Failed to unmarshal event", "error", err)
        return
    }

    alertKey := fmt.Sprintf("alert_sent:%s:%d", event.StudentID, event.Threshold)
    exists, _ := rdb.Exists(ctx, alertKey).Result()
    if exists > 0 {
        return // Already notified within 24h
    }

    _ = notification.SendAbsenceAlert(
        event.StudentID,
        event.CourseCode,
        event.CurrentAbsences,
        event.Threshold,
        event.DaysPeriod,
    )

    // Set alert key with 24h TTL directly
    _ = rdb.Set(ctx, alertKey, "1", 24*time.Hour).Err()
}
}
```

### 6.3 GPA Recalculation Handler

```go
// handlers/gpa.go
package handlers

import (
    "context"
    "encoding/json"
    "ministry-education/config"
    "github.com/redis/go-redis/v9"
    "ministry-education/services/prisma"
)

var ctx = context.Background()
var rdb = config.NewRedisClient()

// Handle GRADE_ENTERED event
func HandleGradeEntered(payload []byte) {
    var event models.GradeEntered
    // unmarshal payload
    
    // Recalculate GPA for student
    gpa, err := prisma.CalculateGPA(event.StudentID)
    if err != nil {
        logger.Error("GPA calculation failed", "error", err, "student_id", event.StudentID)
        return
    }
    
    // Store in cache with TTL
    cacheKey := fmt.Sprintf("gpa:%s", event.StudentID)
    gpaData := map[string]string{
        "gpa":      gpa,
        "updated_at": event.EnteredAt,
        "source_event": event.EventID,
    }
    gpaJSON, _ := json.Marshal(gpaData)
    
    err = rdb.Set(ctx, cacheKey, gpaJSON, config.GPATTL).Err()
    if err != nil {
        logger.Error("Failed to store GPA in cache", "error", err)
    }
    
    // Publish GPA Calculated event (for downstream consumers)
    gpaEvent := models.GPA_Calculated{
        BaseEvent: models.BaseEvent{
            EventID:   uuid.New().String(),
            EventType: "GPA_CALCULATED",
            OccurredAt: time.Now().UTC(),
            Source:    "gpa-calculator",
        },
        StudentID:   event.StudentID,
        OldGPA:      "", // could fetch previous GPA
        NewGPA:      gpa,
        CreditHours: 30, // example
        Term:        "Fall 2024",
    }
    
    _ = Publish("GPA_CALCULATED", gpaEvent)
}
```

## 7. Event Retention & Durability

| Event Type | Retention Period | Storage | Reason |
|------------|-----------------|---------|--------|
| Critical (ENROLLMENT_, GRADE_) | 10 years | PostgreSQL audit_logs | Legal compliance, GDPR |
| Attendance events | 3 years | PostgreSQL | Academic records |
| Notification events | 1 year | Redis + Elasticsearch | Troubleshooting |
| System events (cache, schedule) | 90 days | Redis + Elasticsearch | Monitoring |
| Ephemeral events | 24 hours | Redis only | Performance |

### 7.1 Event Persistence to PostgreSQL

```go
// events/persistence.go
package events

import (
    "context"
    "encoding/json"
    "time"
    "ministry-education/models"
    "ministry-education/db"
)

func PersistEvent(event models.BaseEvent, payload []byte) error {
    ctx := context.Background()
    
    // Serialize full event for storage
    fullEvent := models.PersistedEvent{
        EventID:      event.EventID,
        OccurredAt:   event.OccurredAt,
        EventType:    event.EventType,
        Source:       event.Source,
        CorrelationID: event.CorrelationID,
        Payload:      string(payload),
        CreatedAt:    time.Now(),
    }
    
    // Insert into PostgreSQL
    _, err := db.Pg.Exec(ctx,
        `INSERT INTO event_log (event_id, occurred_at, event_type, source, correlation_id, payload, created_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`,
        fullEvent.EventID, fullEvent.OccurredAt, fullEvent.EventType,
        fullEvent.Source, fullEvent.CorrelationID, fullEvent.Payload, fullEvent.CreatedAt,
    )
    
    return err
}
```

### 7.2 Elasticsearch for Search

```go
// events/elasticsearch.go
package events

import (
    "context"
    "time"
    "ministry-education/db"
)

func IndexEvent(event models.PersistedEvent) error {
    ctx := context.Background()
    
    // Index in Elasticsearch for search
    _, err := db.ES.Index().
        Index("events").
        Id(event.EventID).
        BodyJson(map[string]interface{}{
            "event_type":   event.EventType,
            "occurred_at":  event.OccurredAt,
            "source":       event.Source,
            "student_id":   extractStudentID(event.EventType, event.Payload),
            "course_code":  extractCourseCode(event.EventType, event.Payload),
            "timestamp":    time.Now(),
        }).
        Refresh("wait_for").
        Execute(ctx)
    
    return err
}
```

## 8. Event Flow Diagrams

### 8.1 Student Enrollment Event Flow

```
┌─────────────────┐     ┌─────────────────────────────┐
│  Student Action │     │  API Endpoint: POST /enroll   │
│  (Enroll in      │     │  (use case execution)         │
│   course)         │     └───────▲─────────────────────┘
└───────▲───────────┘          │
        │                      │
        │ 1. Create enrollment │
        │    in PostgreSQL     │
        │                      │
        ├────► [Prisma Save]  │
        │          │           │
        │          ├─ Success ─┤
        │          │           │  2. Publish ENROLLMENT_CREATED event
        │          │           │       → Redis Pub/Sub → events.ENROLLMENT_CREATED
        │          │           │            │
        │          │           ├────► [Notification Handler] → Send welcome email
        │          │           ├────► [Analytics Handler] → Track enrollment stats
        │          │           ├────► [Audit Handler] → Write audit log
        │          │           └────► [Cache Handler] → Invalidate caches
        │          │
        │          └─ Failure ─┤
        │                      │  Return error to client
        ▼                      ▼
┌─────────────────┐  ┌─────────────────────────────┐
│  Response to    │  │  Event Propagation          │
│  Client         │  │  (background/async)         │
└─────────────────┘  └─────────────────────────────┘
```

### 8.2 Absence Threshold Event Flow

```
┌─────────────────┐     ┌─────────────────────────────┐
│  Teacher Action │     │  API Endpoint: POST /attendance │
│  (Mark attendance)│     │  (use case execution)          │
└───────▲───────────┘          └───────▲─────────────────────┘
        │                      │
        │ 1. Mark attendance   │
        │    in PostgreSQL     │
        │                      │
        ├────► [Prisma Save]  │
        │          │           │
        │          ├─ Success ─┤
        │          │           │  1. Calculate absence count
        │          │           │  2. Check threshold (e.g., > 5/week)
        │          │           │  3. If exceeded → Publish ABSENCE_THRESHOLD_EXCEEDED
        │          │           │       → Redis → events.ABSENCE_THRESHOLD_EXCEEDED
        │          │           │            │
        │          │           ├────► [Notification Handler] → Send alert to parent
        │          │           ├────► [Analytics Handler] → Update stats
        │          │           └────► [Cache Handler] → Refresh attendance cache
        │          │
        │          └─ Failure ─┤
        │                      │  Return error
        ▼                      ▼
┌─────────────────┐  ┌─────────────────────────────┐
│  Response to    │  │  Event Propagation          │
│  Client         │  │  (async background processes)│
└─────────────────┘  └─────────────────────────────┘
```

## 9. Event Ordering & Guarantees

### 9.1 At-Least-Once Delivery

```
- Redis Pub/Sub provides at-least-once delivery
- Events may be delivered multiple times
- Handlers must be idempotent
- Use deduplication keys (event_id) for processing

Example deduplication:
  alertKey := fmt.Sprintf("alert_sent:%s:%d", studentID, threshold)
  if rdb.Exists(ctx, alertKey).Val() > 0 {
      // Already processed, skip
  }
```

### 9.2 Event Ordering Guarantees

```
- Within same Redis channel: messages are ordered
- Across channels: no guaranteed ordering
- Use sequence numbers for critical ordering

Example:
type BaseEvent struct {
    EventID      string    `json:"event_id"`
    OccurredAt   time.Time `json:"occurred_at"`
    SequenceID   uint64    `json:"sequence_id"` // monotonically increasing
    ...
}
```

### 9.3 Dead Letter Queue

```
- Events that fail processing N times → moved to DLQ
- DLQ stored in Redis queue: "events.dlq"
- Manual review and reprocessing possible
- Configurable retry count (default: 3)
```

## 10. Monitoring & Alerting

### 9.1 Event Metrics

| Metric | Description | Normal Range | Alert Threshold |
|--------|-------------|--------------|-----------------|
| events.published.total | Total events published per minute | varies by time of day | > 2x normal rate |
| events.consumed.total | Total events consumed per minute | varies | > 50% lagging |
| events.dlq.size | Events in dead letter queue | < 100 | > 1000 |
| events.lag.time | Time between publish and consume | < 5 seconds | > 30 seconds |
| redis.pubsub.channels | Active Pub/Sub channels | 20-30 | > 50 |

### 9.2 Alert Rules

```
# High event lag alert
alert "events_lag_high" {
    expr = (events_published_total - events_consumed_total) / events_published_total > 0.5
    for = 1m
    labels = { severity = "warning" }
    annotations = {
        summary = "Event consumption lagging"
        description = "Lag: {{ $value | humanize }} events behind"
    }
}

# DLQ size alert
alert "events_dlq_high" {
    expr = events_dlq_size > 1000
    for = 5m
    labels = { severity = "critical" }
    annotations = {
        summary = "Event dead letter queue growing"
        description = "DLQ size: {{ $value }} events"
    }
}
```

## 11. Event Versioning

| Version | Release Date | Changes |
|---------|-------------|---------|
| v1.0.0 | 2024-01-15 | Initial event definitions (ENROLLMENT_, ATTENDANCE_, GRADE_) |
| v1.1.0 | 2024-03-01 | Added GPA_CALCULATED, ScheduleUpdated events |
| v1.2.0 | 2024-06-01 | Added CacheInvalidated, improved deduplication |
| v2.0.0 | 2025-01-15 | Breaking changes: new event schema, migration guide required |

### 11.1 Migration Between Versions

```
- Maintain backward compatibility via event aliases
- Old subscribers can listen to both v1 and v2 topics
- Dual-publish during transition period (both v1 and v2)
- Deprecate old events after 3 months grace period

Example alias:
  events.ENROLLMENT_CREATED_v1 → events.ENROLLMENT_CREATED
  events.ENROLLMENT_CREATED_v2 → events.ENROLLMENT_CREATED (new schema)
```