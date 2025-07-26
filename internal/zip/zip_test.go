package zip

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestZipFromPath(t *testing.T) {
	tmpDir := t.TempDir()

	os.Mkdir(filepath.Join(tmpDir, "folder1"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "folder1", "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "folder1", "file2.txt"), []byte("content2"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "folder2"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "folder2", "file3.txt"), []byte("content3"), 0644)

	var buf bytes.Buffer
	err := ZipFromPath(tmpDir, &buf, []string{"folder2"})
	if err != nil {
		t.Fatalf("ZipFromPath() failed: %v", err)
	}

	reader := bytes.NewReader(buf.Bytes())
	zipReader, err := zip.NewReader(reader, int64(buf.Len()))
	if err != nil {
		t.Fatalf("Failed to read zip archive: %v", err)
	}

	rootDir := filepath.Base(tmpDir) // This will be the name of the root directory in the zip archive
	expectedFiles := map[string]bool{
		rootDir + "/folder1/":          true,
		rootDir + "/folder1/file1.txt": true,
		rootDir + "/folder1/file2.txt": true,
	}
	unexpectedFiles := map[string]bool{
		rootDir + "/folder2/":          true,
		rootDir + "/folder2/file3.txt": true,
	}

	foundFiles := make(map[string]bool)

	for _, file := range zipReader.File {
		foundFiles[file.Name] = true
	}

	for expected := range expectedFiles {
		if !foundFiles[expected] {
			t.Errorf("Expected file %s not found in zip archive", expected)
		}
	}

	for unexpected := range unexpectedFiles {
		if foundFiles[unexpected] {
			t.Errorf("Unexpected file %s found in zip archive", unexpected)
		}
	}
}
