package persistence

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestFilePersistence_Append(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int
		wantErr   bool
	}{
		{
			name:      "append valid timestamp",
			timestamp: int(time.Now().Unix()),
			wantErr:   false,
		},
		{
			name:      "append zero timestamp",
			timestamp: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "test_append.log"
			defer os.Remove(filename)

			persister := NewFilePersistence()
			ctx := context.Background()

			err := persister.Append(ctx, tt.timestamp, filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify file was created and contains the timestamp
			if !persister.FileExists(filename) {
				t.Error("File should exist after Append()")
			}

			timestamps, err := persister.ReadAll(ctx, filename)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}

			if len(timestamps) != 1 || timestamps[0] != tt.timestamp {
				t.Errorf("ReadAll() = %v, want [%v]", timestamps, tt.timestamp)
			}
		})
	}
}

func TestFilePersistence_Rewrite(t *testing.T) {
	tests := []struct {
		name       string
		timestamps []int
		wantErr    bool
	}{
		{
			name:       "rewrite with empty slice",
			timestamps: []int{},
			wantErr:    false,
		},
		{
			name:       "rewrite with single timestamp",
			timestamps: []int{int(time.Now().Unix())},
			wantErr:    false,
		},
		{
			name:       "rewrite with multiple timestamps",
			timestamps: []int{int(time.Now().Unix()), int(time.Now().Unix()) - 10, int(time.Now().Unix()) - 20},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "test_rewrite.log"
			defer os.Remove(filename)

			persister := NewFilePersistence()
			ctx := context.Background()

			// First append something
			if err := persister.Append(ctx, 999999, filename); err != nil {
				t.Fatalf("Append() error = %v", err)
			}

			// Then rewrite
			err := persister.Rewrite(ctx, tt.timestamps, filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rewrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify file contains only the rewritten timestamps
			readTimestamps, err := persister.ReadAll(ctx, filename)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}

			if len(readTimestamps) != len(tt.timestamps) {
				t.Errorf("ReadAll() length = %v, want %v", len(readTimestamps), len(tt.timestamps))
			}

			for i, ts := range tt.timestamps {
				if i < len(readTimestamps) && readTimestamps[i] != ts {
					t.Errorf("ReadAll()[%d] = %v, want %v", i, readTimestamps[i], ts)
				}
			}
		})
	}
}

func TestFilePersistence_ReadAll(t *testing.T) {
	tests := []struct {
		name       string
		setupFile  bool
		fileData   []int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "read from existing file",
			setupFile:  true,
			fileData:   []int{int(time.Now().Unix()), int(time.Now().Unix()) - 10},
			wantCount:  2,
			wantErr:    false,
		},
		{
			name:       "read from non-existent file",
			setupFile:  false,
			fileData:   []int{},
			wantCount:  0,
			wantErr:    false,
		},
		{
			name:       "read from empty file",
			setupFile:  true,
			fileData:   []int{},
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := "test_read.log"
			defer os.Remove(filename)

			persister := NewFilePersistence()
			ctx := context.Background()

			if tt.setupFile {
				if err := persister.Rewrite(ctx, tt.fileData, filename); err != nil {
					t.Fatalf("Failed to setup test file: %v", err)
				}
			}

			timestamps, err := persister.ReadAll(ctx, filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(timestamps) != tt.wantCount {
				t.Errorf("ReadAll() length = %v, want %v", len(timestamps), tt.wantCount)
			}

			for i, ts := range tt.fileData {
				if i < len(timestamps) && timestamps[i] != ts {
					t.Errorf("ReadAll()[%d] = %v, want %v", i, timestamps[i], ts)
				}
			}
		})
	}
}

func TestFilePersistence_FileExists(t *testing.T) {
	tests := []struct {
		name     string
		setup    bool
		filename string
		want     bool
	}{
		{
			name:     "file exists",
			setup:    true,
			filename: "test_exists.log",
			want:     true,
		},
		{
			name:     "file does not exist",
			setup:    false,
			filename: "test_not_exists.log",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.filename)

			persister := NewFilePersistence()
			ctx := context.Background()

			if tt.setup {
				if err := persister.Append(ctx, int(time.Now().Unix()), tt.filename); err != nil {
					t.Fatalf("Failed to setup test file: %v", err)
				}
			}

			got := persister.FileExists(tt.filename)
			if got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

