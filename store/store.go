package store

import (
	"context"
	"errors"
)

var (
	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("resource not found")
	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")
)

// Store defines the interface for timestamp storage operations
type Store interface {
	// Store adds a timestamp to the store
	Store(ctx context.Context, timestamp int) error
	// View returns a copy of all timestamps
	View(ctx context.Context) ([]int, error)
	// Count returns the number of stored timestamps
	Count(ctx context.Context) (int, error)
	// Load loads timestamps from persistent storage
	Load(ctx context.Context) error
	// RemoveExpired removes timestamps older than the threshold
	RemoveExpired(ctx context.Context, current, threshold int) error
	// Sync persists the current state to storage
	Sync(ctx context.Context) error
	// Close gracefully closes the store and releases resources
	Close() error
}

