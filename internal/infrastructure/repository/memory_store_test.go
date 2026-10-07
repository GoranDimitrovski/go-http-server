package repository

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
)

func TestMemoryStore_RecordPrunesAndPersists(t *testing.T) {
	ctx := context.Background()
	filename := filepath.Join(t.TempDir(), "ts.log")
	store := NewMemoryStore(filename)

	for i, tc := range []struct{ now, want int }{
		{100, 1},
		{130, 2},
		{159, 3}, // 159-100 = 59 < 60, still inside
		{160, 3}, // 100 expires exactly at the threshold
	} {
		if got, err := store.Record(ctx, tc.now, 60); err != nil || got != tc.want {
			t.Fatalf("step %d: Record(%d) = %d, %v; want %d", i, tc.now, got, err, tc.want)
		}
	}

	reloaded := NewMemoryStore(filename)
	if err := reloaded.Load(ctx); err != nil {
		t.Fatal(err)
	}
	if got, _ := reloaded.Prune(ctx, 160, 60); got != 3 {
		t.Errorf("after reload Prune() = %d, want 3", got)
	}
	if got, _ := reloaded.Prune(ctx, 1000, 60); got != 0 {
		t.Errorf("Prune(far future) = %d, want 0", got)
	}
}

// Run with -race: concurrent Records must neither race nor lose counts.
func TestMemoryStore_ConcurrentRecord(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore(filepath.Join(t.TempDir(), "ts.log"))

	const n = 100
	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := store.Record(ctx, 100, 60); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	if got, _ := store.Prune(ctx, 100, 60); got != n {
		t.Errorf("count = %d, want %d", got, n)
	}
}
