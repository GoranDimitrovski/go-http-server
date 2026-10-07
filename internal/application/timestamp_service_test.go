package application

import (
	"context"
	"errors"
	"testing"
)

type mockRepo struct {
	count            int
	loadErr, saveErr error
}

func (m *mockRepo) Load(ctx context.Context) error { return m.loadErr }
func (m *mockRepo) Record(ctx context.Context, now, threshold int) (int, error) {
	m.count++
	return m.count, m.saveErr
}
func (m *mockRepo) Prune(ctx context.Context, now, threshold int) (int, error) {
	return m.count, m.saveErr
}
func (m *mockRepo) Close() error { return nil }

func TestTimestampService(t *testing.T) {
	ctx := context.Background()
	boom := errors.New("boom")

	if err := NewTimestampService(&mockRepo{loadErr: boom}, 60).Initialize(ctx); !errors.Is(err, boom) {
		t.Errorf("Initialize() with load error = %v, want boom", err)
	}
	if err := NewTimestampService(&mockRepo{saveErr: boom}, 60).Initialize(ctx); !errors.Is(err, boom) {
		t.Errorf("Initialize() with prune error = %v, want boom", err)
	}
	if _, err := NewTimestampService(&mockRepo{saveErr: boom}, 60).RecordTimestamp(ctx); !errors.Is(err, boom) {
		t.Errorf("RecordTimestamp() with save error = %v, want boom", err)
	}
	if got, err := NewTimestampService(&mockRepo{count: 1}, 60).RecordTimestamp(ctx); err != nil || got != 2 {
		t.Errorf("RecordTimestamp() = %d, %v; want 2, nil", got, err)
	}
}
