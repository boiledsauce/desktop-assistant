// app/usecase/file_usecase.go
package usecase

import (
	"desktop-assistant/infra/repository"
	"desktop-assistant/internal/domain"
	"errors"
	"log"
	"os"
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
	mountedDir := os.Getenv("MOUNT_POINT")
	fileName := path

	log.Println("Name of file: ", fileName)
	// Read the file to get its content
	_, err := uc.ReadFile(mountedDir + fileName)
	if err != nil {
		log.Print("Error reading file: ", err)
		return err
	}

	// Get the base path where file is located

	// List all files in the directory
	files, err := uc.repo.ListFiles(mountedDir)
	if err != nil {
		log.Print("Error listing files: ", err)
		return err
	}

	// Find duplicates of the file
	var duplicates []*domain.File
	for _, file := range files {
		// Skip the original file
		if file.Path == mountedDir+fileName {
			continue
		}
		// Compare the name of the path of the files
		// For eaxmple movie(1).mov and movie.mov are duplicates
		// But movie.mov and movie2.mov are not duplicates

		isSameFileNameWithoutSuffix := func(s string, s2 string) bool {
			// Regex pattern to find base name and extension, ignoring "(number)"
			pattern := regexp.MustCompile(`^(.*?)(?: \(\d+\))?(\.\w+)$`)

			// Extract base names and extensions
			match1 := pattern.FindStringSubmatch(s)
			// log.Println("Match1:", match1)

			match2 := pattern.FindStringSubmatch(s2)
			// log.Println("Match2:", match2)

			if len(match1) < 3 || len(match2) < 3 {
				return false // Ensure there are enough groups matched to compare
			}

			// Compare base names and extensions
			return match1[1]+match1[2] == match2[1]+match2[2]
		}

		if isSameFileNameWithoutSuffix(file.Path, filepath.Join(mountedDir, fileName)) {
			duplicates = append(duplicates, file)
		}
	}

	// Log the deletion of the files and the filenames
	for _, file := range duplicates {
		log.Println("Deleting: ", file.Path)
	}

	if err := uc.repo.DeleteFiles(duplicates); err != nil {
		log.Print("Error deleting files: ", err)
		return err
	}

	return nil
}
