package store

import (
	"context"
	"os"
	"simplesurance/persistence"
	"testing"
	"time"
)

func TestMemoryStore_Store(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int
		wantErr   bool
	}{
		{
			name:      "store valid timestamp",
			timestamp: int(time.Now().Unix()),
			wantErr:   false,
		},
		{
			name:      "store zero timestamp",
			timestamp: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryStore("test.log", persistence.NewFilePersistence())
			defer cleanup(t, store, "test.log")

			ctx := context.Background()
			err := store.Store(ctx, tt.timestamp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			count, err := store.Count(ctx)
			if err != nil {
				t.Fatalf("Count() error = %v", err)
			}
			if count != 1 {
				t.Errorf("Count() = %v, want 1", count)
			}
		})
	}
}

func TestMemoryStore_View(t *testing.T) {
	tests := []struct {
		name       string
		timestamps []int
		wantCount  int
	}{
		{
			name:       "view empty store",
			timestamps: []int{},
			wantCount:  0,
		},
		{
			name:       "view with one timestamp",
			timestamps: []int{int(time.Now().Unix())},
			wantCount:  1,
		},
		{
			name:       "view with multiple timestamps",
			timestamps: []int{int(time.Now().Unix()), int(time.Now().Unix()) - 10, int(time.Now().Unix()) - 20},
			wantCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryStore("test.log", persistence.NewFilePersistence())
			defer cleanup(t, store, "test.log")

			ctx := context.Background()
			for _, ts := range tt.timestamps {
				if err := store.Store(ctx, ts); err != nil {
					t.Fatalf("Store() error = %v", err)
				}
			}

			view, err := store.View(ctx)
			if err != nil {
				t.Fatalf("View() error = %v", err)
			}

			if len(view) != tt.wantCount {
				t.Errorf("View() length = %v, want %v", len(view), tt.wantCount)
			}

			// Verify view is a copy (modifying view shouldn't affect store)
			if len(view) > 0 {
				view[0] = 999999
				storeView, _ := store.View(ctx)
				if len(storeView) > 0 && storeView[0] == 999999 {
					t.Error("View() returned reference instead of copy")
				}
			}
		})
	}
}

func TestMemoryStore_Count(t *testing.T) {
	tests := []struct {
		name       string
		timestamps []int
		wantCount  int
	}{
		{
			name:       "count empty store",
			timestamps: []int{},
			wantCount:  0,
		},
		{
			name:       "count single timestamp",
			timestamps: []int{int(time.Now().Unix())},
			wantCount:  1,
		},
		{
			name:       "count multiple timestamps",
			timestamps: []int{int(time.Now().Unix()), int(time.Now().Unix()) - 10},
			wantCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryStore("test.log", persistence.NewFilePersistence())
			defer cleanup(t, store, "test.log")

			ctx := context.Background()
			for _, ts := range tt.timestamps {
				if err := store.Store(ctx, ts); err != nil {
					t.Fatalf("Store() error = %v", err)
				}
			}

			count, err := store.Count(ctx)
			if err != nil {
				t.Fatalf("Count() error = %v", err)
			}

			if count != tt.wantCount {
				t.Errorf("Count() = %v, want %v", count, tt.wantCount)
			}
		})
	}
}

func TestMemoryStore_Load(t *testing.T) {
	tests := []struct {
		name       string
		setupFile  bool
		fileData   []int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "load from existing file",
			setupFile:  true,
			fileData:   []int{int(time.Now().Unix()), int(time.Now().Unix()) - 10},
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "load from non-existent file",
			setupFile:  false,
			fileData:   []int{},
			wantCount:  0,
			wantErr:    false,
		},
		{
			name:       "load from empty file",
			setupFile:  true,
			fileData:   []int{},
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "test_load.log"
			defer os.Remove(filename)

			if tt.setupFile {
				persister := persistence.NewFilePersistence()
				ctx := context.Background()
				if err := persister.Rewrite(ctx, tt.fileData, filename); err != nil {
					t.Fatalf("Failed to setup test file: %v", err)
				}
			}

			store := NewMemoryStore(filename, persistence.NewFilePersistence())
			ctx := context.Background()

			err := store.Load(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			count, err := store.Count(ctx)
			if err != nil {
				t.Fatalf("Count() error = %v", err)
			}

			if count != tt.wantCount {
				t.Errorf("Count() after Load() = %v, want %v", count, tt.wantCount)
			}
		})
	}
}

func TestMemoryStore_RemoveExpired(t *testing.T) {
	now := int(time.Now().Unix())
	tests := []struct {
		name       string
		timestamps []int
		current    int
		threshold  int
		wantCount  int
	}{
		{
			name:       "remove expired timestamps",
			timestamps: []int{now, now - 30, now - 70},
			current:    now,
			threshold:  60,
			wantCount:  2, // now and now-30 should remain
		},
		{
			name:       "keep all timestamps within threshold",
			timestamps: []int{now, now - 30, now - 40},
			current:    now,
			threshold:  60,
			wantCount:  3,
		},
		{
			name:       "remove all expired timestamps",
			timestamps: []int{now - 70, now - 80, now - 90},
			current:    now,
			threshold:  60,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryStore("test.log", persistence.NewFilePersistence())
			defer cleanup(t, store, "test.log")

			ctx := context.Background()
			for _, ts := range tt.timestamps {
				if err := store.Store(ctx, ts); err != nil {
					t.Fatalf("Store() error = %v", err)
				}
			}

			if err := store.RemoveExpired(ctx, tt.current, tt.threshold); err != nil {
				t.Fatalf("RemoveExpired() error = %v", err)
			}

			count, err := store.Count(ctx)
			if err != nil {
				t.Fatalf("Count() error = %v", err)
			}

			if count != tt.wantCount {
				t.Errorf("Count() after RemoveExpired() = %v, want %v", count, tt.wantCount)
			}
		})
	}
}

func TestMemoryStore_Sync(t *testing.T) {
	tests := []struct {
		name       string
		timestamps []int
	}{
		{
			name:       "sync empty store",
			timestamps: []int{},
		},
		{
			name:       "sync with timestamps",
			timestamps: []int{int(time.Now().Unix()), int(time.Now().Unix()) - 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "test_sync.log"
			defer os.Remove(filename)

			store := NewMemoryStore(filename, persistence.NewFilePersistence())
			ctx := context.Background()

			for _, ts := range tt.timestamps {
				if err := store.Store(ctx, ts); err != nil {
					t.Fatalf("Store() error = %v", err)
				}
			}

			if err := store.Sync(ctx); err != nil {
				t.Fatalf("Sync() error = %v", err)
			}

			// Verify file exists and contains correct data
			if _, err := os.Stat(filename); err != nil {
				if len(tt.timestamps) == 0 {
					// Empty file might not exist, which is okay
					return
				}
				t.Fatalf("File should exist after Sync(): %v", err)
			}

			// Load into new store and verify
			newStore := NewMemoryStore(filename, persistence.NewFilePersistence())
			if err := newStore.Load(ctx); err != nil {
				t.Fatalf("Load() after Sync() error = %v", err)
			}

			count, err := newStore.Count(ctx)
			if err != nil {
				t.Fatalf("Count() error = %v", err)
			}

			if count != len(tt.timestamps) {
				t.Errorf("Count() after Load() = %v, want %v", count, len(tt.timestamps))
			}
		})
	}
}

func cleanup(t *testing.T, store *MemoryStore, filename string) {
	t.Helper()
	store.Close()
	os.Remove(filename)
}
