package filetree

import (
	"fmt"
	"io"
	"strings"

	"github.com/Ayushlm10/skim/internal/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// ItemDelegate handles rendering of tree items in the list
type ItemDelegate struct {
	// ShowIndicator shows the selection indicator
	ShowIndicator bool
}

// NewItemDelegate creates a new delegate with default settings
func NewItemDelegate() ItemDelegate {
	return ItemDelegate{
		ShowIndicator: true,
	}
}

// Height returns the height of each item
func (d ItemDelegate) Height() int {
	return 1
}

// Spacing returns the spacing between items
func (d ItemDelegate) Spacing() int {
	return 0
}

// Update handles item-level updates
func (d ItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render renders a single tree item
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(*Item)
	if !ok {
		return
	}

	var b strings.Builder

	// Add indentation based on depth
	indent := strings.Repeat("  ", item.Depth)
	b.WriteString(indent)

	isSelected := index == m.Index()

	// Render based on item type
	if item.IsDir {
		// Directory with expand/collapse indicator
		indicator := styles.TreeCollapsed
		if item.Expanded {
			indicator = styles.TreeExpanded
		}
		b.WriteString(styles.TreeIndicatorStyle.Render(indicator + " "))

		if isSelected {
			b.WriteString(styles.SelectedDirectoryStyle.Render(item.DisplayName()))
			if d.ShowIndicator {
				b.WriteString(" " + styles.TreeIndicatorStyle.Render(styles.SelectedMark))
			}
		} else {
			b.WriteString(styles.DirectoryStyle.Render(item.DisplayName()))
		}
	} else {
		// File with proper alignment
		b.WriteString("  ") // Align with directory indicators

		if isSelected {
			b.WriteString(styles.SelectedItemStyle.Render(item.Name))
			if d.ShowIndicator {
				b.WriteString(" " + styles.TreeIndicatorStyle.Render(styles.SelectedMark))
			}
		} else {
			b.WriteString(styles.FileStyle.Render(item.Name))
		}
	}

	fmt.Fprint(w, b.String())
}

// ShortHelp returns short help text for the list
func (d ItemDelegate) ShortHelp() []string {
	return []string{"j/k: navigate", "enter: open/toggle", "/: filter"}
}

// FullHelp returns full help text for the list
func (d ItemDelegate) FullHelp() [][]string {
	return [][]string{
		{"j/k, ↑/↓: navigate", "enter: open file / toggle dir"},
		{"/: filter", "esc: clear filter"},
	}
}
