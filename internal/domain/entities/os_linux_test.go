package entities

import (
	"os"
	"path/filepath"
	"strings"
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

	// Verify it returns the system temp directory
	expectedTempDir := os.TempDir()
	if tempDir != expectedTempDir {
		t.Errorf("Expected temp dir %s, got %s", expectedTempDir, tempDir)
	}
}

func TestOSLinux_GetInstallationPath(t *testing.T) {
	osLinux := &OSLinux{}
	installPath := osLinux.GetInstallationPath()

	if installPath == "" {
		t.Error("GetInstallationPath should not return empty string")
	}

	// Should contain .local/bin
	if !strings.Contains(installPath, ".local/bin") {
		t.Errorf("Expected installation path to contain '.local/bin', got %s", installPath)
	}

	// Should be an absolute path (contain '/')
	if !strings.Contains(installPath, "/") {
		t.Errorf("Expected installation path to be absolute, got %s", installPath)
	}
}

func TestOSLinux_MakeExecutable(t *testing.T) {
	osLinux := &OSLinux{}

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test_executable_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Initially the file should not be executable
	info, err := os.Stat(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat temp file: %v", err)
	}

	// Make it executable
	err = osLinux.MakeExecutable(tempFile.Name())
	if err != nil {
		t.Errorf("MakeExecutable failed: %v", err)
	}

	// Verify it's now executable
	info, err = os.Stat(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat temp file after making executable: %v", err)
	}

	// Check if owner has execute permission
	mode := info.Mode()
	if mode.Perm()&0o100 == 0 {
		t.Error("File should be executable by owner after MakeExecutable")
	}
}

func TestOSLinux_MakeExecutable_NonExistentFile(t *testing.T) {
	osLinux := &OSLinux{}

	// Try to make a non-existent file executable
	err := osLinux.MakeExecutable("/tmp/non_existent_file_12345")
	if err == nil {
		t.Error("Expected error when making non-existent file executable")
	}

	// Verify error message contains expected text
	if !strings.Contains(err.Error(), "failed to perform change binary permissions") {
		t.Errorf(
			"Expected error to contain 'failed to perform change binary permissions', got: %v",
			err,
		)
	}
}

func TestOSLinux_Move(t *testing.T) {
	osLinux := &OSLinux{}

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test_move_source_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Cleanup in case test fails

	testContent := "test content for move"
	_, err = tempFile.WriteString(testContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Create destination path
	destPath := tempFile.Name() + "_moved"
	defer os.Remove(destPath) // Cleanup

	// Move the file
	err = osLinux.Move(tempFile.Name(), destPath)
	if err != nil {
		t.Errorf("Move failed: %v", err)
	}

	// Verify source file is gone
	if _, err := os.Stat(tempFile.Name()); !os.IsNotExist(err) {
		t.Error("Source file should not exist after move")
	}

	// Verify destination file exists and has correct content
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Error("Destination file should exist after move")
	}

	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content %q, got %q", testContent, string(content))
	}
}

func TestOSLinux_Move_NonExistentSource(t *testing.T) {
	osLinux := &OSLinux{}

	// Try to move a non-existent file
	err := osLinux.Move("/tmp/non_existent_source_12345", "/tmp/destination_12345")
	if err == nil {
		t.Error("Expected error when moving non-existent file")
	}

	// Verify error message contains expected text
	if !strings.Contains(err.Error(), "failed to perform moving folder") {
		t.Errorf("Expected error to contain 'failed to perform moving folder', got: %v", err)
	}
}

func TestOSLinux_Remove(t *testing.T) {
	osLinux := &OSLinux{}

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test_remove_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()

	// Verify file exists
	if _, err := os.Stat(tempFile.Name()); os.IsNotExist(err) {
		t.Fatal("Temp file should exist before removal")
	}

	// Remove the file
	err = osLinux.Remove(tempFile.Name())
	if err != nil {
		t.Errorf("Remove failed: %v", err)
	}

	// Verify file is gone
	if _, err := os.Stat(tempFile.Name()); !os.IsNotExist(err) {
		t.Error("File should not exist after removal")
	}
}

func TestOSLinux_Remove_NonExistentFile(t *testing.T) {
	osLinux := &OSLinux{}

	// Try to remove a non-existent file
	err := osLinux.Remove("/tmp/non_existent_file_12345")
	if err == nil {
		t.Error("Expected error when removing non-existent file")
	}

	// Verify error message contains expected text
	if !strings.Contains(err.Error(), "failed to perform deleting folder") {
		t.Errorf("Expected error to contain 'failed to perform deleting folder', got: %v", err)
	}
}

func TestOSLinux_Extract(t *testing.T) {
	osLinux := &OSLinux{}

	// Create a temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "test_extract_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// We can't easily test actual zip extraction without creating a real zip file
	// So we'll test the error case with a non-existent file
	err = osLinux.Extract("/tmp/non_existent_archive_12345.zip", tempDir)
	if err == nil {
		t.Error("Expected error when extracting non-existent archive")
	}

	// Verify error message contains expected text
	if !strings.Contains(err.Error(), "failed to perform decompressing") {
		t.Errorf("Expected error to contain 'failed to perform decompressing', got: %v", err)
	}
}

func TestGetOS(t *testing.T) {
	osLinux := GetOS()

	if osLinux == nil {
		t.Fatal("GetOS should not return nil")
	}

	// Verify it implements the OS interface by calling a method
	tempDir := osLinux.GetTempDir()
	if tempDir == "" {
		t.Error("GetOS returned object should implement OS interface correctly")
	}
}
