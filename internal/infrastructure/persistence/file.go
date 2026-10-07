package persistence

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

// ReadAll returns the timestamps stored in filename, or none if it doesn't exist.
func ReadAll(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		return []int{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading: %w", err)
	}
	defer file.Close()

	timestamps := []int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
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

// WriteAll replaces filename's contents atomically: write a temp file, then rename over the original,
// so a crash mid-write never leaves a truncated file behind.
func WriteAll(filename string, timestamps []int) error {
	tmp := filename + ".tmp"
	file, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	writer := bufio.NewWriter(file)
	for _, timestamp := range timestamps {
		writer.WriteString(strconv.Itoa(timestamp))
		writer.WriteByte('\n')
	}
	// bufio.Writer keeps the first write error and returns it from Flush.
	if err := errors.Join(writer.Flush(), file.Close()); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	if err := os.Rename(tmp, filename); err != nil {
		return fmt.Errorf("failed to replace file: %w", err)
	}
	return nil
}
