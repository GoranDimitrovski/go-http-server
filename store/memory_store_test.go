package store

import (
	"os"
	"testing"
	"time"
)

var (
	fileName string
    store *TimestampMemoryStore
)

func TestMain(m *testing.M) {
	fileName := "timestamp.log"
	store = New(fileName)
}

func TestStore_Success(t *testing.T) {
	t1 := int(time.Now().Unix())
	store.Store(t1)
	t2 := t1 - 60
	store.Store(t2)

	store.Load()
	count := store.Count()
	expected := 2
	if count != 2 {
		t.Errorf("Expected count %d, got %d", expected, count)
	}

	os.Remove(fileName)
}

func TestStore_Fail(t *testing.T) {
	store.Load()

	count := store.Count()
	expected := 0
	if count != expected {
		t.Errorf("Expected count %d, got %d", expected, count)
	}
}

func TestView_Success(t *testing.T) {
	t1 := int(time.Now().Unix())
	store.Store(t1)
	t2 := t1 - 60
	store.Store(t2)

	store.Load()
	view := store.View()

	if len(view) != 2 {
		t.Errorf("Expected view length 2, got %d", len(view))
	}

	if !contains(view, t1) || !contains(view, t2) {
		t.Errorf("The view does not contain the timestamps")
	}
}

func TestView_Fail(t *testing.T) {
	view := store.View()

	if len(view) != 0 {
		t.Errorf("Expected view length 0, got %d", len(view))
	}
}

func TestLoad_Success(t *testing.T) {
	t1 := int(time.Now().Unix())
	store.Store(t1)
	t2 := t1 - 60
	store.Store(t2)

	store.Load()

	count := store.Count()
	expected := 2
	if count != expected {
		t.Errorf("Expected count %d after Load, got %d", expected, count)
	}
}

func TestLoad_Fail(t *testing.T) {
	fileName := "nonexistent.log"
	store := New(fileName)
	store.Load()

	count := store.Count()
	expected := 0
	if count != expected {
		t.Errorf("Expected count %d after Load, got %d", expected, count)
	}
}

func TestRemoveExpired_Success(t *testing.T) {
    t1 := int(time.Now().Unix())
	store.Store(t1)
	t2 := t1 - 60
	store.Store(t2)

	store.RemoveExpired(t1, 30)

	count := store.Count()
	expected := 1
	if count != expected {
		t.Errorf("Expected count %d after RemoveExpired, got %d", expected, count)
	}
}

func TestRemoveExpired_Fail(t *testing.T) {
	t1 := int(time.Now().Unix())
	store.Store(t1)
	store.RemoveExpired(t1, 120)

	count := store.Count()
	expected := 1
	if count != expected {
		t.Errorf("Expected count %d after RemoveExpired, got %d", expected, count)
	}
}

func TestSync_Success(t *testing.T) {
	t1 := int(time.Now().Unix())
	store.Store(t1)
	t2 := t1 - 60
	store.Store(t2)
	store.Sync()

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Errorf("File %s does not exist after Sync", fileName)
	}

	os.Remove(fileName)
}

func TestSync_Fail(t *testing.T) {
	fileName := "nonexistent.log"
	store := New(fileName)
	store.Sync()

	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		t.Errorf("File %s exists after Sync", fileName)
	}
}

func contains(slice []int, value int) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
