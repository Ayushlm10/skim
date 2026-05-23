package preview

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// Renderer wraps Glamour for markdown rendering
type Renderer struct {
	renderer *glamour.TermRenderer
	width    int
}

// NewRenderer creates a new markdown renderer
func NewRenderer(width int) (*Renderer, error) {
	r, err := glamour.NewTermRenderer(rendererOptions(width)...)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		renderer: r,
		width:    width,
	}, nil
}

// Render renders markdown content to styled terminal output
func (r *Renderer) Render(content string) (string, error) {
	return r.renderer.Render(content)
}

func rendererOptions(width int) []glamour.TermRendererOption {
	style := "light"
	if lipgloss.HasDarkBackground() {
		style = "dark"
	}

	return []glamour.TermRendererOption{
		glamour.WithStandardStyle(style),
		glamour.WithWordWrap(width),
	}
}

// SetWidth updates the word wrap width and recreates the renderer
func (r *Renderer) SetWidth(width int) error {
	if r.width == width {
		return nil
	}

	return r.reset(width)
}

func (r *Renderer) RefreshTheme() error {
	return r.reset(r.width)
}

func (r *Renderer) reset(width int) error {
	newRenderer, err := glamour.NewTermRenderer(rendererOptions(width)...)
	if err != nil {
		return err
	}

	r.renderer = newRenderer
	r.width = width
	return nil
}

// Width returns the current word wrap width
func (r *Renderer) Width() int {
	return r.width
}
