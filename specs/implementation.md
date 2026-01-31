# Local MD Viewer - Implementation Status

## Current Status: Phase 2 Complete

Last Updated: 2026-01-31

---

## Phase 1: Foundation

**Status:** Complete

| Task | Status | Notes |
|------|--------|-------|
| Initialize Go module | Done | `go mod init github.com/athakur/local-md` |
| CLI argument parsing | Done | Accepts path arg, defaults to current dir |
| Main model structure | Done | `internal/app/model.go` with panel focus |
| Basic file tree (flat) | Done | Placeholder with sample structure |
| Basic preview panel | Done | Placeholder welcome content |
| Panel layout | Done | 25/75 split with Lip Gloss |
| Window resize handling | Done | Responds to `tea.WindowSizeMsg` |
| Quit functionality | Done | `q` or `Ctrl+C` to quit |

**Blockers:** None

---

## Phase 2: File Tree Component

**Status:** Complete

| Task | Status | Notes |
|------|--------|-------|
| Directory scanner | Done | `internal/components/filetree/scanner.go` - scans for .md files |
| Tree item struct | Done | `internal/components/filetree/item.go` - with depth tracking |
| Custom list delegate | Done | `internal/components/filetree/delegate.go` - custom rendering |
| Expand/collapse logic | Done | Toggle directories with Enter key |
| Selection tracking | Done | Via bubbles list component |
| Keyboard navigation | Done | j/k, arrows, Enter to toggle/select |

**Blockers:** None

---

## Phase 3: Markdown Preview

**Status:** Not Started

| Task | Status | Notes |
|------|--------|-------|
| Glamour integration | Pending | |
| Preview component | Pending | |
| File loading | Pending | |
| Word wrapping | Pending | |
| Scroll functionality | Pending | |
| Scroll indicator | Pending | |

**Blockers:** None (Phase 2 complete)

---

## Phase 4: Filter/Search

**Status:** Not Started

| Task | Status | Notes |
|------|--------|-------|
| Enable list filtering | Pending | |
| Filter mode activation | Pending | |
| Filter mode exit | Pending | |
| Filter input styling | Pending | |
| Status bar updates | Pending | |

**Blockers:** None (Phase 2 complete)

---

## Phase 5: File Watching

**Status:** Not Started

| Task | Status | Notes |
|------|--------|-------|
| fsnotify wrapper | Pending | |
| FileChangedMsg | Pending | |
| Watch current file | Pending | |
| Debounce logic | Pending | |
| Re-render on change | Pending | |
| Error handling | Pending | |

**Blockers:** Depends on Phase 3

---

## Phase 6: Polish & Styling

**Status:** Not Started

| Task | Status | Notes |
|------|--------|-------|
| Color palette refinement | Pending | |
| Focus indicators | Pending | Done (basic) - accent border on focused panel |
| Status bar info | Pending | Done (basic) - shows key hints |
| Help overlay | Pending | |
| Empty states | Pending | Done (basic) - placeholder messages |
| Loading indicators | Pending | |
| Error messages | Pending | |

**Blockers:** Depends on all previous phases

---

## Files Created

| File | Status | Description |
|------|--------|-------------|
| `specs/design.md` | Done | Design specification |
| `specs/plan.md` | Done | Implementation plan |
| `specs/implementation.md` | Done | This file |
| `main.go` | Done | Entry point with CLI parsing |
| `go.mod` | Done | Module definition |
| `go.sum` | Done | Dependency checksums |
| `internal/app/model.go` | Done | Main model with panel management |
| `internal/app/update.go` | Done | Update logic and key handling |
| `internal/app/view.go` | Done | View rendering for all panels |
| `internal/app/messages.go` | Done | Message types for app communication |
| `internal/styles/styles.go` | Done | Centralized Lip Gloss styles |
| `internal/components/filetree/filetree.go` | Done | File tree component |
| `internal/components/filetree/item.go` | Done | Tree items |
| `internal/components/filetree/delegate.go` | Done | List delegate |
| `internal/components/filetree/scanner.go` | Done | Dir scanner |
| `internal/components/preview/preview.go` | Pending | Preview component |
| `internal/components/preview/renderer.go` | Pending | Glamour wrapper |
| `internal/components/statusbar/statusbar.go` | Pending | Status bar component |
| `internal/watcher/watcher.go` | Pending | File watcher |

---

## Known Issues

None yet.

---

## Decisions Made

| Decision | Rationale | Date |
|----------|-----------|------|
| Use Charm stack | Best-in-class Go TUI libraries | 2026-01-31 |
| Minimal/editorial aesthetic | Clean, professional look | 2026-01-31 |
| Accept CLI path argument | More flexible for different workflows | 2026-01-31 |
| Show helpful message for empty dirs | Better UX than failing | 2026-01-31 |
| 25/75 panel split | More space for content, less for nav | 2026-01-31 |
| Inline rendering in view.go | Simpler for Phase 1; will refactor to components in Phase 2+ | 2026-01-31 |

---

## Next Steps

1. Integrate Glamour for markdown rendering
2. Create preview component with viewport scrolling
3. Load file content on selection
4. Wire file tree selection to preview panel

---

## Change Log

### 2026-01-31
- Created initial specification documents
- Defined project structure
- Planned implementation phases
- **Phase 1 Complete**: Basic dual-panel TUI with placeholder content
  - CLI argument parsing with path validation
  - Main model with focus management
  - Dual-panel layout (25/75 split)
  - Minimal/editorial styling with Lip Gloss
  - Keyboard handling for navigation, Tab switching, quit
  - Status bar with key hints
- **Phase 2 Complete**: File tree component with expand/collapse
  - Directory scanner that filters for .md files only
  - Tree item struct with depth tracking and parent references
  - Custom list delegate for tree-style rendering
  - Expand/collapse directories with Enter key
  - Lazy loading of children on expand
  - Built on bubbles list component for filtering support
