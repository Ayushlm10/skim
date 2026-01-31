package app

// Panel represents which panel has focus
type Panel int

const (
	FileTreePanel Panel = iota
	PreviewPanel
)

// Custom message types for the application

// FileSelectedMsg is sent when a file is selected in the tree
type FileSelectedMsg struct {
	Path string
}

// DirectoryToggleMsg is sent when a directory is expanded/collapsed
type DirectoryToggleMsg struct {
	Path string
}

// FileLoadedMsg is sent when a file's content has been loaded
type FileLoadedMsg struct {
	Path    string
	Content string
}

// FileErrorMsg is sent when there's an error loading a file
type FileErrorMsg struct {
	Path string
	Err  error
}

// FileChangedMsg is sent when a watched file changes
type FileChangedMsg struct {
	Path string
}

// WatchStartMsg is sent when file watching starts
type WatchStartMsg struct {
	Path string
}

// WatchStopMsg is sent when file watching stops
type WatchStopMsg struct{}

// FocusChangedMsg is sent when focus changes between panels
type FocusChangedMsg struct {
	Panel Panel
}

// FilterActiveMsg is sent when filter mode is toggled
type FilterActiveMsg struct {
	Active bool
	Value  string
}

// WindowSizeMsg wraps the terminal size for components
type WindowSizeMsg struct {
	Width  int
	Height int
}
