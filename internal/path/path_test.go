package path

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsurePath(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		isDir   bool
		wantErr bool
	}{
		{"Create directory", filepath.Join(tmpDir, "testdir"), true, false},
		{"Create file path", filepath.Join(tmpDir, "testdir", "file.txt"), false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsurePath(tt.path, tt.isDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsurePath() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if tt.isDir {
					if _, err := os.Stat(tt.path); os.IsNotExist(err) {
						t.Errorf("Directory %s was not created", tt.path)
					}
				} else {
					dir := filepath.Dir(tt.path)
					if _, err := os.Stat(dir); os.IsNotExist(err) {
						t.Errorf("Directory %s was not created for file path", dir)
					}
				}
			}
		})
	}
}

func TestMustEnsurePath(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("Valid path", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustEnsurePath() panicked unexpectedly: %v", r)
			}
		}()
		MustEnsurePath(filepath.Join(tmpDir, "testdir"), true)
	})

	t.Run("Invalid path", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustEnsurePath() did not panic on error")
			}
		}()
		MustEnsurePath("/invalid/path", true)
	})
}
