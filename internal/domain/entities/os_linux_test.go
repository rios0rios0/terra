package entities

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSLinux_Download(t *testing.T) {
	testDownloadSuccess(t, &OSLinux{}, "test_download_linux")
}

func TestOSLinux_DownloadHTTPError(t *testing.T) {
	testDownloadHTTPError(t, &OSLinux{}, "test_download_linux")
}

func TestOSLinux_GetTempDir(t *testing.T) {
	osLinux := &OSLinux{}
	tempDir := osLinux.GetTempDir()

	// Verify temp directory is not empty
	if tempDir == "" {
		t.Error("GetTempDir() returned empty string")
	}

	// Verify temp directory exists and is writable
	tempFile := filepath.Join(tempDir, "terra_test_write_permission")
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Cannot write to temp directory %s: %v", tempDir, err)
	}
	file.Close()

	// Clean up
	err = os.Remove(tempFile)
	if err != nil {
		t.Logf("Failed to clean up test file %s: %v", tempFile, err)
	}
}
