// internal/domain/models.go
package domain

import (
	"fmt"
	"time"
)

const TIME_UNTIL_DOWNLOADED = 5 * time.Second

type File struct {
	Path      string
	Hash      string
	LastWrite time.Time
}

type FileRepository interface {
	Remove(file File) error
	Save(file File) error
	FindByID(id int) (File, error)
}

func (f *File) IsDownloadFinished(currentTime time.Time) bool {
	// If it's been more than 5 seconds since the last write, assume the download has finished
	// print how long it is left
	fmt.Println("TIME UNTIL IT IS DONE: ", TIME_UNTIL_DOWNLOADED-time.Since(f.LastWrite))
	return time.Since(f.LastWrite) > TIME_UNTIL_DOWNLOADED
}