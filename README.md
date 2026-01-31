# skim

A terminal-based markdown viewer built with the [Charm](https://charm.sh) Go TUI stack.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue)

## Features

- **Dual-panel layout** - File tree (25%) and markdown preview (75%)
- **Beautiful rendering** - Glamour-powered markdown with automatic light/dark terminal adaptation
- **File tree navigation** - Expand/collapse directories, filter files with fuzzy search
- **In-preview search** - Search within content with match highlighting and navigation
- **Live reload** - Automatic re-render when files change on disk
- **Keyboard-driven** - Vim-style navigation with full mouse support
- **Minimal aesthetic** - Clean, editorial design with muted colors

## Installation

```bash
go install github.com/Ayushlm10/skim@latest
```

Or build from source:

```bash
git clone https://github.com/Ayushlm10/skim.git
cd skim
go build -o skim
```

## Usage

```bash
# View markdown files in current directory
skim

# View markdown files in a specific directory
skim ~/docs

# View a project's documentation
skim ./specs
```

## Keyboard Shortcuts

### Navigation

| Key | Action |
|-----|--------|
| `j` / `k` / `↑` / `↓` | Move up/down |
| `Enter` | Open file / Toggle directory |
| `Tab` | Switch between panels |

### File Tree

| Key | Action |
|-----|--------|
| `/` | Filter files (fuzzy search) |
| `Esc` | Clear filter |

### Preview

| Key | Action |
|-----|--------|
| `j` / `k` | Scroll line by line |
| `PgUp` / `Ctrl+u` | Scroll half page up |
| `PgDn` / `Ctrl+d` | Scroll half page down |
| `g` | Go to top |
| `G` | Go to bottom |
| `/` | Search in content |
| `n` / `N` | Next/previous match |
| `Esc` | Clear search |

### General

| Key | Action |
|-----|--------|
| `?` | Toggle help overlay |
| `q` / `Ctrl+c` | Quit |

## Tech Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [fsnotify](https://github.com/fsnotify/fsnotify) - File watching

## License

MIT
