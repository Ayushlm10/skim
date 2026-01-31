# AGENTS.md Distribution System - Technical Plan

## Overview

This document describes HOW we'll build the system. For WHAT and WHY, see `spec.md`.

---

## Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Language | TypeScript | Type safety, ecosystem, AI integration |
| Runtime | Bun | Fast, native TypeScript |
| GitLab API | `@gitbeaker/rest` | Mature, typed client |
| Concurrency | `p-limit` | Simple rate limiting |
| Data Store | SQLite | Portable, no infra needed |
| Dashboard | React + Vite + Cloudflare Workers + D1 | Fast dev, easy deploy |
| Generation | OpenCode with custom subagent | AI-powered, flexible |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                              CLI                                 │
│         discover | generate | deliver | status                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │
│  │   Discovery  │───▶│  Generation  │───▶│   Delivery   │       │
│  │  (GitLab API)│    │  (OpenCode)  │    │  (GitLab API)│       │
│  └──────────────┘    └──────────────┘    └──────────────┘       │
│         │                   │                   │                │
│         └───────────────────┴───────────────────┘                │
│                             │                                    │
│                      ┌──────────────┐                           │
│                      │   SQLite DB  │                           │
│                      └──────────────┘                           │
│                             │                                    │
│                      ┌──────────────┐                           │
│                      │  Dashboard   │                           │
│                      │  (D1 sync)   │                           │
│                      └──────────────┘                           │
└─────────────────────────────────────────────────────────────────┘
```

---

## Key Design Decisions

### 1. OpenCode for Generation (not custom extractors)

Instead of building language-specific extractors, we use OpenCode with a specialized `agents-md-generator` subagent. Benefits:
- AI can read and understand any project type
- Easier to iterate (change prompt, not code)
- Already configured with AI Gateway

Subagent location: `~/.config/opencode/agent/agents-md-generator.md`

### 2. GitLab API for Everything (no local clones for discovery)

- Discovery: Use GitLab API to check files exist
- Generation: Clone repo → run OpenCode → cleanup
- Delivery: Use GitLab API to create branch/commit/MR

### 3. Cleanup After Processing

Cloned repos are deleted after MR creation to avoid disk space issues. Use `--keep-repo` flag for debugging.

### 4. Reference: ci_component.go Patterns

We port these patterns from the existing Go tool:
- CLI flags: `--group`, `--repo`, `--dry-run`, `--limit`
- Concurrent processing with rate limiting
- Progress reporting: `[1/100] ✅ repo-name: message`
- Summary at end with counts

### 5. Testing Requirements

Every feature implementation must include corresponding unit tests:
- Use Bun's built-in test runner (`bun test`)
- Mock external dependencies (GitLab API) for fast, reliable tests
- Test edge cases and error conditions
- Run `bun test` before committing to ensure no regressions

---

## Project Structure

```
agents-md-distribution/
├── src/
│   ├── config.ts              # Configuration
│   ├── types.ts               # Type definitions
│   ├── gitlab/                # GitLab API integration
│   ├── discovery/             # Eligibility checking
│   │   └── __tests__/         # Unit tests for discovery
│   ├── generation/            # OpenCode runner + validator
│   ├── delivery/              # MR creation
│   ├── db/                    # SQLite operations
│   ├── runner.ts              # Concurrent processing
│   └── cli.ts                 # CLI entry point
├── dashboard/
│   ├── worker/                # Cloudflare Worker backend
│   └── frontend/              # React frontend
├── data/                      # Local data (gitignored)
├── spec.md
└── plan.md
```

---

## Phases & Tasks

### Phase 1: Discovery + Dashboard

**Goal:** See how many repos we need to work with

| Task | Description | Acceptance Criteria |
|------|-------------|---------------------|
| 1.1 | Project setup | `bun install` works, TypeScript compiles |
| 1.2 | GitLab client | Can authenticate and fetch current user |
| 1.3 | Project fetching | Fetch projects with pagination, supports `--group`, `--repo`, `--limit` |
| 1.4 | File checking | Check if AGENTS.md/CLAUDE.md exist via API |
| 1.5 | Eligibility logic | All rules: activity, archived, size, existing files |
| 1.6 | Database | SQLite schema and CRUD operations |
| 1.7 | Discovery command | `agents-md discover --group devtools` works |
| 1.8 | Runner infrastructure | Concurrent processing with progress output |
| 1.9 | Dashboard backend | Worker with /api/stats, /api/repositories, /api/groups |
| 1.10 | Dashboard frontend | Summary page, repository list, group breakdown |
| 1.10a | URL state sync | Filters/pagination persisted in URL for shareable links |
| 1.10b | Text search | Search repositories by name/path with debounced input |
| 1.10c | Click group to filter | Click namespace in GroupsTable to filter RepositoriesList |
| 1.10d | Data visualizations | Charts for eligibility breakdown and ineligibility reasons |
| 1.10e | Theme support | Dark/light/system mode toggle with localStorage persistence |
| 1.10f | Dashboard UX improvements | See details below |
| 1.11 | Dashboard deployment | Deployed behind Access |
| 1.12 | Hierarchical group navigation | Sidebar tree view for 6K+ repos at scale |

**Deliverable:** Know exactly how many eligible repos exist, with breakdown

---

### Phase 2: Generation

**Goal:** Generate high-quality AGENTS.md files

| Task | Description | Acceptance Criteria |
|------|-------------|---------------------|
| 2.1 | Test subagent | Run on TS, Go, Terraform, Python repos manually |
| 2.2 | OpenCode runner | Invoke CLI, parse JSON output |
| 2.3 | Validator | Check size limits, required sections |
| 2.4 | Generate command | `agents-md generate --repo path` works |

**Deliverable:** Can generate quality AGENTS.md for any repo type

---

### Phase 3: Delivery

**Goal:** Create MRs at scale

| Task | Description | Acceptance Criteria |
|------|-------------|---------------------|
| 3.1 | MR operations | Check existing MR, delete orphaned branches |
| 3.2 | Commit operations | Create branch + commit via API |
| 3.3 | MR template | Clear description explaining AGENTS.md |
| 3.4 | Delivery command | `agents-md deliver --group devtools --dry-run` works |
| 3.5 | Repo cleanup | Delete cloned repo after MR creation |
| 3.6 | Status command | `agents-md status` shows MR breakdown |

**Deliverable:** Can create MRs for all eligible repos with tracking

---

### Phase 4: Pilot & Maintenance

**Goal:** Validate with real users, set up ongoing maintenance

| Task | Description | Acceptance Criteria |
|------|-------------|---------------------|
| 4.1 | Run pilot | Create MRs for devtools group |
| 4.2 | Gather feedback | Track merge rate, modifications, comments |
| 4.3 | Reviewer enhancement | Update CODE_REVIEWER.md to suggest AGENTS.md updates |

**Deliverable:** Pilot complete, maintenance integrated

---

## CLI Reference

```bash
# Discovery
agents-md discover --group devtools          # Discover eligible repos
agents-md discover --group cloudflare --limit 100
agents-md discover --dry-run                 # Preview without DB writes

# Generation
agents-md generate --repo cloudflare/devtools/opencode
agents-md generate --group devtools
agents-md generate --preview                 # Show output, don't save

# Delivery
agents-md deliver --group devtools --dry-run # Preview MRs
agents-md deliver --group devtools           # Create MRs
agents-md deliver --keep-repo                # Don't cleanup cloned repos

# Status
agents-md status                             # Show MR status breakdown
```

---

## Current Status

| Phase | Status | Notes |
|-------|--------|-------|
| Subagent | ✅ Created | `~/.config/opencode/agent/agents-md-generator.md` |
| Phase 1: Discovery | ✅ Done | All tasks complete, deployment ready |
| Phase 2: Generation | ⏳ Next | Can start now |
| Phase 3: Delivery | ⏳ Blocked | Needs Phase 2 |
| Phase 4: Pilot | ⏳ Blocked | Needs Phase 3 |

### Phase 1 Task Progress

| Task | Status | Notes |
|------|--------|-------|
| 1.1 Project setup | ✅ Done | Bun, TypeScript, dependencies |
| 1.2 GitLab client | ✅ Done | Auth, typed client |
| 1.3 Project fetching | ✅ Done | Pagination, group/repo support |
| 1.4 File checking | ✅ Done | fileExists, getFileContent, getRepositoryTree |
| 1.5 Eligibility logic | ✅ Done | All rules: activity, archived, size, existing files |
| 1.6 Database | ✅ Done | SQLite schema and CRUD |
| 1.7 Discovery command | ✅ Done | Full command with DB persistence, --reset flag, DB stats display |
| 1.8 Runner infrastructure | ✅ Done | Concurrent processing with p-limit, --concurrency flag |
| 1.9 Dashboard backend | ✅ Done | Worker API with /api/stats, /api/groups, /api/repositories |
| 1.10 Dashboard frontend | ✅ Done | React + Vite + TanStack Query, industrial terminal aesthetic |
| 1.10a URL state sync | ✅ Done | useUrlState hook with useSyncExternalStore, 14 tests |
| 1.10b Text search | ✅ Done | Worker API search param, debounced search input, 6 API tests + 7 hook tests |
| 1.10c Click group filter | ✅ Done | GroupsTable onClick sets namespace URL param, 8 tests |
| 1.10d Data visualizations | ✅ Done | EligibilityChart (pie) + ReasonsChart (bar), 13 tests |
| 1.10e Theme support | ✅ Done | ThemeProvider + ThemeToggle, 29 tests (16 hook + 13 component) |
| 1.10f Dashboard UX improvements | ✅ Done | Primary metrics, chart interactivity, table sorting |
| 1.11 Dashboard deployment | ✅ Done | Code ready, manual deployment steps needed |
| 1.12 Hierarchical group navigation | ✅ Done | All sub-tasks complete |

---

## Task 1.12: Hierarchical Group Navigation

**Problem:** With 6,292 repos across 200+ namespaces, a flat GroupsTable is unusable. Users can't navigate effectively.

**Solution:** Replace flat table with a collapsible tree sidebar that shows 2 levels of hierarchy.

**Data Context:**
- Total repositories: 6,292
- Eligible: 4,168 (66%)
- Unique namespaces: 200+ (estimated)
- Depth distribution: Level 1 (3 repos), Level 2 (5,000+ repos), Level 3 (1,000+ repos)

### Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Location | Left sidebar | Always visible for navigation |
| Hierarchy depth | Top 2 levels, then flatten | Balance between organization and simplicity |
| Click behavior | Filter direct children only | Click `cloudflare/devtools` → show 144 repos directly in that namespace, not all 230 including subgroups |
| Expand/collapse | Per-node toggle | Progressive disclosure |
| Selected state | Visual highlight + URL sync | Shareable filtered views |

### Sub-tasks

| Task | Description | Status |
|------|-------------|--------|
| 1.12a | **API: Hierarchical groups endpoint** | ✅ Done |
|       | New `GET /api/groups/tree` endpoint that returns nested structure |
|       | Compute hierarchy from namespace paths on-the-fly |
|       | Return: `{ name, path, total, eligible, children: [...] }` |
| 1.12b | **Component: GroupTree** | ✅ Done |
|       | Recursive tree component with expand/collapse |
|       | Show counts (total/eligible) per node |
|       | Terminal-style aesthetic matching existing design |
|       | 20 tests covering rendering, expand/collapse, filtering, accessibility |
| 1.12c | **Layout: Sidebar + Main content** | ✅ Done |
|       | Update App.tsx to have sidebar layout |
|       | Sidebar: GroupTree (sticky, scrollable) |
|       | Main: StatsOverview, Charts, RepositoriesList |
|       | Remove old GroupsTable (replaced by sidebar) |
|       | Responsive design with collapsible sidebar on mobile |
|       | 13 tests for layout structure, sidebar toggle, accessibility |
| 1.12d | **Integration: URL state & filtering** | ✅ Done (in 1.12b) |
|       | Click tree node → set `namespace` URL param |
|       | Highlight selected node in tree |
|       | Clear filter button when namespace is set |
| 1.12e | **Testing** | ✅ Done |
|       | API tests for hierarchical endpoint (in 1.12a) |
|       | Component tests for GroupTree (20 tests in 1.12b) |
|       | Integration tests for filtering (included in component tests) |

### API Design: `/api/groups/tree`

```typescript
interface GroupTreeNode {
  name: string;          // "devtools"
  path: string;          // "cloudflare/devtools" (full path for filtering)
  total: number;         // Repos directly in this namespace
  eligible: number;      // Eligible repos directly in this namespace
  children: GroupTreeNode[];
}

// Response
{
  tree: GroupTreeNode[]  // Top-level groups
}
```

### UI Mockup

```
┌─────────────────────────────────────────────────────────────────────┐
│ Cloudflare / AGENTS.md                    Distribution Dashboard    │
├──────────────┬──────────────────────────────────────────────────────┤
│ GROUPS       │  Stats Overview                                      │
│ ─────────────│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐        │
│ ▼ cloudflare │  │ Total  │ │Eligible│ │ Merge  │ │Progress│        │
│   ▶ devtools │  │ 6,292  │ │ 4,168  │ │  Rate  │ │  0/4K  │        │
│   ▶ platform │  └────────┘ └────────┘ └────────┘ └────────┘        │
│   ▶ workers  │                                                      │
│   ▶ network  │  Charts                                              │
│   ▶ ...      │  [Eligibility Pie]  [Reasons Bar]                   │
│              │                                                      │
│              │  Repositories (filtered by selection)                │
│              │  ┌─────────────────────────────────────────────────┐ │
│              │  │ Name          │ Status │ Reason │ Activity     │ │
│              │  ├─────────────────────────────────────────────────┤ │
│              │  │ opencode      │ ✓      │        │ 2 days ago   │ │
│              │  │ ai-gateway    │ ✓      │        │ 1 week ago   │ │
│              │  └─────────────────────────────────────────────────┘ │
└──────────────┴──────────────────────────────────────────────────────┘
```

---

## Task 1.10f: Dashboard UX Improvements (Pre-deployment)

Planned enhancements to improve dashboard usability and align with spec.md success metrics.

### High Priority

| Improvement | Description | Status |
|-------------|-------------|--------|
| **Primary Metrics Alignment** | Replace "MRs Created" card with "Merge Rate" as primary metric (per spec.md success criteria). Add color coding: green (>70%), yellow (40-70%), red (<40%). Display MR Progress with merged count. | ✅ Done |
| **Chart Interactivity** | Make charts clickable to filter data. Click eligibility pie chart segment → filter repos by status. Click reasons bar → filter repos by that reason. Sync with URL state for shareable links. | ✅ Done |
| **Repository Table Sorting** | Add column sorting to repositories table (currently only groups table has sorting). Sort by name, status, reason, activity. Persist sort preference in URL. | ✅ Done |

### Medium Priority

| Improvement | Description |
|-------------|-------------|
| **Progress/Funnel Visualization** | Add rollout progress section showing conversion funnel: Eligible → Generated → MR Created → MR Merged. Display end-to-end conversion rate prominently. |
| **Migration Tracking** | Add category tabs per spec.md: Generate (no AGENTS.md, no CLAUDE.md), Migrate (has CLAUDE.md, no AGENTS.md), Skip (has AGENTS.md). Surface migration candidates distinctly. |
| **Actionability & Quick Actions** | Add row-level actions: Preview generated AGENTS.md, trigger generation, view in GitLab. Add bulk actions for selected repositories. |
| **Visual Hierarchy Polish** | Add section dividers, improve disabled select styling (show it's intentionally disabled), add show more/less toggle to groups table. |

### Lower Priority

| Improvement | Description |
|-------------|-------------|
| **Keyboard Navigation** | Add power user shortcuts: `/` to search, `j`/`k` to navigate, `?` for help modal. |

---

## Next Step

**Phase 1: Discovery - ✅ COMPLETE**

All code changes are complete. To deploy the dashboard:

```bash
# 1. Create D1 database
cd dashboard/worker
npx wrangler d1 create agents-md-dashboard
# Copy database_id and update wrangler.jsonc

# 2. Apply schema
npx wrangler d1 execute agents-md-dashboard --file=schema.sql

# 3. Build frontend (if needed)
cd ../frontend && bun run build

# 4. Deploy worker
cd ../worker && npx wrangler deploy

# 5. Configure Cloudflare Access (manual via dashboard)

# 6. Sync data to D1
cd ../.. && bun run cli sync
cd dashboard/worker
npx wrangler d1 execute agents-md-dashboard --file=../../data/sync.sql
```

**Next:**
- Phase 2: Generation - Start with Task 2.1 (Test subagent on various repo types)
