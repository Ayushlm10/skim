# Skim

Terminal-based markdown viewer with dual-panel layout (file tree + preview).

**Tech Stack:** Go 1.25+, Charm TUI stack (Bubble Tea, Lip Gloss, Glamour, Bubbles)

## Directory Structure

internal/app/       - Main application state (Model, Update, View)
internal/components/ - UI components (filetree, preview, help)
internal/styles/    - Lipgloss styles and color palette
internal/watcher/   - File watching with fsnotify
main.go             - CLI entry point

## Commands

```bash
go build -o skim     # Build binary
go run .             # Run in current directory
go run . ~/docs      # Run with target directory
```

## Code Patterns

### Bubble Tea Architecture

- Model-Update-View pattern in internal/app/
- Components are self-contained in internal/components/*/
- Each component has its own Model, Update(), View(), and message types
- Parent-child communication via custom tea.Msg types (e.g., FileSelectedMsg, FilterChangedMsg)

### Component Structure

```go
type Model struct { ... }
func New(...) Model { ... }
func (m Model) Init() tea.Cmd { ... }
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) { ... }
func (m Model) View() string { ... }
```

See internal/components/filetree/filetree.go for reference implementation.

### Styling

- All styles in internal/styles/styles.go
- Use lipgloss.AdaptiveColor for light/dark terminal support
- Panel ratio: 25% file tree, 75% preview

### Message Flow

1. User input → app.Update() routes to focused panel
2. Component returns tea.Cmd producing message
3. Message bubbles up to parent via custom Msg types
4. Parent handles cross-component coordination

### Boundaries

✅ Always: Follow existing Bubble Tea patterns in internal/app/
✅ Always: Use adaptive colors from internal/styles/styles.go
⚠️ Requires Approval: Adding new dependencies to go.mod
🚫 Never: Put styles inline - add to styles.go
🚫 Never: Handle keys globally that should be panel-specific
