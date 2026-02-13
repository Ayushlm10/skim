package app

import (
	"strings"

	"github.com/Ayushlm10/skim/internal/styles"
	"github.com/charmbracelet/lipgloss"
)

// View renders the application UI
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var b strings.Builder

	if m.fullscreen {
		// Fullscreen: preview content + status bar only
		b.WriteString(m.renderFullscreenPreview())
		b.WriteString("\n")
		b.WriteString(m.renderStatusBar())
	} else {
		// Normal mode: header + panels + status bar
		b.WriteString(m.renderHeader())
		b.WriteString("\n")
		b.WriteString(m.renderPanels())
		b.WriteString("\n")
		b.WriteString(m.renderStatusBar())
	}

	baseView := b.String()

	// Overlay help if visible
	if m.help.IsVisible() {
		return m.overlayHelp(baseView)
	}

	return baseView
}

// renderHeader renders the top header bar
func (m Model) renderHeader() string {
	title := styles.HeaderStyle.Render("skim")

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

// renderFullscreenPreview renders the preview taking the full terminal area
func (m Model) renderFullscreenPreview() string {
	content := m.preview.View()

	fullHeight := m.FullscreenContentHeight()
	lines := strings.Split(content, "\n")
	for len(lines) < fullHeight {
		lines = append(lines, "")
	}
	if len(lines) > fullHeight {
		lines = lines[:fullHeight]
	}

	return strings.Join(lines, "\n")
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
	// Note: SetSize is called in Update() on WindowSizeMsg
	// SetFocused changes are lost here (value receiver) but focus is visual only

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
	// Note: SetSize is called in Update() on WindowSizeMsg
	// SetFocused changes are lost here (value receiver) but focus is visual only

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
	// Check if we're in fullscreen mode
	if m.fullscreen {
		return m.renderFullscreenStatusBar()
	}

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
			{"i", "ignored"},
			{"f", "fullscreen"},
			{"Tab", "switch"},
			{"?", "help"},
			{"q", "quit"},
		}
	} else {
		// Preview panel - show search mode or normal hints
		if m.preview.IsSearchMode() {
			hints = []struct {
				key  string
				desc string
			}{
				{"⏎", "search"},
				{"Esc", "cancel"},
			}
		} else if m.preview.HasActiveSearch() {
			hints = []struct {
				key  string
				desc string
			}{
				{"n/N", "next/prev match"},
				{"Esc", "clear search"},
				{"/", "new search"},
				{"?", "help"},
			}
		} else {
			hints = []struct {
				key  string
				desc string
			}{
				{"↑↓", "scroll"},
				{"g/G", "top/bottom"},
				{"/", "search"},
				{"f", "fullscreen"},
				{"Tab", "switch"},
				{"?", "help"},
				{"q", "quit"},
			}
		}
	}

	var parts []string
	for _, h := range hints {
		part := styles.HelpKeyStyle.Render(h.key) + " " + styles.HelpDescStyle.Render(h.desc)
		parts = append(parts, part)
	}

	separator := styles.HelpSeparatorStyle.Render("  │  ")
	statusContent := strings.Join(parts, separator)

	// Build right-side status info
	var rightInfo string

	// Show error if any
	if m.lastError != "" {
		// Truncate long errors
		errMsg := m.lastError
		if len(errMsg) > 30 {
			errMsg = errMsg[:27] + "..."
		}
		rightInfo = styles.StatusErrorStyle.Render("error: " + errMsg)
	} else if m.loading {
		// Show loading indicator
		rightInfo = styles.StatusLoadingStyle.Render("loading...")
	} else if m.preview.IsSearchMode() {
		// Show searching indicator (takes priority)
		searchIndicator := styles.SearchPromptStyle.Render("[searching]")
		if m.preview.FileName() != "" {
			rightInfo = styles.StatusValueStyle.Render(m.preview.FileName()) + " " + searchIndicator
		} else {
			rightInfo = searchIndicator
		}
	} else if m.preview.HasActiveSearch() || m.preview.HasSearchNoMatches() {
		// Show search results (visible regardless of focused panel)
		var searchIndicator string
		if m.preview.HasActiveSearch() {
			matchInfo := itoa(m.preview.CurrentMatchIndex()) + "/" + itoa(m.preview.MatchCount())
			searchIndicator = styles.SearchMatchStyle.Render("[" + m.preview.SearchQuery() + ": " + matchInfo + "]")
		} else {
			searchIndicator = styles.SearchNoMatchStyle.Render("[" + m.preview.SearchQuery() + ": no matches]")
		}
		fileName := styles.StatusValueStyle.Render(m.preview.FileName())
		scrollPct := int(m.preview.ScrollPercent() * 100)
		scrollIndicator := styles.HelpDescStyle.Render("[" + itoa(scrollPct) + "%]")
		rightInfo = fileName + " " + searchIndicator + " " + scrollIndicator
	} else if m.filterText != "" {
		// Show active filter
		rightInfo = styles.FilterPromptStyle.Render("filter: ") +
			styles.StatusValueStyle.Render(m.filterText)
	} else if m.FocusedPanel == PreviewPanel && m.preview.FilePath() != "" {
		// Show file info when preview focused (no active search)
		scrollPct := int(m.preview.ScrollPercent() * 100)

		// Build status parts
		fileName := styles.StatusValueStyle.Render(m.preview.FileName())
		scrollIndicator := styles.HelpDescStyle.Render("[" + itoa(scrollPct) + "%]")

		// Add watch indicator if file is being watched
		watchIndicator := ""
		if m.watchedFile != "" && m.watchedFile == m.preview.FilePath() {
			watchIndicator = styles.StatusWatchingStyle.Render(" [watching]")
		}

		rightInfo = fileName + watchIndicator + " " + scrollIndicator
	} else if m.FocusedPanel == FileTreePanel && m.showIgnored {
		// Show indicator when ignored directories are visible
		rightInfo = styles.StatusIgnoredStyle.Render("[showing ignored]")
	}

	// Calculate spacing and add right info
	if rightInfo != "" {
		statusWidth := lipgloss.Width(statusContent)
		rightWidth := lipgloss.Width(rightInfo)
		spacerWidth := m.Width - statusWidth - rightWidth - 4
		if spacerWidth > 0 {
			statusContent = statusContent + strings.Repeat(" ", spacerWidth) + rightInfo
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

// renderFullscreenStatusBar renders the status bar during fullscreen mode
func (m Model) renderFullscreenStatusBar() string {
	var hints []struct {
		key  string
		desc string
	}

	if m.preview.IsSearchMode() {
		hints = []struct {
			key  string
			desc string
		}{
			{"⏎", "search"},
			{"Esc", "cancel"},
		}
	} else if m.preview.HasActiveSearch() {
		hints = []struct {
			key  string
			desc string
		}{
			{"n/N", "next/prev match"},
			{"Esc", "clear search"},
			{"/", "new search"},
			{"f", "exit fullscreen"},
		}
	} else {
		hints = []struct {
			key  string
			desc string
		}{
			{"↑↓", "scroll"},
			{"g/G", "top/bottom"},
			{"/", "search"},
			{"f/Esc", "exit fullscreen"},
			{"?", "help"},
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

	// Build right-side status info
	var rightInfo string

	if m.lastError != "" {
		errMsg := m.lastError
		if len(errMsg) > 30 {
			errMsg = errMsg[:27] + "..."
		}
		rightInfo = styles.StatusErrorStyle.Render("error: " + errMsg)
	} else if m.preview.IsSearchMode() {
		searchIndicator := styles.SearchPromptStyle.Render("[searching]")
		if m.preview.FileName() != "" {
			rightInfo = styles.StatusValueStyle.Render(m.preview.FileName()) + " " + searchIndicator
		} else {
			rightInfo = searchIndicator
		}
	} else if m.preview.HasActiveSearch() || m.preview.HasSearchNoMatches() {
		var searchIndicator string
		if m.preview.HasActiveSearch() {
			matchInfo := itoa(m.preview.CurrentMatchIndex()) + "/" + itoa(m.preview.MatchCount())
			searchIndicator = styles.SearchMatchStyle.Render("[" + m.preview.SearchQuery() + ": " + matchInfo + "]")
		} else {
			searchIndicator = styles.SearchNoMatchStyle.Render("[" + m.preview.SearchQuery() + ": no matches]")
		}
		fileName := styles.StatusValueStyle.Render(m.preview.FileName())
		scrollPct := int(m.preview.ScrollPercent() * 100)
		scrollIndicator := styles.HelpDescStyle.Render("[" + itoa(scrollPct) + "%]")
		rightInfo = fileName + " " + searchIndicator + " " + scrollIndicator
	} else if m.preview.FilePath() != "" {
		scrollPct := int(m.preview.ScrollPercent() * 100)
		fileName := styles.StatusValueStyle.Render(m.preview.FileName())
		scrollIndicator := styles.HelpDescStyle.Render("[" + itoa(scrollPct) + "%]")
		fsIndicator := styles.StatusWatchingStyle.Render("[fullscreen]")
		rightInfo = fileName + " " + fsIndicator + " " + scrollIndicator
	}

	if rightInfo != "" {
		statusWidth := lipgloss.Width(statusContent)
		rightWidth := lipgloss.Width(rightInfo)
		spacerWidth := m.Width - statusWidth - rightWidth - 4
		if spacerWidth > 0 {
			statusContent = statusContent + strings.Repeat(" ", spacerWidth) + rightInfo
		}
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

// overlayHelp renders the help overlay on top of the base view
func (m Model) overlayHelp(baseView string) string {
	// Set help size based on current window
	m.help.SetSize(m.Width, m.Height)

	// Get help overlay content
	helpView := m.help.View()

	// Split base view into lines
	baseLines := strings.Split(baseView, "\n")

	// Split help view into lines
	helpLines := strings.Split(helpView, "\n")

	// Overlay help on top of base
	result := make([]string, len(baseLines))
	copy(result, baseLines)

	for i, helpLine := range helpLines {
		if i < len(result) && helpLine != "" {
			// Replace base line with help line where help has content
			// This creates a simple overlay effect
			result[i] = helpLine
		}
	}

	return strings.Join(result, "\n")
}
