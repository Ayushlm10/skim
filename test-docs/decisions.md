# Architecture Decisions

This document records important architectural and implementation decisions made during development.

---

## Decision 1: Use `@gitbeaker/rest` for GitLab API

**Date:** 2025-01-19  
**Status:** Accepted

### Context
Need a typed GitLab API client for Node.js/Bun.

### Options Considered
1. `@gitbeaker/rest` - Mature, typed, well-maintained
2. Raw fetch with custom types - More control but more work
3. `gitlab` npm package - Less maintained

### Decision
Use `@gitbeaker/rest` v43.x for its TypeScript support and comprehensive API coverage.

### Consequences
- API responses use camelCase (transformed from snake_case)
- Need to handle both formats in mappers for flexibility
- Good pagination support built-in

---

## Decision 2: Group paths must be fully qualified

**Date:** 2025-01-19  
**Status:** Accepted

### Context
GitLab API requires full group paths (e.g., `cloudflare/devtools`) not just group names (`devtools`).

### Decision
- Added `groups` command to help users find correct group paths
- CLI shows helpful error messages when groups not found
- Examples in help text use full paths

### Consequences
- Better UX with group discovery
- Users need to know or discover the full path before running `discover`

---

## Decision 3: Flexible project field mapping

**Date:** 2025-01-19  
**Status:** Accepted

### Context
GitBeaker can return both camelCase (transformed) and snake_case (raw) field names depending on the API endpoint and version.

### Decision
`mapProject()` handles both formats by checking for either field name:
```typescript
const path = project.pathWithNamespace ?? project.path_with_namespace ?? "";
```

### Consequences
- More robust against API inconsistencies
- Slightly more verbose mapping code
- Works with any gitbeaker version

---

## Decision 4: SQLite for local state with Bun's built-in SQLite

**Date:** 2025-01-19  
**Status:** Accepted

### Context
Need to track repository discovery state, eligibility, MR status across runs.

### Decision
Use SQLite via Bun's built-in `bun:sqlite` module, stored in `data/` directory (gitignored):
- Repository state tracking with full schema
- Eligibility results caching
- MR status tracking
- WAL mode enabled for better concurrent access
- In-memory database option for fast testing

### Implementation
Located in `src/db/`:
- `schema.ts` - Table definitions and indexes
- `database.ts` - CRUD operations
- `index.ts` - Module exports
- `__tests__/database.test.ts` - 66 unit tests

### Consequences
- No external database or npm dependencies needed (Bun has native SQLite)
- Portable - can move/backup easily
- Will sync to D1 for dashboard (Phase 1.9)
- Fast tests using in-memory database

---

## Decision 5: Concurrent processing with p-limit

**Date:** 2025-01-19  
**Status:** Accepted

### Context
Need to process many repositories without overwhelming GitLab API.

### Decision
Use `p-limit` for concurrency control, defaulting to 50 concurrent requests (matching Go reference implementation pattern).

### Implementation
Located in `src/runner.ts`:
- `runConcurrent()` - returns results in completion order
- `runConcurrentOrdered()` - returns results in original input order
- `formatDuration()` and `calculateRate()` utilities

The `checkEligibilityBatch()` function now uses `runConcurrentOrdered()` to process repos in parallel while preserving result order. CLI supports `--concurrency` flag.

### Consequences
- Predictable API load
- Progress reporting works correctly (callbacks fire as items complete)
- Easy to tune via `--concurrency` flag (default: 50)
- Results maintain original order for database persistence

---

## Decision 6: Eligibility thresholds and code file detection

**Date:** 2025-01-19  
**Status:** Accepted

### Context
Need to define what makes a repository "eligible" for AGENTS.md generation.

### Decision
Implemented these thresholds in `src/discovery/eligibility.ts`:

| Rule | Threshold | Rationale |
|------|-----------|-----------|
| Activity | 6 months | Focus on active repos where AI tooling matters |
| Review candidates | 6-12 months | Flagged but currently marked inactive |
| Minimum files | >5 files | Skip trivial/placeholder repos |
| Code detection | Extension-based | Recognize common languages |

Code extensions recognized: `.ts`, `.js`, `.py`, `.go`, `.rs`, `.java`, `.kt`, `.c`, `.cpp`, `.cs`, `.rb`, `.php`, `.swift`, `.sh`, `.tf`, `.yaml`, plus `Dockerfile`, `Makefile`.

### Consequences
- Repos with only docs/config files are excluded (no_code)
- 6-12 month repos treated same as >12 month for now (could add "review" status later)
- Extension list may need expansion for niche languages

---

## Decision 7: Testing strategy with mocked dependencies

**Date:** 2025-01-19  
**Status:** Accepted

### Context
Need reliable, fast tests that don't depend on external services (GitLab API).

### Decision
- Use Bun's built-in test runner (`bun test`)
- Mock GitLab client responses for unit tests
- Tests live in `__tests__/` directories alongside source code
- Every feature must have corresponding tests before commit

### Test Structure
```
src/
├── discovery/
│   ├── eligibility.ts
│   ├── index.ts
│   └── __tests__/
│       ├── eligibility.test.ts    # Core logic tests
│       └── code-detection.test.ts # Code file detection
```

### Consequences
- Fast test execution (~50ms for 34 tests)
- No flaky tests from network issues
- Can run tests in CI without GitLab credentials
- Need to keep mocks in sync with actual API behavior

---

## Decision 8: Real-time database persistence during discovery

**Date:** 2025-01-19  
**Status:** Accepted

### Context
When processing many repositories, we need to decide when to persist results to the database - all at once at the end, or as each repository is processed.

### Options Considered
1. Batch persistence at end - Simpler, but loses progress if interrupted
2. Real-time persistence - More resilient, slightly more I/O

### Decision
Persist each repository to the database immediately after eligibility check via the `onProgress` callback:
```typescript
onProgress: (completed, total, checkResult) => {
  const status = checkResult.result.eligible ? "eligible" : "ineligible";
  upsertRepository(db, checkResult.repo, status, checkResult.result.reason);
}
```

### Consequences
- Progress is preserved if discovery is interrupted (Ctrl+C, error, etc.)
- Can resume or re-run discovery - upsert handles existing records
- Slightly more database writes but SQLite with WAL mode handles this well
- `--reset` flag added to clear database before discovery when a fresh start is needed

---

## Decision 9: Dashboard backend with D1 and sync mechanism

**Date:** 2025-01-19  
**Status:** Accepted

### Context
Need to serve dashboard data via API. Local SQLite works for CLI but needs to be accessible from a web dashboard.

### Options Considered
1. Direct SQLite access via API - Requires server with file access
2. D1 with manual sync - Cloudflare Workers + D1, sync data from local SQLite
3. Separate PostgreSQL/MySQL - More infrastructure to manage

### Decision
Use Cloudflare Workers with D1 database, synced from local SQLite:

- **Worker API**: `dashboard/worker/` with endpoints for stats, groups, repositories
- **D1 Schema**: Matches SQLite schema exactly for seamless sync
- **Sync Mechanism**: `agents-md sync` CLI command exports SQL INSERT statements
- **Testing**: Use `@cloudflare/vitest-pool-workers` for realistic D1 testing

### API Endpoints
| Endpoint | Description |
|----------|-------------|
| `GET /api/stats` | Overall statistics (total, eligible, by reason) |
| `GET /api/groups` | Statistics grouped by namespace |
| `GET /api/repositories` | List with filtering (status, namespace, reason) and pagination |
| `GET /api/repositories/:id` | Single repository details |

### Sync Workflow
```bash
# After discovery
agents-md discover --group cloudflare/devtools

# Export to SQL
agents-md sync

# Import to D1
wrangler d1 execute agents-md-dashboard --file=data/sync.sql
```

### Consequences
- Decoupled: CLI works locally, dashboard works on edge
- Simple sync: Just SQL INSERT statements, no complex migration
- Testable: Vitest with workers pool provides realistic D1 environment
- Scalable: D1 handles read traffic, no server to manage

---

## Future Decisions to Make

- [ ] How to handle repos with CLAUDE.md (skip for now, decide post-pilot)
- [ ] Branch naming convention for MRs
- [ ] MR labels strategy
- [ ] Retry strategy for API failures
