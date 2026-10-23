# HORSE STABLE MANAGEMENT SYSTEM — SPEC

**Project**: Horse Stable Management (Activity, Inventory, Schedule, Overtime)  
**Location**: PT Sumberdaya Kuda Asia / `app/`  
**Organization**: eri-co-id (Git owner)  
**Author**: nurafandito (nurafandito@gmail.com)  
**Support**: admin@eri.co.id  
**Date**: 2026-10-23  
**Status**: Spec finalized, ready for Phase 1 implementation

**Objective**: Centralized operations platform for horse training facility: activity logging (with photo proof), consumable & equipment inventory tracking, operational scheduling, and overtime records. Phased rollout, resource-conscious VPS deployment. No admin approval required at log entry.

---

## 1. OBJECTIVE & SCOPE

### Core Mission

Enable trainers, grooms, grass cutters, and cleaners to log daily operations (training, horse activity, inventory, cleaning, feeding) with photo proof. Admin dashboard for full visibility, data management, equipment tracking, and reporting. Zero approval friction at entry point.

### Why This Matters

- **Transparency**: Every activity logged with timestamp, photo, who did it
- **Equipment tracking**: Know which saddle, bridle, blanket assigned to which horse, usage history
- **Consumable management**: Feed/supplement/grass consumption visibility
- **Data-driven decisions**: Reports on horse training patterns, inventory consumption, equipment utilization
- **Auditability**: Immutable activity log for disputes, performance reviews

### Phased Delivery

#### **Phase 1: Activity Logging + Equipment Inventory (MVP)** ⭐

Core: Trainers + Grooms log activities with photo proof. Admin manages consumable & equipment inventory.

**Inclusions:**

- **Authentication**: OTP-based sign-in (no passwords, no signup). Super Admin creates users by email.
- **Schema-driven UI**: Generic form builder + table renderer (all record types)
- **Schema-driven Dashboards**: Dashboard configs define widgets (stats, charts, tables) from queries, rendered generically per object
- Activity record types: **Riding**, **Horse Training**, **Horse Activity** (feed, bathe, walk, groom, etc.), **Cleaning**, **Grass Delivery** (defined via schemas)
- Consumable inventory tracking (serbuk, rumput, vital, equijewel, hay, chaff): in/out only
- **Equipment inventory per horse** (saddle, bridle, halter, blanket, reins, etc.): assign/reassign, track which horse has what
- User roles: Super Admin, Admin, Trainer, Groom, Arit (Grass Cutter), Cleaner (all frontend-native workflows)
- Photo proof upload per activity (mandatory, 1-5 files, max 10MB photo / 50MB video)
- Reports via schema-driven dashboards: Training summary, horse activity timeline, equipment assignment per horse, consumable balance
- **Timeline**: 8-10 weeks (1 full-stack dev) — schema-driven UI + dashboards + OTP auth eliminate custom code per object type

**Not in Phase 1:**

- Schedule calendar, overtime tracking, approval workflows, inventory expiry tracking, horse health alerts

#### **Phase 2: Schedule + Advanced Inventory**

- Full consumable inventory management: stock tracking with reorder levels
- Equipment maintenance tracking (cleaning, repair history)
- Training schedule (calendar view with conflict detection)
- Groom schedule + night piket assignment
- Trainer/groom earned hours tracking
- Inventory reports: consumption trends, reorder alerts
- **Timeline**: 4-5 weeks (post-Phase 1)

#### **Phase 3: Advanced Operations**

- Overtime record & payroll support
- Horse health alerts (derived from activity patterns)
- Inventory expiry/batch tracking
- Equipment depreciation/lifecycle tracking
- Multi-stable support
- Advanced analytics dashboard
- **Timeline**: As needed

---

## 2. USER ROLES & PERMISSIONS

| Role                    | Primary Functions                                                                                     | Devices           | Full Control |
| ----------------------- | ----------------------------------------------------------------------------------------------------- | ----------------- | ------------ |
| **Super Admin**         | All functions: manage users, roles, schemas, objects, data, reports, import/export, retroactive edits | Desktop preferred | ✅ YES       |
| **Admin**               | Manage horses, equipment, lessons, generate reports, edit/delete entries, data import/export          | Desktop preferred | ❌ NO        |
| **Trainer**             | Log riding + horse training, request consumables, view assigned equipment                             | Desktop/Mobile    | ❌ NO        |
| **Groom**               | Log horse activity (care), request consumables, log equipment use, record piket duty                  | Desktop/Mobile    | ❌ NO        |
| **Grass Cutter (Arit)** | Log rumput drop-off activities, confirm delivery quantities                                           | Mobile-first      | ❌ NO        |
| **Cleaner**             | Log facility cleaning activities, confirm areas cleaned                                               | Mobile-first      | ❌ NO        |

**Permission Model:**

- **Super Admin**: Full system access (user management, role assignment, schema edits, data modifications)
- **Admin**: Data management only (cannot manage users/roles/schemas)
- **Trainer/Groom/Arit/Cleaner**: Record creation only (immutable after submission)
- All records immutable once submitted; only Super Admin + Admin can edit retroactively
- Role enforcement in Go middleware, frontend respects permissions

**Frontend-first operations**: All roles operate entirely from app — no paper logs, no approval delays. Entries live immediately.

---

## 3. OBJECTS & DATA MODEL (Schema-Driven)

### 3.1 Activity Records (Schema-Driven Objects)

All records are **append-only** (cannot be deleted by creator). Admin can edit any field or soft-delete if needed.

#### **Phase 1 Object Definitions**

All defined via schema, rendered by generic form builder. Seeded in `internal/db/migrations/009_seed_phase1_objects.sql`; per project owner: `fee_amount` is named `fee`, any object with separate `date`+`time` fields uses one `datetime` field instead, and field keys stay English (`rumput_type` → `grass_type`) even where the spec's own prose uses an Indonesian term:

**riding** — Riding Log

```
date (date, required)
rider (relation: users, required)
horse (relation: horses, required)
code_lesson (relation: lesson_codes, required)
instructor (relation: users)
fee (number, optional)
duration_hours (number, optional)
notes (text, optional)
equipment_used (multiselect: saddle, bridle, halter, blanket, optional)
photos (file: 1-5 max 10MB, required)
```

**training** — Horse Training Log

```
date (date, required)
horse (relation: horses, required)
code_lesson (relation: lesson_codes, required)
instructor (relation: users, required)
duration_hours (number, optional)
fee (number, optional)
notes (text, optional)
equipment_used (multiselect, optional)
photos (file: 1-5 max 10MB, required)
```

**horse_activity** — Horse Activity Log

```
datetime (datetime, required)
horse (relation: horses, required)
activity_type (select: feed, bathe, walk, groom, vet_check, trim, stretch, required)
duration_minutes (number, optional)
notes (text, optional)
photos (file: 1-5 max 10MB, required)
```

**consumable_inventory** — Consumable Inventory Log

```
datetime (datetime, required)
item_type (select: serbuk, rumput, vital, equijewel, hay, chaff, required)
quantity_in (number, optional)
quantity_out (number, optional)
unit (select: kg, liter, required)
current_stock (number, read-only, calculated)
horse_allocated (relation: horses, optional)
notes (text, optional)
```

**equipment_log** — Equipment Assignment/Return

```
datetime (datetime, required)
equipment (relation: equipment, required)
action (select: assign, return, move, inspect, required)
from_horse (relation: horses, optional)
to_horse (relation: horses, optional)
condition (select: good, worn, damaged, required)
notes (text, optional)
photos (file: 1-5 max 10MB, optional)
```

**cleaning** — Facility Cleaning Log

```
datetime (datetime, required)
cleaned_by (relation: users, required)
area (select: stable, arena, tack_room, pasture, wash_area, required)
status (select: completed, in_progress, default: completed)
notes (text, optional)
photos (file: 1-5 max 10MB, optional)
```

**grass_delivery** — Grass Delivery Log

```
datetime (datetime, required)
delivered_by (relation: users, required)
grass_type (select: species/cut-type, required)
quantity (number, required)
unit (select: kg, required)
horse_allocated (relation: horses, optional)
location (select: hay_barn, feeding_area, optional)
photos (file: 1-5 max 10MB, required)
```

**Benefits of schema-driven approach:**

- All objects use same form builder (no custom component per type)
- New activity type = define schema, no frontend code
- Phase 2: Admin UI to add/edit object schemas
- Consistent validation, error handling across all objects

### 3.2 Master Data

#### **Horse**

```
Fields: Name, Code/ID, Breed, DOB, Color, Owner Name, Current Trainer (User),
        Health Status (active, injured, retired), Notes, Photo, Created Date

Used in: All activity logs, inventory allocation, equipment assignment, schedule
Admin-managed: Create, edit, soft-delete, deactivate
```

#### **Equipment** _(Phase 1 NEW)_

```
Fields: Name (saddle, bridle, halter, blanket, reins, etc.), Type (category),
        Code/ID, Serial/Batch, Condition (new, good, worn, damaged),
        Current Horse (optional), Location (tack room, stable, in-use),
        Purchase Date, Notes, Photo, Active (Y/N)

Admin-managed: Create, edit, assign to horse via Equipment Log, mark damaged/retired
Reports: Current assignment per horse, equipment audit trail
```

#### **Lesson Codes**

```
Fields: Code (e.g., "DR-01"), Lesson Type (Dressage, Jump, Beginner, Advanced, Grooming, etc.),
        Default Fee (IDR), Duration (hours, optional), Description, Active (Y/N)

Lookup source for: Riding logs, training logs
Admin-managed: Create, edit, deactivate (keep history)
```

#### **Users**

```
Fields: Name, Email, Phone, Role (Admin/Trainer/Groom/Arit/Cleaner),
        Status (Active/Inactive), Created Date

Permissions: Trainer sees only their own logged activities (except Admin override)
            Groom sees all horse activities they're involved with
            Admin sees everything, can edit/soft-delete any entry
Admin-managed: Create, edit, deactivate
```

#### **Consumable Items** _(Phase 2, but seed in Phase 1)_

```
Fields: Name (serbuk, rumput, vital, equijewel, hay, chaff), Unit (kg, liter),
        Reorder Level (kg), Current Stock (calculated from logs), Last Updated
```

#### **Schedules** _(Phase 2)_

```
Training Schedule: Horse, Date, Time, Instructor, Type, Status (planned/completed)
Groom Schedule: Date, Groom, Shift (morning/afternoon/night), Horses Assigned, Notes
Piket Schedule: Date, Person, Type (night/day), Status, Notes
```

---

## 4. CORE FEATURES BY PHASE

### Phase 1 MVP: Activity Logging + Equipment Inventory

#### **Frontend for Trainers/Grooms/Arit/Cleaners**

**Quick Log Screen** (mobile-optimized)

```
Activity Type → Horse → Date → Notes → Equipment (if applicable) → Photos/Videos → Submit
- Minimal fields, camera-first for photos
- Multiple files: 1-5 photos/videos per activity
- Equipment multi-select (riding: saddle + bridle + halter)
- File size validation shown before upload
```

**View My Activities**

- List of personal logs (last 30 days)
- Filter by activity type, horse, date range
- Edit notes only (photo, date, type locked)
- Equipment used visible

**Horse Timeline**

- One horse's all activities (chronological, newest-first)
- Grouped by date, activity type color-coded
- Photos in gallery view, equipment used noted
- Admin can click → edit modal

**Equipment Assignment** (Groom/Trainer view)

- "Equipment on Horse XYZ" — current assignments
- "Return Equipment" button → log return
- Equipment condition visible

#### **Frontend for Admin**

**Activity Browser**

- Table: Date | Who | What | Horse | Equipment | Notes | Actions
- Filters: Date range, User, Activity Type, Horse, Status
- Pagination (50 per page)
- Click row → Edit modal (photo gallery, all fields editable)

**Equipment Dashboard** _(Phase 1 NEW)_

- "Equipment by Horse" table: Horse | Equipment | Condition | Assigned Date | Logged By
- "Equipment Audit Trail" table: Date | Equipment | From Horse | To Horse | Condition | Notes
- Add Equipment button (+ New Saddle, etc.)
- Edit/Return Equipment modal

**Reports** (export to CSV)

- **Training Summary**: Date range + filters (horse, instructor, type) → table + pie chart (fees by type)
- **Horse Activity Timeline**: One horse + date range → activity types + frequency
- **Consumable Balance**: Item type, in/out summary, current stock
- **Equipment Assignment**: All horses + current equipment + condition

**Master Data Management**

- **Horses**: CRUD table (name, code, breed, trainer, status)
- **Equipment**: CRUD table (name, type, condition, current horse, location)
- **Users**: CRUD table (name, role, email, status) + bulk import CSV
- **Lesson Codes**: CRUD table (code, type, default fee, active)

**Bulk Import** _(Phase 1 critical)_

- CSV upload: Horses | Users | Lesson Codes | Equipment
- Mapping wizard (preview before confirm)
- Validation errors (duplicates, required fields)
- Import → success count + sample

#### **Features NOT in Phase 1**

- ❌ Approval workflow (entries live immediately)
- ❌ Offline sync (mobile requires online)
- ❌ Photo compression/optimization (uploaded as-is)
- ❌ Schedule calendar
- ❌ Overtime tracking
- ❌ Inventory expiry/batch tracking
- ❌ Equipment maintenance scheduling
- ❌ Email/SMS notifications

---

## 5. SCHEMA-DRIVEN UI ARCHITECTURE

### Overview

**Full schema-driven approach** (like Twenty CRM): Backend defines object schemas, frontend renders generically. No custom UI code per object type.

**Benefits:**

- Single form builder + table renderer handles all object types
- Adding new activity types = define schema, no frontend code
- Consistent UX across all objects
- Admin-configurable objects (Phase 2)
- Easier to add custom fields later

### Schema Structure

**Objects** (activities, horses, equipment, users, etc.) are defined by:

- `objectName`: "riding", "training", "horse_activity", "equipment", etc.
- `fields[]`: Array of field definitions
  - `fieldId`: "date", "rider", "horse", "duration", "notes", "photos", etc.
  - `type`: "text" | "number" | "date" | "datetime" | "select" | "multiselect" | "file" | "relation" | "boolean"
  - `label`: Display name
  - `required`: boolean
  - `validation`: Regex, min/max, allowed values
  - `metadata`: Display hints (dateFormat, selectOptions, fileMaxSize, etc.)

### Database Schema (New)

```sql
CREATE TABLE object_definitions (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) UNIQUE NOT NULL, -- "riding", "training", "horse_activity", etc.
  label VARCHAR(255),
  plural_label VARCHAR(255),
  icon VARCHAR(50),
  description TEXT,
  active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE field_definitions (
  id BIGSERIAL PRIMARY KEY,
  object_id BIGINT NOT NULL REFERENCES object_definitions(id),
  name VARCHAR(100) NOT NULL,
  label VARCHAR(255) NOT NULL,
  type VARCHAR(50) NOT NULL, -- text, number, date, select, file, relation, boolean
  required BOOLEAN DEFAULT FALSE,
  position INT, -- field order in form
  validation_rules JSONB, -- {regex, minLength, maxLength, allowedValues, etc.}
  metadata JSONB, -- {dateFormat, selectOptions, fileMaxSize, relation_to_object, etc.}
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(object_id, name)
);

-- Dashboard definitions (schema-driven reports)
CREATE TABLE dashboard_definitions (
  id BIGSERIAL PRIMARY KEY,
  object_id BIGINT REFERENCES object_definitions(id),
  name VARCHAR(100) NOT NULL,
  label VARCHAR(255),
  description TEXT,
  layout JSONB, -- {columns, rows, grid spacing}
  widgets JSONB NOT NULL, -- [{id, type, title, query, metric, position}, ...]
  active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(object_id, name)
);

-- Generic activity/record table (stores data for any object type)
CREATE TABLE records (
  id BIGSERIAL PRIMARY KEY,
  object_id BIGINT NOT NULL REFERENCES object_definitions(id),
  data JSONB NOT NULL, -- {date, rider, horse, duration, notes, etc.} - varies per object
  created_by_id BIGINT NOT NULL REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);

CREATE INDEX idx_records_object_id ON records(object_id);
CREATE INDEX idx_records_created_by_id ON records(created_by_id);
CREATE INDEX idx_records_created_at ON records(created_at DESC);

-- Files table (same, references records generically)
CREATE TABLE files (
  id BIGSERIAL PRIMARY KEY,
  record_id BIGINT NOT NULL REFERENCES records(id) ON DELETE CASCADE,
  field_id BIGINT REFERENCES field_definitions(id),
  file_name VARCHAR(255) NOT NULL,
  file_url VARCHAR(2048) NOT NULL,
  file_size_bytes BIGINT NOT NULL,
  mime_type VARCHAR(50) NOT NULL,
  file_type VARCHAR(20) NOT NULL, -- 'photo' or 'video'
  uploaded_at TIMESTAMP DEFAULT NOW(),
  uploaded_by_id BIGINT NOT NULL REFERENCES users(id)
);

CREATE INDEX idx_files_record_id ON files(record_id);
CREATE INDEX idx_files_uploaded_at ON files(uploaded_at DESC);

-- Master Data Tables (non-schema-driven, fixed structure)
CREATE TABLE horses (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  code VARCHAR(50) UNIQUE,
  breed VARCHAR(100),
  dob DATE,
  color VARCHAR(100),
  owner_name VARCHAR(255),
  current_trainer_id BIGINT REFERENCES users(id),
  health_status VARCHAR(50) DEFAULT 'active',
  notes TEXT,
  photo_url VARCHAR(2048),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);

CREATE TABLE equipment (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(100) NOT NULL,
  code VARCHAR(50) UNIQUE,
  serial_batch VARCHAR(100),
  condition VARCHAR(50) DEFAULT 'good',
  current_horse_id BIGINT REFERENCES horses(id),
  location VARCHAR(255),
  purchase_date DATE,
  notes TEXT,
  photo_url VARCHAR(2048),
  active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  phone VARCHAR(20),
  role VARCHAR(50) NOT NULL, -- 'super_admin', 'admin', 'trainer', 'groom', 'arit', 'cleaner'
  status VARCHAR(50) DEFAULT 'active', -- 'active', 'inactive'
  created_by_id BIGINT REFERENCES users(id), -- Super Admin who created this user
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  deleted_at TIMESTAMP
);

-- OTP table (transient, one row per active OTP request)
CREATE TABLE otp_sessions (
  id BIGSERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL,
  otp_code VARCHAR(10) NOT NULL, -- 6-digit code
  attempt_count INT DEFAULT 0,
  max_attempts INT DEFAULT 5,
  created_at TIMESTAMP DEFAULT NOW(),
  expires_at TIMESTAMP NOT NULL, -- 15 minutes from creation
  verified_at TIMESTAMP
);

CREATE INDEX idx_otp_sessions_email ON otp_sessions(email) WHERE verified_at IS NULL;

-- Role-based access control
CREATE TABLE role_permissions (
  id BIGSERIAL PRIMARY KEY,
  role VARCHAR(50) NOT NULL,
  resource VARCHAR(100) NOT NULL, -- 'objects', 'users', 'schemas', 'reports', etc.
  action VARCHAR(50) NOT NULL, -- 'create', 'read', 'update', 'delete'
  UNIQUE(role, resource, action)
);

CREATE TABLE lesson_codes (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(50) UNIQUE NOT NULL,
  type VARCHAR(100) NOT NULL,
  default_fee BIGINT,
  duration_hours INT,
  description TEXT,
  active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_records_created_at ON records(created_at DESC);
CREATE INDEX idx_records_created_by ON records(created_by_id);
CREATE INDEX idx_horses_active ON horses WHERE deleted_at IS NULL;
CREATE INDEX idx_equipment_horse ON equipment(current_horse_id);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(status) WHERE deleted_at IS NULL;

-- Seed: Default role permissions (Phase 1 startup)
INSERT INTO role_permissions (role, resource, action) VALUES
  ('super_admin', 'users', 'create'), ('super_admin', 'users', 'read'), ('super_admin', 'users', 'update'), ('super_admin', 'users', 'delete'),
  ('super_admin', 'objects', 'create'), ('super_admin', 'objects', 'read'), ('super_admin', 'objects', 'update'), ('super_admin', 'objects', 'delete'),
  ('super_admin', 'schemas', 'create'), ('super_admin', 'schemas', 'read'), ('super_admin', 'schemas', 'update'), ('super_admin', 'schemas', 'delete'),
  ('super_admin', 'reports', 'read'), ('super_admin', 'reports', 'export'),
  ('admin', 'objects', 'create'), ('admin', 'objects', 'read'), ('admin', 'objects', 'update'), ('admin', 'objects', 'delete'),
  ('admin', 'horses', 'create'), ('admin', 'horses', 'read'), ('admin', 'horses', 'update'), ('admin', 'horses', 'delete'),
  ('admin', 'equipment', 'create'), ('admin', 'equipment', 'read'), ('admin', 'equipment', 'update'), ('admin', 'equipment', 'delete'),
  ('admin', 'reports', 'read'), ('admin', 'reports', 'export'),
  ('trainer', 'objects', 'create'), ('trainer', 'objects', 'read'),
  ('groom', 'objects', 'create'), ('groom', 'objects', 'read'),
  ('arit', 'objects', 'create'), ('arit', 'objects', 'read'),
  ('cleaner', 'objects', 'create'), ('cleaner', 'objects', 'read');
```

### Frontend Architecture

**Generic Form Builder**

- Reads schema from `/api/objects/{name}` endpoint
- Renders form fields based on type (text input, date picker, file uploader, select dropdown, etc.)
- Validates input against schema rules before submit
- Reusable across all object types

**Generic Table Renderer**

- Reads schema, renders columns
- Sortable, filterable by field type
- Pagination, bulk actions (delete, export)
- Row detail modal for editing

**Generic Dashboard Renderer**

- Reads dashboard config from `/api/dashboards/{objectName}` endpoint
- Renders widgets (stat cards, charts, tables) from widget definitions
- Widget types: `stat` (count/sum), `chart` (bar/line/pie), `table` (query results), `timeline` (activity stream)
- Queries stored in dashboard config, no frontend hard-coded queries
- Admin can redesign dashboard without code deployment (Phase 2)

**Benefits:**

- One form component instead of 5-6 per object type
- One table component instead of multiple views
- One dashboard renderer handles all reports
- Adding new object type = schema + dashboard config, no code

---

## 6. TECH STACK & ARCHITECTURE

**Pattern**: Based on proven Fidelitas monorepo (Go backend + Vite + React frontend).

### Stack Layers

| Layer            | Technology                                           | Rationale                                               |
| ---------------- | ---------------------------------------------------- | ------------------------------------------------------- |
| **Backend**      | Go 1.27+ (stdlib `net/http`)                         | Fast, low-footprint, proven on Fidelitas VPS            |
| **Frontend**     | Vite + React 19 + TypeScript                         | Fast bundler, component-driven, reuse admin-ui patterns |
| **State Mgmt**   | TanStack Query (data) + React Context (UI)           | Caching, deduplication, sync                            |
| **Data Tables**  | TanStack React Table                                 | Admin reports, bulk views, filtering                    |
| **Charts**       | Recharts                                             | Training summary, consumption trends                    |
| **UI**           | shadcn/ui + Tailwind CSS 4                           | Accessibility, responsive, reuse existing               |
| **Database**     | PostgreSQL 16 (shared)                               | ACID, JSON for flexible metadata                        |
| **Cache/Queue**  | Redis (shared)                                       | Reporting aggregates, photo processing                  |
| **File Storage** | Garage (S3-compatible, shared)                       | Photos, attached to activities                          |
| **Deployment**   | Docker Compose (dev), Docker on Fidelitas VPS (prod) | Single container per service                            |

### Architecture Principles

1. **Monolithic backend**: Single Go binary, modular by package
2. **Schema-driven**: Objects defined via schemas, not hardcoded types
3. **Role-based access**: JWT + role middleware, permission checks per endpoint
4. **Super Admin segregation**: Only super_admin can manage users, roles, schemas
5. **Immutable records**: User logs cannot be edited (super_admin + admin can edit retroactively)
6. **Photo proof first-class**: Every activity record has `files` relationship
7. **Performance**: Eager-loaded queries, Redis cache for reports, pagination
8. **No approval friction**: Records saved immediately; Super Admin/Admin edits retroactively

### Query Optimization & N+1 Prevention

**N+1 Prevention Strategy:**

- **Eager loading**: Join `files`, `users`, `horses`, `equipment` at query time (not loop fetches)
- **Example**: `SELECT a.*, f.*, u.name, h.name, e.name FROM activities a LEFT JOIN files f ON f.activity_id = a.id LEFT JOIN users u ON u.id = a.created_by_id LEFT JOIN horses h ON h.id = a.horse_id LEFT JOIN equipment e ON...` (batch in one query)
- **Batch loading** (for complex cases): Fetch activity IDs first, then fetch related data in one pass per relation

**Indexing Strategy:**

- `idx_activities_occurred_at DESC` — most queries filter by date range
- `idx_activities_horse_id` — horse timeline queries
- `idx_activities_created_by_id` — user's own logs
- `idx_activities_type` — activity type filters
- `idx_files_activity_id` — fetch files for activity
- `idx_files_uploaded_at DESC` — recent uploads
- `idx_equipment_current_horse_id` — equipment by horse

**Pagination & Limits:**

- All list endpoints paginated (default 50/page, max 500)
- No unbounded result sets (enforced in handler)
- Cursor-based pagination for large tables (offset inefficient on 1M+ rows)

**Caching (Redis):**

- Report aggregates cached 1-hour TTL (training fees by type, consumable balance)
- Invalidated on new activity insert
- User-specific dashboard caches 30-min TTL
- Equipment assignment cache 5-min TTL (changes frequently)

**Query Patterns (Go):**

- Use parameterized queries (prepared statements) to prevent injection + plan caching
- Helper function: `db.GetRecordWithRelations(ctx, recordID)` — one call, returns record + files + creator + relations
- Avoid loading files for list views — lazy-load on detail view only

### Authentication & Authorization

**User Management (Super Admin Only)**

- **No self-registration**: Users created by Super Admin only
- **Email-based access**: Super Admin assigns email → user invited
- **User table**: `name`, `email` (unique, required), `role`, `status` (active/inactive), `created_by_id`

**Sign-In Flow**

1. User enters email on login screen
2. System checks if email exists in `users` table + status = active
3. If exists: Send OTP (one-time password) to email via SMTP
4. User enters OTP code (6 digits, expires 15 minutes)
5. On valid OTP: Issue JWT token (24-hour TTL, refresh token 7-day)
6. Token stored in `HttpOnly` cookie (CSRF protection)
7. Role + permissions cached in token claims

**No Google Sign-In in Phase 1** (planned for Phase 2 if needed; currently OTP-only)

**JWT Token:**

- Payload: `{sub: user_id, email, role, permissions: [...]}`
- Issued after OTP validation
- Refresh: 7-day token refreshes access token
- Logout: Clear HttpOnly cookie

**Middleware:**

- `AuthMiddleware`: Verify JWT, extract user + role, reject if expired
- `RoleMiddleware`: Check role_permissions for resource + action
- `OwnerMiddleware`: For field-level access (trainer sees only own records)

**Super Admin Enforcement:**

- User creation/deletion/role assignment → requires `super_admin` role
- Schema edits → requires `super_admin` role
- User invitation: Super Admin sets email → system sends invite link
- No self-demotion (prevent accidental lockout)

### VPS Resource Constraints (Fidelitas VPS)

**Conservative Phase 1 approach:**

- Go binary: ~50-100 MB RAM idle
- Database: ~500 MB for 6 months logs + OTP sessions (auto-cleaned)
- Photos: ~100 GB/year
- Redis: <100 MB (reporting aggregates)
- SMTP: For OTP delivery (configured via Fidelitas mail server or Gmail SMTP)

**VPS can handle Phase 1 comfortably** once migrations out complete.

**SMTP Configuration (OTP Delivery):**

- Use Fidelitas mail server or Gmail SMTP
- OTP sent in seconds; no queuing needed
- Email template: "Your OTP code: [code] (expires 15 minutes). Do not share."
- Logging: Log OTP delivery status (no OTP codes in logs)

**Database Connection Pooling:**

- PgBouncer or pgx pool: 10-20 connections (shared with other projects)
- Go's `sql.DB` with `SetMaxOpenConns(20)`, `SetMaxIdleConns(5)`
- Query timeout: 30 seconds (admin override, 60s for bulk reports)
- Prepared statement caching: Auto (pgx stdlib)

### Security & Data Protection

**Authentication & Sessions**

- OTP-based: 6-digit code, single-use, 15-min expiry, max 5 attempts per email
- JWT tokens: 24-hour access, 7-day refresh
- HttpOnly cookies: CSRF protection, secure flag on HTTPS
- No password storage: Users created by Super Admin only

**Authorization & Access Control**

- Role-based: Super Admin > Admin > Trainer/Groom/Arit/Cleaner
- Field-level: Trainer sees own records only (created_by_id = user_id)
- Super Admin-only: User management, schema edits, role assignment
- Middleware: All endpoints require valid JWT + role check

**Data Protection**

- Parameterized queries: All database queries use prepared statements (`$1`, `$2`)
- Input validation: Schema-based validation on all record fields
- File upload: MIME type + size validation (photos ≤10MB, videos ≤50MB)
- File storage: Garage (S3-compatible), signed URLs (no direct access)
- Metadata stripping: EXIF data removed from photos before storage
- Rate limiting: OTP (5 attempts/15min), API (100 req/min per user)

**Secrets Management**

- `.env` never committed (use `.env.example`)
- SMTP credentials: Environment variables only
- Database password: VPS secrets manager or `.env`
- No secrets in logs: OTP codes, tokens, passwords never logged
- `gitleaks detect --redact` runs before every commit

**API Security**

- HTTPS/TLS: Required in production (VPS nginx)
- CORS: Restrict to frontend domain only
- Request validation: Schema + size limits (max 10MB body)
- Response sanitization: No sensitive data in error messages
- API keys: Not used (JWT + RBAC instead)

**Audit & Compliance**

- Immutable audit trail: Records never deleted by users, only soft-deleted by Admin/Super Admin
- Timestamps: All records include `created_at`, `updated_at`, `deleted_at`, `created_by_id`
- Activity log: All user actions logged (create, edit, delete)
- Data export: Available to Admin + Super Admin only (CSV, includes all fields)

**Threat Prevention**

- SQL injection: Parameterized queries + prepared statements
- XSS: React auto-escapes, no innerHTML, sanitize HTML in text fields
- CSRF: HttpOnly cookies + SameSite=Lax
- Brute force: OTP rate limit (5 attempts), JWT expiry
- Dependency vulnerabilities: `go mod tidy`, `npm audit` before each commit

---

## 7. API ENDPOINTS (Go Backend)

### Objects & Schemas

```
GET    /api/objects                   — List all object definitions (all roles)
GET    /api/objects/{name}            — Get schema for object type (all roles)
POST   /api/objects                   — Create object definition (super_admin only, Phase 2)
PATCH  /api/objects/{name}            — Update schema (super_admin only, Phase 2)
DELETE /api/objects/{name}            — Delete object definition (super_admin only, Phase 2)
```

### Records (Generic CRUD for all object types)

```
GET    /api/objects/{name}/records    — List records of type (paginated, filters by field; all roles)
POST   /api/objects/{name}/records    — Create record (validate against schema; all roles can create own)
GET    /api/objects/{name}/records/{id} — Get single record (all roles)
PATCH  /api/objects/{name}/records/{id} — Update record (super_admin + admin only, retroactive edits)
DELETE /api/objects/{name}/records/{id} — Soft-delete record (super_admin + admin only)

Role-based field access (Phase 2):
- Trainer sees only their own logged records (created_by_id = current user)
- Groom sees all horse activities they're involved with
- Super Admin + Admin see all records
- Salary/fee fields hidden from non-admin

Examples:
  GET  /api/objects/riding/records               → all riding logs (admin view)
  POST /api/objects/horse_activity/records       → new horse activity (trainer/groom creates)
  GET  /api/objects/equipment_log/records?horse_id=5 → equipment changes for horse 5
```

### Files

```
POST   /api/objects/{name}/records/{id}/files  — Upload photo/video
                                                 (validate: 10MB photos, 50MB videos, 1-5 per record)
GET    /api/objects/{name}/records/{id}/files  — List files for record
DELETE /api/files/{id}                          — Delete file (admin only)
```

### Equipment (Master Data)

```
GET    /api/equipment                 — List all equipment
POST   /api/equipment                 — Create equipment
PATCH  /api/equipment/{id}            — Edit equipment
GET    /api/equipment/{id}            — Get single equipment
DELETE /api/equipment/{id}            — Soft-delete (admin only)
GET    /api/equipment/horse/{horseId} — Current equipment on horse
```

### Dashboards (Schema-Driven Reports)

```
GET    /api/dashboards                            — List all dashboard definitions (all roles)
GET    /api/dashboards/{objectName}               — Get dashboard config for object (all roles)
POST   /api/dashboards/{objectName}               — Create/update dashboard (super_admin + admin, Phase 2)
DELETE /api/dashboards/{objectName}               — Delete dashboard (super_admin only, Phase 2)

GET    /api/dashboards/{objectName}/widgets/{id}  — Execute widget query (all roles)
GET    /api/dashboards/{objectName}/widgets/{id}?format=csv — Export widget data as CSV (admin + super_admin)
```

**Example Widget Definition:**

Widgets are **structured descriptors, never a raw SQL string.** Storing and executing admin-authored SQL text directly (as an earlier draft of this section showed) is a SQL-injection / arbitrary-query-execution risk (CWE-89) and contradicts §12/CONSTRAINTS.md's own "parameterized, no SQL injection" requirement for this feature. `metric`/`field`/`group_by` are validated against the object's real `field_definitions` before any query is built — see `internal/dashboard` (`validate` logic in `service.go`, identifier allow-list in `repository.go`). Seeded in `internal/db/migrations/012_seed_dashboards.sql`.

```json
{
  "id": "fees_by_lesson",
  "type": "chart",
  "title": "Fees by Lesson Code",
  "metric": "sum",
  "field": "fee",
  "group_by": "code_lesson",
  "chart_type": "pie"
}
```

### Reports (Legacy - Phase 2+)

```
GET    /api/reports/training-summary              — By date range, horse, instructor, type (admin + super_admin) [Legacy: use /api/dashboards/riding instead]
GET    /api/reports/horse-timeline                — Activities for horse (all roles, filtered by visibility)
GET    /api/reports/consumable-balance            — Inventory in/out summary (admin + super_admin)
GET    /api/reports/equipment-assignment          — Current equipment per horse (all roles)
GET    /api/reports/equipment-audit               — Equipment change history (admin + super_admin)
GET    /api/reports/{objectName}/{fieldName}      — Custom report (count/sum by field; admin + super_admin)
GET    /api/reports/export/{reportName}           — Export to CSV (admin + super_admin)
```

### Auth (OTP-based, no password/signup)

```
POST   /api/auth/send-otp           — Send OTP to email (no auth required)
                                     Request: {email}
                                     Response: {message, expires_in_seconds: 900}

POST   /api/auth/verify-otp         — Verify OTP, issue JWT (no auth required)
                                     Request: {email, otp_code}
                                     Response: {access_token, refresh_token, user: {id, email, role, name}}

POST   /api/auth/refresh            — Refresh access token (refresh token in cookie)
                                     Response: {access_token}

POST   /api/auth/logout             — Clear session (all roles)
GET    /api/auth/me                 — Current user info (all roles, requires JWT)
```

### Master Data & User Management

```
GET    /api/users                   — List users (super_admin only)
POST   /api/users                   — Create user + send invite email (super_admin only)
                                     Request: {email, name, role}
                                     Response: {id, email, role, status, created_at}
                                     (Email sent: "You've been invited to {app}. Sign in at {url}")

PATCH  /api/users/{id}              — Edit user (super_admin only)
PATCH  /api/users/{id}/role         — Change user role (super_admin only)
DELETE /api/users/{id}              — Soft-delete user (super_admin only)
GET    /api/users/{id}              — Get user details (super_admin only)

GET    /api/horses                  — List all horses (all roles)
POST   /api/horses                  — Create horse (admin + super_admin)
PATCH  /api/horses/{id}             — Edit horse (admin + super_admin)
DELETE /api/horses/{id}             — Soft-delete horse (admin + super_admin)

GET    /api/equipment               — List equipment (all roles)
POST   /api/equipment               — Create equipment (admin + super_admin)
PATCH  /api/equipment/{id}          — Edit equipment (admin + super_admin)
DELETE /api/equipment/{id}          — Soft-delete equipment (admin + super_admin)

GET    /api/lesson-codes            — List lesson codes (all roles)
POST   /api/lesson-codes            — Create lesson code (admin + super_admin)
PATCH  /api/lesson-codes/{id}       — Edit lesson code (admin + super_admin)
DELETE /api/lesson-codes/{id}       — Soft-delete lesson code (admin + super_admin)
```

---

## 8. UI/UX DESIGN PRINCIPLES

### Mobile-First

- Quick log: 4 taps (Type → Horse → Photo → Submit)
- Camera-first: Large button, full-screen capture
- Multiple files: Drag/select up to 5 photos or videos per activity
- File size validation: Show error if file >10MB (photo) or >50MB (video)
- Retry on upload failure: Auto-retry with backoff, user-visible progress
- Dark mode support
- Touch-friendly: 44px+ tap targets

### Admin Dashboard

- Activity browser with filters
- **Equipment dashboard** (NEW): Current assignment per horse, audit trail
- Reports: Training, activity, consumable, equipment (tabbed)
- Master data CRUD (Horses, Equipment, Users, Lessons)

### Design System

- Reuse Fidelitas: shadcn/ui, Tailwind, Lucide icons
- Color coding: Activity types + equipment categories
- Accessibility: WCAG 2.1 AA

---

## 9. PROJECT STRUCTURE

```
app/
├── backend/
│   ├── cmd/server/
│   ├── internal/
│   │   ├── objects/            — Schema definitions, generic record CRUD
│   │   ├── records/            — Generic record handler (all activity types)
│   │   ├── files/              — File upload, validation, storage
│   │   ├── dashboards/         — Dashboard definitions, widget query execution
│   │   ├── equipment/          — Equipment master data CRUD
│   │   ├── horses/             — Horse master data CRUD
│   │   ├── users/              — User management
│   │   ├── reports/            — Report generation (schema-aware)
│   │   ├── import/             — Bulk CSV import
│   │   ├── shared/
│   │   │   ├── db/             — Query helpers, eager-load patterns, N+1 prevention
│   │   │   ├── schema/         — Field validator, type mapper
│   │   │   └── auth/           — JWT, role enforcement
│   ├── migrations/
│   ├── go.mod
│   └── CLAUDE.md
├── frontend/
│   ├── src/
│   │   ├── pages/
│   │   │   ├── RecordForm.tsx       — Generic form for all record types
│   │   │   ├── RecordList.tsx       — Generic table for all record types
│   │   │   ├── Dashboard.tsx        — Quick log, recent records
│   │   │   ├── Reports.tsx
│   │   │   ├── Equipment.tsx        — Equipment master data
│   │   │   ├── Horses.tsx           — Horse master data
│   │   │   ├── Users.tsx
│   │   │   └── AdminSchemas.tsx     — Manage object schemas (Phase 2)
│   │   ├── components/
│   │   │   ├── SchemaFormBuilder/    — Generic form renderer
│   │   │   ├── SchemaTableRenderer/  — Generic table renderer
│   │   │   ├── SchemaDashboard/      — Generic dashboard renderer (charts, stats, tables)
│   │   │   ├── FileUpload/
│   │   │   ├── ReportViewer/
│   │   │   └── common/
│   │   ├── hooks/
│   │   ├── types/
│   │   ├── utils/
│   │   └── styles/
│   ├── vite.config.ts
│   ├── package.json
│   └── CLAUDE.md
├── docker-compose.yml
├── SPEC.md
├── README.md
└── .gitignore
```

---

## 10. CODE STYLE & CONVENTIONS

### Go

- Naming: PascalCase (types), camelCase (functions)
- Structure: handler.go, service.go, repository.go, models.go
- Error handling: Explicit wrapping with context
- Formatting: `gofmt` standard, 80-char width

### React

- Naming: PascalCase (components), camelCase (hooks)
- Hooks: TanStack Query (`useQuery`, `useMutation`), React Context
- Formatting: Prettier, 100-char width, single quotes

---

## 11. TESTING STRATEGY

### Go (must-have)

- Unit: Activity creation, file upload, report queries
- Integration: API endpoints (POST/PATCH), equipment assignment
- **Performance tests**: List activities with 1000 entries — must be <100ms, verify no N+1
- **Query audit**: Run `EXPLAIN ANALYZE` on report queries, ensure index usage
- Run: `go test ./...`

### React (critical paths)

- E2E: Trainer logs activity → photo → view in timeline
- E2E: Admin imports equipment.csv → assigns to horse → audits changes
- E2E: Admin generates equipment report → downloads CSV
- Run: `npm run e2e` (Playwright)

### Manual Checklist

- [ ] Trainer logs 5 activities with photos on mobile
- [ ] Admin assigns 10 pieces of equipment to horses
- [ ] Admin imports horses.csv + equipment.csv
- [ ] Report exports to CSV
- [ ] Dark mode readable
- [ ] Keyboard navigation works
- [ ] Concurrent logins OK

---

## 12. KNOWN CONSTRAINTS

### Phase 1 **NOT** Included

- ❌ Approval workflow, schedule calendar, overtime, expiry tracking, health alerts
- ❌ Equipment maintenance scheduling, depreciation tracking
- ❌ Offline mode, photo compression
- ❌ Email/SMS notifications

### Always Required

- ✅ Photo proof (mandatory)
- ✅ Immutable audit trail
- ✅ Role-based access (enforced)
- ✅ Admin override
- ✅ Data export (CSV)

### File Upload Constraints

- **Formats supported**: JPEG, PNG, WebP (photos); MP4, MOV, WebM (videos)
- **Per-file limits**: 10 MB max for photos, 50 MB max for videos
- **Per-activity limit**: 1-5 files per activity log
- **Validation**: Enforced on upload (frontend + backend), rejected if exceeded
- **Metadata stripping**: Location/EXIF data removed before storage (privacy)

### Boundaries

- **Trainer earnings**: Admin visible only
- **No self-registration** (Phase 1)
- **Photo metadata**: Location stripped (privacy)
- **Not payroll system**: Phase 3+
- **VPS capacity**: Monitor disk, archive old activities if needed

---

## 13. ACCEPTANCE CRITERIA (Phase 1)

**Auth:**

1. ✅ Super Admin creates user with email → user receives invite email
2. ✅ User enters email on login → system sends OTP code (expires 15min)
3. ✅ User enters invalid OTP 5x → account locked for 1 minute
4. ✅ User enters valid OTP → JWT issued, logged in, redirected to dashboard
5. ✅ User logs out → session cleared, redirected to login
6. ✅ Unauthenticated request → 401, redirected to login
7. ✅ JWT expires → refresh token renews it, or 401 if expired

**Core Workflows:** 8. ✅ Trainer logs riding activity with photo → appears in Admin dashboard 9. ✅ Groom logs horse activity with photo → visible in horse timeline 10. ✅ Admin assigns equipment (saddle, bridle) to horse → visible in equipment dashboard 11. ✅ Admin bulks-imports horses.csv + equipment.csv without manual entry 12. ✅ Admin generates training summary report (horse, instructor, type, date, fees) → pie chart 13. ✅ Admin generates horse activity report (activities by date, frequency) → table 14. ✅ Admin generates equipment assignment report (horse | equipment | condition | date) → CSV export 15. ✅ Grass cutter logs rumput delivery with photo 16. ✅ Cleaner logs facility cleaning with photo

**Permissions:** 17. ✅ Role-based access enforced (trainer ≠ admin, can't delete) 18. ✅ E2E test passes: Super Admin creates user → User logs in via OTP → Trainer logs + photo → Admin edits → CSV export 19. ✅ Super Admin only: User creation, role assignment, schema edits (verified via API 403 errors) 20. ✅ Trainer sees only own records (other trainers' records hidden)

**Performance & Quality:** 21. ✅ Performance: <1s page loads, <100ms form submit, <3s photo upload, OTP send <2s 22. ✅ Query efficiency: No N+1 queries in record list (verified via query logs), all eager-loaded joins 23. ✅ Database indexes: All filter columns indexed (created_at, object_id, created_by_id, email) 24. ✅ Schema-driven UI: Generic form builder + table renderer handle all 6 activity types (no custom components) 25. ✅ Mobile responsive (320px+), WCAG 2.1 AA 26. ✅ OTP security: Code rotates per request, not reusable, expires after use

---

## 14. NEXT STEPS

### Pre-Implementation (Week 1)

1. ✅ Spec + Constraints approved
2. Initialize git, docker-compose.yml, Go module, React project
3. PostgreSQL migrations (schema tables, users, otp_sessions)
4. OTP auth endpoint (send-otp, verify-otp, refresh, logout)

### Phase 1 Implementation (Weeks 2-10)

1. **Auth & Schema Infrastructure** (Weeks 2-3)
   - OTP auth endpoints (send-otp, verify-otp, refresh)
   - User management endpoints (Super Admin only)
   - Object definitions, field definitions, records tables
   - Dashboard definitions table (for reports/widgets)
   - Field validator, type mapper, schema retrieval endpoints
2. **Generic Record API** (Weeks 3-4)
   - POST/GET/PATCH `/api/objects/{name}/records`
   - Schema validation, field filtering, pagination
3. **Master Data APIs** (Week 4)
   - Horses, Equipment, Users, Lesson Codes CRUD
4. **File Upload & Processing** (Week 4-5)
   - Photo/video upload, validation, Garage storage, EXIF stripping
5. **Schema-Driven Dashboard Infrastructure** (Week 5)
   - Dashboard definitions, widget query execution
   - `/api/dashboards/{objectName}` endpoints
   - Dashboard seeding (riding, training, equipment, consumable)
6. **Schema-Driven Frontend** (Weeks 6-7)
   - Generic form builder (render from schema)
   - Generic table renderer (sort, filter, paginate)
   - Generic dashboard renderer (charts, stats, tables from widget configs)
   - Quick log screen (mobile optimized)
7. **Reports via Dashboards** (Week 7-8)
   - Training summary, horse timeline, consumable balance, equipment assignment via dashboard widgets
   - Widget types: stat cards, charts (bar/pie), tables, timelines
   - CSV export for widget data
8. **Polish, Testing & Accessibility** (Weeks 8-10)
   - E2E tests (form submit, photo upload, dashboard widget export)
   - Performance audit (N+1 queries, load times, widget query optimization)
   - WCAG 2.1 AA audit
   - Mobile responsiveness

### Phase 2: Admin Schema Management + Custom Dashboards (Post-Phase 1)

- Schema editor UI: Admin can add/edit object definitions and fields
- Dashboard designer UI: Admin can create/edit dashboard widgets without code
- Custom object types (beyond Phase 1 activity types)
- Field validation rules UI
- Dashboard widget editor (queries, chart types, metrics)
- Schema versioning & migration
- Custom field types (computed fields, derived metrics)

### Post-Phase 1 (Phase 1 completion)

- Deploy to Fidelitas VPS
- Gather feedback
- Plan Phase 2 based on actual usage patterns

---

## How to Extend (Phase 2+)

See **SKILL.md** for how to add new activity types and dashboards without writing code:

- Add new object definition (activity type) → automatic form + table rendering
- Add new dashboard widget → automatic chart/stat/table rendering
- Customize field validation → no code needed, schema-driven

---

## References

- **SKILL.md**: How to extend without code (activity types, dashboards, Phase 2+)
- **CONSTRAINTS.md**: Quality bar (testing, security, performance)
- **CLAUDE.md**: Code style and dev patterns
- **tasks/plan.md**: Phase 1 implementation plan with 23 tasks
- **TanStack**: [https://tanstack.com](https://tanstack.com)
- **shadcn/ui**: [https://ui.shadcn.com](https://ui.shadcn.com)
