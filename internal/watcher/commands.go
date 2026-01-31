package watcher

import (
	tea "github.com/charmbracelet/bubbletea"
)

// FileChangedMsg is sent when a watched file changes
type FileChangedMsg struct {
	Path string
}

// WatchErrorMsg is sent when there's a watcher error
type WatchErrorMsg struct {
	Err error
}

// WaitForChange creates a command that waits for a file change event
// This should be called after each FileChangedMsg to continue listening
func WaitForChange(w *Watcher) tea.Cmd {
	return func() tea.Msg {
		select {
		case path, ok := <-w.Events:
			if !ok {
				return nil
			}
			return FileChangedMsg{Path: path}
		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			return WatchErrorMsg{Err: err}
		}
	}
}

// StartWatching creates a command to start watching a file and wait for changes
func StartWatching(w *Watcher, path string) tea.Cmd {
	return func() tea.Msg {
		err := w.Watch(path)
		if err != nil {
			return WatchErrorMsg{Err: err}
		}
		// Return immediately to trigger the first WaitForChange
		return watchStartedMsg{path: path}
	}
}

// watchStartedMsg is an internal message indicating watch has started
type watchStartedMsg struct {
	path string
}

// WatchStartedMsg returns true if the msg is a watch started message
func IsWatchStartedMsg(msg tea.Msg) (string, bool) {
	if m, ok := msg.(watchStartedMsg); ok {
		return m.path, true
	}
	return "", false
}
