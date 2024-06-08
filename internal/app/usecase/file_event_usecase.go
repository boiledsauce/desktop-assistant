package usecase

import (
	"crypto/sha256"
	"desktop-assistant/infra/repository"
	"desktop-assistant/internal/domain"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"
)

type FileEventUseCase struct {
	watcher              *repository.FileWatcherPolling
	files                map[string]*domain.File
	fileDownloadFinished chan *domain.File
	once                 sync.Once
	mu                   sync.Mutex
}

func NewFileEventUseCase(watcher *repository.FileWatcherPolling) *FileEventUseCase {
	return &FileEventUseCase{
		watcher:              watcher,
		files:                make(map[string]*domain.File),
		fileDownloadFinished: make(chan *domain.File, 10), // Buffered channel to avoid blocking
	}
}

func (feuc *FileEventUseCase) FileDownloadFinishedChannel() <-chan *domain.File {
	return feuc.fileDownloadFinished
}

func (feuc *FileEventUseCase) HandleFileEvents() {
	go func() {
		for {
			select {
			case event, ok := <-feuc.watcher.Events():
				if !ok {
					return // Channel closed, exit the loop
				}
				fmt.Println("Event: ", event)
				feuc.handleEvent(event)
			case err, ok := <-feuc.watcher.Errors():
				if !ok {
					return // Channel closed, exit the loop
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()
}

func (feuc *FileEventUseCase) handleEvent(event domain.FileEvent) {
	log.Printf("Detected file change: %s, Event: %v", event.Path, event.Type)
	feuc.once.Do(func() {
		go feuc.checkDownloads()
	})
	// You can expand this function to perform different actions based on the event type
	switch event.Type {
	case domain.Create:
		log.Printf("Create event: %s", event.Path)
		feuc.create(event)
	case domain.Write:
		log.Printf("Write event: %s", event.Path)
		feuc.write(event)
	case domain.Remove:
		feuc.remove(event)
	case domain.Rename:
		feuc.rename(event)
	case domain.Unknown:
		feuc.unknown(event)
	default:
		log.Println("Default event type")
	}
}

func (feuc *FileEventUseCase) checkDownloads() {
	// Make goroutine only once
	for {
		time.Sleep(1 * time.Second) // Check every second
		log.Println("Checking for files...")
		feuc.mu.Lock()
		for _, file := range feuc.files {
			if file.IsDownloadFinished(time.Now()) {
				log.Println("Download finished:", file.Path)
				if isDuplicateFilename(file.Path) {
					log.Println("Duplicate found:", file.Path)
					select {
					case feuc.fileDownloadFinished <- file:
						log.Printf("Sent file download finished event: %s", file.Path)
					default:
						log.Printf("Channel is full, dropping event: %s", file.Path)
					}
				}
				delete(feuc.files, file.Path)
			}
		}
		feuc.mu.Unlock()
	}
}

func (feuc *FileEventUseCase) create(event domain.FileEvent) *domain.File {
	file := &domain.File{Path: event.Path, LastWrite: event.Timestamp, Hash: feuc.hashFile(event.Path)}
	feuc.mu.Lock()
	defer feuc.mu.Unlock()

	for _, existingFile := range feuc.files {
		if isDuplicateFilename(existingFile.Path) {
			log.Printf("Duplicate file created: %s", file.Path)
			// Send duplicate file event
		}
	}

	feuc.files[event.Path] = file
	return file
}

func (feuc *FileEventUseCase) hashFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read file for hashing: %v", err)
		return ""
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

func (feuc *FileEventUseCase) write(event domain.FileEvent) error {
	feuc.mu.Lock()
	defer feuc.mu.Unlock()

	file, ok := feuc.files[event.Path]
	if !ok {
		return fmt.Errorf("file not found: %s", event.Path)
	}
	file.LastWrite = event.Timestamp

	return nil
}

func (feuc *FileEventUseCase) remove(event domain.FileEvent) {
	feuc.mu.Lock()
	defer feuc.mu.Unlock()

	delete(feuc.files, event.Path)
}

func (feuc *FileEventUseCase) rename(event domain.FileEvent) {
	log.Println("Renamed event:", event.Path)
}

func (feuc *FileEventUseCase) unknown(event domain.FileEvent) {
	log.Println("Unknown event type")
}

func isDuplicateFilename(filename string) bool {
	// Regex pattern to find "(number)" before the file extension
	pattern := regexp.MustCompile(`\(\d+\)\.\w+$`)
	return pattern.MatchString(filename)
}
