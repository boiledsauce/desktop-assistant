package repository

import (
	"os"
	"time"

	"desktop-assistant/internal/domain"
)

type FileWatcherPolling struct {
	events    chan domain.FileEvent
	errors    chan error
	directory string
	ticker    *time.Ticker
	files     map[string]os.FileInfo
}

// NewFileWatcherPolling creates a new file watcher that uses polling to detect file changes.
func NewFileWatcherPolling(directory string) *FileWatcherPolling {
	return &FileWatcherPolling{
		directory: directory,
		events:    make(chan domain.FileEvent, 10), // Buffered channel to handle bursts of changes
		errors:    make(chan error, 10),
		files:     make(map[string]os.FileInfo),
	}
}

func (fw *FileWatcherPolling) Start(interval time.Duration) {
	fw.ticker = time.NewTicker(interval)
	go func() {
		for range fw.ticker.C {
			fw.pollDirectory()
		}
	}()
}

func (fw *FileWatcherPolling) pollDirectory() {
	currentFiles, err := os.ReadDir(fw.directory)
	if err != nil {
		fw.errors <- err
		return
	}

	currentMap := make(map[string]os.FileInfo)
	for _, entry := range currentFiles {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		currentMap[entry.Name()] = info
		oldInfo, found := fw.files[entry.Name()]
		if !found {
			// New file
			fw.events <- domain.FileEvent{Path: entry.Name(), Timestamp: info.ModTime(), Type: domain.Create}
		} else if info.ModTime() != oldInfo.ModTime() {
			// Modified file
			fw.events <- domain.FileEvent{Path: entry.Name(), Timestamp: info.ModTime(), Type: domain.Write}
		}
		delete(fw.files, entry.Name())
	}

	// Check for deleted files
	for filename, oldInfo := range fw.files {
		fw.events <- domain.FileEvent{Path: filename, Timestamp: oldInfo.ModTime(), Type: domain.Remove}
	}

	// Update the stored files info
	fw.files = currentMap
}

func (fw *FileWatcherPolling) Stop() {
	if fw.ticker != nil {
		fw.ticker.Stop()
	}
	close(fw.events)
	close(fw.errors)
}

func (fw *FileWatcherPolling) Events() <-chan domain.FileEvent {
	return fw.events
}

func (fw *FileWatcherPolling) Errors() <-chan error {
	return fw.errors
}

func (fw *FileWatcherPolling) Directory() string {
	return fw.directory
}
