package repository

import (
	"log"
	"time"

	"desktop-assistant/internal/domain"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	events    chan domain.FileEvent
	watcher   *fsnotify.Watcher
	directory string
}

// Events returns a receive-only channel of domain events.
func (fw *FileWatcher) Events() <-chan domain.FileEvent {
	return fw.events
}

func (fw *FileWatcher) Errors() <-chan error {
	return fw.watcher.Errors
}

// Directory returns the directory being watched.
func (fw *FileWatcher) Directory() string {
	return fw.directory
}

// NewFileWatcher creates a new file watcher for the specified directory.
func NewFileWatcher(directory string) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create a new Watcher: %v", err)
		return nil, err
	}

	if err := watcher.Add(directory); err != nil {
		watcher.Close()
		log.Printf("Failed to add directory '%s' to watcher: %v", directory, err)
		return nil, err
	}

	return &FileWatcher{watcher: watcher, directory: directory, events: make(chan domain.FileEvent)}, nil
}

// Start begins watching file events in the directory.
func (fw *FileWatcher) Start() {
	go func() {
		for {
			select {
			case event, ok := <-fw.watcher.Events:
				if !ok {
					return // Channel closed, exit the loop
				}
				fw.handleEvent(event)
			case err, ok := <-fw.watcher.Errors:
				if !ok {
					return // Channel closed, exit the loop
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()
}

// handleEvent processes file system events by logging them and performing actions.
func (fw *FileWatcher) handleEvent(event fsnotify.Event) {
	// Convert the fsnotify.Event to a domain.Event and send it to the Events channel
	fw.events <- domain.FileEvent{
		Path:      event.Name,
		Timestamp: time.Now(),
		Type:      mapFsnotifyOpToEventType(event.Op),
	}
}

// mapFsnotifyOpToEventType maps fsnotify.Op values to domain.EventType values.
func mapFsnotifyOpToEventType(op fsnotify.Op) domain.EventType {
	switch op {
	case fsnotify.Create:
		return domain.Create
	case fsnotify.Write:
		return domain.Write
	case fsnotify.Remove:
		return domain.Remove
	case fsnotify.Rename:
		return domain.Rename
	case fsnotify.Chmod:
		return domain.Close
	default:
		return domain.Unknown
	}
}

// Stop stops the file watcher and closes all resources.
func (fw *FileWatcher) Stop() {
	fw.watcher.Close()
}
