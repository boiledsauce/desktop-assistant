// app/usecase/file_usecase.go
package usecase

import (
	"desktop-assistant/infra/repository"
	"errors"
	"path/filepath"
)

type FileSystemUseCase struct {
	repo repository.FileSystemRepository
}

func NewFileSystemUseCase(repo repository.FileSystemRepository) *FileSystemUseCase {
	return &FileSystemUseCase{
		repo: repo,
	}
}

// CreateFile ensures business rules and validations before creating a file.
func (uc *FileSystemUseCase) CreateFile(path string, data []byte) error {
	// Validate the path or data based on business rules
	if path == "" {
		return errors.New("path cannot be empty")
	}
	if data == nil {
		return errors.New("data cannot be nil")
	}

	// Normalize the path to avoid issues like directory traversal attacks
	cleanPath := filepath.Clean(path)

	// Check for business-specific conditions, e.g., file size, file type, naming conventions
	// Example: Let's say business rules dictate that files cannot be created in a certain directory
	if filepath.Dir(cleanPath) == "restricted_directory" {
		return errors.New("cannot create files in restricted directory")
	}

	// Call the repository function to create the file
	return uc.repo.CreateFile(cleanPath, data)
}

func (uc *FileSystemUseCase) ReadFile(path string) ([]byte, error) {
	return uc.repo.ReadFile(path)
}

func (uc *FileSystemUseCase) DeleteFile(path string) error {
	return uc.repo.DeleteFile(path)
}
