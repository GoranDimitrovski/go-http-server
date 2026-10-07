package domain

import "context"

// TimestampRepository holds the timestamps inside a sliding window.
// Record and Prune drop timestamps where now-ts >= threshold, persist, and return the remaining count.
type TimestampRepository interface {
	Load(ctx context.Context) error
	Record(ctx context.Context, now, threshold int) (int, error)
	Prune(ctx context.Context, now, threshold int) (int, error)
	Close() error
}
