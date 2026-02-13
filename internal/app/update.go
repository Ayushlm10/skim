package app

import (
	"github.com/Ayushlm10/skim/internal/components/filetree"
	"github.com/Ayushlm10/skim/internal/components/preview"
	"github.com/Ayushlm10/skim/internal/watcher"
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

		if m.fullscreen {
			// In fullscreen, preview gets full terminal dimensions
			m.preview.SetSize(m.Width-2, m.FullscreenContentHeight())
		} else {
			// Normal mode: update component sizes with panel split
			fileTreeWidth, previewWidth := m.PanelWidths()
			contentHeight := m.ContentHeight()
			m.fileTree.SetSize(fileTreeWidth-2, contentHeight)
			m.preview.SetSize(previewWidth-2, contentHeight)
		}

		return m, nil

	// Keyboard input
	case tea.KeyMsg:
		return m.handleKeypress(msg)

	// Mouse input - route to appropriate panel based on X coordinate
	case tea.MouseMsg:
		return m.handleMouse(msg)

	// Custom messages
	case FileSelectedMsg:
		// Load the file content
		m.loading = true
		m.lastError = ""
		return m, preview.LoadFile(msg.Path)

	case FocusChangedMsg:
		m.FocusedPanel = msg.Panel
		return m, nil

	case FilterActiveMsg:
		m.filterActive = msg.Active
		m.filterText = msg.Value
		return m, nil

	// File tree filter state change
	case filetree.FilterChangedMsg:
		m.filterActive = msg.Active
		m.filterText = msg.Value
		return m, nil

	// File tree ignored directories toggled
	case filetree.IgnoredDirsToggledMsg:
		m.showIgnored = msg.ShowIgnored
		return m, nil

	// File tree component messages
	case filetree.FileSelectedMsg:
		// Load the file content when a file is selected in the tree
		m.loading = true
		m.lastError = ""
		return m, preview.LoadFile(msg.Path)

	case filetree.DirectoryToggledMsg:
		// Directory was toggled, tree already updated
		return m, nil

	// Preview component messages
	case preview.FileLoadedMsg:
		// Forward to preview component
		var cmd tea.Cmd
		m.preview, cmd = m.preview.Update(msg)
		m.loading = false

		// Handle errors
		if msg.Error != nil {
			m.lastError = msg.Error.Error()
			return m, cmd
		}

		// Start watching the newly loaded file
		m.lastError = ""
		if m.watcher != nil {
			m.watchedFile = msg.Path
			return m, tea.Batch(cmd, watcher.StartWatching(m.watcher, msg.Path))
		}
		return m, cmd

	// Watcher messages (Phase 5)
	case watcher.FileChangedMsg:
		// File changed, reload it
		if msg.Path == m.watchedFile {
			return m, tea.Batch(
				preview.LoadFile(msg.Path),
				watcher.WaitForChange(m.watcher),
			)
		}
		return m, watcher.WaitForChange(m.watcher)

	case watcher.WatchErrorMsg:
		// Log error but continue watching
		// TODO: Show error in status bar in Phase 6
		return m, watcher.WaitForChange(m.watcher)
	}

	// Handle internal watch started message
	if path, ok := watcher.IsWatchStartedMsg(msg); ok {
		m.watchedFile = path
		return m, watcher.WaitForChange(m.watcher)
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
	// If help is visible, handle help keys first
	if m.help.IsVisible() {
		switch msg.String() {
		case "?", "esc", "enter", "q":
			m.help.Hide()
			return m, nil
		}
		// Ignore other keys when help is visible
		return m, nil
	}

	// Global keys (work regardless of focus/mode)
	switch msg.String() {
	case "ctrl+c", "q":
		// Clean up watcher before quitting
		if m.watcher != nil {
			_ = m.watcher.Close()
		}
		return m, tea.Quit

	case "?":
		// Toggle help overlay
		m.help.Toggle()
		return m, nil

	case "tab":
		// Switch panel focus (no-op in fullscreen since preview is always focused)
		if !m.fullscreen {
			if m.FocusedPanel == FileTreePanel {
				m.FocusedPanel = PreviewPanel
			} else {
				m.FocusedPanel = FileTreePanel
			}
		}
		return m, nil

	case "f":
		// Don't toggle fullscreen if user is typing in search or filter
		if m.preview.IsSearchMode() || m.filterActive {
			break
		}
		m.fullscreen = !m.fullscreen
		if m.fullscreen {
			m.FocusedPanel = PreviewPanel
			m.preview.SetSize(m.Width-2, m.FullscreenContentHeight())
		} else {
			fileTreeWidth, previewWidth := m.PanelWidths()
			contentHeight := m.ContentHeight()
			m.fileTree.SetSize(fileTreeWidth-2, contentHeight)
			m.preview.SetSize(previewWidth-2, contentHeight)
		}
		return m, nil

	case "esc":
		// Exit fullscreen if active (and no search/filter is consuming Esc)
		if m.fullscreen && !m.preview.IsSearchMode() && !m.preview.HasActiveSearch() {
			m.fullscreen = false
			fileTreeWidth, previewWidth := m.PanelWidths()
			contentHeight := m.ContentHeight()
			m.fileTree.SetSize(fileTreeWidth-2, contentHeight)
			m.preview.SetSize(previewWidth-2, contentHeight)
			return m, nil
		}
	}

	// Panel-specific keys
	if m.fullscreen {
		// In fullscreen, all keys go to preview
		return m.handlePreviewKeys(msg)
	}

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
	// Delegate to preview component
	var cmd tea.Cmd
	m.preview, cmd = m.preview.HandleKey(msg)
	return m, cmd
}

// handleMouse routes mouse events to the appropriate panel based on X coordinate
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Only handle mouse wheel events for scrolling
	if msg.Button != tea.MouseButtonWheelUp && msg.Button != tea.MouseButtonWheelDown {
		return m, nil
	}

	// In fullscreen, all mouse events go to preview
	if m.fullscreen {
		var cmd tea.Cmd
		m.preview, cmd = m.preview.HandleMouse(msg)
		return m, cmd
	}

	// Calculate panel boundary (file tree width + left border)
	fileTreeWidth, _ := m.PanelWidths()
	panelBoundary := fileTreeWidth + 2 // +2 for left panel border

	// Route to appropriate panel based on mouse X position
	if msg.X < panelBoundary {
		// Mouse is over file tree panel
		var cmd tea.Cmd
		m.fileTree, cmd = m.fileTree.Update(msg)
		return m, cmd
	}

	// Mouse is over preview panel
	var cmd tea.Cmd
	m.preview, cmd = m.preview.HandleMouse(msg)
	return m, cmd
}
