package watcher

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher wraps fsnotify with debouncing and Bubble Tea integration
type Watcher struct {
	watcher *fsnotify.Watcher

	// Currently watched file
	watchedPath string

	// Debounce settings
	debounceDelay time.Duration

	// Channel for file change events (debounced)
	Events chan string

	// Channel for errors
	Errors chan error

	// Mutex for thread-safe operations
	mu sync.Mutex

	// Track if we're running
	running bool

	// Done channel for shutdown
	done chan struct{}
}

// New creates a new file watcher
func New() (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		watcher:       fsWatcher,
		debounceDelay: 100 * time.Millisecond, // 100ms debounce
		Events:        make(chan string, 10),
		Errors:        make(chan error, 10),
		done:          make(chan struct{}),
	}

	// Start the event processing goroutine
	go w.processEvents()

	w.running = true
	return w, nil
}

// Watch starts watching a file, removing any previously watched file
func (w *Watcher) Watch(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Remove old watch if exists
	if w.watchedPath != "" {
		_ = w.watcher.Remove(w.watchedPath)
	}

	// Add new watch
	err := w.watcher.Add(path)
	if err != nil {
		return err
	}

	w.watchedPath = path
	return nil
}

// Unwatch stops watching the current file
func (w *Watcher) Unwatch() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.watchedPath != "" {
		_ = w.watcher.Remove(w.watchedPath)
		w.watchedPath = ""
	}
}

// WatchedPath returns the currently watched path
func (w *Watcher) WatchedPath() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.watchedPath
}

// Close stops the watcher and cleans up resources
func (w *Watcher) Close() error {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return nil
	}
	w.running = false
	w.mu.Unlock()

	close(w.done)
	return w.watcher.Close()
}

// processEvents handles fsnotify events with debouncing
func (w *Watcher) processEvents() {
	var (
		// Debounce timer
		timer *time.Timer
		// Pending event path
		pendingPath string
	)

	for {
		select {
		case <-w.done:
			if timer != nil {
				timer.Stop()
			}
			return

		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// Only handle Write and Create events
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				w.mu.Lock()
				currentPath := w.watchedPath
				w.mu.Unlock()

				// Only process events for the file we're watching
				if event.Name == currentPath {
					// Debounce: reset timer on each event
					if timer != nil {
						timer.Stop()
					}
					pendingPath = event.Name
					timer = time.AfterFunc(w.debounceDelay, func() {
						// Non-blocking send
						select {
						case w.Events <- pendingPath:
						default:
						}
					})
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			// Non-blocking send
			select {
			case w.Errors <- err:
			default:
			}
		}
	}
}
