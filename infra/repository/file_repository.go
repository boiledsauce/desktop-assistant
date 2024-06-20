package repository

import (
	"desktop-assistant/internal/domain"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileSystemRepository interface {
	CreateFile(path string, data []byte) (*domain.File, error)
	ReadFile(path string) (*domain.File, error)
	GetFileEntity(file multipart.File) (*domain.File, error)

	ListFiles(directory string) ([]*domain.File, error)
	DeleteFile(file *domain.File) error
	DeleteFiles(files []*domain.File) error
}

type FileSystemRepositoryImpl struct{}

func NewFileSystemRepository() FileSystemRepository {
	return &FileSystemRepositoryImpl{}
}

func (fsr *FileSystemRepositoryImpl) CreateFile(path string, data []byte) (*domain.File, error) {
	file := &domain.File{
		Path: path,
		Hash: "",
	}

	return file, nil
}

func (fsr *FileSystemRepositoryImpl) GetFileEntity(file multipart.File) (*domain.File, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Print("Error reading file: ", err)
		return nil, err
	}

	return &domain.File{
		Path:    "",
		Hash:    "",
		Content: fileBytes,
		Size:    int64(len(fileBytes)),
	}, nil
}

func (fsr *FileSystemRepositoryImpl) ReadFile(path string) (*domain.File, error) {
	// Open the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err // Handle errors such as file not found or read permission issues
	}

	// Retrieve file info to get additional details like size
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err // Handle potential errors in retrieving file info
	}

	// Create and return the File entity
	file := &domain.File{
		Path:    path,
		Hash:    "",
		Content: data,
		Size:    fileInfo.Size(),
	}
	log.Printf("File read successfully. Path: %s, Size: %d bytes\n", path, fileInfo.Size())
	return file, nil
}

func (fsr *FileSystemRepositoryImpl) DeleteFile(file *domain.File) error {
	return os.Remove(file.Path)
}

func (fsr *FileSystemRepositoryImpl) DeleteFiles(files []*domain.File) error {
	for _, file := range files {
		err := fsr.DeleteFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fsr *FileSystemRepositoryImpl) ListFiles(directory string) ([]*domain.File, error) {
	var files []*domain.File
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, &domain.File{
				Path:      path,
				Hash:      info.Name(),
				LastWrite: info.ModTime(),
			})
		}
		return nil
	})
	return files, err
}
