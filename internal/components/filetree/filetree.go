package filetree

import (
	"github.com/Ayushlm10/skim/internal/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Messages for communication with parent

// FileSelectedMsg is sent when a file is selected
type FileSelectedMsg struct {
	Path string
}

// DirectoryToggledMsg is sent when a directory is expanded/collapsed
type DirectoryToggledMsg struct {
	Path     string
	Expanded bool
}

// FilterChangedMsg is sent when the filter state changes
type FilterChangedMsg struct {
	Active bool
	Value  string
}

// Model is the file tree component model
type Model struct {
	// Root path being scanned
	RootPath string

	// All items in flattened display order
	items []*Item

	// The underlying list component
	list list.Model

	// Scan options
	scanOptions ScanOptions

	// Dimensions
	width  int
	height int

	// Focus state
	focused bool
}

// New creates a new file tree component
func New(rootPath string, width, height int) Model {
	// Create delegate
	delegate := NewItemDelegate()

	// Create list with empty items initially
	l := list.New([]list.Item{}, delegate, width, height)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true) // Show filter input when filtering
	l.Styles.NoItems = styles.EmptyStateStyle
	l.DisableQuitKeybindings()

	// Apply custom styles
	l.Styles = treeListStyles()

	// Customize key map to avoid conflicts
	l.KeyMap.Quit.SetEnabled(false)

	return Model{
		RootPath:    rootPath,
		items:       nil,
		list:        l,
		scanOptions: DefaultScanOptions(),
		width:       width,
		height:      height,
		focused:     true,
	}
}

// Init initializes the component and starts scanning
func (m Model) Init() tea.Cmd {
	return m.scanRoot()
}

// scanRoot scans the root directory
func (m Model) scanRoot() tea.Cmd {
	return func() tea.Msg {
		items, err := ScanDirectory(m.RootPath, m.scanOptions)
		if err != nil {
			return scanErrorMsg{err}
		}
		return scanCompleteMsg{items}
	}
}

// scanErrorMsg is sent when scanning fails
type scanErrorMsg struct {
	err error
}

// scanCompleteMsg is sent when scanning completes
type scanCompleteMsg struct {
	items []*Item
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(m.width, m.height)

	case scanCompleteMsg:
		m.items = msg.items
		m.rebuildList()

	case scanErrorMsg:
		// Handle error - could show in status
		return m, nil

	case tea.KeyMsg:
		if m.focused {
			return m.handleKey(msg)
		}
	}

	// Forward to list for filtering etc.
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// handleKey handles keyboard input
func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	// When filtering, delegate most keys to the list
	if m.IsFiltering() {
		switch msg.String() {
		case "esc":
			// Exit filter mode and clear filter
			m.list.ResetFilter()
			return m, func() tea.Msg {
				return FilterChangedMsg{Active: false, Value: ""}
			}
		case "enter":
			// Accept filter and select item
			if m.list.FilterState() == list.Filtering {
				// Let list handle the enter to accept filter
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, tea.Batch(cmd, func() tea.Msg {
					return FilterChangedMsg{Active: false, Value: m.list.FilterValue()}
				})
			}
			return m.handleSelect()
		default:
			// Forward to list for filter input
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			// Notify parent of filter changes
			return m, tea.Batch(cmd, func() tea.Msg {
				return FilterChangedMsg{Active: true, Value: m.list.FilterValue()}
			})
		}
	}

	switch msg.String() {
	case "enter":
		return m.handleSelect()

	case "up", "k":
		m.list.CursorUp()
		return m, nil

	case "down", "j":
		m.list.CursorDown()
		return m, nil

	case "/":
		// Enter filter mode - delegate to list's filtering
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, tea.Batch(cmd, func() tea.Msg {
			return FilterChangedMsg{Active: true, Value: ""}
		})

	case "esc":
		// Clear filter if there is one
		if m.list.FilterValue() != "" {
			m.list.ResetFilter()
			return m, func() tea.Msg {
				return FilterChangedMsg{Active: false, Value: ""}
			}
		}
		return m, nil
	}

	return m, nil
}

// handleSelect handles enter key - opens file or toggles directory
func (m Model) handleSelect() (Model, tea.Cmd) {
	selected := m.list.SelectedItem()
	if selected == nil {
		return m, nil
	}

	item, ok := selected.(*Item)
	if !ok {
		return m, nil
	}

	if item.IsDir {
		// Toggle directory expand/collapse
		item.Toggle()

		// Load children if expanding and not yet loaded
		if item.Expanded && !item.HasChildren() {
			_ = ScanChildren(item, m.scanOptions)
		}

		m.rebuildList()

		return m, func() tea.Msg {
			return DirectoryToggledMsg{
				Path:     item.Path,
				Expanded: item.Expanded,
			}
		}
	}

	// File selected - send message to parent
	return m, func() tea.Msg {
		return FileSelectedMsg{Path: item.Path}
	}
}

// rebuildList reconstructs the flattened list from tree structure
func (m *Model) rebuildList() {
	var flatItems []list.Item

	var flatten func(items []*Item)
	flatten = func(items []*Item) {
		for _, item := range items {
			flatItems = append(flatItems, item)
			if item.IsDir && item.Expanded && item.HasChildren() {
				flatten(item.Children)
			}
		}
	}

	flatten(m.items)
	m.list.SetItems(flatItems)
}

// View renders the component
func (m Model) View() string {
	if len(m.items) == 0 {
		return m.renderEmptyState()
	}

	return m.list.View()
}

// renderEmptyState renders a helpful message when no files are found
func (m Model) renderEmptyState() string {
	title := styles.EmptyStateTitleStyle.Render("No Markdown Files")
	hint := styles.EmptyStateHintStyle.Render("This directory contains no .md files.\nTry navigating to a different directory.")

	// Center the content
	content := title + "\n\n" + hint

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

// SetSize updates the component size
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

// SetFocused sets the focus state
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

// IsFocused returns the focus state
func (m Model) IsFocused() bool {
	return m.focused
}

// SelectedItem returns the currently selected item
func (m Model) SelectedItem() *Item {
	selected := m.list.SelectedItem()
	if selected == nil {
		return nil
	}
	item, _ := selected.(*Item)
	return item
}

// ItemCount returns the total number of visible items
func (m Model) ItemCount() int {
	return len(m.list.Items())
}

// FilterState returns the current filter text and active state
func (m Model) FilterState() (string, bool) {
	return m.list.FilterValue(), m.list.FilterState() == list.Filtering
}

// IsFiltering returns true if the list is in filtering mode
func (m Model) IsFiltering() bool {
	return m.list.FilterState() == list.Filtering
}

// FilterValue returns the current filter value
func (m Model) FilterValue() string {
	return m.list.FilterValue()
}

// HasActiveFilter returns true if there's a non-empty filter applied
func (m Model) HasActiveFilter() bool {
	return m.list.FilterValue() != ""
}

// treeListStyles returns custom styles for the tree list
func treeListStyles() list.Styles {
	s := list.DefaultStyles()

	s.Title = lipgloss.NewStyle().
		Foreground(styles.Highlight).
		Bold(true)

	// Filter input styling - minimal/editorial aesthetic
	s.FilterPrompt = styles.FilterPromptStyle
	s.FilterCursor = styles.FilterCursorStyle

	// Pagination and help styling
	s.PaginationStyle = lipgloss.NewStyle().Foreground(styles.Muted)
	s.HelpStyle = lipgloss.NewStyle().Foreground(styles.Muted)

	// No items message
	s.NoItems = styles.EmptyStateStyle

	return s
}
