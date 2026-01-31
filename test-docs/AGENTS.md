# AGENTS.md

Automated AGENTS.md generation and distribution system for AI coding tools, using GitLab API for discovery and Cloudflare Workers for the dashboard.

**Tech Stack:** TypeScript, Bun, React, Cloudflare Workers + D1
**Package Manager:** bun (root, frontend) | pnpm for wrangler CLI

## Directory Structure

```
src/               - CLI and core logic (discovery, GitLab client, database)
dashboard/worker/  - Cloudflare Worker API backend (D1 database)
dashboard/frontend/- React + Vite dashboard UI
data/              - Local SQLite database and sync files (gitignored)
```

## Commands

```bash
# Install dependencies
bun install

# CLI operations
bun run cli discover --group cloudflare/devtools  # Discover eligible repos
bun run cli sync                                   # Export SQLite → D1 sync SQL

# Testing
bun test ./src/                    # CLI/core tests
bun run test:all                   # All tests (CLI + worker + frontend)
bun run test:worker                # Worker API tests (vitest)
bun run test:frontend              # Frontend tests (vitest)

# Development
bun run dev:frontend               # Start frontend dev server
bun run dev:worker                 # Start worker dev server

# Type checking
bun run typecheck                  # Root project
cd dashboard/worker && bun run typecheck
cd dashboard/frontend && bun run lint
```

## Environment Variables

```bash
GITLAB_TOKEN       # Required: GitLab API token
GITLAB_HOST        # Required: GitLab host URL
```

## Local Development Workflow

### Setting up the dashboard locally
```bash
# 1. Delete existing local data (clean slate)
rm -f data/agents.db
rm -rf dashboard/worker/.wrangler

# 2. Run discovery to populate SQLite
bun run cli discover --group cloudflare/devtools

# 3. Generate sync SQL and apply to local D1
bun run cli sync
cd dashboard/worker
pnpm dlx wrangler d1 execute agents-md-dashboard --local --file=schema.sql
pnpm dlx wrangler d1 execute agents-md-dashboard --local --file=../../data/sync.sql

# 4. Start worker (terminal 1) and frontend (terminal 2)
bun run dev          # in dashboard/worker/
bun run dev          # in dashboard/frontend/
```

## Code Patterns

- **GitLab client:** See `src/gitlab/client.ts` for API patterns using @gitbeaker/rest
- **Eligibility logic:** `src/discovery/eligibility.ts` defines all rules
- **Worker API:** REST endpoints at `/api/stats`, `/api/groups`, `/api/repositories`
- **Frontend hooks:** Custom hooks in `dashboard/frontend/src/hooks/` (useUrlState, useDebounce, useTheme)
- **Tests colocated:** `__tests__/` directories next to source files

## Boundaries

✅ **Always:** Run `bun test ./src/` before committing CLI changes
✅ **Always:** Use wrangler.jsonc (not wrangler.toml) for worker config
⚠️ **Requires Approval:** Deploying worker or creating GitLab MRs
🚫 **Never:** Commit data/ directory contents (gitignored)
🚫 **Never:** Use npm or yarn (this project uses bun)
