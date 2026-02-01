package app

import (
	"github.com/Ayushlm10/skim/internal/components/filetree"
	"github.com/Ayushlm10/skim/internal/components/help"
	"github.com/Ayushlm10/skim/internal/components/preview"
	"github.com/Ayushlm10/skim/internal/styles"
	"github.com/Ayushlm10/skim/internal/watcher"
	tea "github.com/charmbracelet/bubbletea"
)

// Model is the main application model
type Model struct {
	// Root path being viewed
	RootPath string

	// Window dimensions
	Width  int
	Height int

	// Panel focus
	FocusedPanel Panel

	// File tree component (Phase 2)
	fileTree filetree.Model

	// Preview component (Phase 3)
	preview preview.Model

	// Help overlay (Phase 6)
	help help.Model

	// File watcher (Phase 5)
	watcher     *watcher.Watcher
	watchedFile string

	// UI state
	ready        bool
	filterActive bool
	filterText   string
	loading      bool
	lastError    string
	showIgnored  bool // Whether ignored directories are visible
}

// New creates a new application model
func New(rootPath string) Model {
	// Create file tree with initial dimensions (will be resized)
	ft := filetree.New(rootPath, 30, 20)

	// Create preview component with initial dimensions
	pv := preview.New(60, 20)

	// Create help overlay
	h := help.New()

	// Create file watcher
	w, _ := watcher.New()

	return Model{
		RootPath:     rootPath,
		FocusedPanel: FileTreePanel,
		fileTree:     ft,
		preview:      pv,
		help:         h,
		watcher:      w,
		ready:        false,
	}
}

// Init initializes the model and returns an initial command
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("skim"),
		m.fileTree.Init(),
		m.preview.Init(),
	)
}

// PanelWidths calculates the width of each panel based on total width
func (m Model) PanelWidths() (fileTree, preview int) {
	// Account for borders (2 chars each panel)
	usableWidth := m.Width - 4

	fileTree = int(float64(usableWidth) * styles.FileTreeRatio)
	preview = usableWidth - fileTree

	// Ensure minimum widths
	if fileTree < 20 {
		fileTree = 20
		preview = usableWidth - fileTree
	}
	if preview < 30 {
		preview = 30
		fileTree = usableWidth - preview
	}

	return fileTree, preview
}

// ContentHeight returns the height available for panel content (inner height for lipgloss)
func (m Model) ContentHeight() int {
	// Total height minus:
	// - header (1 line)
	// - newline after header (1 line)
	// - panel borders (2 lines: top + bottom, added by lipgloss on top of inner height)
	// - newline after panels (1 line)
	// - status bar (1 line)
	// Total overhead: 6 lines
	// lipgloss.Height() sets inner height, border adds 2 more
	return m.Height - 6
}
