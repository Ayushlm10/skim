package app

import (
	"path/filepath"
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
	var b strings.Builder

	// Show filter input if active
	if m.filterActive {
		filterLine := styles.FilterPromptStyle.Render("Filter: ") +
			styles.FilterInputStyle.Render(m.filterText+"_")
		b.WriteString(filterLine)
		b.WriteString("\n\n")
	}

	// Placeholder content for Phase 1
	if len(m.fileTreeItems) == 0 {
		// Show a sample tree structure as placeholder
		placeholderItems := []struct {
			name     string
			isDir    bool
			depth    int
			expanded bool
			selected bool
		}{
			{"docs/", true, 0, true, false},
			{"design.md", false, 1, false, true},
			{"plan.md", false, 1, false, false},
			{"implementation.md", false, 1, false, false},
			{"notes/", true, 0, false, false},
		}

		for i, item := range placeholderItems {
			line := m.renderTreeItem(item.name, item.isDir, item.depth, item.expanded, i == 1)
			b.WriteString(line)
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(styles.EmptyStateStyle.Render("(Phase 1 placeholder)"))
	}

	// Ensure content fills the available height
	content := b.String()
	lines := strings.Split(content, "\n")
	for len(lines) < height {
		lines = append(lines, "")
	}

	return strings.Join(lines[:height], "\n")
}

// renderTreeItem renders a single item in the file tree
func (m Model) renderTreeItem(name string, isDir bool, depth int, expanded, selected bool) string {
	var b strings.Builder

	// Add indentation
	indent := strings.Repeat("  ", depth)
	b.WriteString(indent)

	// Add tree indicator for directories
	if isDir {
		indicator := styles.TreeCollapsed
		if expanded {
			indicator = styles.TreeExpanded
		}
		b.WriteString(styles.TreeIndicatorStyle.Render(indicator + " "))
		b.WriteString(styles.DirectoryStyle.Render(name))
	} else {
		b.WriteString("  ") // Align with directory indicators
		style := styles.FileStyle
		if selected {
			style = styles.SelectedItemStyle
		}
		b.WriteString(style.Render(name))
		if selected {
			b.WriteString(" " + styles.TreeIndicatorStyle.Render(styles.SelectedMark))
		}
	}

	return b.String()
}

// renderPreview renders the markdown preview content
func (m Model) renderPreview(width, height int) string {
	if m.previewPath == "" {
		// Show placeholder for Phase 1
		placeholder := []string{
			"# Welcome to Local MD Viewer",
			"",
			"Select a markdown file from the left panel to preview it here.",
			"",
			"## Quick Start",
			"",
			"- Use **j/k** or **arrow keys** to navigate",
			"- Press **Enter** to open a file",
			"- Press **Tab** to switch panels",
			"- Press **/** to filter files",
			"- Press **q** to quit",
			"",
			"---",
			"",
			"_Phase 1 placeholder content_",
		}

		content := strings.Join(placeholder, "\n")
		return styles.PreviewStyle.Width(width).Render(content)
	}

	// Render actual content (to be enhanced in Phase 3 with Glamour)
	content := m.previewContent
	if content == "" {
		content = styles.NoPreviewStyle.Render("Loading " + filepath.Base(m.previewPath) + "...")
	}

	return styles.PreviewStyle.Width(width).Render(content)
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	// Build help hints
	hints := []struct {
		key  string
		desc string
	}{
		{"↑↓", "navigate"},
		{"⏎", "open"},
		{"/", "filter"},
		{"Tab", "switch"},
		{"q", "quit"},
	}

	var parts []string
	for _, h := range hints {
		part := styles.HelpKeyStyle.Render(h.key) + " " + styles.HelpDescStyle.Render(h.desc)
		parts = append(parts, part)
	}

	separator := styles.HelpSeparatorStyle.Render("  │  ")
	statusContent := strings.Join(parts, separator)

	return styles.StatusBarStyle.
		Width(m.Width).
		Render(statusContent)
}
