package preview

import (
	"github.com/charmbracelet/glamour"
)

// Renderer wraps Glamour for markdown rendering
type Renderer struct {
	renderer *glamour.TermRenderer
	width    int
}

// NewRenderer creates a new markdown renderer
func NewRenderer(width int) (*Renderer, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
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

// SetWidth updates the word wrap width and recreates the renderer
func (r *Renderer) SetWidth(width int) error {
	if r.width == width {
		return nil
	}

	newRenderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
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
