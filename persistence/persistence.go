package persistence

import (
	"context"
	"errors"
	"io"
)

var (
	// ErrFileNotFound is returned when a file doesn't exist
	ErrFileNotFound = errors.New("file not found")
	// ErrFileOperationFailed is returned when a file operation fails
	ErrFileOperationFailed = errors.New("file operation failed")
)

// FilePersistence defines the interface for file operations
type FilePersistence interface {
	// Append appends a timestamp to the file
	Append(ctx context.Context, timestamp int, filename string) error
	// Rewrite rewrites the entire file with the given timestamps
	Rewrite(ctx context.Context, timestamps []int, filename string) error
	// ReadAll reads all timestamps from the file
	ReadAll(ctx context.Context, filename string) ([]int, error)
	// FileExists checks if a file exists
	FileExists(filename string) bool
}

// FileReader defines the interface for reading files
type FileReader interface {
	Read(ctx context.Context, filename string) (io.ReadCloser, error)
}

// FileWriter defines the interface for writing files
type FileWriter interface {
	Write(ctx context.Context, filename string) (io.WriteCloser, error)
	Append(ctx context.Context, filename string) (io.WriteCloser, error)
}

