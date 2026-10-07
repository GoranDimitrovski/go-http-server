package repository

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"simplesurance/internal/infrastructure/persistence"
)

// MemoryStore keeps timestamps in memory and rewrites the file on every change.
// Rewriting the whole file per request is O(n) disk I/O; switch to a periodic flush if throughput matters.
type MemoryStore struct {
	mu         sync.Mutex // also serializes file writes
	timestamps []int
	fileName   string
}

func NewMemoryStore(fileName string) *MemoryStore {
	return &MemoryStore{timestamps: []int{}, fileName: fileName}
}

func (s *MemoryStore) Load(ctx context.Context) error {
	timestamps, err := persistence.ReadAll(s.fileName)
	if err != nil {
		return fmt.Errorf("failed to load timestamps: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.timestamps = timestamps
	return nil
}

func (s *MemoryStore) Record(ctx context.Context, now, threshold int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.timestamps = append(s.timestamps, now)
	return s.pruneAndSave(now, threshold)
}

func (s *MemoryStore) Prune(ctx context.Context, now, threshold int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pruneAndSave(now, threshold)
}

// pruneAndSave must be called with s.mu held.
func (s *MemoryStore) pruneAndSave(now, threshold int) (int, error) {
	s.timestamps = slices.DeleteFunc(s.timestamps, func(ts int) bool { return now-ts >= threshold })
	if err := persistence.WriteAll(s.fileName, s.timestamps); err != nil {
		return 0, fmt.Errorf("failed to persist timestamps: %w", err)
	}
	return len(s.timestamps), nil
}

func (s *MemoryStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return persistence.WriteAll(s.fileName, s.timestamps)
}
