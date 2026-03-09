package application

import (
	"context"
	"fmt"
	"time"

	"simplesurance/internal/domain"
)
type TimestampService struct {
	repo      domain.TimestampRepository
	threshold int
}
func NewTimestampService(repo domain.TimestampRepository, threshold int) *TimestampService {
	return &TimestampService{
		repo:      repo,
		threshold: threshold,
	}
}
func (s *TimestampService) Initialize(ctx context.Context) error {
	if err := s.repo.Load(ctx); err != nil {
		return fmt.Errorf("failed to load timestamps: %w", err)
	}

	current := int(time.Now().Unix())
	if err := s.repo.RemoveExpired(ctx, current, s.threshold); err != nil {
		return fmt.Errorf("failed to remove expired timestamps: %w", err)
	}

	if err := s.repo.Sync(ctx); err != nil {
		return fmt.Errorf("failed to sync after initialization: %w", err)
	}

	return nil
}
func (s *TimestampService) RecordTimestamp(ctx context.Context) (int, error) {
	current := int(time.Now().Unix())
	if err := s.repo.RemoveExpired(ctx, current, s.threshold); err != nil {
		return 0, fmt.Errorf("failed to remove expired timestamps: %w", err)
	}
	if err := s.repo.Store(ctx, current); err != nil {
		return 0, fmt.Errorf("failed to store timestamp: %w", err)
	}
	if err := s.repo.Sync(ctx); err != nil {
		return 0, fmt.Errorf("failed to sync timestamp: %w", err)
	}
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get count: %w", err)
	}

	return count, nil
}
