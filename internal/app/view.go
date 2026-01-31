package app

import (
	"strings"

	"github.com/athakur/local-md/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

// View renders the application UI
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Render header
	b.WriteString(m.renderHeader())
	b.WriteString("\n")

	// Render main content (file tree + preview)
	b.WriteString(m.renderPanels())

	// Render status bar
	b.WriteString(m.renderStatusBar())

	return b.String()
}

// renderHeader renders the top header bar
func (m Model) renderHeader() string {
	title := styles.HeaderStyle.Render("Local MD Viewer")

	// Show the current path (truncated if needed)
	path := m.RootPath
	maxPathLen := m.Width - lipgloss.Width(title) - 4
	if len(path) > maxPathLen && maxPathLen > 10 {
		// Truncate from the beginning with ~
		path = "~" + path[len(path)-maxPathLen+1:]
	}
	pathStr := styles.HeaderPathStyle.Render(path)

	// Create header with title left, path right
	spacer := strings.Repeat(" ", m.Width-lipgloss.Width(title)-lipgloss.Width(pathStr))
	if len(spacer) < 0 {
		spacer = " "
	}

	return title + spacer + pathStr
}

// renderPanels renders the file tree and preview panels side by side
func (m Model) renderPanels() string {
	fileTreeWidth, previewWidth := m.PanelWidths()
	contentHeight := m.ContentHeight()

	// Render file tree panel
	fileTreeContent := m.renderFileTree(fileTreeWidth-2, contentHeight)
	fileTreePanel := m.stylePanelBox(fileTreeContent, fileTreeWidth, contentHeight, m.FocusedPanel == FileTreePanel)

	// Render preview panel
	previewContent := m.renderPreview(previewWidth-2, contentHeight)
	previewPanel := m.stylePanelBox(previewContent, previewWidth, contentHeight, m.FocusedPanel == PreviewPanel)

	// Join panels horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, fileTreePanel, previewPanel)
}

// stylePanelBox applies panel styling with border
func (m Model) stylePanelBox(content string, width, height int, focused bool) string {
	style := styles.PanelStyle
	if focused {
		style = styles.FocusedPanelStyle
	}

	return style.
		Width(width).
		Height(height).
		Render(content)
}

// renderFileTree renders the file tree content
func (m Model) renderFileTree(width, height int) string {
	// Update file tree size and focus state
	m.fileTree.SetSize(width, height)
	m.fileTree.SetFocused(m.FocusedPanel == FileTreePanel)

	// Render the file tree component
	content := m.fileTree.View()

	// Ensure content fills the available height
	lines := strings.Split(content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}

	return strings.Join(lines[:height], "\n")
}

// renderPreview renders the markdown preview content
func (m Model) renderPreview(width, height int) string {
	// Update preview size and focus state
	m.preview.SetSize(width, height)
	m.preview.SetFocused(m.FocusedPanel == PreviewPanel)

	// Render the preview component
	content := m.preview.View()

	// Ensure content fills the available height
	lines := strings.Split(content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}

	// Truncate if too long
	if len(lines) > height {
		lines = lines[:height]
	}

	return strings.Join(lines, "\n")
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	// Check if we're in filter mode
	if m.filterActive {
		return m.renderFilterStatusBar()
	}

	// Build help hints based on focused panel
	var hints []struct {
		key  string
		desc string
	}

	if m.FocusedPanel == FileTreePanel {
		hints = []struct {
			key  string
			desc string
		}{
			{"↑↓", "navigate"},
			{"⏎", "open"},
			{"/", "filter"},
			{"Tab", "switch"},
			{"q", "quit"},
		}
	} else {
		hints = []struct {
			key  string
			desc string
		}{
			{"↑↓", "scroll"},
			{"g/G", "top/bottom"},
			{"Tab", "switch"},
			{"q", "quit"},
		}
	}

	var parts []string
	for _, h := range hints {
		part := styles.HelpKeyStyle.Render(h.key) + " " + styles.HelpDescStyle.Render(h.desc)
		parts = append(parts, part)
	}

	separator := styles.HelpSeparatorStyle.Render("  │  ")
	statusContent := strings.Join(parts, separator)

	// Add filter indicator if there's an active filter
	if m.filterText != "" {
		filterInfo := styles.FilterPromptStyle.Render("filter: ") +
			styles.StatusValueStyle.Render(m.filterText)
		// Calculate spacing
		statusWidth := lipgloss.Width(statusContent)
		filterWidth := lipgloss.Width(filterInfo)
		spacerWidth := m.Width - statusWidth - filterWidth - 4
		if spacerWidth > 0 {
			statusContent = statusContent + strings.Repeat(" ", spacerWidth) + filterInfo
		}
	} else if m.FocusedPanel == PreviewPanel && m.preview.FilePath() != "" {
		// Add scroll indicator if in preview and file is loaded
		scrollPct := int(m.preview.ScrollPercent() * 100)
		// Add watch indicator if file is being watched
		watchIndicator := ""
		if m.watchedFile != "" && m.watchedFile == m.preview.FilePath() {
			watchIndicator = " [watching]"
		}
		scrollInfo := styles.StatusValueStyle.Render(
			m.preview.FileName() + watchIndicator + " " + styles.HelpDescStyle.Render(
				"["+itoa(scrollPct)+"%]",
			),
		)
		// Calculate spacing
		statusWidth := lipgloss.Width(statusContent)
		scrollWidth := lipgloss.Width(scrollInfo)
		spacerWidth := m.Width - statusWidth - scrollWidth - 4
		if spacerWidth > 0 {
			statusContent = statusContent + strings.Repeat(" ", spacerWidth) + scrollInfo
		}
	}

	return styles.StatusBarStyle.
		Width(m.Width).
		Render(statusContent)
}

// renderFilterStatusBar renders the status bar during filter mode
func (m Model) renderFilterStatusBar() string {
	hints := []struct {
		key  string
		desc string
	}{
		{"⏎", "accept"},
		{"Esc", "cancel"},
	}

	var parts []string
	for _, h := range hints {
		part := styles.HelpKeyStyle.Render(h.key) + " " + styles.HelpDescStyle.Render(h.desc)
		parts = append(parts, part)
	}

	separator := styles.HelpSeparatorStyle.Render("  │  ")
	statusContent := strings.Join(parts, separator)

	// Add filtering indicator
	filterIndicator := styles.FilterPromptStyle.Render("FILTERING")
	statusWidth := lipgloss.Width(statusContent)
	filterWidth := lipgloss.Width(filterIndicator)
	spacerWidth := m.Width - statusWidth - filterWidth - 4
	if spacerWidth > 0 {
		statusContent = statusContent + strings.Repeat(" ", spacerWidth) + filterIndicator
	}

	return styles.StatusBarStyle.
		Width(m.Width).
		Render(statusContent)
}

// itoa converts int to string without importing strconv
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + itoa(-i)
	}
	var digits []byte
	for i > 0 {
		digits = append([]byte{byte('0' + i%10)}, digits...)
		i /= 10
	}
	return string(digits)
}
