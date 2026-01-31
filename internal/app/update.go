package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Window resize
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.ready = true
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
	}

	return m, nil
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
	// If filter mode is active, handle filter input
	if m.filterActive {
		switch msg.String() {
		case "esc":
			m.filterActive = false
			m.filterText = ""
			return m, nil
		case "enter":
			m.filterActive = false
			return m, nil
		case "backspace":
			if len(m.filterText) > 0 {
				m.filterText = m.filterText[:len(m.filterText)-1]
			}
			return m, nil
		default:
			// Add character to filter (if printable)
			if len(msg.String()) == 1 {
				m.filterText += msg.String()
			}
			return m, nil
		}
	}

	// Normal file tree navigation
	switch msg.String() {
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}
		return m, nil

	case "down", "j":
		if m.selectedIndex < len(m.fileTreeItems)-1 {
			m.selectedIndex++
		}
		return m, nil

	case "enter":
		// Open file / toggle directory (to be implemented in Phase 2)
		return m, nil

	case "/":
		// Enter filter mode
		m.filterActive = true
		m.filterText = ""
		return m, nil

	case "esc":
		// Clear selection / reset
		m.selectedIndex = 0
		return m, nil
	}

	return m, nil
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
