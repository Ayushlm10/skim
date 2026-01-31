package filetree

import (
	"path/filepath"
	"strings"
)

// Item represents a file or directory in the tree
type Item struct {
	// Path is the full path to the item
	Path string

	// Name is the display name (base name)
	Name string

	// IsDir indicates if this is a directory
	IsDir bool

	// Depth is the nesting level (0 = root level)
	Depth int

	// Expanded indicates if directory is expanded (only relevant for dirs)
	Expanded bool

	// Children contains child items (only populated when expanded)
	Children []*Item

	// Parent points to the parent directory (nil for root items)
	Parent *Item

	// Visible indicates if item should be shown (for filtering)
	Visible bool
}

// NewItem creates a new tree item
func NewItem(path string, isDir bool, depth int) *Item {
	return &Item{
		Path:     path,
		Name:     filepath.Base(path),
		IsDir:    isDir,
		Depth:    depth,
		Expanded: false,
		Children: nil,
		Parent:   nil,
		Visible:  true,
	}
}

// Toggle expands or collapses a directory
func (i *Item) Toggle() {
	if i.IsDir {
		i.Expanded = !i.Expanded
	}
}

// HasChildren returns true if the directory has children loaded
func (i *Item) HasChildren() bool {
	return len(i.Children) > 0
}

// FilterText returns the text used for filtering
func (i Item) FilterText() string {
	return strings.ToLower(i.Name)
}

// Title implements list.Item interface for bubbles list
func (i Item) Title() string {
	return i.Name
}

// Description implements list.Item interface for bubbles list
func (i Item) Description() string {
	return i.Path
}

// FilterValue implements list.Item interface for bubbles list
func (i Item) FilterValue() string {
	return i.Name
}

// IsMarkdown returns true if the file is a markdown file
func (i *Item) IsMarkdown() bool {
	if i.IsDir {
		return false
	}
	ext := strings.ToLower(filepath.Ext(i.Name))
	return ext == ".md" || ext == ".markdown"
}

// DisplayName returns the name with directory indicator if needed
func (i *Item) DisplayName() string {
	if i.IsDir {
		return i.Name + "/"
	}
	return i.Name
}
