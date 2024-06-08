// app/usecase/file_usecase.go
package usecase

import (
	"desktop-assistant/infra/repository"
	"desktop-assistant/internal/domain"
	"errors"
	"path/filepath"
	"regexp"
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
func (uc *FileSystemUseCase) CreateFile(path string, data []byte) (*domain.File, error) {
	// Validate the path or data based on business rules
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	// Normalize the path to avoid issues like directory traversal attacks
	cleanPath := filepath.Clean(path)

	// Check for business-specific conditions, e.g., file size, file type, naming conventions
	// Example: Let's say business rules dictate that files cannot be created in a certain directory
	if filepath.Dir(cleanPath) == "restricted_directory" {
		return nil, errors.New("cannot create files in restricted directory")
	}

	// Call the repository function to create the file
	return uc.repo.CreateFile(cleanPath, data)
}

func (uc *FileSystemUseCase) ReadFile(path string) (*domain.File, error) {
	return uc.repo.ReadFile(path)
}

func (uc *FileSystemUseCase) DeleteFile(path string) error {
	file := &domain.File{
		Path: path,
	}
	return uc.repo.DeleteFile(file)
}

func (uc *FileSystemUseCase) DeleteDuplicateFiles(path string) error {
	// Read the file to get its content
	_, err := uc.ReadFile(path)
	if err != nil {
		return err
	}

	// Get the base path where file is located
	basePath := filepath.Dir(path)

	// List all files in the directory
	files, err := uc.repo.ListFiles(basePath)
	if err != nil {
		return err
	}

	// Find duplicates of the file
	var duplicates []*domain.File
	for _, file := range files {
		// Skip the original file
		if file.Path == path {
			continue
		}
		// Compare the name of the path of the files
		// For eaxmple movie(1).mov and movie.mov are duplicates
		// But movie.mov and movie2.mov are not duplicates
		if regexp.MustCompile(`\(\d+\)`).MatchString(file.Path) {
			// Check if the file name without the (1) or (2) is the same
			originalFileName := regexp.MustCompile(`\(\d+\)`).ReplaceAllString(path, "")
			if originalFileName == file.Path {
				duplicates = append(duplicates, file)
			}
		}
	}

	// Delete the duplicates
	return uc.repo.DeleteFiles(duplicates)

}
