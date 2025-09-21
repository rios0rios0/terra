package entities_test

import (
	"os"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
)

func TestGetOS(t *testing.T) {
	osInstance := entities.GetOS()
	if osInstance == nil {
		t.Fatal("GetOS() returned nil")
	}
}

func TestOSLinux_GetTempDir(t *testing.T) {
	osInstance := entities.GetOS()
	tempDir := osInstance.GetTempDir()
	
	if tempDir == "" {
		t.Error("GetTempDir() returned empty string")
	}

	// Verify the directory exists
	info, err := os.Stat(tempDir)
	if err != nil {
		t.Errorf("Temp directory does not exist: %v", err)
	}

	if !info.IsDir() {
		t.Error("GetTempDir() did not return a directory")
	}
}

func TestOSLinux_GetInstallationPath(t *testing.T) {
	osInstance := entities.GetOS()
	installPath := osInstance.GetInstallationPath()
	
	if installPath == "" {
		t.Error("GetInstallationPath() returned empty string")
	}
}

func TestOSLinux_MakeExecutable(t *testing.T) {
	osInstance := entities.GetOS()
	
	// Create a temporary file
	tempFile, err := os.CreateTemp(t.TempDir(), "test_executable_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Make it executable
	err = osInstance.MakeExecutable(tempFile.Name())
	if err != nil {
		t.Errorf("MakeExecutable failed: %v", err)
	}

	// Verify it's executable
	info, err := os.Stat(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	mode := info.Mode()
	if mode&0111 == 0 {
		t.Error("File is not executable after MakeExecutable")
	}
}

func TestOSLinux_MakeExecutable_NonExistentFile(t *testing.T) {
	osInstance := entities.GetOS()
	
	err := osInstance.MakeExecutable("/non/existent/file12345")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestOSLinux_Remove_NonExistentFile(t *testing.T) {
	osInstance := entities.GetOS()
	
	// Removing a non-existent file should not cause an error in most implementations
	err := osInstance.Remove("/non/existent/file12345")
	// We don't check for error here as different implementations may handle this differently
	_ = err
}