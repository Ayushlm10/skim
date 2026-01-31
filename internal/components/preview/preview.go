package preview

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/athakur/local-md/internal/styles"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Messages for communication with parent

// FileLoadedMsg is sent when a file has been loaded and rendered
type FileLoadedMsg struct {
	Path    string
	Content string
	Error   error
}

// Model is the preview component model
type Model struct {
	// Current file being previewed
	filePath string

	// Raw markdown content
	rawContent string

	// Rendered content
	renderedContent string

	// Viewport for scrolling
	viewport viewport.Model

	// Renderer for markdown
	renderer *Renderer

	// Dimensions
	width  int
	height int

	// Focus state
	focused bool

	// Ready state (viewport initialized)
	ready bool

	// Error state
	err error
}

// New creates a new preview component
func New(width, height int) Model {
	// Create viewport
	vp := viewport.New(width, height)
	vp.Style = lipgloss.NewStyle()
	vp.SetContent("")

	// Create renderer (with padding for viewport)
	renderer, _ := NewRenderer(width - 4)

	return Model{
		viewport: vp,
		renderer: renderer,
		width:    width,
		height:   height,
		focused:  false,
		ready:    true,
	}
}

// Init initializes the component
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil

	case FileLoadedMsg:
		if msg.Error != nil {
			m.err = msg.Error
			m.renderedContent = ""
			m.viewport.SetContent(m.renderError(msg.Error))
		} else {
			m.filePath = msg.Path
			m.rawContent = msg.Content
			m.err = nil

			// Render the content
			rendered, err := m.renderer.Render(msg.Content)
			if err != nil {
				m.err = err
				m.viewport.SetContent(m.renderError(err))
			} else {
				m.renderedContent = rendered
				m.viewport.SetContent(rendered)
				m.viewport.GotoTop()
			}
		}
		return m, nil
	}

	// Forward to viewport when focused
	if m.focused {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// HandleKey handles keyboard input when focused
func (m Model) HandleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.viewport.LineUp(1)
		return m, nil

	case "down", "j":
		m.viewport.LineDown(1)
		return m, nil

	case "pgup", "ctrl+u":
		m.viewport.HalfViewUp()
		return m, nil

	case "pgdown", "ctrl+d":
		m.viewport.HalfViewDown()
		return m, nil

	case "g":
		m.viewport.GotoTop()
		return m, nil

	case "G":
		m.viewport.GotoBottom()
		return m, nil

	case "home":
		m.viewport.GotoTop()
		return m, nil

	case "end":
		m.viewport.GotoBottom()
		return m, nil
	}

	return m, nil
}

// View renders the component
func (m Model) View() string {
	if !m.ready {
		return "Initializing preview..."
	}

	if m.filePath == "" {
		return m.renderWelcome()
	}

	return m.viewport.View()
}

// renderWelcome renders the welcome message when no file is selected
func (m Model) renderWelcome() string {
	welcome := []string{
		"",
		"  # Welcome to Local MD Viewer",
		"",
		"  Select a markdown file from the left panel",
		"  to preview it here.",
		"",
		"  ## Quick Start",
		"",
		"  - Use **j/k** or **arrow keys** to navigate",
		"  - Press **Enter** to open a file",
		"  - Press **Tab** to switch panels",
		"  - Press **/** to filter files",
		"  - Press **q** to quit",
		"",
	}

	// Render welcome message through glamour for consistent styling
	content := strings.Join(welcome, "\n")
	if m.renderer != nil {
		rendered, err := m.renderer.Render(content)
		if err == nil {
			return rendered
		}
	}

	return styles.PreviewStyle.Render(content)
}

// renderError renders an error message
func (m Model) renderError(err error) string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Padding(1, 2)

	return errorStyle.Render(fmt.Sprintf("Error loading file:\n\n%s", err.Error()))
}

// SetSize updates the component size
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height

	// Update renderer width for word wrap
	if m.renderer != nil {
		_ = m.renderer.SetWidth(width - 4) // Account for padding
	}

	// Re-render content if we have any
	if m.rawContent != "" && m.renderer != nil {
		rendered, err := m.renderer.Render(m.rawContent)
		if err == nil {
			m.renderedContent = rendered
			m.viewport.SetContent(rendered)
		}
	}
}

// SetFocused sets the focus state
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

// IsFocused returns the focus state
func (m Model) IsFocused() bool {
	return m.focused
}

// FilePath returns the current file path
func (m Model) FilePath() string {
	return m.filePath
}

// FileName returns just the file name
func (m Model) FileName() string {
	if m.filePath == "" {
		return ""
	}
	return filepath.Base(m.filePath)
}

// ScrollPercent returns the current scroll position as a percentage
func (m Model) ScrollPercent() float64 {
	return m.viewport.ScrollPercent()
}

// AtTop returns true if scrolled to the top
func (m Model) AtTop() bool {
	return m.viewport.AtTop()
}

// AtBottom returns true if scrolled to the bottom
func (m Model) AtBottom() bool {
	return m.viewport.AtBottom()
}

// TotalLines returns the total number of lines in the content
func (m Model) TotalLines() int {
	return m.viewport.TotalLineCount()
}

// VisibleLines returns the number of visible lines
func (m Model) VisibleLines() int {
	return m.viewport.VisibleLineCount()
}

// LoadFile creates a command to load a file
func LoadFile(path string) tea.Cmd {
	return func() tea.Msg {
		content, err := os.ReadFile(path)
		if err != nil {
			return FileLoadedMsg{
				Path:  path,
				Error: err,
			}
		}

		return FileLoadedMsg{
			Path:    path,
			Content: string(content),
		}
	}
}
