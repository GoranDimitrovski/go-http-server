package domain

import "context"
type TimestampRepository interface {
	Store(ctx context.Context, timestamp int) error
	View(ctx context.Context) ([]int, error)
	Count(ctx context.Context) (int, error)
	Load(ctx context.Context) error
	RemoveExpired(ctx context.Context, current, threshold int) error
	Sync(ctx context.Context) error
	Close() error
}
