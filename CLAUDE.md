# Horse Stable Management System — Development Guide

**Org**: eri-co-id · **Author**: nurafandito (nurafandito@gmail.com) · **Support**: admin@eri.co.id

**Read before coding:** `SPEC.md` (what to build), `CONSTRAINTS.md` (quality bar), `tasks/plan.md` (current task).

**Status**: Phase 1 in progress. Auth (real JWT + real Logout revocation) + master data CRUD + user management + Generic Record API + file upload (local disk, no Garage/S3 creds yet) + SMTP mailer (real net/smtp, LogMailer dev fallback) + dashboards/reports (structured widgets, not stored raw SQL) all working end-to-end against a real DB. No commits on `master` yet. Frontend is still unbuilt (explicitly deferred).

---

## Project Map

```
cmd/server/main.go          — routes + wiring (stdlib http.ServeMux, Go 1.22+ method patterns)
internal/auth/              — OTP sign-in + JWT: handler → service → repository, models, jwt, otp, middleware
internal/user/               — User management CRUD (super_admin only) + SeedSuperAdmin boot-time bootstrap
internal/horse/              — Horse master data CRUD: handler → service → repository, models
internal/equipment/          — Equipment master data CRUD (+ ListByHorse)
internal/lessoncode/         — Lesson code CRUD (soft-"delete" = active=false, history kept)
internal/schema/              — object_definitions/field_definitions reader + ValidateData (schema validation)
internal/record/              — Generic Record CRUD: schema-validated, ownership-scoped (non-admin sees only
                               their own records — CONSTRAINTS.md data-access floor), admin retroactive edits
internal/file/                — Photo/video upload: schema-aware size/count limits, sanitized filenames,
                               reuses record.Service.Get for ownership before attaching a file
internal/storage/             — Storage interface; LocalStorage is the dev stand-in for Garage/S3 (no creds yet)
internal/mailer/               — Mailer: SMTPMailer (stdlib net/smtp) for real, LogMailer dev fallback
internal/dashboard/            — Reports: widgets are structured {metric,field,group_by} descriptors
                               validated against real field_definitions, NEVER raw stored SQL (CWE-89) —
                               deliberate deviation from SPEC.md §7's literal example
internal/db/                — Connect() + embedded SQL migrations (001–013)
internal/db/migrations/     — users, otp_sessions, object_definitions, records+files, role_permissions,
                               horses, equipment, lesson_codes, phase-1 object/field seed, role_permissions seed,
                               files.uploaded_by_id, phase-1 dashboard seed, users.tokens_invalidated_at
frontend/src/               — React scaffold only (App.tsx, main.tsx, index.css) — unbuilt
tasks/plan.md, tasks/todo.md — task breakdown + checklist
```

Each domain package has its own `*_test.go` for service-layer validation (no DB needed). Two real integration suites run against a live Postgres (skip gracefully if `DATABASE_URL` is unreachable — see `internal/testutil`): `internal/auth/otp_integration_test.go` (OTP happy path, 5-attempt lockout, expiry) and `internal/record/ownership_integration_test.go` (CONSTRAINTS.md's data-access floor: trainer isolation, admin sees all). Everything else is still verified by hand against the local dev Postgres, not by an automated suite.

Load only the area you're working on. Everything beyond `internal/auth` and `internal/db` is unbuilt.

---

## Decision-Making & Assumptions

**Do not assume.** Before implementing, surface:

```
ASSUMPTIONS I'M MAKING:
1. [assumption]
2. [assumption]
→ Confirm or correct, and I'll proceed.
```

Clarify tech choices, ambiguous or conflicting specs, and architectural trade-offs. Silent assumptions → rework.

---

## Commit & Task Sync

Every commit maps to a task in `tasks/plan.md` or TodoWrite. No silent commits.

1. Claim the task (TodoWrite or `tasks/plan.md`)
2. Code + test + verify
3. `git commit` — imperative mood, references the task
4. Mark complete

---

## Stack (verified against `go.mod` / `frontend/package.json`)

| Layer | Actual | Notes |
|-------|--------|-------|
| Backend | Go 1.27.1, stdlib `net/http` | Layered monolith: handler → service → repository |
| DB driver | `database/sql` + `github.com/lib/pq` | **Not pgx.** Only non-stdlib Go dep in the project |
| Database | PostgreSQL 16, JSONB for flexible schemas | Migrations embedded via `embed` in `internal/db/migrate.go` |
| Frontend | React 19, TypeScript 6, Vite 8 | `@base-ui/react` + Tailwind 4 + `@tabler/icons-react` |
| Frontend lint | `oxlint` (`npm run lint`) | No eslint, no prettier |
| Auth | Hand-rolled HMAC-SHA256 JWT (`internal/auth/jwt.go`) | No JWT library |
| Storage | Garage (S3-compatible) + imgproxy | Planned, not wired yet |

**Dependency rule:** stdlib first. Adding any Go module or npm package needs a reason stated in the commit message — the project deliberately runs on one Go dep.

Not installed (do not import without adding + justifying): TanStack Query, shadcn/ui runtime, vitest, Playwright, Redis client, SMTP library.

---

## Go Conventions

Per-domain layout under `internal/{domain}/`: `handler.go` (HTTP in/out), `service.go` (business logic), `repository.go` (SQL only), `models.go` (types/constants), `*_test.go`.

- Handler: decode, validate presence, call service, map error → status. See `internal/auth/handler.go`.
- Service: business rules, no SQL, no `http`. See `internal/auth/service.go`.
- Repository: parameterized SQL only (`$1`, `$2`), no validation, no business logic. See `internal/auth/repository.go`.
- Errors: wrap with context — `fmt.Errorf("get user: %w", err)`. Return, don't panic (except unrecoverable env at boot).
- Naming: exported `PascalCase`, unexported `camelCase`. `gofmt` standard, 100 cols soft.
- No suppression comments (`//nolint`). Handle the error or explain why ignoring is safe, in code.

## Frontend Conventions

- Import alias: `@/*` → `frontend/src/*` (declared in `vite.config.ts` and `tsconfig.app.json`).
- Components `PascalCase`, hooks `use*` `camelCase`, truly-constant values `CONSTANT_CASE`.
- Tailwind 4 via `@tailwindcss/vite`; global styles live in `src/index.css` (no `tailwind.config` file, no `styles/globals.css`).
- `tsc -b` must pass (`npm run build`) — `noUnusedLocals`/`noUnusedParameters` are on.
- No `@ts-ignore`. Type it, or use `unknown` plus a `// TODO: type this` line.
- Structure (`pages/`, `hooks/`, `types/`, `utils/`) does not exist yet — create directories as the first real feature needs them, per `SPEC.md`'s schema-driven renderer design.

---

## Comments

Default: no comment. Add only when WHY is non-obvious. Max 2 lines, placed above the code. Never comment out code — delete it. No past tense ("fixed", "added") — git has that.

---

## Database Patterns

- **No N+1.** One query with `LEFT JOIN`s over a loop of per-row fetches.
- **Soft delete:** filter `deleted_at IS NULL` on record reads.
- **Index** columns used in `WHERE` / `ORDER BY` / `JOIN`: `records.created_at DESC`, `records.created_by_id`, `records.object_id`, `files.record_id`, `users.role`. No index on JSONB `data` unless a specific query needs it.
- **Pool:** `SetMaxOpenConns(20)`, `SetMaxIdleConns(5)` — the VPS Postgres is shared across projects. Query contexts carry timeouts.
- Schema changes go in a new numbered migration in `internal/db/migrations/`; never edit an applied one.

---

## Security

- Parameterized queries always — never string-concatenated SQL.
- Uploads: validate MIME type, extension, and size.
- Record `data` JSONB: validate against `field_definitions` before storing.
- JWT in HttpOnly + Secure + SameSite cookie; role middleware on every endpoint; super_admin gate for user/schema management.
- OTP: 6 digits from `crypto/rand`, 15-min expiry, single-use, max 5 attempts per email. Response is identical whether the email exists or not.
- **Never log OTP codes or tokens.** Resolved: OTP delivery goes through `internal/auth.Mailer` (real SMTP via `internal/mailer.SMTPMailer`, stdlib `net/smtp`, no library). Only `mailer.LogMailer` — the dev fallback when `SMTP_*` env vars are unset — logs the code, and it says so explicitly in the log line.
- Secrets in env only (see `.env.example`); never commit `.env`. `gitleaks detect --redact` runs in `make check-fast`.
- `gitleaks`/`semgrep`/`osv-scanner` are **not installed** in this dev environment (installing them is a system-wide `brew`/`pip` change — ask before doing it). Last manual review pass (no tool installs): no hardcoded secrets found; all 46 SQL calls across every repository are parameterized except `internal/dashboard/repository.go`'s identifier interpolation, which is whitelisted against real `field_definitions` names plus a regex defense-in-depth (tested — injection payloads rejected with a clean 404, no DB error). One real gap found and fixed: `/uploads/` was served with zero auth (anyone with the URL could fetch any file) — now requires `RequireAuth`; still not scoped to per-record ownership (that lands with Garage/imgproxy signed URLs per SPEC.md §6). `npm audit`: 7 HIGH findings, all transitive through `shadcn`'s own CLI tooling (devDependency, not shipped code) — tracked as `CONSTRAINTS.md` exception E5, not fixed (the auto-fix force-downgrades `shadcn` to `1.0.0`).

---

## Testing

Go: `make test` (`go test -v -race ./...`). Unit tests mock the repository; integration tests use a real test DB. Target in `CONSTRAINTS.md`: coverage ≥80%.

Frontend: no test runner installed. If a task requires frontend tests, raise it before adding vitest/Playwright.

A skipped test in committed code is a defect — delete it or fix it.

---

## Performance Targets (from `SPEC.md`)

| Operation | Target | Approach |
|-----------|--------|----------|
| Page load | <1s | lazy-load routes, paginate tables |
| Form submit | <100ms | optimistic update, debounced validation |
| Photo upload (10MB) | <3s | stream upload + progress |
| List 1000 records | <500ms | 50/page, indexed filters, eager joins |
| Report | <3s | cache; async if >10s |

Measure locally when a task touches these. Load testing comes after Phase 1.

---

## Common Mistakes

1. Business logic in the repository layer (belongs in service)
2. Unwrapped or swallowed errors
3. N+1 queries
4. Suppression comments instead of fixes
5. Missing index on a new filter column
6. Importing a package that isn't in `package.json` / `go.mod`
7. Editing an already-applied migration
8. Skipped tests left in the commit

---

## Before Committing

```bash
make check-fast          # gofmt -l . && go vet ./... && gitleaks detect --redact
make check-task          # check-fast + go test -v -race ./...
cd frontend && npm run lint && npm run build
git diff --stat
git commit               # imperative, references the task
```

If a constraint blocks you, say so in the commit message with the reason — don't weaken `CONSTRAINTS.md`.

---

## Reference

[base-ui](https://base-ui.com) · [Tailwind 4](https://tailwindcss.com/docs) · [PostgreSQL 16](https://www.postgresql.org/docs/16/) · [net/http](https://pkg.go.dev/net/http) · [database/sql](https://pkg.go.dev/database/sql) · [lib/pq](https://pkg.go.dev/github.com/lib/pq) · [Vite](https://vite.dev)
