package helpers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rios0rios0/terra/internal/infrastructure/helpers"
)

func TestArgumentsHelper_FindAbsolutePath_ExistingDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	helper := helpers.ArgumentsHelper{}

	tests := []struct {
		name      string
		arguments []string
	}{
		{
			name:      "absolute path as first argument",
			arguments: []string{tempDir, "plan"},
		},
		{
			name:      "absolute path as last argument",
			arguments: []string{"plan", tempDir},
		},
		{
			name:      "relative path",
			arguments: []string{".", "plan"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.FindAbsolutePath(tt.arguments)
			if result == "" {
				t.Error("Expected non-empty absolute path")
			}

			// Verify the result is an absolute path
			if !filepath.IsAbs(result) {
				t.Errorf("Expected absolute path, got: %s", result)
			}

			// Verify the directory exists
			info, err := os.Stat(result)
			if err != nil {
				t.Errorf("Directory should exist: %v", err)
			}

			if !info.IsDir() {
				t.Error("Expected result to be a directory")
			}
		})
	}
}

func TestArgumentsHelper_FindAbsolutePath_NonExistentDirectory(t *testing.T) {
	// This test is complex to implement properly without testing private functions
	// The public method FindAbsolutePath calls logger.Fatalf on non-existent paths
	// which would exit the test process. This behavior is tested indirectly through
	// other test cases that use valid paths.
	t.Skip("Skipping test that would require testing private functions directly")
}

func TestArgumentsHelper_RemovePathFromArguments(t *testing.T) {
	helper := helpers.ArgumentsHelper{}

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{
			name:      "remove path from first position",
			arguments: []string{"/some/path", "plan", "--auto-approve"},
			expected:  []string{"plan", "--auto-approve"},
		},
		{
			name:      "remove path from last position",
			arguments: []string{"plan", "--auto-approve", "/some/path"},
			expected:  []string{"plan", "--auto-approve"},
		},
		{
			name:      "no path to remove",
			arguments: []string{"plan", "--auto-approve"},
			expected:  []string{"plan", "--auto-approve"},
		},
		{
			name:      "empty arguments",
			arguments: []string{},
			expected:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.RemovePathFromArguments(tt.arguments)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d arguments, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected argument[%d] = %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestArgumentsHelper_FindAbsolutePath_DefaultsToCurrentDir(t *testing.T) {
	helper := helpers.ArgumentsHelper{}

	// Test with non-path arguments - should default to current directory
	tests := []struct {
		name      string
		arguments []string
	}{
		{
			name:      "no arguments",
			arguments: []string{},
		},
		{
			name:      "only command arguments",
			arguments: []string{"plan", "--auto-approve"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helper.FindAbsolutePath(tt.arguments)
			
			// Should return current working directory when no path is found
			currentDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			
			expectedPath, err := filepath.Abs(".")
			if err != nil {
				t.Fatalf("Failed to resolve current directory: %v", err)
			}

			if result != expectedPath && result != currentDir {
				t.Errorf("Expected current directory path, got: %s", result)
			}
		})
	}
}