package app

import (
	"github.com/athakur/local-md/internal/components/filetree"
	"github.com/athakur/local-md/internal/styles"
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

	// Preview state (placeholder for Phase 3)
	previewContent string
	previewPath    string

	// UI state
	ready        bool
	filterActive bool
	filterText   string
}

// New creates a new application model
func New(rootPath string) Model {
	// Create file tree with initial dimensions (will be resized)
	ft := filetree.New(rootPath, 30, 20)

	return Model{
		RootPath:       rootPath,
		FocusedPanel:   FileTreePanel,
		fileTree:       ft,
		previewContent: "",
		ready:          false,
	}
}

// Init initializes the model and returns an initial command
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Local MD Viewer"),
		m.fileTree.Init(),
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

// ContentHeight returns the height available for panel content
func (m Model) ContentHeight() int {
	// Total height minus header (1) and status bar (1) and borders (4)
	return m.Height - 6
}
