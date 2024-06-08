package handler

import (
	"desktop-assistant/internal/app/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FileSystemHandler struct {
	filesystemUc *usecase.FileSystemUseCase
}

func NewFileSystemHandler(filesystemUc *usecase.FileSystemUseCase) *FileSystemHandler {
	return &FileSystemHandler{
		filesystemUc: filesystemUc,
	}
}

func (fh *FileSystemHandler) CreateFile(c echo.Context) error {
	if file, err := fh.filesystemUc.CreateFile(c.QueryParam("path"), []byte(c.QueryParam("data"))); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	} else {
		return c.JSON(http.StatusCreated, file)
	}
}

func (fh *FileSystemHandler) ReadFile(c echo.Context) ([]byte, error) {
	path := c.QueryParam("path")
	if _, err := fh.filesystemUc.ReadFile(path); err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "File not found")
	} else {
		return nil, nil
	}
}

func (fh *FileSystemHandler) DeleteFile(c echo.Context) error {
	path := c.QueryParam("path")
	err := fh.filesystemUc.DeleteFile(path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "File not found")
	}
	return c.String(http.StatusOK, "File deleted successfully")
}

// Define the Request type
type Request struct {
	Path string `json:"path"`
	// Define the fields of the Request type
}

// Deletes duplicates of the file
func (fh *FileSystemHandler) DeleteDuplicateFiles(c echo.Context) error {
	var body Request

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	err := fh.filesystemUc.DeleteDuplicateFiles(body.Path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "File not found")
	}
	return c.String(http.StatusOK, "File deleted successfully")
}
