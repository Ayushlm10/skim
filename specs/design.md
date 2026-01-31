# Local MD Viewer - Design Specification

## Overview

A terminal-based markdown viewer built with the Charm Go TUI stack, designed for developers who work extensively with markdown files for planning, specifications, and documentation.

## Problem Statement

Developers often maintain numerous markdown files for:
- Project specifications and requirements
- Architecture Decision Records (ADRs)
- Meeting notes and planning documents
- Technical documentation

Existing solutions require either:
- Opening a full IDE/editor
- Using web-based preview tools
- Switching context away from the terminal

## Solution

A lightweight, fast, terminal-native markdown viewer that:
- Lives where developers already work (the terminal)
- Provides instant, beautiful markdown rendering
- Enables quick navigation through document collections
- Watches for changes and updates automatically

## User Experience

### Visual Layout

```
┌─────────────────────────────────────────────────────────────────┐
│  Local MD Viewer                                    ~/docs      │
├──────────────────────┬──────────────────────────────────────────┤
│                      │                                          │
│  ▾ docs/             │  # Project Specification                 │
│    ├── spec.md  ◀    │                                          │
│    ├── plan.md       │  This document outlines the key          │
│    └── decisions/    │  requirements and decisions for...       │
│        ▸ adr-001.md  │                                          │
│        ▸ adr-002.md  │  ## Goals                                │
│                      │                                          │
│  ▾ notes/            │  - Fast local viewing                    │
│    └── meeting.md    │  - Clean typography                      │
│                      │  - Minimal resource usage                │
│                      │                                          │
│  ┌────────────────┐  │  ## Implementation                       │
│  │ Filter: spec   │  │                                          │
│  └────────────────┘  │  The implementation follows...           │
│                      │                                          │
├──────────────────────┴──────────────────────────────────────────┤
│  ↑↓ navigate  │  ⏎ open  │  / filter  │  Tab switch  │  q quit  │
└─────────────────────────────────────────────────────────────────┘
```

### Aesthetic Direction: Minimal/Editorial

The visual design follows a minimal, editorial aesthetic:

**Typography & Spacing**
- Clean, readable markdown rendering
- Generous whitespace and padding
- Clear visual hierarchy between headings, body, and code

**Color Palette**
- Muted, sophisticated colors
- Adaptive to terminal background (light/dark)
- Subtle accent colors for selection and focus states
- No garish or overwhelming colors

**Visual Elements**
- Thin, refined borders (or borderless panels)
- Subtle focus indicators
- Clean tree indicators (▸ ▾ for directories)
- Understated selection highlighting

### Interaction Model

**Navigation**
- Arrow keys or vim-style (j/k) for list navigation
- Enter to open files or expand/collapse directories
- Tab to switch focus between file tree and preview
- Page Up/Down to scroll preview content
- Mouse wheel scrolls whichever panel is hovered (no click required)

**File Filtering (in file tree)**
- `/` activates filter mode in file tree
- Fuzzy matching on filenames
- Esc to clear filter and return to normal mode

**Content Search (in preview)**
- `/` activates search mode when preview is focused
- Case-insensitive search within markdown content
- `n` / `N` to navigate between matches
- Status bar shows current match position (e.g., "match 2/5")
- Esc to clear search and return to normal mode

**File Watching**
- Automatic refresh when viewed file changes
- Debounced updates to prevent flickering
- Visual indicator when file is being watched

## Technical Architecture

### Component Hierarchy

```
App (main model)
├── FileTree (left panel)
│   ├── List (bubbles/list)
│   ├── TreeItems (custom delegate)
│   └── FilterInput (built-in)
├── Preview (right panel)
│   ├── Viewport (bubbles/viewport)
│   ├── GlamourRenderer
│   └── SearchInput (textinput for content search)
├── StatusBar (bottom)
└── FileWatcher (background)
```

### Data Flow

```
User Input → App.Update() → Route to focused component
                         → Update state
                         → Return commands (file read, watch, etc.)

File Change → Watcher → FileChangedMsg → App.Update() → Refresh preview

Window Resize → tea.WindowSizeMsg → Recalculate panel sizes
                                  → Update child components
```

### Key Design Decisions

1. **Split-Pane Layout**
   - Fixed ratio: 25% file tree, 75% preview
   - Responsive to terminal size
   - Minimum widths to prevent unusable states

2. **File Tree as Flattened List**
   - Directories and files in single list
   - Indentation indicates depth
   - Expand/collapse changes list contents
   - Simpler than true tree widget, sufficient for this use case

3. **Glamour for Markdown**
   - Battle-tested markdown renderer
   - Built-in themes that adapt to terminal
   - Proper word wrapping support

4. **fsnotify for File Watching**
   - Cross-platform (Linux, macOS, Windows)
   - Event-based, not polling
   - Minimal resource usage

## Constraints & Boundaries

### In Scope
- Viewing `.md` files only
- Local filesystem navigation
- Single file preview at a time
- Read-only viewing

### Out of Scope
- Editing markdown files
- Remote file access
- Multiple preview tabs
- Export/conversion features
- Syntax highlighting for code blocks (beyond Glamour's default)

## Success Criteria

1. **Fast startup** - Under 100ms to first render
2. **Responsive** - No perceptible lag when navigating
3. **Low memory** - Under 50MB for typical usage
4. **Reliable watching** - Changes reflected within 500ms
5. **Intuitive** - Usable without reading documentation
