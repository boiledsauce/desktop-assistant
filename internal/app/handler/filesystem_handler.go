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
	return fh.filesystemUc.CreateFile(c.QueryParam("path"), []byte(c.QueryParam("data")))
}

func (fh *FileSystemHandler) ReadFile(c echo.Context) ([]byte, error) {
	path := c.QueryParam("path")
	data, err := fh.filesystemUc.ReadFile(path)
	return data, err
}

func (fh *FileSystemHandler) DeleteFile(c echo.Context) error {
	path := c.QueryParam("path")
	err := fh.filesystemUc.DeleteFile(path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "File not found")
	}
	return c.String(http.StatusOK, "File deleted successfully")
}
