package store

import (
	"context"
	"fmt"
	"simplesurance/persistence"
	"sync"
)

// MemoryStore implements the Store interface using in-memory storage
type MemoryStore struct {
	timestamps []int
	fileName   string
	persister  persistence.FilePersistence
	mu         sync.RWMutex
}

// NewMemoryStore creates a new memory store instance
func NewMemoryStore(fileName string, persister persistence.FilePersistence) *MemoryStore {
	return &MemoryStore{
		timestamps: make([]int, 0),
		fileName:   fileName,
		persister:  persister,
	}
}

// Store adds a timestamp to the store
func (s *MemoryStore) Store(ctx context.Context, timestamp int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.timestamps = append(s.timestamps, timestamp)
	return nil
}

// View returns a copy of all timestamps
func (s *MemoryStore) View(ctx context.Context) ([]int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]int, len(s.timestamps))
	copy(result, s.timestamps)
	return result, nil
}

// Count returns the number of stored timestamps
func (s *MemoryStore) Count(ctx context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.timestamps), nil
}

// Load loads timestamps from persistent storage
func (s *MemoryStore) Load(ctx context.Context) error {
	timestamps, err := s.persister.ReadAll(ctx, s.fileName)
	if err != nil {
		return fmt.Errorf("failed to load timestamps: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.timestamps = timestamps
	return nil
}

// RemoveExpired removes timestamps older than the threshold
func (s *MemoryStore) RemoveExpired(ctx context.Context, current, threshold int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Pre-allocate with estimated capacity for better performance
	validTimestamps := make([]int, 0, len(s.timestamps))
	for _, timestamp := range s.timestamps {
		if current-timestamp < threshold {
			validTimestamps = append(validTimestamps, timestamp)
		}
	}

	// Release unused capacity by creating a new slice with exact size
	if len(validTimestamps) < cap(validTimestamps) {
		trimmed := make([]int, len(validTimestamps))
		copy(trimmed, validTimestamps)
		s.timestamps = trimmed
	} else {
		s.timestamps = validTimestamps
	}
	return nil
}

// Sync persists the current state to storage
func (s *MemoryStore) Sync(ctx context.Context) error {
	s.mu.RLock()
	timestamps, err := s.View(ctx)
	s.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to get timestamps for sync: %w", err)
	}

	if err := s.persister.Rewrite(ctx, timestamps, s.fileName); err != nil {
		return fmt.Errorf("failed to sync timestamps: %w", err)
	}

	return nil
}

// Close gracefully closes the store and releases resources
func (s *MemoryStore) Close() error {
	// Memory store doesn't need cleanup, but we can sync before closing
	ctx := context.Background()
	return s.Sync(ctx)
}
