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
	return &TimestampService{repo: repo, threshold: threshold}
}

func (s *TimestampService) Initialize(ctx context.Context) error {
	if err := s.repo.Load(ctx); err != nil {
		return fmt.Errorf("failed to load timestamps: %w", err)
	}
	if _, err := s.repo.Prune(ctx, int(time.Now().Unix()), s.threshold); err != nil {
		return fmt.Errorf("failed to remove expired timestamps: %w", err)
	}
	return nil
}

// RecordTimestamp records now and returns the number of timestamps inside the window.
func (s *TimestampService) RecordTimestamp(ctx context.Context) (int, error) {
	count, err := s.repo.Record(ctx, int(time.Now().Unix()), s.threshold)
	if err != nil {
		return 0, fmt.Errorf("failed to record timestamp: %w", err)
	}
	return count, nil
}
