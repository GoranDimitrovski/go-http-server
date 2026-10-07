package persistence

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestWriteAllReadAll(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "ts.log")

	got, err := ReadAll(filename)
	if err != nil || len(got) != 0 {
		t.Fatalf("ReadAll(missing) = %v, %v; want empty, nil", got, err)
	}

	for _, want := range [][]int{{3, 1, 2}, {}} { // second write must fully replace the first
		if err := WriteAll(filename, want); err != nil {
			t.Fatalf("WriteAll(%v) error = %v", want, err)
		}
		got, err := ReadAll(filename)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("ReadAll() = %v, %v; want %v", got, err, want)
		}
	}
	if _, err := os.Stat(filename + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("temp file left behind: %v", err)
	}
}

func TestReadAllRejectsGarbage(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "ts.log")
	os.WriteFile(filename, []byte("1\nnope\n"), 0o644)
	if _, err := ReadAll(filename); err == nil {
		t.Error("ReadAll() error = nil, want parse error")
	}
}
