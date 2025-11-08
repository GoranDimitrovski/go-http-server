package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

// mockStore is a mock implementation of store.Store for testing
type mockStore struct {
	timestamps []int
	storeErr   error
	loadErr    error
	syncErr    error
	countErr   error
	removeErr  error
}

func (m *mockStore) Store(ctx context.Context, timestamp int) error {
	if m.storeErr != nil {
		return m.storeErr
	}
	m.timestamps = append(m.timestamps, timestamp)
	return nil
}

func (m *mockStore) View(ctx context.Context) ([]int, error) {
	result := make([]int, len(m.timestamps))
	copy(result, m.timestamps)
	return result, nil
}

func (m *mockStore) Count(ctx context.Context) (int, error) {
	if m.countErr != nil {
		return 0, m.countErr
	}
	return len(m.timestamps), nil
}

func (m *mockStore) Load(ctx context.Context) error {
	return m.loadErr
}

func (m *mockStore) RemoveExpired(ctx context.Context, current, threshold int) error {
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

func (m *mockStore) Sync(ctx context.Context) error {
	return m.syncErr
}

func (m *mockStore) Close() error {
	return nil
}

func TestTimestampService_Initialize(t *testing.T) {
	tests := []struct {
		name      string
		mockStore *mockStore
		wantErr   bool
	}{
		{
			name: "successful initialization",
			mockStore: &mockStore{
				timestamps: []int{int(time.Now().Unix()) - 100},
			},
			wantErr: false,
		},
		{
			name: "load error",
			mockStore: &mockStore{
				loadErr: errors.New("load failed"),
			},
			wantErr: true,
		},
		{
			name: "sync error",
			mockStore: &mockStore{
				syncErr: errors.New("sync failed"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTimestampService(tt.mockStore, 60)
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
		mockStore *mockStore
		wantCount int
		wantErr   bool
	}{
		{
			name: "successful record",
			mockStore: &mockStore{
				timestamps: []int{now - 30},
			},
			wantCount: 2, // existing + new
			wantErr:   false,
		},
		{
			name: "remove expired error",
			mockStore: &mockStore{
				removeErr: errors.New("remove failed"),
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "store error",
			mockStore: &mockStore{
				storeErr: errors.New("store failed"),
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "sync error",
			mockStore: &mockStore{
				syncErr: errors.New("sync failed"),
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "count error",
			mockStore: &mockStore{
				countErr: errors.New("count failed"),
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTimestampService(tt.mockStore, 60)
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

