package styles

import "github.com/charmbracelet/lipgloss"

// Minimal/Editorial color palette - muted and sophisticated
var (
	// Base colors - adaptive to terminal
	Subtle    = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	Highlight = lipgloss.AdaptiveColor{Light: "#2D2D2D", Dark: "#E0E0E0"}
	Accent    = lipgloss.AdaptiveColor{Light: "#4A7C59", Dark: "#7EB38E"}
	Muted     = lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"}
	Border    = lipgloss.AdaptiveColor{Light: "#DDDDDD", Dark: "#3A3A3A"}

	// Panel widths (ratios)
	FileTreeRatio = 0.25
	PreviewRatio  = 0.75
)

// Header styles
var (
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Highlight).
			Padding(0, 1)

	HeaderPathStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Padding(0, 1)
)

// Panel styles
var (
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border)

	FocusedPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(Accent)
)

// File tree styles
var (
	FileTreeStyle = lipgloss.NewStyle().
			Padding(0, 1)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Highlight).
				Background(lipgloss.AdaptiveColor{Light: "#EEEEEE", Dark: "#333333"})

	DirectoryStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	FileStyle = lipgloss.NewStyle().
			Foreground(Highlight)

	TreeIndicatorStyle = lipgloss.NewStyle().
				Foreground(Subtle)
)

// Preview styles
var (
	PreviewStyle = lipgloss.NewStyle().
			Padding(0, 1)

	NoPreviewStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true).
			Padding(1, 2)
)

// Status bar styles
var (
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Background(lipgloss.AdaptiveColor{Light: "#F5F5F5", Dark: "#1A1A1A"}).
			Padding(0, 1)

	StatusKeyStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	StatusValueStyle = lipgloss.NewStyle().
				Foreground(Highlight)
)

// Filter input styles
var (
	FilterPromptStyle = lipgloss.NewStyle().
				Foreground(Accent)

	FilterInputStyle = lipgloss.NewStyle().
				Foreground(Highlight)
)

// Help styles
var (
	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Accent)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(Muted)

	HelpSeparatorStyle = lipgloss.NewStyle().
				Foreground(Subtle)
)

// Empty state styles
var (
	EmptyStateStyle = lipgloss.NewStyle().
		Foreground(Muted).
		Italic(true).
		Align(lipgloss.Center)
)

// Tree indicators
const (
	TreeExpanded  = "▾"
	TreeCollapsed = "▸"
	TreeBranch    = "├──"
	TreeLastItem  = "└──"
	TreeVertical  = "│  "
	TreeEmpty     = "   "
	SelectedMark  = "◀"
)
