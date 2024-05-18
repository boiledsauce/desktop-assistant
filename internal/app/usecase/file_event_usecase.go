// usecase/file_event_usecase.go

package usecase

import (
	"desktop-assistant/infra/repository"
	"desktop-assistant/internal/domain"
	"fmt"
	"log"
	"sync"
	"time"
)

type FileEventUseCase struct {
	watcher              *repository.FileWatcher
	files                map[string]*domain.File
	fileDownloadFinished chan *domain.File
	once                 sync.Once
	mu                   sync.Mutex
}

func NewFileEventUseCase(watcher *repository.FileWatcher) *FileEventUseCase {
	return &FileEventUseCase{
		watcher:              watcher,
		files:                make(map[string]*domain.File),
		fileDownloadFinished: make(chan *domain.File),
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

	go func() {
		for file := range feuc.fileDownloadFinished {
			// Here you can take actions on the finalized file
			// For example, you can send an HTTP response or notify other parts of your program
			log.Println("Received finished download:", file.Path)
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
				delete(feuc.files, file.Path)
				feuc.fileDownloadFinished <- file
			}
		}
		feuc.mu.Unlock()
	}
}

func (feuc *FileEventUseCase) create(event domain.FileEvent) *domain.File {
	file := &domain.File{Path: event.Path, LastWrite: event.Timestamp, Hash: "placeholder"}
	feuc.mu.Lock()
	defer feuc.mu.Unlock()

	feuc.files[event.Path] = file
	return file
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
	feuc.files[event.Path] = nil
}

func (feuc *FileEventUseCase) rename(event domain.FileEvent) {
	fmt.Println("Renamed")
}

func (feuc *FileEventUseCase) unknown(event domain.FileEvent) {
	log.Println("Unknown event type")
}

func (feuc *FileEventUseCase) def(event domain.FileEvent) {
	log.Println("Default event type")
}
