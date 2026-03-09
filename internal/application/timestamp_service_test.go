package application

import (
	"context"
	"errors"
	"testing"
	"time"
)
type mockRepo struct {
	timestamps []int
	storeErr   error
	loadErr    error
	syncErr    error
	countErr   error
	removeErr  error
}

func (m *mockRepo) Store(ctx context.Context, timestamp int) error {
	if m.storeErr != nil {
		return m.storeErr
	}
	m.timestamps = append(m.timestamps, timestamp)
	return nil
}

func (m *mockRepo) View(ctx context.Context) ([]int, error) {
	result := make([]int, len(m.timestamps))
	copy(result, m.timestamps)
	return result, nil
}

func (m *mockRepo) Count(ctx context.Context) (int, error) {
	if m.countErr != nil {
		return 0, m.countErr
	}
	return len(m.timestamps), nil
}

func (m *mockRepo) Load(ctx context.Context) error {
	return m.loadErr
}

func (m *mockRepo) RemoveExpired(ctx context.Context, current, threshold int) error {
	if m.removeErr != nil {
		return m.removeErr
	}
	var valid []int
	for _, ts := range m.timestamps {
		if current-ts < threshold {
			valid = append(valid, ts)
		}
	}
	m.timestamps = valid
	return nil
}

func (m *mockRepo) Sync(ctx context.Context) error {
	return m.syncErr
}

func (m *mockRepo) Close() error {
	return nil
}

func TestTimestampService_Initialize(t *testing.T) {
	tests := []struct {
		name      string
		mockRepo *mockRepo
		wantErr   bool
	}{
		{
			name: "successful initialization",
			mockRepo: &mockRepo{
				timestamps: []int{int(time.Now().Unix()) - 100},
			},
			wantErr: false,
		},
		{
			name: "load error",
			mockRepo: &mockRepo{
				loadErr: errors.New("load failed"),
			},
			wantErr: true,
		},
		{
			name: "sync error",
			mockRepo: &mockRepo{
				syncErr: errors.New("sync failed"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTimestampService(tt.mockRepo, 60)
			ctx := context.Background()

			err := service.Initialize(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimestampService_RecordTimestamp(t *testing.T) {
	now := int(time.Now().Unix())
	tests := []struct {
		name      string
		mockRepo *mockRepo
		wantCount int
		wantErr   bool
	}{
		{
			name: "successful record",
			mockRepo: &mockRepo{
				timestamps: []int{now - 30},
			},
			wantCount: 2, // existing + new
			wantErr:   false,
		},
		{
			name: "remove expired error",
			mockRepo: &mockRepo{
				removeErr: errors.New("remove failed"),
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "store error",
			mockRepo: &mockRepo{
				storeErr: errors.New("store failed"),
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "sync error",
			mockRepo: &mockRepo{
				syncErr: errors.New("sync failed"),
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "count error",
			mockRepo: &mockRepo{
				countErr: errors.New("count failed"),
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTimestampService(tt.mockRepo, 60)
			ctx := context.Background()

			count, err := service.RecordTimestamp(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && count != tt.wantCount {
				t.Errorf("RecordTimestamp() count = %v, want %v", count, tt.wantCount)
			}
		})
	}
}
