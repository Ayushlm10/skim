package help

import (
	"strings"

	"github.com/Ayushlm10/skim/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyBinding represents a single key binding
type KeyBinding struct {
	Key  string
	Desc string
}

// KeySection represents a group of related key bindings
type KeySection struct {
	Title    string
	Bindings []KeyBinding
}

// Model is the help overlay component
type Model struct {
	visible bool
	width   int
	height  int
}

// New creates a new help overlay
func New() Model {
	return Model{
		visible: false,
	}
}

// Toggle toggles the help overlay visibility
func (m *Model) Toggle() {
	m.visible = !m.visible
}

// Show shows the help overlay
func (m *Model) Show() {
	m.visible = true
}

// Hide hides the help overlay
func (m *Model) Hide() {
	m.visible = false
}

// IsVisible returns whether the overlay is visible
func (m Model) IsVisible() bool {
	return m.visible
}

// SetSize updates the overlay dimensions
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.visible {
			switch msg.String() {
			case "?", "esc", "q", "enter":
				m.visible = false
				return m, nil
			}
		}
	}
	return m, nil
}

// View renders the help overlay
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	// Define key bindings sections
	sections := []KeySection{
		{
			Title: "Navigation",
			Bindings: []KeyBinding{
				{Key: "↑ / k", Desc: "Move up"},
				{Key: "↓ / j", Desc: "Move down"},
				{Key: "Enter", Desc: "Open file / Toggle folder"},
				{Key: "Tab", Desc: "Switch panel focus"},
				{Key: "i", Desc: "Toggle ignored directories"},
			},
		},
		{
			Title: "Preview",
			Bindings: []KeyBinding{
				{Key: "PgUp / Ctrl+u", Desc: "Scroll up half page"},
				{Key: "PgDn / Ctrl+d", Desc: "Scroll down half page"},
				{Key: "g", Desc: "Go to top"},
				{Key: "G", Desc: "Go to bottom"},
				{Key: "f", Desc: "Toggle fullscreen mode"},
			},
		},
		{
			Title: "File Tree Filter",
			Bindings: []KeyBinding{
				{Key: "/", Desc: "Enter filter mode (file tree)"},
				{Key: "Esc", Desc: "Clear filter / Exit"},
				{Key: "Enter", Desc: "Accept filter"},
			},
		},
		{
			Title: "Preview Search",
			Bindings: []KeyBinding{
				{Key: "/", Desc: "Search in content (preview)"},
				{Key: "n", Desc: "Next match"},
				{Key: "N", Desc: "Previous match"},
				{Key: "Esc", Desc: "Clear search"},
			},
		},
		{
			Title: "General",
			Bindings: []KeyBinding{
				{Key: "?", Desc: "Toggle this help"},
				{Key: "q / Ctrl+c", Desc: "Quit"},
			},
		},
	}

	return m.renderOverlay(sections)
}

// renderOverlay renders the help content in a centered overlay
func (m Model) renderOverlay(sections []KeySection) string {
	// Style definitions
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Highlight).
		Bold(true).
		MarginBottom(1)

	sectionTitleStyle := lipgloss.NewStyle().
		Foreground(styles.Accent).
		Bold(true).
		MarginTop(1)

	keyStyle := lipgloss.NewStyle().
		Foreground(styles.Highlight).
		Width(16)

	descStyle := lipgloss.NewStyle().
		Foreground(styles.Muted)

	// Build help content
	var content strings.Builder

	content.WriteString(titleStyle.Render("Keyboard Shortcuts"))
	content.WriteString("\n")

	for _, section := range sections {
		content.WriteString(sectionTitleStyle.Render(section.Title))
		content.WriteString("\n")

		for _, binding := range section.Bindings {
			line := keyStyle.Render(binding.Key) + descStyle.Render(binding.Desc)
			content.WriteString(line)
			content.WriteString("\n")
		}
	}

	content.WriteString("\n")
	content.WriteString(lipgloss.NewStyle().
		Foreground(styles.Subtle).
		Italic(true).
		Render("Press ? or Esc to close"))

	// Create the overlay box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Accent).
		Padding(1, 3).
		Background(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#1A1A1A"})

	helpBox := boxStyle.Render(content.String())

	// Calculate centering
	boxWidth := lipgloss.Width(helpBox)
	boxHeight := lipgloss.Height(helpBox)

	// Center horizontally
	leftPadding := (m.width - boxWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	// Center vertically
	topPadding := (m.height - boxHeight) / 2
	if topPadding < 0 {
		topPadding = 0
	}

	// Create backdrop with semi-transparent effect (using spaces)
	// For terminal, we just position the box
	centeredBox := lipgloss.NewStyle().
		MarginLeft(leftPadding).
		MarginTop(topPadding).
		Render(helpBox)

	return centeredBox
}
