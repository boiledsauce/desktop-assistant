package handler

import (
	"desktop-assistant/internal/app/usecase"
	"desktop-assistant/internal/domain"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	fileEventUseCase *usecase.FileEventUseCase
	files            map[string]*domain.File
}

func NewFileHandler(fileEventUseCase *usecase.FileEventUseCase) *FileHandler {
	return &FileHandler{
		fileEventUseCase: fileEventUseCase,
		files:            make(map[string]*domain.File),
	}
}

func (fh *FileHandler) HandleFileEvents(c echo.Context) error {
	// Set the necessary headers for SSE
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")
	c.Response().WriteHeader(200)
	fh.fileEventUseCase.HandleFileEvents()

	// Create a channel to signal when the client disconnects
	disconnected := make(chan struct{})
	defer close(disconnected)

	// Close the connection when the client disconnects
	go func() {
		select {
		case <-c.Request().Context().Done():
			disconnected <- struct{}{}
		}
	}()

	// Handle the file events
	for {
		select {
		case event := <-fh.fileEventUseCase.FileDownloadFinishedChannel():
			// Format the event data
			data := fmt.Sprintf("data: %v\n\n", event)
			log.Printf("Sending event: %s", data) // Log the event being sent

			// Send the event to the client
			if _, err := c.Response().Write([]byte(data)); err != nil {
				log.Printf("Error writing response: %v", err) // Log the error
				return err
			}
			c.Response().Flush()
		case <-disconnected:
			// The client disconnected, so stop handling events
			log.Println("Client disconnected") // Log the disconnection
			return nil
		}
	}
}
