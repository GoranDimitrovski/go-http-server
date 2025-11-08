package persistence

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
)

// FilePersistenceImpl implements FilePersistence interface
type FilePersistenceImpl struct{}

// NewFilePersistence creates a new file persistence instance
func NewFilePersistence() *FilePersistenceImpl {
	return &FilePersistenceImpl{}
}

// Append appends a timestamp to the file
func (f *FilePersistenceImpl) Append(ctx context.Context, timestamp int, filename string) error {
	return f.WriteToFile(ctx, []int{timestamp}, filename, true)
}

// Rewrite rewrites the entire file with the given timestamps
func (f *FilePersistenceImpl) Rewrite(ctx context.Context, timestamps []int, filename string) error {
	if err := f.truncateFile(ctx, filename); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	return f.WriteToFile(ctx, timestamps, filename, false)
}

// ReadAll reads all timestamps from the file
func (f *FilePersistenceImpl) ReadAll(ctx context.Context, filename string) ([]int, error) {
	if !f.FileExists(filename) {
		return []int{}, nil
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %w", err)
	}
	defer file.Close()

	var timestamps []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		timestamp, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp '%s': %w", line, err)
		}
		timestamps = append(timestamps, timestamp)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while scanning file: %w", err)
	}

	return timestamps, nil
}

// FileExists checks if a file exists
func (f *FilePersistenceImpl) FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// WriteToFile writes timestamps to a file
func (f *FilePersistenceImpl) WriteToFile(ctx context.Context, timestamps []int, filename string, append bool) error {
	flags := os.O_CREATE | os.O_WRONLY
	if append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	file, err := os.OpenFile(filename, flags, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, timestamp := range timestamps {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if _, err := writer.WriteString(fmt.Sprintf("%d\n", timestamp)); err != nil {
			return fmt.Errorf("failed to write timestamp: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func (f *FilePersistenceImpl) truncateFile(ctx context.Context, filename string) error {
	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, nothing to truncate
		}
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	return file.Close()
}

