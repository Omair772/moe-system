# ✅ Todo — Ministry of Education Integrated System (Updated Architecture)

## 1. Project Overview
Comprehensive integrated system for Ministry of Education managing students, courses, attendance, grades, and reporting using Clean Architecture with Go backend, React frontend, PostgreSQL, and native Go ORM.

---

## 2. Phase 1: Documentation & Foundation (Weeks 1-2)

### Week 1
- [ ] Initialize Go module and repository structure
- [ ] Configure Docker Compose (Go, PostgreSQL, Redis)
- [ ] Set up GitHub Actions CI/CD pipeline
- [ ] Create project README and documentation scaffold
- [ ] Establish coding standards (golangci-lint, gofmt)
- [ ] Create docs/ directory and scaffold the 13 markdown files
- [ ] **Tech Stack Decision**: Evaluate ent ORM vs pgx/sqlx (Decision by end of week)

### Week 2
- [ ] Set up PostgreSQL database and initial schema design
- [ ] Initialize chosen Go ORM (ent or pgx) migration system with golang-migrate
- [ ] Configure Redis for caching/TTL and Event Messaging (Redis Streams / RabbitMQ)
- [ ] Create foundational Docker networks
- [ ] **Auth Strategy Decision**: Design Dual-Token mechanism (Access 15-60m + Refresh 24h in Redis)
- [ ] **Cache Key Design**: Design new key structure `cache:api:{role}:{user_id}` and `cache:attendance:{school_id}:{class_id}:{date}`
- [ ] Create initial Go module dependencies

---

## 3. Phase 2: Go Backend & ORM Migration (Weeks 3-5)

### Week 3
- [ ] Implement PostgreSQL Table Partitioning setup for:
  - `ATTENDANCE_RECORDS` (temporal partitioning by month/year)
  - `AUDIT_LOGS` (temporal partitioning by year/quarter)
- [ ] Design partitioning strategy and create partition tables
- [ ] Implement chosen Go ORM (ent ORM or pgx/sqlx) with entity definitions
- [ ] Set up golang-migrate for database schema management
- [ ] Begin implementing API endpoints with native Go ORM
- [ ] **Database Partitioning Task**: Configure partitioning scripts and triggers

### Week 4
- [ ] Implement go-playground/validator for Go Backend validation (replace Zod in backend)
- [ ] Implement Zod validation in React Frontend (Week 8 context)
- [ ] Implement Dual-Token authentication mechanism:
  - Access Token: 15-60 minutes expiry
  - Refresh Token: 24 hours in Redis with Revocation support
- [ ] **Async Notification Queue Task**: Set up Redis Streams / RabbitMQ for Event Messaging
- [ ] Implement event publishing for enrollment/attendance events
- [ ] Configure worker pools for background processing

### Week 5
- [ ] Build Background Worker Pool for processing notifications/SMS
- [ ] Move notification sending outside HTTP Request cycle
- [ ] Implement Redis Streams consumer for event handling
- [ ] Implement SMS service worker integration
- [ ] Unit test suite for ORM operations (>80% coverage)
- [ ] Integration tests for authentication flow
- [ ] API documentation (Swagger/OpenAPI generation)

---

## 4. Phase 3: React Frontend & Integration (Weeks 6-10)

### Week 6
- [ ] React project setup (Vite + TypeScript)
- [ ] Configure Axios with interceptors
- [ ] Create common UI component library
- [ ] Set up React state management (Zustand or Redux)
- [ ] Configure environment variables
- [ ] **Frontend Validation**: Implement Zod schemas for React forms

### Week 7
- [ ] Student dashboard component
- [ ] Course enrollment flow UI
- [ ] Attendance view component
- [ ] Student list with search and pagination
- [ ] Responsive design (mobile-first)
- [ ] API client integration (all student-facing endpoints)

### Week 8
- [ ] Teacher interface (grade entry, attendance marking)
- [ ] Parent portal (view child's data, attendance, grades)
- [ ] Notification system (in-app messages)
- [ ] Form validation with Zod React Hook Form
- [ ] Toast notifications for user feedback
- [ ] API client integration (teacher/parent endpoints)
- [ ] **Zod Focus**: Ensure Zod is used exclusively in frontend (Week 8)

### Week 9
- [ ] Report generation UI
- [ ] Advanced search and filtering
- [ ] Cross-browser testing (Chrome, Firefox, Safari)
- [ ] Performance optimization (Lighthouse >90)
- [ ] State management setup (global stores)
- [ ] API client integration (admin endpoints)

### Week 10
- [ ] End-to-end tests (Cypress)
- [ ] Bug fixes from testing phase
- [ ] Final UI/UX refinements
- [ ] Documentation updates
- [ ] Code review and refactoring

---

## 5. Phase 4: Production Release & Launch (Weeks 11-16)

### Week 11
- [ ] Staging environment deployment
- [ ] Load testing and performance tuning
- [ ] Security penetration testing
- [ ] Database migration validation
- [ ] Backup and disaster recovery validation
- [ ] Monitoring setup (Prometheus + Grafana basic)

### Week 12
- [ ] Production deployment
- [ ] Full monitoring & alerting setup
- [ ] User Acceptance Testing (UAT) with ministry stakeholders
- [ ] Final audit log validation
- [ ] Performance benchmarking (p95 < 200ms)
- [ ] GDPR compliance validation

### Week 13
- [ ] User training sessions (admin, teachers, parents, students)
- [ ] Go-live preparation checklist
- [ ] Production monitoring activation
- [ ] SLA and baseline metrics establishment
- [ ] Handover documentation completion

### Week 14
- [ ] Official system launch
- [ ] Post-launch monitoring (first 2 weeks)
- [ ] Bug fix sprint (critical issues only)
- [ ] Retrospective meeting and lessons learned
- [ ] Project closure and team recognition

### Week 15-16
- [ ] Ongoing maintenance planning
- [ ] Version 2.0 planning and roadmap
- [ ] Dependency updates schedule
- [ ] Community/or knowledge transfer

---

## 6. Daily Task Tracking Format

```
[timestamp] [PRIORITY] Task description
  @project: ministry-education
  @phase: Phase-2
  @tech: go|react|postgres|redis
  @depends: task-id (if blocked)
  @estimate: 2h

Examples:
[2024-01-15 09:00] [HIGH] Implement student create endpoint with go-playground validation
  @project: ministry-education
  @phase: Phase-2
  @tech: go
  @depends: none
  @estimate: 3h

[2024-01-15 14:30] [MEDIUM] Add Redis Streams for event messaging
  @project: ministry-education
  @phase: Phase-2
  @tech: redis
  @depends: 2024-01-15 09:00
  @estimate: 2h
```

---

## 7. Weekly Progress Checklist

### Sprint Week Review
- [ ] Completed tasks from sprint board
- [ ] New tasks identified/added
- [ ] Blockers resolved or escalated
- [ ] Test coverage updated
- [ ] Code review status
- [ ] Documentation updates
- [ ] Metrics captured (velocity, burnup/burndown)

### Key Metrics to Track
- Tasks completed vs. planned
- Defect rate (new bugs vs. resolved)
- Test coverage percentage
- API response time trends
- Cache hit ratio
- Team velocity (story points/week)
- Phase gate completion (each phase must pass before moving to next)
- ORM migration progress (ent/pgx vs Prisma comparison)
- Async queue health (event processing latency, worker count)
- Partition performance (ATTENDANCE_RECORDS, AUDIT_LOGS)

### Architecture Change Milestones
- [ ] Week 1-2: Tech stack evaluation and decision (ent vs pgx/sqlx)
- [ ] Week 2: Auth Dual-Token design and Redis Streams/Basic setup
- [ ] Week 3: Table partitioning setup for ATTENDANCE_RECORDS and AUDIT_LOGS
- [ ] Week 3-4: go-playground/validator implementation (replace Zod in backend)
- [ ] Week 4-5: Background Worker Pool and Async Notification Queue setup
- [ ] Week 5: Full async notification flow verified (no synchronous HTTP sends)
- [ ] Week 5: Cache key refactoring implemented (`cache:api:{role}:{user_id}`)
- [ ] Week 5: Dual-Token auth mechanism validated (Access + Refresh in Redis)