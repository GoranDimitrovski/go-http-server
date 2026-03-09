package persistence

import (
	"context"
	"errors"
	"io"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrFileOperationFailed = errors.New("file operation failed")
)
type FilePersistence interface {
	Append(ctx context.Context, timestamp int, filename string) error
	Rewrite(ctx context.Context, timestamps []int, filename string) error
	ReadAll(ctx context.Context, filename string) ([]int, error)
	FileExists(filename string) bool
}
type FileReader interface {
	Read(ctx context.Context, filename string) (io.ReadCloser, error)
}
type FileWriter interface {
	Write(ctx context.Context, filename string) (io.WriteCloser, error)
	Append(ctx context.Context, filename string) (io.WriteCloser, error)
}
