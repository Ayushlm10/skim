package preview

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/athakur/local-md/internal/styles"
	"github.com/charmbracelet/bubbles/textinput"
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

	// Search state (Phase 7.2)
	searchMode   bool            // Whether search input is active
	searchInput  textinput.Model // Text input for search query
	searchQuery  string          // Current search query (after Enter)
	matches      []int           // Line numbers in rawContent that match
	currentMatch int             // Index into matches slice (0-based)
}

// New creates a new preview component
func New(width, height int) Model {
	// Create viewport
	vp := viewport.New(width, height)
	vp.Style = lipgloss.NewStyle()
	vp.SetContent("")

	// Enable mouse wheel scrolling with snappy delta
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 5

	// Create renderer (with padding for viewport)
	renderer, _ := NewRenderer(width - 4)

	// Create search input
	ti := textinput.New()
	ti.Placeholder = "search..."
	ti.Prompt = "/"
	ti.PromptStyle = styles.FilterPromptStyle
	ti.TextStyle = styles.FilterInputStyle
	ti.Cursor.Style = styles.FilterCursorStyle
	ti.CharLimit = 100
	ti.Width = width - 10

	return Model{
		viewport:     vp,
		renderer:     renderer,
		width:        width,
		height:       height,
		focused:      false,
		ready:        true,
		searchInput:  ti,
		searchMode:   false,
		matches:      nil,
		currentMatch: 0,
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
		// Clear any existing search when loading a new file
		m.clearSearch()
		m.searchMode = false
		m.searchInput.Blur()

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
	// Handle search mode input
	if m.searchMode {
		return m.handleSearchKey(msg)
	}

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

	case "/":
		// Enter search mode
		m.searchMode = true
		m.searchInput.Focus()
		m.searchInput.SetValue("")
		return m, textinput.Blink

	case "n":
		// Next match
		if len(m.matches) > 0 {
			m.currentMatch = (m.currentMatch + 1) % len(m.matches)
			(&m).scrollToCurrentMatch()
		}
		return m, nil

	case "N":
		// Previous match
		if len(m.matches) > 0 {
			m.currentMatch--
			if m.currentMatch < 0 {
				m.currentMatch = len(m.matches) - 1
			}
			(&m).scrollToCurrentMatch()
		}
		return m, nil

	case "esc":
		// Clear search if active
		if m.searchQuery != "" {
			(&m).clearSearch()
		}
		return m, nil
	}

	return m, nil
}

// handleSearchKey handles keys when search input is active
func (m Model) handleSearchKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Execute search and exit search mode
		query := m.searchInput.Value()
		m.searchMode = false
		m.searchInput.Blur()
		if query != "" {
			m.searchQuery = query
			// Use pointer to ensure modifications persist
			(&m).performSearch()
			if len(m.matches) > 0 {
				m.currentMatch = 0
				(&m).scrollToCurrentMatch()
			}
		}
		return m, nil

	case "esc":
		// Cancel search mode without searching
		m.searchMode = false
		m.searchInput.Blur()
		m.searchInput.SetValue("")
		return m, nil
	}

	// Forward other keys to textinput
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

// performSearch searches rawContent for the query and stores matching line numbers
func (m *Model) performSearch() {
	m.matches = nil
	m.currentMatch = 0

	if m.searchQuery == "" || m.rawContent == "" {
		return
	}

	// Case-insensitive search in raw content for line tracking (for scrolling)
	query := strings.ToLower(m.searchQuery)
	lines := strings.Split(m.rawContent, "\n")

	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), query) {
			m.matches = append(m.matches, i)
		}
	}

	// Apply highlighting to rendered content
	m.applySearchHighlight()
}

// countVisibleMatches counts occurrences in visible text (for status bar display)
func (m Model) countVisibleMatches() int {
	if m.searchQuery == "" || m.renderedContent == "" {
		return 0
	}

	// Strip ANSI codes and count matches
	visible := stripANSI(m.renderedContent)
	lowerVisible := strings.ToLower(visible)
	lowerQuery := strings.ToLower(m.searchQuery)

	count := 0
	searchStart := 0
	for {
		idx := strings.Index(lowerVisible[searchStart:], lowerQuery)
		if idx == -1 {
			break
		}
		count++
		searchStart = searchStart + idx + 1
	}
	return count
}

// stripANSI removes ANSI escape sequences from a string
func stripANSI(s string) string {
	result := &strings.Builder{}
	result.Grow(len(s))

	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && s[j] != 'm' {
				j++
			}
			if j < len(s) {
				j++
				i = j
				continue
			}
		}
		result.WriteByte(s[i])
		i++
	}
	return result.String()
}

// applySearchHighlight highlights search matches in the rendered content
func (m *Model) applySearchHighlight() {
	if m.searchQuery == "" || m.renderedContent == "" {
		m.viewport.SetContent(m.renderedContent)
		return
	}

	highlighted := highlightMatches(m.renderedContent, m.searchQuery)
	// Debug: Log that highlighting is being applied
	// Remove this after debugging
	_ = highlighted // Ensure we're using the variable
	m.viewport.SetContent(highlighted)
}

// highlightMatches applies reverse video highlighting to all occurrences of query
// It handles ANSI escape sequences properly by tracking visible text positions
func highlightMatches(content, query string) string {
	if query == "" {
		return content
	}

	// ANSI codes for reverse video (swap fg/bg)
	const reverseOn = "\x1b[7m"
	const reverseOff = "\x1b[27m"

	lowerQuery := strings.ToLower(query)
	queryLen := len(query)

	// First, find all match positions in the visible text (ignoring ANSI codes)
	// Build a map of visible position -> original position
	visibleText := &strings.Builder{}
	visibleToOriginal := make([]int, 0, len(content))

	i := 0
	for i < len(content) {
		// Skip ANSI escape sequences
		if content[i] == '\x1b' && i+1 < len(content) && content[i+1] == '[' {
			j := i + 2
			for j < len(content) && content[j] != 'm' {
				j++
			}
			if j < len(content) {
				j++ // Include the 'm'
				i = j
				continue
			}
		}

		// Track this visible character's original position
		visibleToOriginal = append(visibleToOriginal, i)
		visibleText.WriteByte(content[i])
		i++
	}

	visible := visibleText.String()
	lowerVisible := strings.ToLower(visible)

	// Find all match start positions in visible text
	var matchStarts []int
	searchStart := 0
	for {
		idx := strings.Index(lowerVisible[searchStart:], lowerQuery)
		if idx == -1 {
			break
		}
		matchStarts = append(matchStarts, searchStart+idx)
		searchStart = searchStart + idx + 1
	}

	if len(matchStarts) == 0 {
		return content
	}

	// Build result with highlights inserted at correct positions
	result := &strings.Builder{}
	result.Grow(len(content) + len(matchStarts)*20)

	// Create a set of positions where we need to insert highlight codes
	highlightOn := make(map[int]bool)  // original positions to insert reverseOn
	highlightOff := make(map[int]bool) // original positions to insert reverseOff

	for _, visStart := range matchStarts {
		if visStart < len(visibleToOriginal) {
			highlightOn[visibleToOriginal[visStart]] = true
		}
		visEnd := visStart + queryLen
		if visEnd <= len(visibleToOriginal) {
			// Insert reverseOff AFTER the last character of match
			if visEnd < len(visibleToOriginal) {
				highlightOff[visibleToOriginal[visEnd]] = true
			} else {
				// Match extends to end of visible content
				highlightOff[len(content)] = true
			}
		}
	}

	// Now rebuild content with highlights
	for i := 0; i <= len(content); i++ {
		// Check if we need to turn off highlight here
		if highlightOff[i] {
			result.WriteString(reverseOff)
		}
		// Check if we need to turn on highlight here
		if highlightOn[i] {
			result.WriteString(reverseOn)
		}

		if i < len(content) {
			result.WriteByte(content[i])
		}
	}

	return result.String()
}

// scrollToCurrentMatch scrolls the viewport to show the current match
func (m *Model) scrollToCurrentMatch() {
	if len(m.matches) == 0 || m.currentMatch < 0 || m.currentMatch >= len(m.matches) {
		return
	}

	// Get the line number in raw content
	rawLine := m.matches[m.currentMatch]

	// The viewport shows rendered content, which may have different line counts
	// We need to estimate where in the rendered content this line appears
	// Since Glamour can add blank lines, headers, etc., we estimate by ratio
	rawLineCount := len(strings.Split(m.rawContent, "\n"))
	renderedLineCount := m.viewport.TotalLineCount()

	if rawLineCount == 0 {
		return
	}

	// Estimate the rendered line position
	ratio := float64(rawLine) / float64(rawLineCount)
	targetLine := int(ratio * float64(renderedLineCount))

	// Center the match in the viewport
	halfView := m.viewport.VisibleLineCount() / 2
	scrollTo := targetLine - halfView
	if scrollTo < 0 {
		scrollTo = 0
	}

	m.viewport.SetYOffset(scrollTo)
}

// clearSearch clears the search state
func (m *Model) clearSearch() {
	m.searchQuery = ""
	m.matches = nil
	m.currentMatch = 0
	m.searchInput.SetValue("")

	// Restore original rendered content (without highlights)
	if m.renderedContent != "" {
		m.viewport.SetContent(m.renderedContent)
	}
}

// HandleMouse handles mouse input (scrolling)
func (m Model) HandleMouse(msg tea.MouseMsg) (Model, tea.Cmd) {
	// Forward to viewport - it handles mouse wheel natively
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the component
func (m Model) View() string {
	if !m.ready {
		return "Initializing preview..."
	}

	if m.filePath == "" {
		return m.renderWelcome()
	}

	viewportContent := m.viewport.View()

	// If search mode is active, show the search input at the bottom
	if m.searchMode {
		// Calculate available lines for viewport
		lines := strings.Split(viewportContent, "\n")

		// Reserve one line for search input
		maxLines := m.height - 1
		if len(lines) > maxLines {
			lines = lines[:maxLines]
		}

		// Build search input line
		searchLine := m.renderSearchInput()

		// Pad with empty lines if needed to push search to bottom
		for len(lines) < maxLines {
			lines = append(lines, "")
		}
		lines = append(lines, searchLine)

		return strings.Join(lines, "\n")
	}

	return viewportContent
}

// renderSearchInput renders the search input line
func (m Model) renderSearchInput() string {
	inputStyle := lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#2A2A2A"}).
		Foreground(styles.Highlight).
		Width(m.width - 2)

	return inputStyle.Render(m.searchInput.View())
}

// renderWelcome renders the welcome message when no file is selected
func (m Model) renderWelcome() string {
	welcome := []string{
		"",
		"# Welcome to Local MD Viewer",
		"",
		"Select a markdown file from the left panel to preview it here.",
		"",
		"## Quick Start",
		"",
		"| Key | Action |",
		"|-----|--------|",
		"| `↑` `↓` / `j` `k` | Navigate files |",
		"| `Enter` | Open file or toggle folder |",
		"| `Tab` | Switch between panels |",
		"| `/` | Filter files |",
		"| `?` | Show keyboard shortcuts |",
		"| `q` | Quit |",
		"",
		"---",
		"",
		"*Files are automatically reloaded when modified.*",
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
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.Error).
		Bold(true).
		MarginBottom(1)

	messageStyle := lipgloss.NewStyle().
		Foreground(styles.Muted)

	hintStyle := lipgloss.NewStyle().
		Foreground(styles.Subtle).
		Italic(true).
		MarginTop(1)

	title := titleStyle.Render("Unable to Load File")
	message := messageStyle.Render(err.Error())
	hint := hintStyle.Render("Select another file or check the file path.")

	content := title + "\n\n" + message + "\n\n" + hint

	return lipgloss.NewStyle().
		Padding(2, 3).
		Render(content)
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

	// Update search input width
	m.searchInput.Width = width - 10

	// Re-render content if we have any
	if m.rawContent != "" && m.renderer != nil {
		rendered, err := m.renderer.Render(m.rawContent)
		if err == nil {
			m.renderedContent = rendered
			// If there's an active search, apply highlighting
			if m.searchQuery != "" {
				m.applySearchHighlight()
			} else {
				m.viewport.SetContent(rendered)
			}
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

// IsSearchMode returns whether search input is active
func (m Model) IsSearchMode() bool {
	return m.searchMode
}

// SearchQuery returns the current search query
func (m Model) SearchQuery() string {
	return m.searchQuery
}

// MatchCount returns the number of matches found (in visible rendered text)
func (m Model) MatchCount() int {
	return m.countVisibleMatches()
}

// CurrentMatchIndex returns the current match index (1-based for display)
func (m Model) CurrentMatchIndex() int {
	if len(m.matches) == 0 {
		return 0
	}
	return m.currentMatch + 1
}

// HasActiveSearch returns whether there's an active search with results
func (m Model) HasActiveSearch() bool {
	return m.searchQuery != "" && len(m.matches) > 0
}

// HasSearchNoMatches returns whether there's a search query with no results
func (m Model) HasSearchNoMatches() bool {
	return m.searchQuery != "" && len(m.matches) == 0
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
