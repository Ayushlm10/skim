package styles

import "github.com/charmbracelet/lipgloss"

// Minimal/Editorial color palette - muted and sophisticated
var (
	// Base colors - adaptive to terminal with improved contrast
	Subtle    = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#626262"}
	Highlight = lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#F0F0F0"}
	Accent    = lipgloss.AdaptiveColor{Light: "#2D6A4F", Dark: "#95D5B2"}
	Muted     = lipgloss.AdaptiveColor{Light: "#6B6B6B", Dark: "#9A9A9A"}
	Border    = lipgloss.AdaptiveColor{Light: "#E0E0E0", Dark: "#404040"}

	// Secondary colors for visual hierarchy
	AccentDim  = lipgloss.AdaptiveColor{Light: "#40916C", Dark: "#74C69D"}
	Warning    = lipgloss.AdaptiveColor{Light: "#B07D2B", Dark: "#F4A261"}
	Error      = lipgloss.AdaptiveColor{Light: "#C94C4C", Dark: "#FF6B6B"}
	Success    = lipgloss.AdaptiveColor{Light: "#2D6A4F", Dark: "#52B788"}
	Background = lipgloss.AdaptiveColor{Light: "#FAFAFA", Dark: "#1A1A1A"}

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
			Background(lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#252525"}).
			Padding(0, 1)

	StatusKeyStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	StatusValueStyle = lipgloss.NewStyle().
				Foreground(Highlight)

	StatusWatchingStyle = lipgloss.NewStyle().
				Foreground(Success).
				Bold(true)

	StatusLoadingStyle = lipgloss.NewStyle().
				Foreground(AccentDim).
				Italic(true)

	StatusErrorStyle = lipgloss.NewStyle().
				Foreground(Error).
				Bold(true)
)

// Filter input styles
var (
	FilterPromptStyle = lipgloss.NewStyle().
				Foreground(Accent).
				Bold(true)

	FilterInputStyle = lipgloss.NewStyle().
				Foreground(Highlight)

	FilterCursorStyle = lipgloss.NewStyle().
				Foreground(Accent)

	FilterTextStyle = lipgloss.NewStyle().
			Foreground(Highlight).
			Background(lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#2A2A2A"})
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

	EmptyStateTitleStyle = lipgloss.NewStyle().
				Foreground(Subtle).
				Bold(true).
				Align(lipgloss.Center).
				MarginBottom(1)

	EmptyStateHintStyle = lipgloss.NewStyle().
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
