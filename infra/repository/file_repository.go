package repository

import (
	"os"
)

type FileSystemRepository interface {
	CreateFile(path string, data []byte) error
	ReadFile(path string) ([]byte, error)
	DeleteFile(path string) error
}

type fileSystemRepository struct{}

func NewFileSystemRepository() FileSystemRepository {
	return &fileSystemRepository{}
}

func (fsr *fileSystemRepository) CreateFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func (fsr *fileSystemRepository) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fsr *fileSystemRepository) DeleteFile(path string) error {
	return os.Remove(path)
}
