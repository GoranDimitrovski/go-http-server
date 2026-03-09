package repository

import (
	"context"
	"fmt"
	"sync"

	"simplesurance/internal/infrastructure/persistence"
)
type MemoryStore struct {
	timestamps []int
	fileName   string
	persister  persistence.FilePersistence
	mu         sync.RWMutex
}
func NewMemoryStore(fileName string, persister persistence.FilePersistence) *MemoryStore {
	return &MemoryStore{
		timestamps: make([]int, 0),
		fileName:   fileName,
		persister:  persister,
	}
}
func (s *MemoryStore) Store(ctx context.Context, timestamp int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.timestamps = append(s.timestamps, timestamp)
	return nil
}
func (s *MemoryStore) View(ctx context.Context) ([]int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]int, len(s.timestamps))
	copy(result, s.timestamps)
	return result, nil
}
func (s *MemoryStore) Count(ctx context.Context) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.timestamps), nil
}
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
func (s *MemoryStore) RemoveExpired(ctx context.Context, current, threshold int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	validTimestamps := make([]int, 0, len(s.timestamps))
	for _, timestamp := range s.timestamps {
		if current-timestamp < threshold {
			validTimestamps = append(validTimestamps, timestamp)
		}
	}
	if len(validTimestamps) < cap(validTimestamps) {
		trimmed := make([]int, len(validTimestamps))
		copy(trimmed, validTimestamps)
		s.timestamps = trimmed
	} else {
		s.timestamps = validTimestamps
	}
	return nil
}
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
func (s *MemoryStore) Close() error {
	ctx := context.Background()
	return s.Sync(ctx)
}
