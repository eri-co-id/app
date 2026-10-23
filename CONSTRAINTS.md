# Constraints

**Project**: Horse Stable Management System  
**Organization**: eri-co-id  
**Author**: nurafandito (nurafandito@gmail.com)  
**Support**: admin@eri.co.id  
**Last reviewed**: 2026-10-23  
**Enforced by**: Claude Code agent

---

## Floor (always enforced, no setup required)

- ✅ No new suppression comments: `@ts-ignore`, `eslint-disable`, `# noqa`, `// nolint`
- ✅ No unimplemented stubs: `panic("not implemented")`, empty `catch {}`, `throw new Error("TODO")`
- ✅ No skipped tests: `.skip`, `.only` left in committed code (unless reason in commit message)
- ✅ No deleted tests without reason in commit message
- ✅ No secrets in source: passwords, API keys, tokens, DB URLs
- ✅ This file does not get weakened to make a change pass

---

## Enforced with Numbers

| Dimension | Rule | Checked by | Runs at | Budget |
|-----------|------|-----------|---------|--------|
| **Types: Go** | Zero build errors | `go build ./...` | edit + task | 5s |
| **Types: React** | Zero type errors | `tsc --noEmit` | edit + task | 5s |
| **Lint: Go** | `gofmt` compliant, `go vet` clean | `gofmt -l . && go vet ./...` | edit + task | 3s |
| **Lint: React** | Prettier + ESLint pass | `prettier --check . && eslint .` | edit + task | 5s |
| **Secrets** | No secrets in source, no OTP codes in logs | `gitleaks detect --redact --no-banner` | task + CI | 10s |
| **Security: code** | No HIGH semgrep findings (OWASP Top 10, CWE) | `semgrep scan --config p/owasp-top-ten` | task end | 30s |
| **Security: dependencies** | No high/critical CVE in Go + JS | `go mod tidy && npm audit` | task end | 10s |
| **Security: OTP auth** | OTP valid only 15min, max 5 attempts, single-use | Integration tests: expired OTP fails, 6th attempt locked | task end | 20s |
| **Security: data access** | Trainer cannot fetch/edit other trainer's records | Unit test: `GET /api/records` filtered by `created_by_id` | task end | 10s |
| **Coverage: new code** | Changed lines ≥80% covered | `go test ./... -cover` + `vitest run --coverage` | task end | 90s |
| **Performance: queries** | Zero N+1 patterns detected, eager-loaded joins verified | `EXPLAIN ANALYZE` on report queries (manual audit) | task end | N/A |
| **Performance: endpoints** | POST/PATCH <100ms, GET list <500ms (local), OTP send <2s | Load test on feature completionᵃ | task end | N/A |

ᵃ *Performance testing scope: Only run detailed load tests when task involves record creation, list endpoints, or report generation. Simple CRUD operations pass if response time feels reasonable locally.*

---

## Measured, not yet Enforced

| Metric | Today | Direction | When enforced |
|--------|-------|-----------|----------------|
| Project test coverage | 0% (new project) | Must not fall below 75% | After Phase 1 complete |
| Bundle size (React) | 0 kB (new project) | Track & alert if >500 kB | After Phase 1 complete |
| Database query time (p95) | baseline at Phase 1 deploy | Must not degrade | Post-deployment monitoring |
| Semgrep HIGH findings | 0 (baseline) | Must stay at 0 | Every task (already enforced) |
| Dependency vulnerabilities | 0 (baseline) | Must stay at 0 | Every task (already enforced) |
| OTP success rate | TBD after Phase 1 | Track in monitoring (alert if <99%) | Phase 2 monitoring |

---

## Exceptions

| ID | Rule | Path | Reason | Owner | Expires |
|----|------|------|--------|-------|---------|
| E1 | Coverage ≥80% on new code | `internal/reports/` | Report queries need complex setup; cover with integration tests first | @team | 2026-12-23 |
| E2 | Performance audit | File upload endpoints | 50MB video processing is I/O bound; acceptable if <5s locally | @team | 2026-12-23 |
| E3 | Semgrep: SQL injection (CWE-89) | `internal/db/eager_load.go` | Dynamic table/column names in `EXPLAIN ANALYZE` helper (not user input); validated via schema definitions | @team | 2027-01-23 |
| E4 | HTTPS enforcement | Local development only | HTTP allowed in localhost dev; required in staging/prod (nginx enforces) | @team | Permanent (dev-only) |
| E5 | Dependency vulnerabilities: no high/critical CVE | `frontend/package.json` (`shadcn`) | 7 HIGH `npm audit` findings (stack-exhaustion DoS in `braces`, GHSA-vfj7-8cjw-p6xm) are transitive through `shadcn`'s own CLI tooling (`micromatch`→`fast-glob`→`@shadcn/registry`/`ts-morph`) — a devDependency never shipped to users, not runtime/browser code. `npm audit fix --force` downgrades `shadcn` to `1.0.0`, breaking the already-chosen `base-nova` style in `components.json`. Owner declined the forced downgrade; re-check when shadcn ships a patched `fast-glob`. | @team | 2026-12-23 |

---

## Installation & Scripts

### Go Setup

```bash
# Already standard in Go
go build ./...
go vet ./...
gofmt -l .
go test ./... -cover
go mod tidy  # Clean up dependencies
```

### React Setup

```bash
cd frontend
npm install  # Includes vitest, @vitest/ui for coverage
npm install -D prettier eslint gitleaks semgrep
npm audit  # Check for dependency vulnerabilities
```

### Security Tools Installation

```bash
# Secrets scanning (required)
brew install gitleaks  # or: pip install gitleaks

# Code scanning with OWASP rules (required)
pip install semgrep  # or: brew install semgrep

# Dependency vulnerability scanning (required)
brew install osv-scanner
```

### Project Scripts

Add to project root `Makefile` or shell scripts:

```bash
# Fast checks (edit loop, <10s)
check-fast:
	cd backend && go build ./... && go vet ./...
	cd frontend && prettier --check . && eslint .
	gitleaks detect --redact --no-banner

# Task end checks (~120s, includes tests + security)
check-task: check-fast
	cd backend && go test ./... -cover -v
	cd frontend && vitest run --coverage
	semgrep scan --config p/owasp-top-ten --quiet
	go mod tidy && npm audit

# Full checks (CI only, no time limit)
check-full: check-task
	cd backend && go test ./... -race
	cd frontend && npm run build
	osv-scanner scan source -r .
	# Manual: EXPLAIN ANALYZE on report queries in staging
	# Manual: OTP integration tests + data access control verification
```

Or as npm/go scripts in respective `package.json` / `Makefile`.

**Key security commands:**

```bash
# Secrets (runs before every commit)
gitleaks detect --redact --no-banner

# Code scanning (OWASP Top 10)
semgrep scan --config p/owasp-top-ten --config p/cwe-top-25

# Dependency vulnerabilities
osv-scanner scan source -r .

# OTP auth integration tests
go test ./internal/auth -v -run TestOTP

# Data access control tests
go test ./internal/records -v -run TestTrainerCannotFetchOtherTrainerRecords
```

---

## Lifecycle Wiring

| Phase | Command | What Runs | Budget |
|-------|---------|-----------|--------|
| **BUILD** (after edit) | `check-fast` | Types, lint, secrets (Go + React) | <10s |
| **VERIFY** (task end) | `check-task` | All above + unit tests + coverage | ~90s |
| **REVIEW** (before merge) | `check-full` + manual | All above + integration tests + performance audit | minutes |
| **SHIP** (CI/deployment) | Auto-run `check-full` | Same as review, plus E2E smoke tests | CI gates |

---

## Floor Guard Checklist (per PR/commit)

Before considering a task done, verify:

- [ ] No new `@ts-ignore`, `eslint-disable`, `// nolint` (if present, state reason in commit message)
- [ ] No `panic()` or empty `catch {}` without error handling context
- [ ] No `.skip` or `.only` on tests
- [ ] No deleted test files (if needed, mark as deprecated with a comment)
- [ ] `gitleaks detect --redact` passes (no secrets, **especially OTP codes in logs**)
- [ ] No hardcoded credentials (SMTP password, API keys, database URLs → use env vars)
- [ ] Coverage of changed lines ≥80% (Go and React)
- [ ] Parameterized queries used (all database queries with `$1`, `$2` notation)
- [ ] File uploads validated (MIME type + size before storage)
- [ ] `CONSTRAINTS.md` not weakened (compare `git diff` against branch point)

**Security-Specific Checks:**

- [ ] **OTP auth**: Test passes for 6-digit code, 15-min expiry, 5-attempt lockout
- [ ] **Data access**: Trainer sees only own records (created_by_id filter enforced)
- [ ] **Secrets**: No password/OTP/token in logs (only delivery status logged)
- [ ] **Rate limiting**: OTP endpoint returns 429 after 5 attempts per email
- [ ] **Parameterized queries**: All user input goes through prepared statements
- [ ] **EXIF stripping**: Photos uploaded without metadata (validated in file handler)

---

## Security Testing Acceptance Criteria (by Feature)

### OTP Authentication

- [ ] OTP code: 6 digits, cryptographically random (not sequential), single-use
- [ ] Expiry: Exactly 15 minutes, auto-deletes from `otp_sessions` after verification
- [ ] Rate limit: 5 failed attempts per email per 15 minutes, then 429 response
- [ ] Attempt tracking: `otp_sessions.attempt_count` increments on wrong code, resets on correct
- [ ] Email template: Never includes full code in subject; body says "Your code: [code]"
- [ ] Logging: OTP delivery status logged (sent/failed), never the code itself
- [ ] Integration test passes: `TestOTPExpiredAfter15Minutes`, `TestOTP5FailedAttemptsLocked`, `TestOTPSingleUseOnly`

### Data Access Control (Trainer/Groom Records)

- [ ] Trainer token → `GET /api/objects/riding/records` returns only `created_by_id = trainer_id`
- [ ] Trainer cannot fetch other trainer's record via `GET /api/objects/riding/records/{id}` (403 Forbidden)
- [ ] Admin token → `GET /api/objects/riding/records` returns all records (no filter)
- [ ] Parameterized query: Filter added via `WHERE created_by_id = $1` (not string concat)
- [ ] Unit test: `TestTrainerCannotFetchOtherTrainerRecords`, `TestAdminSeesAllRecords`

### File Upload Security

- [ ] MIME type validation: Only JPEG, PNG, WebP (photos); MP4, MOV, WebM (videos)
- [ ] File size: Reject >10MB photos, >50MB videos before storage
- [ ] EXIF stripping: Metadata removed via library (not manual regex)
- [ ] Virus scanning: Optional in Phase 1 (flag for Phase 2 if needed)
- [ ] URL signing: Garage signed URLs expire after 1 hour

### Secret Management

- [ ] `.env` never committed (checked by `gitleaks detect --redact`)
- [ ] SMTP password: Environment variable only, never in logs
- [ ] JWT secret: In `.env`, rotatable
- [ ] OTP codes: Never logged (only "sent to user@email.com")
- [ ] `gitleaks` catches any hardcoded password patterns (test: add fake secret, verify gitleaks blocks it)

---

## Guard: What to Watch in Diffs

✋ **Red flags in code review:**

1. **A threshold moved** — Check if `CONSTRAINTS.md` was edited to relax a rule instead of fixing the code
2. **A test got easier** — `.skip`, `.only`, or deleted test files without reason
3. **A checker silenced** — New suppression comment added (should have a reason)
4. **Work unfinished** — `panic("TODO")`, empty error handler, incomplete function
5. **Exception appeared** — New row in Exceptions table without discussion

**Security-specific red flags:**

6. **String concatenation in SQL** — `query := "SELECT * FROM users WHERE id = " + id` → Should be parameterized
7. **Password/secret in code** — `password := "admin123"` → Must be in `.env` with comment explaining why
8. **OTP code logged** — `log.Printf("OTP: %s", code)` → Only log delivery status
9. **File upload without validation** — Direct `io.Copy` from request → Must validate MIME + size first
10. **Missing EXIF stripping** — Photos stored with metadata → Must use library (piexifgo, exiftool)
11. **No rate limiting** → POST endpoints need rate limit checks (OTP: 5/15min)
12. **Direct SQL query from user input** — Even with good intentions → Always use prepared statements
13. **Semgrep findings ignored** → `# nosemgrep: rule-id` without comment explaining safe case

Tightening the bar is silent. **Loosening it should be loud** (discussed, documented).

---

## Notes for Agents

Read this file before writing code. Do not weaken constraints to make a change pass. If a constraint blocks progress, flag it in the commit message with the reason, and update the file explicitly (it will show in review).

**Performance constraint scope:** You don't need to benchmark every endpoint. For CRUD operations, response time is "good" if it feels instant locally. Run detailed checks when the task explicitly involves optimization or you notice slowness.

**Coverage scope:** Changed lines, not project total. If you edit 50 lines, those 50 need 80% test coverage (40 lines tested). Project coverage is tracked separately and will be enforced after Phase 1.

---

## See Also

- `SPEC.md` — What to build
- `CLAUDE.md` (to be written) — How to build it (code style, conventions)
- `README.md` — Setup and deployment
