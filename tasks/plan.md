# Implementation Plan: Horse Stable Management — Phase 1

## Overview

Build a schema-driven activity logging system with OTP authentication, generic form/table/dashboard renderers, and equipment tracking. Phase 1 delivers trainers/grooms logging activities with photos, admins viewing/editing all records, and schema-driven reports—all from reusable components that work with any object type.

**Timeline:** 8-10 weeks (1 full-stack dev)  
**Quality bar:** CONSTRAINTS.md (types, lint, coverage ≥80%, security scanning, OTP/auth tests)

---

## Architecture Decisions

1. **Schema-driven everything:** Objects, fields, records, dashboards all defined in database (JSONB), not code. Avoids custom components per activity type.
2. **Vertical slicing:** Build Auth → Master Data → One activity type → Reports, not "all database, all API, all UI."
3. **Generic renderers:** One SchemaFormBuilder, one SchemaTableRenderer, one SchemaDashboard handle all objects.
4. **OTP auth:** No passwords, no signup. Super Admin creates users by email only. Simpler, more secure.
5. **Role-based access:** Trainer sees own records; Admin sees all; Super Admin manages everything.
6. **Eager loading:** All queries join related data in one pass (no N+1 queries).

---

## Phase 1: Week 1—Project Initialization & Schema

### Task 1: Initialize git, docker-compose, Go module, React project

**Description:** Set up project structure, Git, PostgreSQL + Redis via Docker Compose, Go module with stdlib net/http, React with Vite + TypeScript.

**Acceptance criteria:**
- [ ] Git initialized with `.gitignore`, `.gitconfig` (author: nurafandito), `.gitmessage`
- [ ] `docker-compose.yml` has PostgreSQL 16, Redis, with volumes persisting data
- [ ] `go mod init` with go.mod, go.sum ready for dependencies
- [ ] React Vite project initialized with TypeScript, Tailwind, shadcn/ui
- [ ] Backend runs: `go run ./cmd/server` (listens :8080, connects to DB)
- [ ] Frontend runs: `npm run dev` (Vite, :5173)
- [ ] Docker compose up/down works without errors

**Verification:**
- [ ] `go build ./... && npm run build` succeeds
- [ ] `docker compose up -d && docker compose logs` shows PostgreSQL + Redis ready
- [ ] `curl http://localhost:8080/healthz` returns 200

**Dependencies:** None

**Files likely touched:**
- `.gitignore`, `.gitconfig`, `.gitmessage`
- `docker-compose.yml`
- `backend/go.mod`, `backend/go.sum`
- `backend/cmd/server/main.go`
- `frontend/package.json`, `vite.config.ts`
- `README.md` (setup instructions)

**Estimated scope:** Medium (5 files)

---

### Task 2: PostgreSQL migrations — Schema foundation (users, otp_sessions, object_definitions, field_definitions, records, files)

**Description:** Create all migration files for database schema: users (with role + status), otp_sessions (for OTP auth), object_definitions, field_definitions, records (generic JSONB), files, plus indexes on common filter columns.

**Acceptance criteria:**
- [ ] Migration files created in `backend/migrations/` (001_users.sql, 002_otp_sessions.sql, 003_schema.sql, 004_records_files.sql, 005_indexes.sql)
- [ ] All tables created: users, otp_sessions, object_definitions, field_definitions, records, files
- [ ] Indexes on: created_at DESC, object_id, created_by_id, email, record_id
- [ ] `psql` can connect and `\dt` shows all tables
- [ ] Migration rollback works (can drop and re-create)

**Verification:**
- [ ] `go run cmd/server` applies migrations automatically (or manual `./migrate up`)
- [ ] `SELECT * FROM pg_tables WHERE schemaname='public'` shows 6 tables
- [ ] `SELECT * FROM pg_indexes WHERE schemaname='public'` shows 8+ indexes

**Dependencies:** Task 1 (Docker Compose up)

**Files likely touched:**
- `backend/migrations/001_users.sql`
- `backend/migrations/002_otp_sessions.sql`
- `backend/migrations/003_schema.sql`
- `backend/migrations/004_records_files.sql`
- `backend/migrations/005_indexes.sql`
- `backend/cmd/server/migrate.go` (migration runner)

**Estimated scope:** Medium (6 files)

---

### Task 3: OTP Auth API endpoints (send-otp, verify-otp, refresh, logout, /me)

**Description:** Implement OTP-based authentication: send 6-digit code to email, verify code with rate limiting (5 attempts/15min), issue JWT tokens, refresh flow, logout, and current user endpoint.

**Acceptance criteria:**
- [ ] `POST /api/auth/send-otp` — Email exists + active → sends OTP code via SMTP, returns 200 with expiry
- [ ] `POST /api/auth/verify-otp` — Valid OTP → issues JWT + refresh token, invalid OTP → 429 after 5 attempts
- [ ] `POST /api/auth/refresh` — Valid refresh token → new access token
- [ ] `POST /api/auth/logout` — Clears HttpOnly cookie
- [ ] `GET /api/auth/me` — Returns current user (requires valid JWT)
- [ ] OTP code: 6 digits, cryptographically random, expires 15 min, single-use
- [ ] Attempt tracking: `otp_sessions.attempt_count` increments, resets on success
- [ ] JWT tokens: access 24h, refresh 7d, stored in HttpOnly cookies

**Verification:**
- [ ] Tests pass: `go test ./internal/auth -v`
  - `TestSendOTPValidEmail` — OTP sent successfully
  - `TestSendOTPInvalidEmail` — 400 Bad Request
  - `TestVerifyOTPExpired` — OTP expired after 15min → 400
  - `TestVerifyOTPTooManyAttempts` — 5 failed attempts → 429
  - `TestVerifyOTPSuccess` — Valid OTP → JWT issued, verified in `/me`
  - `TestRefreshToken` — Expired access token → new one issued
- [ ] `curl -X POST http://localhost:8080/api/auth/send-otp -d 'email=trainer@example.com'` returns valid response
- [ ] Integration test: Full OTP flow (send → verify → refresh → logout) works end-to-end

**Dependencies:** Task 2 (migrations with users + otp_sessions tables)

**Files likely touched:**
- `backend/internal/auth/handler.go`
- `backend/internal/auth/service.go`
- `backend/internal/auth/repository.go`
- `backend/internal/auth/models.go`
- `backend/internal/auth/auth_test.go`
- `backend/cmd/server/routes.go` (register auth routes)

**Estimated scope:** Large (6 files)

---

### Task 4: Seed super_admin user + role_permissions

**Description:** Create initial super_admin user (hardcoded in seed script) and populate `role_permissions` table with default permissions for each role (super_admin, admin, trainer, groom, arit, cleaner).

**Acceptance criteria:**
- [ ] Seed script creates super_admin user with email `admin@equine.local` (or configurable), password bypassed (OTP only)
- [ ] `role_permissions` table seeded with 30+ rows: (role, resource, action) tuples for all roles
- [ ] Super Admin can: create users, manage schemas, edit/delete records
- [ ] Admin can: manage master data (horses, equipment), view reports, edit/delete records
- [ ] Trainer/Groom/Arit/Cleaner can: create records (immutable after creation), view own
- [ ] Seed script is idempotent (safe to run multiple times)

**Verification:**
- [ ] Run seed: `go run ./cmd/seed` (or migration hooks it)
- [ ] `SELECT COUNT(*) FROM users WHERE role='super_admin'` returns 1
- [ ] `SELECT COUNT(*) FROM role_permissions` returns ≥30
- [ ] Query: `SELECT action FROM role_permissions WHERE role='super_admin' AND resource='users'` returns ['create', 'read', 'update', 'delete']

**Dependencies:** Task 2 (users + role_permissions tables created)

**Files likely touched:**
- `backend/cmd/seed/main.go`
- `backend/migrations/006_seed_roles.sql`

**Estimated scope:** Small (2 files)

---

## Checkpoint: Foundation Ready

- [ ] Git initialized, Docker Compose running (PostgreSQL + Redis)
- [ ] Database schema created with all tables + indexes
- [ ] OTP auth working (send, verify, refresh, logout, /me endpoints)
- [ ] Super Admin user created, role_permissions seeded
- [ ] `go test ./...` passes for auth + schema packages
- [ ] `go build ./...` succeeds

---

## Phase 1: Weeks 2—3: Schema Infrastructure & Master Data CRUD

### Task 5: Schema seeding — Object and field definitions for all 6 activity types

**Description:** Seed `object_definitions` and `field_definitions` for: riding, training, horse_activity, consumable_inventory, equipment_log, cleaning, grass_delivery. Each object includes 5-10 fields (date, horse, notes, photos, etc.) with types (text, date, select, relation, file, etc.).

**Reference:** See **SKILL.md** for how to add new activity types + dashboards without code (useful for Phase 2 when Super Admin can define schemas via UI).

**Acceptance criteria:**
- [ ] 6 objects seeded in `object_definitions` (one row per activity type)
- [ ] ~50 field definitions across all objects (field_definitions)
- [ ] Field types supported: text, number, date, select, multiselect, relation, file, boolean
- [ ] Validation rules and metadata (max length, allowed values, file size) stored in JSONB
- [ ] Query: `GET /api/objects` returns all 6 objects with labels + icons
- [ ] Query: `GET /api/objects/riding` returns riding schema with all 10 fields

**Verification:**
- [ ] Seed script: `go run ./cmd/seed` creates all objects + fields
- [ ] `SELECT COUNT(*) FROM object_definitions` returns 6
- [ ] `SELECT COUNT(*) FROM field_definitions` returns ~50
- [ ] `curl http://localhost:8080/api/objects/riding` returns valid schema JSON

**Dependencies:** Task 4 (seed script infrastructure)

**Files likely touched:**
- `backend/cmd/seed/seed_objects.go` (seed logic)
- `backend/migrations/007_seed_schemas.sql` (or in Go seed)
- `backend/internal/objects/handler.go` (GET /api/objects endpoints)
- `backend/internal/objects/repository.go`

**Estimated scope:** Medium (4 files)

---

### Task 6: Master Data CRUD — Horses (schema + API + validation)

**Description:** Build complete CRUD for horses: POST (create), GET (list + detail), PATCH (update), DELETE (soft-delete). Include validation, role-based access (Admin only), and tests.

**Acceptance criteria:**
- [ ] `POST /api/horses` — Admin only, validates required fields, returns created horse with ID
- [ ] `GET /api/horses` — List all horses, paginated (50/page), filters by active status
- [ ] `GET /api/horses/{id}` — Get single horse with all fields
- [ ] `PATCH /api/horses/{id}` — Admin only, updates any field
- [ ] `DELETE /api/horses/{id}` — Admin only, soft-delete (sets deleted_at)
- [ ] Validation: name required, code unique, DOB optional, health_status enum
- [ ] Indexes on: code (unique), deleted_at (for soft-delete filtering)

**Verification:**
- [ ] Tests pass: `go test ./internal/horses -v`
  - `TestCreateHorse` — Valid data → success, invalid → error
  - `TestCreateHorseUnauthorized` — Trainer attempts → 403
  - `TestListHorses` — Returns horses, excludes deleted
  - `TestGetHorse` — Returns single horse
  - `TestUpdateHorse` — Admin only, updates correctly
  - `TestDeleteHorse` — Soft-delete works, still queryable via deleted_at
- [ ] `curl -H "Authorization: Bearer $JWT" -X POST http://localhost:8080/api/horses -d '{"name":"Midnight","breed":"Arabian"}'` creates horse
- [ ] E2E: Admin creates horse → lists horses → updates name → soft-deletes

**Dependencies:** Task 4 (role_permissions, JWT auth)

**Files likely touched:**
- `backend/internal/horses/handler.go`
- `backend/internal/horses/service.go`
- `backend/internal/horses/repository.go`
- `backend/internal/horses/models.go`
- `backend/internal/horses/horses_test.go`
- `backend/cmd/server/routes.go`

**Estimated scope:** Large (6 files)

---

### Task 7: Master Data CRUD — Equipment (similar to Horses)

**Description:** CRUD for equipment: create, list, get, update, delete. Each equipment has name, type, condition, current_horse_id, location, serial, notes. Only Admin can manage.

**Acceptance criteria:**
- [ ] `POST /api/equipment` — Create, validates name + type required
- [ ] `GET /api/equipment` — List all, paginated, filters by condition/type
- [ ] `GET /api/equipment/{id}` — Get single with current assignment
- [ ] `PATCH /api/equipment/{id}` — Update (Admin only)
- [ ] `DELETE /api/equipment/{id}` — Soft-delete (Admin only)
- [ ] Foreign key: `current_horse_id` references horses table
- [ ] Condition enum: good, worn, damaged
- [ ] Index on current_horse_id (for "equipment by horse" queries)

**Verification:**
- [ ] Tests pass: `go test ./internal/equipment -v` (similar structure to horses)
- [ ] `curl -X POST http://localhost:8080/api/equipment -d '{"name":"Saddle A","type":"saddle","condition":"good"}'`
- [ ] E2E: Create equipment → assign to horse → view on horse → update condition

**Dependencies:** Task 6 (horses must exist first for foreign key)

**Files likely touched:**
- `backend/internal/equipment/handler.go`, `service.go`, `repository.go`, `models.go`, `*_test.go`
- `backend/cmd/server/routes.go`

**Estimated scope:** Large (6 files)

---

### Task 8: Master Data CRUD — Lesson Codes

**Description:** CRUD for lesson codes: code (e.g., "DR-01"), type, default_fee, duration, description. Used by riding + training logs to populate fee/type fields.

**Acceptance criteria:**
- [ ] `POST /api/lesson-codes` — Create, code must be unique
- [ ] `GET /api/lesson-codes` — List all
- [ ] `PATCH /api/lesson-codes/{id}` — Update
- [ ] `DELETE /api/lesson-codes/{id}` — Soft-delete
- [ ] Enums: type (Dressage, Jump, Beginner, Advanced, etc.)
- [ ] All Admin-only

**Verification:**
- [ ] Tests: `go test ./internal/lesson_codes -v`
- [ ] Seed default lesson codes (DR-01, JM-01, BG-01, etc.)

**Dependencies:** Task 4 (role_permissions)

**Files likely touched:**
- `backend/internal/lesson_codes/handler.go`, `service.go`, `repository.go`, `models.go`, `*_test.go`
- `backend/cmd/seed/seed_lesson_codes.go`

**Estimated scope:** Medium (5 files)

---

## Checkpoint: Master Data Complete

- [ ] Horses, Equipment, Lesson Codes CRUD all working
- [ ] All endpoints tested + authorized
- [ ] Database soft-deletes working correctly
- [ ] Seed script populates demo data

---

## Phase 1: Week 3—4: Generic Record API

### Task 9: Generic record creation API — POST /api/objects/{name}/records

**Description:** Implement POST endpoint that accepts any object type, validates data against schema, and stores in `records` table as JSONB. Include file upload handling (photo/video validation) and schema validation.

**Acceptance criteria:**
- [ ] `POST /api/objects/riding/records` accepts `{date, rider_id, horse_id, duration, notes, photos: [...]}`
- [ ] Validates against riding schema: required fields, type checking, value constraints
- [ ] All roles can create (immutable after creation)
- [ ] created_by_id automatically set to current user
- [ ] Returns created record with ID
- [ ] File validation: photos ≤10MB, videos ≤50MB, supported MIME types (JPEG/PNG/WebP, MP4/MOV/WebM)
- [ ] EXIF metadata stripped from photos before storage
- [ ] Parameterized query (no SQL injection vulnerability)

**Verification:**
- [ ] Tests pass: `go test ./internal/records -v`
  - `TestCreateRecordValidSchema` — Valid data → success
  - `TestCreateRecordInvalidSchema` — Missing required field → error
  - `TestCreateRecordFileValidation` — 11MB photo → error, 10MB photo → success
  - `TestCreateRecordExifStripping` — Photo metadata removed
  - `TestCreateRecordOwnership` — created_by_id set to logged-in user
- [ ] Semgrep scan: no SQL injection findings
- [ ] E2E: Create riding record with photo → appears in database with metadata stripped

**Dependencies:** Task 5 (schema definitions exist), Task 3 (JWT auth)

**Files likely touched:**
- `backend/internal/records/handler.go`
- `backend/internal/records/service.go`
- `backend/internal/records/repository.go`
- `backend/internal/files/handler.go` (file upload/validation)
- `backend/internal/files/exif.go` (metadata stripping)
- `backend/internal/records/records_test.go`

**Estimated scope:** Large (6 files)

---

### Task 10: Generic record list API — GET /api/objects/{name}/records

**Description:** Implement GET with pagination, filtering by schema fields, eager-loading related data (creator, horse, equipment, etc.), and role-based field visibility (trainer sees own only).

**Acceptance criteria:**
- [ ] `GET /api/objects/riding/records` returns all records, paginated (50/page)
- [ ] Query params: `?page=1&created_by_id={user_id}&horse_id={horse_id}` — filters by schema field
- [ ] Eager loading: joins users (creator), horses, equipment in one query (no N+1)
- [ ] Trainer role: filtered to created_by_id = trainer_id automatically
- [ ] Admin role: sees all records
- [ ] Returns JSON: `{data: [...], total: N, page: 1, per_page: 50}`
- [ ] Indexes used: created_by_id, object_id, created_at (verified via EXPLAIN ANALYZE)

**Verification:**
- [ ] Tests: `go test ./internal/records -v`
  - `TestListRecordsTrainerFiltered` — Trainer sees only own records
  - `TestListRecordsAdminUnfiltered` — Admin sees all
  - `TestListRecordsEagerLoading` — Query count = 1 (no N+1)
  - `TestListRecordsPagination` — Returns correct page
- [ ] EXPLAIN ANALYZE shows index scans (no full table scans)
- [ ] Load test: 1000 records list in <500ms (local)

**Dependencies:** Task 9 (records exist)

**Files likely touched:**
- `backend/internal/records/handler.go` (GET endpoint)
- `backend/internal/records/repository.go` (eager-load query)
- `backend/internal/records/records_test.go`

**Estimated scope:** Medium (3 files)

---

### Task 11: Generic record detail + update API — GET/PATCH /api/objects/{name}/records/{id}

**Description:** Get single record with all relations, and update (Admin only, retroactive edits). Soft-delete via PATCH status = deleted_at.

**Acceptance criteria:**
- [ ] `GET /api/objects/riding/records/123` returns record + files + creator + relations
- [ ] `PATCH /api/objects/riding/records/123` — Admin only, updates any field, returns updated
- [ ] `DELETE /api/objects/riding/records/123` — Admin only, soft-delete (sets deleted_at)
- [ ] Trainer cannot update their own record (immutable after creation)
- [ ] Query: single eager-loaded fetch (join files, creator, relations)

**Verification:**
- [ ] Tests: `go test ./internal/records -v`
  - `TestGetRecordWithRelations` — Returns all data
  - `TestUpdateRecordAdminOnly` — Trainer → 403, Admin → success
  - `TestTrainerCannotEditOwnRecord` — Attempt to PATCH own record → 403
  - `TestSoftDelete` — Record still exists (deleted_at set), not returned in list
- [ ] E2E: Admin edits trainer's record, changes date + notes, updates correctly

**Dependencies:** Task 10 (list working)

**Files likely touched:**
- `backend/internal/records/handler.go`
- `backend/internal/records/repository.go`
- `backend/internal/records/records_test.go`

**Estimated scope:** Medium (3 files)

---

## Checkpoint: Generic Record API Complete

- [ ] All CRUD endpoints working for any object type
- [ ] Eager loading verified (no N+1 queries)
- [ ] Role-based access enforced (trainer immutable, admin editable)
- [ ] File upload with validation + EXIF stripping
- [ ] Tests pass, load test <500ms for 1000 records

---

## Phase 1: Weeks 5—7: Frontend — Schema-Driven Components

### Task 12: Frontend — SchemaFormBuilder component (React)

**Description:** Generic form component that renders fields based on schema definition. Handles validation, file upload, and submission.

**Acceptance criteria:**
- [ ] Component accepts schema (ObjectDefinition) + optional initialData
- [ ] Renders form with fields: text input, date picker, select dropdown, multiselect, file upload, relation (select horse/trainer)
- [ ] File upload: shows progress, validates size/MIME, handles errors
- [ ] Validation: real-time error messages, prevents submit if invalid
- [ ] On submit: calls API via TanStack Query (useQuery/useMutation)
- [ ] Responsive: works on mobile (320px+)
- [ ] Accessible: labels, ARIA, keyboard navigation (WCAG 2.1 AA)

**Verification:**
- [ ] E2E tests pass: `npm run e2e`
  - Render riding form → fill fields → upload photo → submit → success
  - Invalid email → error message → fix → resubmit → success
  - Large photo (>10MB) → size warning → choose valid → submit
- [ ] Accessibility audit: axe-core passes (zero critical)
- [ ] Mobile test: iPhone 12 viewport, touch-friendly (44px tap targets)

**Dependencies:** Task 5 (schema definitions), Task 3 (JWT auth, /api/auth/me)

**Files likely touched:**
- `frontend/src/components/SchemaFormBuilder/index.tsx`
- `frontend/src/components/SchemaFormBuilder/useFormState.ts`
- `frontend/src/components/FileUpload/index.tsx`
- `frontend/src/hooks/useMutation.ts` (API calls)
- `frontend/src/types/schema.ts`

**Estimated scope:** Large (5 files)

---

### Task 13: Frontend — SchemaTableRenderer component (React)

**Description:** Generic table component that renders records from any object type. Includes sorting, filtering, pagination, and row detail modal for editing.

**Acceptance criteria:**
- [ ] Component accepts schema + records array
- [ ] Renders columns based on schema fields (excluding files/relations for table view)
- [ ] Sorting: clickable column headers, ascending/descending
- [ ] Filtering: filter by field values (date range, select dropdown)
- [ ] Pagination: 50 rows/page, next/prev buttons
- [ ] Row click → detail modal with full record + edit form (Admin only)
- [ ] Responsive: horizontal scroll on mobile (table doesn't overflow)
- [ ] Accessible: keyboard navigation, screen reader support

**Verification:**
- [ ] E2E tests:
  - Render riding records table → sort by date → paginate to page 2
  - Click row → detail modal → edit notes → save → table updates
  - Filter by horse → shows only matching records
- [ ] Load test: 1000 rows rendered in <2s

**Dependencies:** Task 12 (form builder), Task 10 (list API)

**Files likely touched:**
- `frontend/src/components/SchemaTableRenderer/index.tsx`
- `frontend/src/components/SchemaTableRenderer/useSorting.ts`
- `frontend/src/components/SchemaTableRenderer/useFiltering.ts`
- `frontend/src/hooks/useQuery.ts` (fetch data)

**Estimated scope:** Large (4 files)

---

### Task 14: Frontend — SchemaFormBuilder + SchemaTableRenderer integration (full log pages)

**Description:** Build quick-log page (mobile-optimized form) + activity list page using the generic components.

**Acceptance criteria:**
- [ ] Quick log page: form for riding/training/horse_activity (mobile-first, camera button for photo)
- [ ] Activity list page: table of user's own records (trainers), all records (admin)
- [ ] Routes: `/log/{activity_type}`, `/activities`
- [ ] Form submit → redirect to success page or back to list
- [ ] Mobile: 4-tap flow (type → horse → date → submit)

**Verification:**
- [ ] E2E: Trainer logs riding activity on mobile → photo → submit → appears in activity list
- [ ] E2E: Admin views all activities → filters by horse → exports CSV (mock)

**Dependencies:** Task 12 + Task 13 (components exist)

**Files likely touched:**
- `frontend/src/pages/QuickLog.tsx`
- `frontend/src/pages/ActivityList.tsx`
- `frontend/src/types/api.ts`

**Estimated scope:** Medium (3 files)

---

### Task 15: Frontend — Dashboard infrastructure (SchemaDashboard component)

**Description:** Generic dashboard component that renders widgets (stat cards, charts, tables) from dashboard schema definitions. Queries backend for widget data and renders charts.

**Acceptance criteria:**
- [ ] Component accepts dashboard definition (name, widgets array)
- [ ] Widget types: stat (count/sum), chart (bar/pie/line, via Recharts), table (query results), timeline (activity stream)
- [ ] Each widget has: query (SQL executed on backend), metric (field to display), chart options
- [ ] Responsive grid layout (2-3 columns, stacks on mobile)
- [ ] Refresh button to re-query data
- [ ] Export widget data to CSV (Admin only)

**Verification:**
- [ ] E2E: Admin views riding dashboard → stat card shows "45 hours", pie chart shows "fees by type", table shows top horses
- [ ] Refresh data → updates widgets
- [ ] Export CSV → downloads training_summary.csv

**Dependencies:** Task 5 (dashboard schema), backend dashboard API (Task 16)

**Files likely touched:**
- `frontend/src/components/SchemaDashboard/index.tsx`
- `frontend/src/components/SchemaDashboard/useWidgetQuery.ts`
- `frontend/src/components/widgets/StatCard.tsx`
- `frontend/src/components/widgets/ChartWidget.tsx`

**Estimated scope:** Large (4 files)

---

## Checkpoint: Frontend Complete

- [ ] SchemaFormBuilder, SchemaTableRenderer, SchemaDashboard all working
- [ ] Quick log page functional on mobile + desktop
- [ ] Activity list page with filters
- [ ] All E2E tests passing
- [ ] WCAG 2.1 AA accessibility verified

---

## Phase 1: Weeks 5—8: Backend — Dashboard & Reports

### Task 16: Dashboard infrastructure — API + seed (backend)

**Description:** Dashboard definitions seeded for each activity type (riding, training, horse_activity, etc.). Each dashboard has 3-4 widgets with pre-configured queries.

**Reference:** See **SKILL.md** for how to add new dashboard widgets without code (useful for Phase 2 when Super Admin can design dashboards via UI).

**Acceptance criteria:**
- [ ] `dashboard_definitions` table has rows for each activity type
- [ ] Each dashboard has widgets: [stat (total), chart (by type/horse), table (recent)]
- [ ] `GET /api/dashboards/{objectName}` returns dashboard config
- [ ] `GET /api/dashboards/{objectName}/widgets/{id}` executes widget query, returns data
- [ ] Widget queries: aggregations (SUM, COUNT, GROUP BY) on records table
- [ ] Example: riding dashboard → widgets: [total hours, fees by type pie chart, recent rides table]

**Verification:**
- [ ] Seed creates riding, training, equipment, consumable dashboards
- [ ] `curl http://localhost:8080/api/dashboards/riding` returns dashboard JSON
- [ ] `curl http://localhost:8080/api/dashboards/riding/widgets/total_hours` returns `{value: 123.5}`

**Dependencies:** Task 5 (schema definitions), Task 9-11 (records exist)

**Files likely touched:**
- `backend/cmd/seed/seed_dashboards.go`
- `backend/internal/dashboards/handler.go`
- `backend/internal/dashboards/repository.go`
- `backend/internal/dashboards/widget_executor.go` (execute widget queries)

**Estimated scope:** Large (4 files)

---

### Task 17: Widget query execution + caching (backend)

**Description:** Execute dashboard widget queries against records table. Cache results in Redis (1-hour TTL for slow aggregations, 5-min for frequently-changing equipment assignments).

**Acceptance criteria:**
- [ ] Widget executor reads query from dashboard_definition, executes against records
- [ ] Caching: stat widgets cached 1-hour (aggregate), table widgets 5-min (fresh data)
- [ ] Cache invalidated on new record insert (via trigger or app-level)
- [ ] Query safety: parameterized, no user input in query string
- [ ] Performance: complex aggregation (<3s first run, <100ms cached)

**Verification:**
- [ ] Test: execute widget query → response time <3s
- [ ] Test: second call → <100ms (cached)
- [ ] Test: create new riding record → cache invalidated, next widget call is fresh
- [ ] Semgrep: no SQL injection findings in widget executor

**Dependencies:** Task 16 (dashboard definitions)

**Files likely touched:**
- `backend/internal/dashboards/widget_executor.go`
- `backend/internal/shared/cache/cache.go` (Redis helpers)

**Estimated scope:** Medium (2 files)

---

## Checkpoint: Reports & Dashboards Complete

- [ ] All dashboard widgets executing + caching
- [ ] Frontend displays dashboards correctly
- [ ] Performance: <3s first load, <100ms cached

---

## Phase 1: Weeks 8—10: Polish, Testing, Accessibility, Security

### Task 18: File upload — Full integration (backend + frontend)

**Description:** Complete file upload flow: frontend form → backend validation → Garage storage → EXIF stripping → signed URL response. Handle retries, progress display, error handling.

**Acceptance criteria:**
- [ ] Frontend file input: drag/drop + click, shows preview
- [ ] Progress bar: upload % shown during transfer
- [ ] Validation: MIME type + size checked before upload
- [ ] Backend: receives file, validates, strips EXIF, stores in Garage
- [ ] Response: signed URL (expires 1 hour) for secure retrieval
- [ ] Error handling: size exceeded, unsupported format, upload failed → user-friendly error
- [ ] Retry: on network failure, auto-retry with backoff

**Verification:**
- [ ] E2E: Upload 10MB photo → 2-3s → preview displayed → record saved with file URL
- [ ] E2E: Upload 11MB photo → "File too large" error → choose 10MB → success
- [ ] E2E: Upload .txt file → "Unsupported format" error
- [ ] Load test: concurrent uploads (5 x 10MB photos) → all succeed

**Dependencies:** Task 9 (record creation), Task 12 (form builder)

**Files likely touched:**
- `backend/internal/files/handler.go`
- `backend/internal/files/exif.go`
- `backend/internal/files/garage_client.go`
- `frontend/src/components/FileUpload/index.tsx`

**Estimated scope:** Large (4 files)

---

### Task 19: E2E tests — Full user workflows

**Description:** Playwright E2E tests covering: Super Admin creates user → Trainer logs in via OTP → logs riding activity → Admin edits → generates report → exports CSV.

**Acceptance criteria:**
- [ ] Test: Super Admin flow
  - Login as super_admin (seeded user)
  - Create new trainer user (email: trainer@example.com)
  - Email received with invite
  - Trainer logs in (OTP sent, verified)
  - Trainer accesses dashboard
- [ ] Test: Trainer logs activity
  - Select activity type (riding)
  - Fill form (date, horse, duration, photo)
  - Submit → success, activity visible in list
- [ ] Test: Admin edits
  - Admin views all activities
  - Clicks trainer's activity → detail modal
  - Changes notes → saves
  - Activity list reflects change
- [ ] Test: Reports
  - Admin views riding dashboard
  - Stat widget shows total hours
  - Pie chart shows fees by type
  - Export CSV → file downloaded

**Verification:**
- [ ] All E2E tests pass: `npm run e2e` (headless)
- [ ] No flaky tests (run 3x, all pass)
- [ ] Logs show no unhandled errors

**Dependencies:** Tasks 3, 9-15, 16-17 (all functionality exists)

**Files likely touched:**
- `frontend/e2e/full-workflow.spec.ts`
- `frontend/playwright.config.ts`

**Estimated scope:** Large (2 files, but comprehensive)

---

### Task 20: Performance audit — N+1 queries, load tests, bundle size

**Description:** Verify no N+1 queries, endpoint response times <targets (form submit <100ms, list <500ms), React bundle <500kB.

**Acceptance criteria:**
- [ ] Query audit: run test suite with DB query logging, verify no N+1 (list endpoint queries <2 per page)
- [ ] Load test: 1000 riding records → list endpoint <500ms (p95)
- [ ] Form submit: POST new riding record <100ms (p95)
- [ ] OTP send: <2s (SMTP + DB write)
- [ ] React bundle (minified + gzipped): <500kB
- [ ] EXPLAIN ANALYZE on slow queries: confirms index use

**Verification:**
- [ ] Go test with query logging: `go test ./internal/records -v -run TestLoadRecords`
- [ ] Query log shows ~2 queries per request (no N+1)
- [ ] Load test: `wrk -t4 -c100 http://localhost:8080/api/objects/riding/records` shows <500ms p95
- [ ] Bundle analysis: `npm run analyze` shows total <500kB

**Dependencies:** All implementation tasks

**Files likely touched:**
- `backend/internal/records/records_test.go`
- `frontend/vite.config.ts` (bundle analysis)
- `frontend/package.json` (build scripts)

**Estimated scope:** Medium (3 files)

---

### Task 21: Security audit — Semgrep, gitleaks, OTP/auth tests, data access tests

**Description:** Full security review: code scanning for OWASP Top 10, dependency vulnerabilities, OTP rate limiting + attempt tracking tests, trainer data access control tests.

**Acceptance criteria:**
- [ ] Semgrep scan: zero HIGH findings (OWA…P Top 10, CWE)
- [ ] gitleaks: no secrets leaked
- [ ] osv-scanner: no high/critical CVE in dependencies
- [ ] OTP tests: 5-attempt lockout verified, code expires after 15min
- [ ] Data access: trainer cannot fetch/edit other trainer's records (unit tests)
- [ ] Parameterized queries: all database code uses `$1`, `$2` (no string concat)
- [ ] File validation: MIME + size enforced before storage
- [ ] EXIF stripping: photo metadata confirmed removed

**Verification:**
- [ ] `semgrep scan --config p/owasp-top-ten` → zero HIGH
- [ ] `gitleaks detect --redact` → no findings
- [ ] `osv-scanner scan source -r .` → no high/critical
- [ ] Tests: `go test ./internal/auth -v -run OTP`, `go test ./internal/records -v -run DataAccess`
- [ ] EXIF test: upload photo with metadata → verify removed in storage

**Dependencies:** All implementation tasks

**Files likely touched:**
- Tests throughout (auth, records, files packages)
- No new files needed (use existing test files)

**Estimated scope:** Medium (0 new files, augment existing tests)

---

### Task 22: Accessibility audit — WCAG 2.1 AA

**Description:** Full accessibility review: keyboard navigation, screen reader support, color contrast, ARIA labels, form validation messages.

**Acceptance criteria:**
- [ ] axe-core scan: zero critical or serious violations
- [ ] Keyboard: all interactive elements focusable, logical tab order, Enter/Space activate buttons
- [ ] Screen reader: form labels associated, error messages announced, buttons have text
- [ ] Color: contrast ratio ≥4.5:1 for text, ≥3:1 for large text
- [ ] Mobile: 44px+ tap targets, zoom works (no max-scale=1)
- [ ] Dark mode: all colors readable in light + dark themes

**Verification:**
- [ ] axe-core scan: `npm run audit:a11y` on all pages → zero violations
- [ ] Manual keyboard test: tab through quick log form → all elements reachable
- [ ] Screen reader test (NVDA/JAWS): navigate form → labels read correctly, errors announced
- [ ] Color contrast: check with Chrome DevTools (Accessibility tab)
- [ ] Mobile: test on iPhone 12, pinch-zoom works

**Dependencies:** All frontend tasks (12-15)

**Files likely touched:**
- React component files (add ARIA labels, semantic HTML)
- CSS files (contrast, dark mode)

**Estimated scope:** Medium (3-5 files affected, but small changes per file)

---

### Task 23: Documentation — README, API reference, developer guide

**Description:** Write comprehensive docs: local setup (docker-compose, npm, go), API endpoint reference, schema-driven patterns, troubleshooting.

**Acceptance criteria:**
- [ ] README.md: quick start (git clone → docker compose up → npm run dev)
- [ ] API.md: endpoint reference for auth, objects, records, dashboards with curl examples
- [ ] DEVELOPER.md: codebase overview, schema extension guide, adding new activity types
- [ ] TROUBLESHOOTING.md: common issues (DB migration failures, JWT expired, file upload fails)

**Verification:**
- [ ] README instructions: tested on fresh machine (or simulated)
- [ ] API examples: copy-pasteable curl commands that work
- [ ] Developer guide: clear enough for new dev to add a new activity type

**Dependencies:** All implementation tasks

**Files likely touched:**
- `README.md`, `API.md`, `DEVELOPER.md`, `TROUBLESHOOTING.md`

**Estimated scope:** Small (4 markdown files)

---

## Checkpoint: Phase 1 Complete

- [ ] All core functionality implemented
- [ ] All tests passing (unit + integration + E2E)
- [ ] Security audit passed (semgrep, gitleaks, osv-scanner)
- [ ] Performance targets met (<100ms form, <500ms list)
- [ ] Accessibility audit passed (WCAG 2.1 AA)
- [ ] Documentation complete
- [ ] Ready for deployment to VPS

---

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| OTP email delivery fails (SMTP misconfigured) | High | Test SMTP early (Task 3). Use mock SMTP in dev, real in staging. Alert on failures. |
| N+1 queries in eager loading | Medium | Audit queries early (Task 10). Verify EXPLAIN ANALYZE on all report queries. Load test (Task 20). |
| Schema validation too strict/lenient | Medium | Test with real data early. Adjust validation rules based on feedback. Validation is in JSONB, easy to tweak. |
| File upload security (virus, executable) | Medium | Validate MIME type + size (Task 18). Skip virus scanning Phase 1, plan for Phase 2 if needed. |
| React bundle bloat | Low | Monitor bundle size early (add size-limit to Task 20). Use code splitting if needed. |
| Super Admin password reset broken | High | No password reset needed (OTP only). If super admin loses access to email, manually add new super admin user to DB. Plan recovery process. |

---

## Open Questions

1. **SMTP provider**: Use Fidelitas mail server or Gmail SMTP? Confirm credentials for Task 3.
2. **Garage S3 endpoint**: Confirm endpoint URL + credentials for file storage (Task 18).
3. **VPS deployment**: When deploying to Fidelitas VPS, confirm domain + SSL cert path for nginx reverse proxy.
4. **Monitoring**: Plan for OTP success rate monitoring (Phase 2). Where should alerts go (Slack, email)?

---

## Implementation Notes

- **Task dependencies**: Follow the order; don't skip foundation tasks.
- **Testing strategy**: Each task includes unit tests. E2E tests run after all features complete (Task 19).
- **Checkpoints**: Human review recommended after each checkpoint to catch issues early.
- **Time buffer**: 8-10 weeks assumes no major blockers. Adjust phase timing if needed.

---

**Author**: nurafandito  
**Organization**: eri-co-id  
**Support**: admin@eri.co.id  
**Last updated**: 2026-10-23
