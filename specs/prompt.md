# Local MD Viewer - Prompts & Working Instructions

## How We Work in This Repository

This repository follows a **specification-first, iterative development** approach. Before writing code, we document the design, plan, and track implementation progress.

### Repository Structure

```
specs/
├── design.md          # WHAT we're building and WHY (high-level)
├── plan.md            # HOW we plan to build it (detailed tasks)
├── implementation.md  # WHERE we are (living status tracker)
└── prompt.md          # This file - prompts and instructions
```

### Workflow

1. **Read the specs first** - Before making changes, read `design.md` to understand the vision and `plan.md` to understand the approach.

2. **Check implementation status** - Look at `implementation.md` to see what's done, in progress, and pending.

3. **Work in phases** - Follow the phases defined in `plan.md`. Complete one phase before moving to the next.

4. **Update as you go** - After completing tasks, update `implementation.md` with status changes.

5. **Document decisions** - Add significant decisions to the "Decisions Made" table in `implementation.md`.

### For AI Assistants

When working on this project:
- Load and reference the spec files before making changes
- Follow the established patterns and aesthetic direction
- Update `implementation.md` after completing tasks
- Ask clarifying questions rather than making assumptions
- Use the Go TUI skills (bubbletea, lipgloss, bubbles, glamour) for implementation guidance
- **Create a git commit at the end of each phase** with a descriptive message
- **Check git history** (`git log --oneline`) to understand what's been done
- **Update AGENTS.md** when discovering patterns or instructions that would help future AI assistants

### Maintaining AGENTS.md

The `AGENTS.md` file at the repository root contains instructions specifically for AI coding assistants. When you discover something that would help future assistants work on this codebase, add it to AGENTS.md:

**What to add:**
- Code patterns specific to this project (e.g., how to add a new component)
- Common pitfalls or gotchas encountered during implementation
- Project-specific conventions that differ from general best practices
- Boundaries and constraints (what NOT to do)
- Useful commands or workflows

**Format:**
- Keep it concise and scannable
- Use the existing section structure
- Add to "Code Patterns" for implementation patterns
- Add to "Boundaries" for constraints
- Add to "Workflow" for process instructions

### Git History

The git history serves as an additional source of truth for this project. Each phase completion is marked with a commit. Use `git log --oneline` to see the progression:

```bash
git log --oneline --grep="Phase"
```

---

## Stored Prompts

### Initial Project Prompt

```
I work with a lot of markdown files for planning, spec, decisions. Can we build 
a local markdown viewer which I can use for viewing all these files? Use 
frontend design skill to understand how to build the UI and build it in Go. 
You can use the go-tui skills.
```

### Feature Requirements (from Q&A)

```
Project Location: Use current directory (local-md)

Features requested:
- File browser sidebar
- Rendered markdown preview
- Fuzzy search/filter
- Watch mode / live reload

Aesthetic: Minimal/Editorial
- Clean typography, generous whitespace, refined details

CLI Behavior:
- Accept path argument (e.g., `mdview ./docs`)
- Default to current directory if no argument

Edge Cases:
- Show helpful message for empty directories or no markdown files
```

### Design Direction Prompt

```
Visual style: Minimal/Editorial

This means:
- Muted, sophisticated color palette
- Adaptive to terminal background (light/dark)
- Subtle accent colors for selection and focus
- Thin, refined borders (or borderless)
- Generous whitespace and padding
- Clean tree indicators (▸ ▾)
- No garish or overwhelming colors

Typography:
- Clean, readable markdown rendering
- Clear visual hierarchy
- Proper code block formatting
```

### Layout Specification

```
┌─────────────────────────────────────────────────────────────────┐
│  Local MD Viewer                                    ~/docs      │
├──────────────────────┬──────────────────────────────────────────┤
│                      │                                          │
│  ▾ docs/             │  # Project Specification                 │
│    ├── spec.md  ◀    │                                          │
│    ├── plan.md       │  This document outlines...               │
│    └── decisions/    │                                          │
│        ▸ adr-001.md  │  ## Goals                                │
│                      │                                          │
│  ┌────────────────┐  │  - Fast local viewing                    │
│  │ Filter: ___    │  │  - Clean typography                      │
│  └────────────────┘  │                                          │
│                      │                                          │
├──────────────────────┴──────────────────────────────────────────┤
│  ↑↓ navigate  │  ⏎ open  │  / filter  │  Tab switch  │  q quit  │
└─────────────────────────────────────────────────────────────────┘

Panel ratio: 25% file tree, 75% preview
```

### Key Bindings Prompt

```
Navigation:
- ↑/k: Move up
- ↓/j: Move down
- Enter: Open file / Toggle directory
- Tab: Switch panel focus

Search:
- /: Enter filter mode
- Esc: Exit filter, clear

Preview scrolling:
- PgUp/Ctrl+u: Scroll up
- PgDn/Ctrl+d: Scroll down
- g: Go to top
- G: Go to bottom

Global:
- ?: Toggle help
- q/Ctrl+c: Quit
```

### Tech Stack Prompt

```
Build a Go TUI application using the Charm stack:

- github.com/charmbracelet/bubbletea - Elm-style TUI framework
- github.com/charmbracelet/bubbles - Pre-built components (list, viewport)
- github.com/charmbracelet/lipgloss - Styling and layout
- github.com/charmbracelet/glamour - Markdown rendering
- github.com/fsnotify/fsnotify - File system watching

Follow the Model/Update/View pattern from Bubble Tea.
Use bubbles/list for the file tree with a custom delegate.
Use bubbles/viewport for scrollable markdown preview.
Use Glamour with auto style for terminal-adaptive rendering.
```

---

## Phase Completion Commits

Each phase should end with a git commit. The commit message format:

```
Phase N: [Phase Name] - [Brief description]

- Bullet point of key changes
- Another key change
```

### Phase 1 Commit (2026-01-31)

```
Phase 1: Foundation - Basic dual-panel TUI with Charm stack

- Initialize Go module with Bubble Tea, Lip Gloss, Glamour, fsnotify
- CLI argument parsing (accepts path, defaults to current dir)
- Main model with panel focus management
- Dual-panel layout (25% file tree / 75% preview)
- Minimal/editorial styling with adaptive colors
- Keyboard handling (Tab, j/k, arrows, q to quit)
- Status bar with key hints
- Placeholder content for file tree and preview
```

### Phase 2 Commit (2026-01-31)

```
Phase 2: File Tree Component - Directory scanning with expand/collapse

- Add filetree component package (scanner, item, delegate, filetree)
- Implement directory scanner for .md files with depth tracking
- Custom list delegate for tree-style rendering with ▸/▾ indicators
- Expand/collapse directories via Enter key with lazy child loading
- Wire component into main app model with message passing
- Built on bubbles list for filtering support
```

### Phase 3 Commit (2026-01-31)

```
Phase 3: Markdown Preview - Glamour rendering with viewport scrolling

- Add preview component package (renderer, preview)
- Integrate Glamour with auto style for light/dark terminal adaptation
- Implement viewport-based scrolling (j/k, PgUp/PgDn, g/G, Ctrl+u/d)
- Wire file selection to preview via LoadFile command
- Adaptive word wrapping based on panel width
- Status bar shows filename and scroll percentage when preview focused
```

---

## Future Prompts

_Add new prompts here as the project evolves._

### Template for Adding Prompts

```
### [Prompt Name]

Date: YYYY-MM-DD
Context: [Why this prompt was needed]

\`\`\`
[The actual prompt text]
\`\`\`
```
