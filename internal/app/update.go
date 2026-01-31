package app

import (
	"github.com/athakur/local-md/internal/components/filetree"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	// Window resize
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.ready = true

		// Update file tree size
		fileTreeWidth, _ := m.PanelWidths()
		contentHeight := m.ContentHeight()
		m.fileTree.SetSize(fileTreeWidth-2, contentHeight)

		return m, nil

	// Keyboard input
	case tea.KeyMsg:
		return m.handleKeypress(msg)

	// Mouse input (for future use)
	case tea.MouseMsg:
		return m, nil

	// Custom messages
	case FileSelectedMsg:
		m.previewPath = msg.Path
		return m, nil

	case FileLoadedMsg:
		m.previewContent = msg.Content
		m.previewPath = msg.Path
		return m, nil

	case FileErrorMsg:
		m.previewContent = "Error loading file: " + msg.Err.Error()
		return m, nil

	case FocusChangedMsg:
		m.FocusedPanel = msg.Panel
		return m, nil

	case FilterActiveMsg:
		m.filterActive = msg.Active
		if !msg.Active {
			m.filterText = ""
		}
		return m, nil

	// File tree component messages
	case filetree.FileSelectedMsg:
		m.previewPath = msg.Path
		// TODO: Load file content in Phase 3
		return m, nil

	case filetree.DirectoryToggledMsg:
		// Directory was toggled, tree already updated
		return m, nil
	}

	// Forward messages to file tree when focused
	if m.FocusedPanel == FileTreePanel {
		var cmd tea.Cmd
		m.fileTree, cmd = m.fileTree.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// handleKeypress processes keyboard input
func (m Model) handleKeypress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys (work regardless of focus/mode)
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "?":
		// Toggle help (to be implemented in Phase 6)
		return m, nil

	case "tab":
		// Switch panel focus
		if m.FocusedPanel == FileTreePanel {
			m.FocusedPanel = PreviewPanel
		} else {
			m.FocusedPanel = FileTreePanel
		}
		return m, nil
	}

	// Panel-specific keys
	switch m.FocusedPanel {
	case FileTreePanel:
		return m.handleFileTreeKeys(msg)
	case PreviewPanel:
		return m.handlePreviewKeys(msg)
	}

	return m, nil
}

// handleFileTreeKeys handles keys when file tree is focused
func (m Model) handleFileTreeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Delegate to file tree component
	var cmd tea.Cmd
	m.fileTree, cmd = m.fileTree.Update(msg)
	return m, cmd
}

// handlePreviewKeys handles keys when preview is focused
func (m Model) handlePreviewKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		// Scroll up (to be implemented in Phase 3)
		return m, nil

	case "down", "j":
		// Scroll down (to be implemented in Phase 3)
		return m, nil

	case "pgup", "ctrl+u":
		// Page up (to be implemented in Phase 3)
		return m, nil

	case "pgdown", "ctrl+d":
		// Page down (to be implemented in Phase 3)
		return m, nil

	case "g":
		// Go to top (to be implemented in Phase 3)
		return m, nil

	case "G":
		// Go to bottom (to be implemented in Phase 3)
		return m, nil
	}

	return m, nil
}
