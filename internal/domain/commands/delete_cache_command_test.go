package commands_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDeleteCacheCommand(t *testing.T) {
	cmd := NewDeleteCacheCommand()

	if cmd == nil {
		t.Fatal("NewDeleteCacheCommand returned nil")
	}
}

func TestDeleteCacheCommand_Execute(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Change to temp directory for the test
	t.Chdir(tempDir)

	// Create test directories to be deleted
	testDirs := []string{
		".terraform",
		"module1/.terraform",
		"module2/.terragrunt-cache",
		"nested/deep/.terraform",
		"nested/deep/.terragrunt-cache",
	}

	for _, dir := range testDirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}

		// Create a test file in each directory
		testFile := filepath.Join(dir, "test.txt")
		err = os.WriteFile(testFile, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", testFile, err)
		}
	}

	// Also create some directories that should NOT be deleted
	keepDirs := []string{
		"src",
		"docs",
		"other/.terraform-backup", // Similar name but should not match
	}

	for _, dir := range keepDirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create keep directory %s: %v", dir, err)
		}
	}

	// Verify directories exist before deletion
	for _, dir := range testDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Fatalf("Test directory %s was not created", dir)
		}
	}

	// Execute the delete command
	cmd := NewDeleteCacheCommand()
	cmd.Execute([]string{".terraform", ".terragrunt-cache"})

	// Verify that target directories were deleted
	for _, dir := range testDirs {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Errorf("Directory %s should have been deleted but still exists", dir)
		}
	}

	// Verify that other directories were NOT deleted
	for _, dir := range keepDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Directory %s should not have been deleted but was removed", dir)
		}
	}
}

func TestDeleteCacheCommand_ExecuteWithEmptyList(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Change to temp directory for the test
	t.Chdir(tempDir)

	// Create a test directory
	testDir := ".terraform"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Execute with empty list
	cmd := NewDeleteCacheCommand()
	cmd.Execute([]string{})

	// Verify directory still exists (should not be deleted)
	if _, statErr := os.Stat(testDir); os.IsNotExist(statErr) {
		t.Error("Directory should not have been deleted when no targets specified")
	}
}

func TestDeleteCacheCommand_ExecuteWithNonExistentDirectories(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Change to temp directory for the test
	t.Chdir(tempDir)

	// Execute with non-existent directory names
	cmd := NewDeleteCacheCommand()

	// This should not crash or error - it should just find no directories to delete
	cmd.Execute([]string{".nonexistent", ".alsononexistent"})

	// Test passed if no panic occurred
}

func TestDeleteCacheCommand_ExecuteWithSpecificDirectory(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Change to temp directory for the test
	t.Chdir(tempDir)

	// Create test directories
	terraformDir := "project/.terraform"
	terragruntDir := "project/.terragrunt-cache"

	err := os.MkdirAll(terraformDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create terraform directory: %v", err)
	}

	err = os.MkdirAll(terragruntDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create terragrunt directory: %v", err)
	}

	// Execute command to delete only terraform directories
	cmd := NewDeleteCacheCommand()
	cmd.Execute([]string{".terraform"})

	// Verify terraform directory was deleted
	if _, statErr := os.Stat(terraformDir); !os.IsNotExist(statErr) {
		t.Error("Terraform directory should have been deleted")
	}

	// Verify terragrunt directory still exists
	if _, statErr := os.Stat(terragruntDir); os.IsNotExist(statErr) {
		t.Error("Terragrunt directory should not have been deleted")
	}
}
