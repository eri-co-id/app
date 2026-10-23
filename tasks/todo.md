# Phase 1 Implementation — Task List

**Status**: Ready to start  
**Timeline**: 8-10 weeks (1 full-stack dev)  
**Quality bar**: See `CONSTRAINTS.md`

---

## Phase 1: Week 1 — Project Initialization & Schema

- [ ] **Task 1**: Initialize git, docker-compose, Go module, React project
  - [ ] Git initialized with `.gitignore`, `.gitconfig`, `.gitmessage`
  - [ ] Docker Compose running (PostgreSQL 16, Redis)
  - [ ] Go module setup + build works
  - [ ] React Vite project initialized
  - [ ] Both run locally without errors

- [ ] **Task 2**: PostgreSQL migrations — Schema foundation
  - [ ] All migration files created (001-005)
  - [ ] Tables created: users, otp_sessions, object_definitions, field_definitions, records, files
  - [ ] Indexes created on common filter columns
  - [ ] Migration rollback tested

- [ ] **Task 3**: OTP Auth API endpoints
  - [ ] POST /api/auth/send-otp implemented + tested
  - [ ] POST /api/auth/verify-otp with rate limiting (5 attempts) implemented + tested
  - [ ] POST /api/auth/refresh implemented + tested
  - [ ] POST /api/auth/logout implemented + tested
  - [ ] GET /api/auth/me implemented + tested
  - [ ] OTP: 6 digits, random, 15-min expiry, single-use
  - [ ] Integration tests pass
  - [ ] Semgrep scan clean (no auth vulnerabilities)

- [ ] **Task 4**: Seed super_admin user + role_permissions
  - [ ] Super Admin user created (seeded)
  - [ ] role_permissions table seeded (30+ rows)
  - [ ] Seed script idempotent

### Checkpoint: Foundation Ready ✓
- [ ] Git initialized, Docker running
- [ ] Database schema complete with all tables + indexes
- [ ] OTP auth working end-to-end
- [ ] Super Admin + role_permissions seeded
- [ ] `go test ./...` passes
- [ ] Ready for human review

---

## Phase 1: Weeks 2-3 — Schema Infrastructure & Master Data CRUD

- [ ] **Task 5**: Schema seeding — Object + field definitions for all 6 activity types
  - [ ] 6 objects seeded (riding, training, horse_activity, consumable_inventory, equipment_log, cleaning, grass_delivery)
  - [ ] ~50 field definitions across all objects
  - [ ] Field types supported: text, number, date, select, multiselect, relation, file, boolean
  - [ ] GET /api/objects returns all 6 objects
  - [ ] GET /api/objects/{name} returns schema with fields

- [ ] **Task 6**: Master Data CRUD — Horses
  - [ ] POST /api/horses (Admin only, with validation)
  - [ ] GET /api/horses (list, paginated)
  - [ ] GET /api/horses/{id} (detail)
  - [ ] PATCH /api/horses/{id} (update, Admin only)
  - [ ] DELETE /api/horses/{id} (soft-delete, Admin only)
  - [ ] Tests pass (create, list, update, delete, unauthorized checks)
  - [ ] E2E: Admin creates horse → lists → updates → deletes

- [ ] **Task 7**: Master Data CRUD — Equipment
  - [ ] Full CRUD (POST, GET list/detail, PATCH, DELETE)
  - [ ] Admin-only authorization enforced
  - [ ] Foreign key to horses (current_horse_id)
  - [ ] Tests pass
  - [ ] E2E: Create equipment → assign to horse → view

- [ ] **Task 8**: Master Data CRUD — Lesson Codes
  - [ ] Full CRUD (POST, GET list/detail, PATCH, DELETE)
  - [ ] Admin-only authorization
  - [ ] Seed default lesson codes
  - [ ] Tests pass

### Checkpoint: Master Data Complete ✓
- [ ] Horses, Equipment, Lesson Codes all CRUD working
- [ ] All endpoints tested + authorized
- [ ] Soft-deletes working correctly
- [ ] Seed script populates demo data
- [ ] Ready for human review

---

## Phase 1: Weeks 3-4 — Generic Record API

- [ ] **Task 9**: Generic record creation — POST /api/objects/{name}/records
  - [ ] Create record with schema validation
  - [ ] All roles can create (immutable after)
  - [ ] created_by_id auto-set
  - [ ] File validation (photos ≤10MB, videos ≤50MB)
  - [ ] EXIF metadata stripped from photos
  - [ ] Parameterized query (SQL injection protected)
  - [ ] Tests pass (valid/invalid schema, file validation, EXIF stripping)
  - [ ] Semgrep scan clean

- [ ] **Task 10**: Generic record list — GET /api/objects/{name}/records
  - [ ] List with pagination (50/page)
  - [ ] Filters by schema fields
  - [ ] Eager loading (no N+1 queries)
  - [ ] Trainer: sees own records only
  - [ ] Admin: sees all records
  - [ ] Tests pass (filtering, eager loading, pagination, trainer access)
  - [ ] EXPLAIN ANALYZE confirms index use
  - [ ] Load test: 1000 records < 500ms

- [ ] **Task 11**: Record detail + update — GET/PATCH /api/objects/{name}/records/{id}
  - [ ] GET single record with all relations
  - [ ] PATCH (Admin only, retroactive edits)
  - [ ] DELETE (soft-delete, Admin only)
  - [ ] Trainer cannot edit own record
  - [ ] Tests pass (get/update/delete, access control)
  - [ ] E2E: Admin edits trainer's record successfully

### Checkpoint: Generic Record API Complete ✓
- [ ] All CRUD endpoints working for any object type
- [ ] Eager loading verified (no N+1)
- [ ] Role-based access enforced
- [ ] File upload with validation + EXIF stripping
- [ ] Tests pass, load test <500ms
- [ ] Ready for human review

---

## Phase 1: Weeks 5-7 — Frontend — Schema-Driven Components

- [ ] **Task 12**: Frontend — SchemaFormBuilder component
  - [ ] Generic form rendering from schema
  - [ ] Field types: text, date, select, multiselect, file, relation
  - [ ] File upload with progress + validation
  - [ ] Real-time validation + error messages
  - [ ] API integration (TanStack Query)
  - [ ] Responsive (mobile 320px+)
  - [ ] Accessible (WCAG 2.1 AA)
  - [ ] E2E: Fill form → upload photo → submit → success
  - [ ] axe-core audit passes

- [ ] **Task 13**: Frontend — SchemaTableRenderer component
  - [ ] Generic table from schema
  - [ ] Sorting, filtering, pagination (50/page)
  - [ ] Row click → detail modal with edit (Admin only)
  - [ ] Responsive (horizontal scroll on mobile)
  - [ ] Accessible (keyboard navigation, screen reader)
  - [ ] E2E: Sort → filter → paginate → edit row
  - [ ] Load test: 1000 rows < 2s

- [ ] **Task 14**: SchemaFormBuilder + SchemaTableRenderer integration
  - [ ] Quick log page (mobile-optimized form)
  - [ ] Activity list page (with filters)
  - [ ] Routes: /log/{activity_type}, /activities
  - [ ] Mobile: 4-tap quick log
  - [ ] E2E: Trainer logs riding → appears in list (mobile)
  - [ ] E2E: Admin filters by horse → sees results

- [ ] **Task 15**: Frontend — Dashboard infrastructure (SchemaDashboard component)
  - [ ] Generic dashboard from schema
  - [ ] Widget types: stat, chart (Recharts), table, timeline
  - [ ] Responsive grid layout (2-3 cols, stacks mobile)
  - [ ] Refresh button to re-query
  - [ ] CSV export (Admin only)
  - [ ] E2E: View dashboard → stat card + chart + table
  - [ ] Export CSV → file downloads

### Checkpoint: Frontend Complete ✓
- [ ] All schema-driven components working
- [ ] Quick log + list pages functional
- [ ] Dashboard widgets rendering correctly
- [ ] E2E tests passing
- [ ] WCAG 2.1 AA verified
- [ ] Ready for human review

---

## Phase 1: Weeks 5-8 — Backend — Dashboard & Reports

- [ ] **Task 16**: Dashboard infrastructure — API + seed
  - [ ] dashboard_definitions table seeded for each activity type
  - [ ] GET /api/dashboards/{objectName} returns config
  - [ ] GET /api/dashboards/{objectName}/widgets/{id} executes query
  - [ ] Widget queries: SUM, COUNT, GROUP BY aggregations
  - [ ] Example: riding dashboard with 3-4 widgets

- [ ] **Task 17**: Widget query execution + caching
  - [ ] Execute widget queries against records table
  - [ ] Redis caching (1h for aggregates, 5m for tables)
  - [ ] Cache invalidation on new record insert
  - [ ] Query safety (parameterized, no SQL injection)
  - [ ] Performance: <3s first run, <100ms cached
  - [ ] Semgrep scan clean

### Checkpoint: Reports & Dashboards Complete ✓
- [ ] All dashboard widgets working + caching
- [ ] Frontend displays dashboards correctly
- [ ] Performance: <3s first load, <100ms cached
- [ ] Ready for human review

---

## Phase 1: Weeks 8-10 — Polish, Testing, Accessibility, Security

- [ ] **Task 18**: File upload — Full integration
  - [ ] Frontend: drag/drop, click, preview
  - [ ] Progress bar during transfer
  - [ ] Validation: MIME + size checked before upload
  - [ ] Backend: validates, strips EXIF, stores in Garage
  - [ ] Response: signed URL (1h expiry)
  - [ ] Error handling: size exceeded, format, network → user-friendly
  - [ ] Retry: auto-retry on failure
  - [ ] E2E: Upload 10MB photo → success, 11MB → error
  - [ ] Load test: 5 concurrent uploads → all succeed

- [ ] **Task 19**: E2E tests — Full user workflows
  - [ ] Super Admin creates user, trainer logs in via OTP
  - [ ] Trainer logs riding activity with photo
  - [ ] Admin edits trainer's record
  - [ ] Admin generates report + exports CSV
  - [ ] All tests pass (headless), no flaky tests

- [ ] **Task 20**: Performance audit
  - [ ] Query audit: no N+1 (list endpoint <2 queries/page)
  - [ ] Load test: 1000 records → list <500ms (p95)
  - [ ] Form submit: <100ms (p95)
  - [ ] OTP send: <2s (p95)
  - [ ] React bundle: <500kB (gzipped)
  - [ ] EXPLAIN ANALYZE confirms index use
  - [ ] `wrk` load test results documented

- [ ] **Task 21**: Security audit
  - [ ] Semgrep scan: zero HIGH findings (OWASP Top 10, CWE)
  - [ ] gitleaks: no secrets detected
  - [ ] osv-scanner: no high/critical CVE
  - [ ] OTP tests: 5-attempt lockout verified, 15-min expiry
  - [ ] Data access: trainer cannot fetch other trainer's records
  - [ ] Parameterized queries verified (no string concat)
  - [ ] File validation: MIME + size enforced
  - [ ] EXIF stripping verified

- [ ] **Task 22**: Accessibility audit — WCAG 2.1 AA
  - [ ] axe-core scan: zero critical/serious violations
  - [ ] Keyboard navigation tested (tab order logical)
  - [ ] Screen reader tested (labels, errors announced)
  - [ ] Color contrast ≥4.5:1 for text
  - [ ] Mobile: 44px+ tap targets, zoom works
  - [ ] Dark mode readable (light + dark themes)

- [ ] **Task 23**: Documentation
  - [ ] README.md: local setup (docker-compose, npm, go)
  - [ ] API.md: endpoint reference with curl examples
  - [ ] DEVELOPER.md: codebase overview, schema extension guide
  - [ ] TROUBLESHOOTING.md: common issues + solutions
  - [ ] Examples tested (instructions work on fresh machine)

### Checkpoint: Phase 1 Complete ✓
- [ ] All core functionality implemented
- [ ] All tests passing (unit + integration + E2E)
- [ ] Security audit passed (semgrep, gitleaks, osv-scanner)
- [ ] Performance targets met
- [ ] Accessibility audit passed (WCAG 2.1 AA)
- [ ] Documentation complete
- [ ] **Ready for deployment to VPS**

---

## Notes

- **Check task dependencies** before starting a task (listed in plan.md)
- **Update this list** as you complete tasks (check off the box `- [x]`)
- **Checkpoints**: After each checkpoint, pause for human review before continuing
- **Issues**: If blocked, check `CONSTRAINTS.md` (quality bar), `SPEC.md` (requirements), or `CLAUDE.md` (code style)
- **Risks**: See `tasks/plan.md` for risk mitigation strategies
- **Support**: admin@eri.co.id

---

**Author**: nurafandito  
**Organization**: eri-co-id  
**Created**: 2026-10-23  
**Last updated**: 2026-10-23
