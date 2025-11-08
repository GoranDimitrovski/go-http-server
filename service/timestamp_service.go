package service

import (
	"context"
	"fmt"
	"simplesurance/store"
	"time"
)

// TimestampService handles timestamp-related business logic
type TimestampService struct {
	store     store.Store
	threshold int
}

// NewTimestampService creates a new timestamp service
func NewTimestampService(store store.Store, threshold int) *TimestampService {
	return &TimestampService{
		store:     store,
		threshold: threshold,
	}
}

// Initialize loads existing timestamps and removes expired ones
func (s *TimestampService) Initialize(ctx context.Context) error {
	if err := s.store.Load(ctx); err != nil {
		return fmt.Errorf("failed to load timestamps: %w", err)
	}

	current := int(time.Now().Unix())
	if err := s.store.RemoveExpired(ctx, current, s.threshold); err != nil {
		return fmt.Errorf("failed to remove expired timestamps: %w", err)
	}

	if err := s.store.Sync(ctx); err != nil {
		return fmt.Errorf("failed to sync after initialization: %w", err)
	}

	return nil
}

// RecordTimestamp records a new timestamp and returns the current count
func (s *TimestampService) RecordTimestamp(ctx context.Context) (int, error) {
	current := int(time.Now().Unix())

	// Remove expired timestamps first
	if err := s.store.RemoveExpired(ctx, current, s.threshold); err != nil {
		return 0, fmt.Errorf("failed to remove expired timestamps: %w", err)
	}

	// Store the new timestamp
	if err := s.store.Store(ctx, current); err != nil {
		return 0, fmt.Errorf("failed to store timestamp: %w", err)
	}

	// Sync to persistent storage
	if err := s.store.Sync(ctx); err != nil {
		return 0, fmt.Errorf("failed to sync timestamp: %w", err)
	}

	// Get the count
	count, err := s.store.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get count: %w", err)
	}

	return count, nil
}

