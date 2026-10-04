# 📝 Audit Specification — Ministry of Education System

## 1. Audit Purpose & Scope

This document defines the comprehensive audit framework for the Ministry of Education integrated system, ensuring compliance, data integrity, and regulatory adherence.

### 1.1 Audit Objectives
- **Data Integrity**: Ensure all student/teacher data is accurate and consistent
- **Compliance**: Meet educational regulations and data protection laws
- **Security**: Monitor access patterns and detect unauthorized activities
- **Accountability**: Track who did what and when
- **Quality**: Validate business rule enforcement

### 1.2 Regulatory Framework
- **GDPR**: Student and staff data privacy
- **FERPA**: Educational record protection
- **Local Education Laws**: Country-specific requirements
- **ISO 27001**: Information security management
- **SOC 2**: Service organization control requirements

## 2. Audit Log Structure

### 2.1 Required Audit Entity Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `id` | UUID | Primary key | `550e8400-e29b-41d4-a716-446655440000` |
| `timestamp` | datetime UTC | When the action occurred | `2024-01-15T10:30:00Z` |
| `user_id` | UUID | Who performed the action | `11111111-1111-41d4-a716-446655440000` |
| `user_role` | enum | Role of the user | `STUDENT`, `TEACHER`, `ADMIN` |
| `action` | enum | Type of action | `CREATE`, `UPDATE`, `DELETE`, `LOGIN` |
| `entity_type` | enum | Which entity affected | `STUDENT`, `COURSE`, `ENROLLMENT` |
| `entity_id` | UUID | Specific record ID | `22222222-2222-41d4-a716-446655440000` |
| `changes` | JSON | What changed (before/after) | `{ "grade": "A→B" }` |
| `ip_address` | string | Source IP | `192.168.1.100` |
| `user_agent` | string | Client information | `Mozilla/5.0` |
| `correlation_id` | string | Request tracing | `req-abc123` |
| `success` | boolean | Operation success status | `true` |
| `failure_reason` | string | Why it failed (if applicable) | `VALIDATION_ERROR` |

### 2.2 Audit Action Types

```
AUTHENTICATION:
├─ LOGIN_SUCCESS
├─ LOGIN_FAILURE
├─ LOGOUT
├─ TOKEN_REFRESH
└─ PASSWORD_RESET

ADMINISTRATIVE:
├─ CREATE_STUDENT
├─ UPDATE_STUDENT
├─ DELETE_STUDENT (soft-delete)
├─ CREATE_TEACHER
├─ UPDATE_TEACHER
├─ ASSIGN_ROLE
└─ CHANGE_PERMISSIONS

ACADEMIC:
├─ GRADE_ENTRY
├─ UPDATE_GRADE
├─ ATTENDANCE_MARK
├─ COURSE_ENROLLMENT
├─ COURSE_DROP
└─ REPORT_GENERATION

SYSTEM:
├─ API_CALL
├─ CACHE_ACCESS
├─ DATABASE_QUERY
├─ FILE_UPLOAD
└─ EXPORT_DATA
```

## 3. Audit Implementation

### 3.1 Audit Logger (Go Middleware)

```go
// Middleware to capture all API requests
func AuditMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Generate correlation ID
        correlationID := generateCorrelationID()
        c.Set("correlation_id", correlationID)
        
        // Start timing
        start := time.Now()
        
        // Process request
        c.Next()
        
        // Log audit entry
        duration := time.Since(start)
        logAuditEntry(&models.AuditLog{
            CorrelationID: correlationID,
            UserID:        getUserID(c),
            UserRole:      getUserRole(c),
            Action:        determineAction(c),
            EntityType:    determineEntityType(c),
            EntityID:      getEntityID(c),
            Changes:       getRequestChanges(c),
            IPAddress:     c.ClientIP(),
            UserAgent:     c.Request.UserAgent(),
            Duration:      duration,
            Success:       c.Writer.Status() < 400,
        })
    }
}
```

### 3.2 Audit Storage Strategy

| Storage | Purpose | TTL |
|---------|---------|-----|
| **PostgreSQL** | Long-term audit retention (7 years) | N/A (permanent) |
| **Redis** | Recent audits (last 30 days) | 30 days |
| **Elasticsearch** | Searchable audit logs | 1 year |
| **S3/MinIO** | Immutable backup | Indefinite |

### 3.3 Audit Query Examples

```sql
-- Get all admin actions in last 24 hours
SELECT * FROM audit_logs
WHERE user_role = 'ADMIN'
  AND timestamp >= NOW() - INTERVAL '24 hours'
  AND action LIKE 'ADMIN:%';

-- Get student data changes
SELECT * FROM audit_logs
WHERE entity_type = 'STUDENT'
  AND action IN ('UPDATE', 'DELETE')
  AND timestamp >= NOW() - INTERVAL '30 days';

-- Get failed authentication attempts
SELECT * FROM audit_logs
WHERE action = 'LOGIN_FAILURE'
  AND success = false
  AND timestamp >= NOW() - INTERVAL '1 hour';
```

## 4. Data Retention Policy

| Data Type | Retention Period | Storage | Reason |
|-----------|-----------------|---------|--------|
| **Audit Logs** | 7 years | PostgreSQL + S3 | Legal compliance, GDPR |
| **Session Data** | 30 days | Redis | Performance, cache |
| **Error Logs** | 1 year | Elasticsearch | Troubleshooting |
| **Access Logs** | 1 year | PostgreSQL | Security monitoring |
| **Backup Logs** | Indefinite | S3 | Disaster recovery |

## 5. Audit Reports

### 5.1 Standard Reports

| Report Type | Frequency | Audience | Content |
|-------------|-----------|----------|---------|
| **Daily Access Report** | Daily | IT Security | API calls, error rates, top users |
| **Weekly Activity Summary** | Weekly | Administration | Student/teacher activity metrics |
| **Monthly Compliance Report** | Monthly | Legal/Regulatory | Data access patterns, violations |
| **Quarterly Security Review** | Quarterly | CISO | Threat assessment, vulnerabilities |
| **Annual Audit Report** | Annual | Board of Directors | Full system audit, compliance status |

### 5.2 Report Metrics

```
┌─────────────────────────────────────────────────┐
│        AUDIT METRICS DASHBOARD                 │
├─────────────────────────────────────────────────┤
│ Total API Calls:         2,458,369            │
│ Successful:              2,412,147 (98.1%)     │
│ Failed:                    46,222 (1.9%)        │
│ Unique Users:              3,245                │
│ Admin Actions:               1,234              │
│ Data Exports:                  45              │
│ Security Incidents:            3 (critical)    │
└─────────────────────────────────────────────────┘
```

## 6. Audit Compliance Checklist

### 6.1 Required Validations
- [ ] Every API endpoint logs entry and exit
- [ ] All data modifications have before/after images
- [ ] Authentication events are fully logged
- [ ] Admin actions require additional verification
- [ ] No PII exposed in log outputs (masked/scrubbed)
- [ ] Audit logs are tamper-evident (immutable)
- [ ] Log integrity verified via hash/chains
- [ ] Retention policy automated enforcement
- [ ] Report generation available on demand
- [ ] Alerting for suspicious patterns

### 6.2 Alerting Thresholds

| Metric | Threshold | Action |
|--------|-----------|--------|
| Failed login attempts | > 100/hour | IP block, security review |
| Admin privilege escalations | > 50/day | Notify security team |
| Data exports | > 10/day | Require approval, log detail |
| Failed validation errors | > 500/hour | Investigate input issues |
| Database query timeouts | > 5% of queries | Performance review |
| Cache miss rate | > 30% | Review caching strategy |

## 7. GDPR & Data Subject Rights

### 7.1 Rights Supported
- **Right to Access**: Student/parent can request all stored data
- **Right to Rectification**: Corrections to inaccurate data
- **Right to Erasure**: "Right to be forgotten" (with legal exceptions)
- **Right to Restriction**: Limit data processing
- **Right to Data Portability**: Export data in machine-readable format
- **Right to Object**: Opt-out of certain processing

### 7.2 Implementation
```go
// Example: Data export functionality
func (h *AuditHandler) ExportData(c *gin.Context) {
    subjectID := c.Param("subject_id")
    format := c.Query("format") // json, csv, pdf
    
    // Verify user permission
    if !userHasPermission(c, "data:export") {
        c.AbortWithStatusJSON(403, ...)
        return
    }
    
    // Collect data with proper filtering
    data := auditService.ExportSubjectData(subjectID, format)
    
    // Set proper headers
    c.Header("Content-Disposition", 
        fmt.Sprintf(`attachment; filename="audit-%s.%s"`, subjectID, format))
    c.Header("Content-Type", "application/octet-stream")
    
    c.Data(http.StatusOK, format, data)
}
```

## 8. Audit Tooling

### 8.1 Logging Infrastructure
- **Elasticsearch**: Log search and analysis
- **Kibana**: Dashboard and visualization
- **Logstash**: Log preprocessing and forwarding
- **Filebeat**: Shipping logs to central store

### 8.2 Alerting & Notification
- **PagerDuty**: Critical security alerts
- **Slack**: Daily/weekly summary reports
- **Email**: Monthly compliance digests
- **Custom Webhook**: Integration with ticketing systems

### 8.3 Query Examples for Common Scenarios

```sql
-- Find users who accessed records they shouldn't
SELECT user_id, COUNT(*) as access_count
FROM audit_logs
WHERE entity_type = 'GRADE'
  AND user_role != 'TEACHER'
GROUP BY user_id
HAVING COUNT(*) > 10;

-- Track data modification pattern
SELECT DATE(timestamp) as date,
       action,
       COUNT(*) as count
FROM audit_logs
WHERE timestamp >= NOW() - INTERVAL '30 days'
GROUP BY DATE(timestamp), action
ORDER BY date, count DESC;
```