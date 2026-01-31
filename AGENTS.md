A terminal-based markdown viewer built with the Charm Go TUI stack for browsing local markdown files.
**Tech Stack:** Go 1.25.2, Bubble Tea, Lip Gloss
**Binary:** `skim`
 Directory Structure
internal/app/         - Main application model, update, view logic
internal/styles/      - Centralized Lip Gloss styles
internal/components/  - UI components (planned: filetree, preview, statusbar, watcher)
specs/                - Design docs and implementation tracking
test-docs/            - Test markdown files for development
## Commands
```bash
go run main.go [path]  # Run viewer with optional directory path (defaults to .)
go build -o skim       # Build binary
go mod tidy            # Update dependencies
Code Patterns
Bubble Tea Architecture: Model-Update-View pattern in internal/app/
- model.go - Application state and panel management
- update.go - Event handling and state updates
- view.go - Rendering logic
- messages.go - Custom message types
Panel System: Split-pane layout with focus management
- 25% file tree, 75% preview (configurable ratio in internal/styles/styles.go)
- FocusedPanel enum tracks active panel
- PanelWidths() calculates responsive dimensions with minimums
Component Pattern (upcoming phases): Custom components wrap Bubble Tea models
- See specs/implementation.md for planned component structure
- Components will go in internal/components/
Project Status
Phase 1 Complete: Foundation with dual-panel layout and placeholder content
Next Phase: Directory scanner and file tree component with expand/collapse
See specs/implementation.md for detailed phase breakdown and status.
Workflow
Adding a Component
1. Create in internal/components/<name>/
2. Implement Bubble Tea model interface (Init, Update, View)
3. Wire into main model in internal/app/model.go
Testing Changes
1. Use test-docs/ directory for sample markdown files
2. Run: go run main.go test-docs
Boundaries
✅ Always: Match minimal/editorial aesthetic defined in specs/design.md
✅ Always: Use Charm libraries (Bubble Tea, Lip Gloss, Bubbles, Glamour)
⚠️ Requires Approval: Breaking changes to panel ratio or keyboard mappings
🚫 Never: Add file editing features - read-only viewer only
🚫 Never: Support non-markdown files
