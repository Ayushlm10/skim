# Local MD Viewer - Implementation Plan

## Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| TUI Framework | Bubble Tea | Elm-style Model/Update/View loop |
| Components | Bubbles | List, viewport, text input |
| Styling | Lip Gloss | Colors, borders, layout |
| Markdown | Glamour | Terminal markdown rendering |
| File Watching | fsnotify | Cross-platform FS events |
| Language | Go 1.21+ | Performance, single binary |

## Project Structure

```
local-md/
├── main.go                     # Entry point, CLI parsing
├── go.mod                      # Module definition
├── go.sum                      # Dependency checksums
│
├── internal/
│   ├── app/
│   │   ├── model.go            # Main Bubble Tea model
│   │   ├── update.go           # Update logic (key handling, messages)
│   │   ├── view.go             # View rendering
│   │   └── messages.go         # Custom message types
│   │
│   ├── components/
│   │   ├── filetree/
│   │   │   ├── filetree.go     # File tree component
│   │   │   ├── item.go         # Tree item (file/directory)
│   │   │   ├── delegate.go     # Custom list delegate
│   │   │   └── scanner.go      # Directory scanning logic
│   │   │
│   │   ├── preview/
│   │   │   ├── preview.go      # Markdown preview component
│   │   │   └── renderer.go     # Glamour wrapper
│   │   │
│   │   └── statusbar/
│   │       └── statusbar.go    # Bottom status bar
│   │
│   ├── styles/
│   │   └── styles.go           # Centralized Lip Gloss styles
│   │
│   └── watcher/
│       └── watcher.go          # File system watcher
│
├── specs/
│   ├── design.md               # Design specification
│   ├── plan.md                 # This file
│   └── implementation.md       # Implementation status
│
└── README.md                   # User documentation
```

## Implementation Phases

### Phase 1: Foundation (Core Structure)

**Goal:** Basic working application with dual-panel layout

**Tasks:**
1. Initialize Go module with dependencies
2. Create main.go with CLI argument parsing
3. Implement main model with focus management
4. Create basic file tree (flat list, no tree structure yet)
5. Create basic preview with hardcoded content
6. Implement panel layout with Lip Gloss
7. Handle window resize events

**Deliverable:** App launches, shows two panels, can quit with `q`

### Phase 2: File Tree Component

**Goal:** Navigable file tree with directory support

**Tasks:**
1. Implement directory scanner (recursive, .md only)
2. Create tree item struct (path, name, depth, isDir, expanded)
3. Implement custom list delegate for tree rendering
4. Handle expand/collapse for directories
5. Track selected item and notify parent
6. Implement basic keyboard navigation

**Deliverable:** Navigate directories, see tree structure, select files

### Phase 3: Markdown Preview

**Goal:** Beautiful markdown rendering in viewport

**Tasks:**
1. Integrate Glamour renderer with auto style
2. Create preview component with viewport
3. Implement file loading on selection
4. Handle word wrapping based on width
5. Implement scroll with Page Up/Down
6. Show scroll position indicator

**Deliverable:** Select file in tree, see rendered markdown

### Phase 4: Filter/Search

**Goal:** Fuzzy search through file list

**Tasks:**
1. Enable built-in list filtering
2. Implement `/` to enter filter mode
3. Handle Esc to exit filter mode
4. Style the filter input
5. Update status bar to show filter state

**Deliverable:** Type `/spec` to filter to files containing "spec"

### Phase 5: File Watching

**Goal:** Auto-refresh when files change

**Tasks:**
1. Implement fsnotify watcher wrapper
2. Create watch command that returns FileChangedMsg
3. Watch currently selected file
4. Debounce rapid changes (250ms)
5. Re-render preview on change
6. Handle watcher errors gracefully

**Deliverable:** Edit file externally, preview updates automatically

### Phase 6: Polish & Styling

**Goal:** Production-ready visual quality

**Tasks:**
1. Refine color palette (adaptive light/dark)
2. Implement focus indicators
3. Add status bar with file info (path, size, modified)
4. Create help overlay (toggle with `?`)
5. Handle empty states gracefully
6. Add loading indicator for large files
7. Error handling and user feedback

**Deliverable:** Polished, professional-looking application

## Key Bindings Specification

| Key | Context | Action |
|-----|---------|--------|
| `↑` / `k` | File tree | Move selection up |
| `↓` / `j` | File tree | Move selection down |
| `Enter` | File tree | Open file / Toggle directory |
| `Tab` | Global | Switch focus between panels |
| `/` | File tree | Enter filter mode |
| `Esc` | Filter mode | Exit filter, clear filter |
| `Esc` | Normal | Clear selection / Reset view |
| `PgUp` / `Ctrl+u` | Preview | Scroll up |
| `PgDn` / `Ctrl+d` | Preview | Scroll down |
| `g` | Preview | Go to top |
| `G` | Preview | Go to bottom |
| `?` | Global | Toggle help overlay |
| `q` / `Ctrl+c` | Global | Quit application |

## Message Types

```go
// Custom messages for the application
type (
    // File tree events
    FileSelectedMsg    struct{ Path string }
    DirectoryToggleMsg struct{ Path string }
    
    // File operations
    FileLoadedMsg   struct{ Path, Content string }
    FileErrorMsg    struct{ Path string; Err error }
    FileChangedMsg  struct{ Path string }
    
    // Watcher
    WatchStartMsg struct{ Path string }
    WatchStopMsg  struct{}
    
    // UI state
    FocusChangedMsg struct{ Panel string }
    FilterActiveMsg struct{ Active bool }
)
```

## Dependencies

```go
module github.com/user/local-md

go 1.21

require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/bubbles v0.18.0
    github.com/charmbracelet/lipgloss v0.10.0
    github.com/charmbracelet/glamour v0.6.0
    github.com/fsnotify/fsnotify v1.7.0
)
```

## Testing Strategy

### Unit Tests
- Directory scanner (various directory structures)
- Tree item rendering
- Message handling in update functions

### Integration Tests
- Full app initialization
- Navigation flow
- File selection and preview

### Manual Testing Checklist
- [ ] Launch with no arguments (current directory)
- [ ] Launch with path argument
- [ ] Navigate deep directory structure
- [ ] Filter files with various patterns
- [ ] View large markdown files
- [ ] Edit file and verify auto-refresh
- [ ] Resize terminal window
- [ ] Test in light and dark terminals

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Large directories slow scanning | Lazy loading, limit depth |
| Large files slow rendering | Truncate/paginate, show loading |
| Rapid file changes overwhelm | Debounce watcher events |
| Terminal doesn't support colors | Graceful degradation via Lip Gloss |
| Non-UTF8 files | Detect and show error message |

## Timeline Estimate

| Phase | Estimated Time |
|-------|----------------|
| Phase 1: Foundation | 1-2 hours |
| Phase 2: File Tree | 2-3 hours |
| Phase 3: Preview | 1-2 hours |
| Phase 4: Filter | 1 hour |
| Phase 5: Watching | 1-2 hours |
| Phase 6: Polish | 2-3 hours |
| **Total** | **8-13 hours** |
