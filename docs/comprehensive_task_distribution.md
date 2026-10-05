# 📋 Document Comprehensive Task Distribution and Unified Technical Path
## Ministry of Education System

---

## 1. Introduction and Unified Work Strategy

### 1.1 Project Philosophy: Full-Stack & Phase Rotation Model

This project follows a **Complete Integration (Full-Stack) & Phase Rotation Model** where all 10 team members participate in every phase of the project. This approach ensures:

- **Shared Ownership**: Every student understands the entire system, not just their isolated component
- **Enhanced Pair Programming**: Knowledge transfer is maximized through constant rotation
- **Unified Understanding**: All team members can discuss any component with confidence
- **Supervisor Readiness**: Prepared for discussions with **Eng. Salah Al-Sayani** (Project Supervisor) at any level

### 1.2 Pair Programming Strategy

- **Dynamic Pairing**: Partners rotate regularly (typically every 2-3 days) to ensure broad knowledge distribution
- **Role Rotation**: Each session has a "Driver" (writing code) and "Navigator" (reviewing architecture), roles switch frequently
- **Quality Focus**: Pair programming significantly raises code quality through real-time review and collaboration
- **Discussion Preparation**: Pairs are always ready to discuss their module with **Eng. Salah Al-Sayani**, understanding both implementation details and architectural context

### 1.3 Project Overview

- **Project**: Ministry of Education Integrated System
- **Supervisor**: **Eng. Salah Al-Sayani**
- **Project Manager & Solutions Architect**: **Omair Sadiq Al-Dedaa**
- **Total Team Members**: 10 students
- **Methodology**: Full-Stack & Phase Rotation Model
- **Project Phases**: 5 complete phases (all students participate in each)

---

## 2. Unified Project Roadmap (5 Phases) & Student Roles

### Phase 1: Analysis & Design (Figma, Wireframing, User Flow)

**Participating Students**: All 10 students

**Activities**:
- Figma wireframing and UI/UX design
- User flow mapping and requirements gathering
- Database entity-relationship diagramming
- API specification drafting
- Project setup and repository initialization

**Student Activities**:
- Collaborative design sessions
- Requirements documentation
- Initial architecture sketches
- Setup of development environments

---

### Phase 2: Database Engineering (PostgreSQL, Schema, ent ORM, Migrations)

**Participating Students**: All 10 students

**Activities**:
- Database schema design and implementation
- ent ORM model creation
- Migration script development with golang-migrate
- Index optimization and performance tuning
- Seed data setup

**Student Activities**:
- Schema design workshops
- ent model generation
- Migration testing and validation
- Performance optimization reviews

---

### Phase 3: Backend Development & Security (Go APIs, Redis Cache, RabbitMQ Queues)

**Participating Students**: All 10 students

**Activities**:
- Go API development with Fiber/Chi framework
- Redis caching strategy implementation
- RabbitMQ/Rabbit Streams event handling
- Authentication & authorization (Dual-Token Strategy)
- Security hardening and OWASP compliance

**Student Activities**:
- Go service development
- API endpoint creation
- Cache integration and TTL configuration
- Event queue setup and consumer implementation
- Security feature implementation

---

### Phase 4: Frontend Development & Integration (React, Vite, Tailwind, Zustand, Axios)

**Participating Students**: All 10 students

**Activities**:
- React 18 component development
- Vite build setup and optimization
- React Router v6 navigation implementation
- Zustand state management
- Tailwind CSS styling
- Axios configuration with token interceptors
- react-helmet-async for SEO & Open Graph

**Student Activities**:
- Component development
- Page routing setup
- State management implementation
- Styling with Tailwind
- API integration with Axios
- SEO metadata setup

---

### Phase 5: Testing, Security & Infrastructure (Docker, Nginx, k6, Postman, CI/CD)

**Participating Students**: All 10 students

**Activities**:
- Docker containerization
- Docker Compose multi-service setup
- Nginx reverse proxy configuration
- k6 load testing
- Playwright/Postman API testing
- CI/CD pipeline setup with GitHub Actions
- Security scanning and vulnerability assessment

**Student Activities**:
- Dockerfile creation
- Multi-container setup
- Testing script development
- Pipeline configuration
- Security assessment and remediation

---

## 3. Student Module Distribution (10 Students)

Each student is responsible for specific modules in partnership with a colleague. Below are the detailed assignments:

### 1. **عمير صادق الدعداع** (Project Manager & Solutions Architect)

- **Module Responsibility**: Overall system architecture, integration pathways, supervisory modules
- **Figma/React Screens**: Project dashboard, architecture overview, supervisory controls
- **DB Tables & Schema**: Master schema, migration coordination, cross-module relationships
- **Go APIs & Redis/RabbitMQ**: Authentication orchestration, system-wide events, monitoring APIs
- **Technologies**: Go (Golang), Ent ORM, Project management tools, Documentation tools
- **Project Manager Additional Tasks**: Team coordination, progress tracking, supervisor meeting preparation, resource allocation, quality assurance oversight

---

### 2. **صقر علي سعيد محمد**

- **Module Responsibility**: Authentication & Authorization module, User Management
- **Figma/React Screens**: Login/Register screens, User profile management, Role-based access interfaces
- **DB Tables & Schema**: Users table, Refresh tokens, Permission tables, Audit logs
- **Go APIs & Redis/RabbitMQ**: JWT generation/validation, Token refresh endpoints, Session management APIs
- **Technologies**: Go (Golang), go-playground/validator, Redis, JWT libraries
- **Special Focus**: Dual-Token Strategy implementation, Redis-based revocation system

---

### 3. **هلال بلال ابراهيم عبدالرحمن**

- **Module Responsibility**: Attendance Tracking system, Student absence monitoring
- **Figma/React Screens**: Attendance dashboard, Absence alerts, Student attendance history
- **DB Tables & Schema**: Attendance records, Absence thresholds, Student absence history
- **Go APIs & Redis/RabbitMQ**: Attendance marking APIs, Threshold event publishing, Cache invalidation
- **Technologies**: Go (Golang), go-playground/validator, Redis Streams, RabbitMQ
- **Special Focus**: Async notification triggers, Cache key structure `cache:attendance:{school_id}:{class_id}:{date}`

---

### 4. **مهند عبده سيف زيد**

- **Module Responsibility**: Student Management & Academic Records
- **Figma/React Screens**: Student profiles, Academic transcripts, Grade entry forms
- **DB Tables & Schema**: Students table, Enrollments, Grades, Academic history
- **Go APIs & Redis/RabbitMQ**: Student CRUD APIs, Grade submission, Academic reporting
- **Technologies**: Go (Golang), Ent ORM, PostgreSQL, JSON Spec
- **Special Focus**: Academic integrity, Grade calculation algorithms, Transcript generation

---

### 5. **عبدالله محمد قاسم نصر**

- **Module Responsibility**: Course Management & Scheduling
- **Figma/React Screens**: Course catalog, Schedule viewer, Prerequisite verification
- **DB Tables & Schema**: Courses table, Prerequisites, Schedule conflicts, Credits
- **Go APIs & Redis/RabbitMQ**: Course creation/updation, Schedule conflict checking, Credits management
- **Technologies**: Go (Golang), Ent ORM, PostgreSQL, Scheduling algorithms
- **Special Focus**: Course capacity management, Prerequisite validation, Schedule optimization

---

### 6. **عبدالله عبدالملك ردمان عقلان**

- **Module Responsibility**: Notification & Reporting System
- **Figma/React Screens**: Notification center, Report generation, Export utilities
- **DB Tables & Schema**: Notification logs, Report metadata, Export history
- **Go APIs & Redis/RabbitMQ**: Notification publishing, Report generation triggers, Export pipelines
- **Technologies**: Go (Golang), Redis Streams, RabbitMQ, Report generation libraries
- **Special Focus**: Async notification processing (RabbitMQ/Redis Streams), Report formatting, Export pipelines

---

### 7. **عصام احمد علي**

- **Module Responsibility**: Database & Infrastructure Management
- **Figma/React Screens**: Dashboard, System monitoring, Infrastructure overview
- **DB Tables & Schema**: System metrics, Connection pooling, Performance statistics
- **Go APIs & Redis/RabbitMQ**: Health checks, Metrics endpoints, Infrastructure APIs
- **Technologies**: Go (Golang), PostgreSQL performance optimization, Docker, Nginx
- **Special Focus**: Database optimization, Infrastructure monitoring, Container orchestration

---

### 8. **وليد علي احمد**

- **Module Responsibility**: API Design & Documentation
- **Figma/React Screens**: API explorer, Developer documentation, Test interfaces
- **DB Tables & Schema**: API logs, Request/response schemas, Deprecation records
- **Go APIs & Redis/RabbitMQ**: API gateway, Documentation serving, Deprecation management
- **Technologies**: Go (Golang), Swagger/OpenAPI, API gateway patterns
- **Special Focus**: API design consistency, Documentation standards, Deprecation workflows

---

### 9. **عدي مهيوب عبدالجليل هزبر**

- **Module Responsibility**: Testing & Quality Assurance
- **Figma/React Screens**: Test reports, Quality metrics, Debug interfaces
- **DB Tables & Schema**: Test results, Bug reports, Quality metrics history
- **Go APIs & Redis/RabbitMQ**: Test trigger endpoints, Metrics collection, Quality dashboards
- **Technologies**: Go (Golang), k6, Playwright/Postman, Test frameworks
- **Special Focus**: Load testing (k6), API testing (Playwright/Postman), Quality metric calculation

---

### 10. **عمار محمد حسن ناجي**

- **Module Responsibility**: Security & Compliance
- **Figma/React Screens**: Security dashboard, Compliance reports, Vulnerability scans
- **DB Tables & Schema**: Security logs, Compliance records, Vulnerability history
- **Go APIs & Redis/RabbitMQ**: Security monitoring APIs, Compliance checks, Audit logs
- **Technologies**: Go (Golang), OWASP Security, Bcrypt (Cost 12), Security scanning tools
- **Special Focus**: OWASP Security compliance, Bcrypt password hashing (Cost 12), Vulnerability assessment, Security audit logging

---

## 4. Pairs Module Matrix (Dual Module Assignment)

| Module | Primary Owner | Co-Lead Owner | Collaboration Focus |
|--------|--------------|---------------|--------------------|
| Authentication & Authorization | Omair Sadiq - Sakr Ali- Helal Belal | Students 2 & 3 | JWT, RBAC, Dual-Token Strategy |
| Attendance & Absence Tracking | Mohannad | Abdullah | Real-time marking, Threshold alerts |
| Student Management | Students 6 & 7 | Students 6 & 7 | CRUD operations, Academic records |
| Course Management | Students 8 & 9 | Students 8 & 9 | Catalog, Scheduling, Prerequisites |
| Notifications & Reports | Students 6 & 10 | Students 6 & 10 | Async events, RabbitMQ/Redis Streams |
| DB & Infrastructure | Students 7 & 10 | Students 7 & 10 | PostgreSQL, ent ORM, Docker, Nginx |
| API Design & Docs | Students 8 & 9 | Students 8 & 9 | Swagger/OpenAPI, Developer experience |
| Testing & QA | Students 9 & 10 | Students 9 & 10 | k6, Playwright, Quality metrics |
| Security & Compliance | Students 10 & 1 (Amer) | Students 10 & 1 (Amer) | OWASP, Bcrypt, Audits |

*Note: Each pair works end-to-end on their module from Figma design through to deployment and testing.*

---

## 5. RACI Matrix and Code Management (Git Flow)

### 5.1 RACI Matrix (Responsibility Assignment)

| Project Element | Responsible (R) | Accountable (A) | Consulted (C) | Informed (I) |
|----------------|-----------------|-----------------|---------------|--------------|
| Code Commit | All team members | Project Manager | Lead developers | Supervisor |
| Feature Development | Module owners | Project Manager | Pair programming partner | Team |
| Bug Fixes | Module owner | Project Manager | QA team | Supervisor |
| Code Review | All team members | Lead developer | Module owner | Team |
| Deployment | DevOps pair | Project Manager | All team members | Supervisor |
| Requirements | All team members | Project Manager | Supervisor | Team |
| Design Reviews | Design leads | Project Manager | All team members | Supervisor |

### 5.2 Git Flow Branching Strategy

```
main              ← Production-ready (protected branch)
develop           ← Integration branch (protected)
feature/*         ← New feature branches (from develop)
hotfix/*          ← Production fix branches (from main)
release/*         ← Release branches (from develop)
```

**Branch Naming Conventions**:
```
feature/student-api-v1
feature/course-enrollment
fix/jwt-auth-timeout
release/v1.0.0-beta
```

### 5.3 Commit Message Convention

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

### 5.4 Daily Standup Format

```
Yesterday:
- What I completed (feature/xyz, bug/abc)

Today:
- What I will work on
- Any dependencies needed

Blockers:
- Issue 1: description
- Issue 2: description

@project: ministry-education
@phase: Phase-2
@tech: go|react|postgres|redis
@depends: task-id (if blocked)
@estimate: 2h
```

### 5.5 Daily Communication Channels

- **Slack/Discord**: Quick questions, informal discussions
- **GitHub Issues**: Task tracking and bug reports
- **Pull Request Discussions**: Code review and feedback
- **Weekly Sync Meetings**: Progress review with supervisor
- **Phase End Reviews**: End-of-phase presentations to **Eng. Salah Al-Sayani**

---

## 6. Technical Stack Summary

### 6.1 Languages
- Go (Golang) - Backend services
- TypeScript - Frontend development
- JavaScript (ES6+) - Frontend interactivity
- SQL - Database queries
- Bash - Shell scripting

### 6.2 Frontend
- React 18 - UI library
- Vite - Build tool and dev server
- React Router v6 - Navigation
- Zustand - State management
- Tailwind CSS - Styling
- Zod & React Hook Form - Validation (Frontend only)
- Axios - API client with token interceptors
- react-helmet-async - SEO & Meta Tags

### 6.3 Backend
- Go (Golang) - Service implementation
- Fiber/Chi - Web framework
- go-playground/validator - Input validation (Backend)
- REST APIs - JSON API specification
- JSON Spec - API documentation format

### 6.4 Database & Cache
- PostgreSQL 15+ - Primary database
- ent ORM / pgx - Database ORM
- Redis - Session revocation, Access/Refresh Token caching, Attendance caching

### 6.5 Asynchronous Processing
- RabbitMQ - Event queue for notifications
- Redis Streams - Real-time event processing

### 6.6 DevOps & Security
- Docker - Containerization
- Docker Compose - Multi-container setup
- Nginx - Reverse proxy
- GitHub Actions - CI/CD
- k6 - Load testing
- Playwright/Postman - API testing
- Bcrypt (Cost 12) - Password hashing
- OWASP Security - Security guidelines

### 6.7 Design Tools
- Figma - UI/UX design
- Draw.io / Enterprise Architect - System architecture
- Postman - API testing
- Jira / GitHub Projects - Task management

---

## 7. Success Criteria and Phase Gates

### Phase Gate Criteria:

| Gate | Criteria | Responsible |
|------|----------|-------------|
| Gate 1: Analysis Complete | All requirements documented, Design approved | Project Manager + Supervisor |
| Gate 2: DB Schema Complete | Schema migrated, Tests passing | DB Team + Project Manager |
| Gate 3: Backend Ready | APIs documented, Security reviewed | Backend Team + Supervisor |
| Gate 4: Frontend Ready | UI complete, SEO optimized | Frontend Team + Supervisor |
| Gate 5: Full System Ready | All tests passing, Deployment ready | All teams + Supervisor |

### Phase Exit Criteria:
- All planned features implemented for the phase
- Unit test coverage > 80% for implemented features
- Security review completed with no CRITICAL findings
- Performance benchmarks met (API response time, load testing)
- Documentation updated and accessible
- Supervisor sign-off obtained

---

## 8. Weekly and Phase Review Cadence

### Weekly Reviews:
- **Every Friday**: 30-minute standup reviewing week progress
- **Bi-weekly**: 1-hour review with supervisor (**Eng. Salah Al-Sayani**)
- **Monthly**: Comprehensive phase review and planning

### Phase Reviews:
- **After each of 5 phases**: 2-hour comprehensive review
- **Deliverables**: Phase deliverables document, lessons learned, next phase plan
- **Supervisor Feedback**: Direct feedback from **Eng. Salah Al-Sayani**

### Phase Transition:
- **Transition Meeting**: 30-minute meeting between phase leads
- **Documentation Update**: Update phase documentation before starting next phase
- **Retrospective**: 15-minute retrospective to capture improvements

---

## 8. Contact and Communication Tree

### Leadership:
- **Project Supervisor**: **Eng. Salah Al-Sayani**
- **Project Manager**: **Amer Sadiq Al-Daqqah**

### Module Leads:
- Authentication: Students 2 & 3
- Attendance: Students 4 & 5
- Student Management: Students 6 & 7
- Course Management: Students 8 & 9
- Notifications: Students 6 & 10
- DB & Infrastructure: Students 7 & 10
- API Design: Students 8 & 9
- Testing & QA: Students 9 & 10
- Security & Compliance: Students 10 & 1

### Daily Standup:
- **Time**: Daily at 9:00 AM
- **Format**: 3 questions (yesterday/today/blockers)
- **Location**: Virtual (Teams/Zoom) or physical (designated room)
- **Duration**: Maximum 15 minutes

### Escalation Path:
1. Team member → Module lead
2. Module lead → Project Manager
3. Project Manager → Supervisor
4. Supervisor → External review if needed

---

## 9. Emergency and Contingency Procedures

### Critical Failure Procedures:
1. **Immediate**: Identify the failure type and impact
2. **Immediate**: Notify module leads and Project Manager
3. **Within 1 hour**: Implement rollback if necessary
4. **Within 2 hours**: Root cause analysis begins
5. **Within 24 hours**: Fix deployed and verified
6. **Within 48 hours**: Retrospective and prevention plan

### Backup Procedures:
- Daily database snapshots
- Code repository backed up to remote
- Configuration files versioned
- Emergency contact list maintained

---

---
*Document Version: 1.0*
*Last Updated: *[Current Date]*
*Prepared for: Ministry of Education System Project*
*Supervised by: **Eng. Salah Al-Sayani***
*Project Manager: **Omair Sadiq Al-Dedaa***
*Team Size: 10 Students*

---
