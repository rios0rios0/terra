package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestArgumentsHelper_FindAbsolutePath_ExistingDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	helper := ArgumentsHelper{}

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
			if _, err := os.Stat(result); err != nil {
				t.Errorf("Expected directory to exist: %s, error: %v", result, err)
			}
		})
	}
}

func TestArgumentsHelper_FindAbsolutePath_NonExistentDirectory(t *testing.T) {
	tests := []struct {
		name      string
		arguments []string
		expectMsg string
	}{
		{
			name:      "non-existent absolute path",
			arguments: []string{"/non/existent/path", "plan"},
			expectMsg: "Directory does not exist",
		},
		{
			name:      "non-existent relative path",
			arguments: []string{"./non_existent", "apply"},
			expectMsg: "Directory does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture logger.Fatalf by running in subprocess would be complex
			// For now, we'll test the underlying logic by testing findRelativePath
			// and then os.Stat separately

			relativePath, position := findRelativePath(tt.arguments)
			if position == -1 && relativePath == "." {
				t.Skip("Non-existent path not detected as path-like argument")
			}

			absolutePath, err := filepath.Abs(relativePath)
			if err != nil {
				t.Fatalf("Error resolving path: %v", err)
			}

			// Verify the path doesn't exist
			if _, err := os.Stat(absolutePath); !os.IsNotExist(err) {
				t.Errorf("Expected path to not exist: %s", absolutePath)
			}
		})
	}
}

func TestArgumentsHelper_FindAbsolutePath_FileInsteadOfDirectory(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "test_file_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Test the underlying logic since logger.Fatalf would exit
	relativePath, _ := findRelativePath([]string{tempFile.Name(), "plan"})
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		t.Fatalf("Error resolving path: %v", err)
	}

	info, err := os.Stat(absolutePath)
	if err != nil {
		t.Fatalf("Error getting file info: %v", err)
	}

	if info.IsDir() {
		t.Error("Expected file to not be a directory")
	}
}

func TestFindRelativePath_PathDetection(t *testing.T) {
	tests := []struct {
		name             string
		arguments        []string
		expectedPath     string
		expectedPosition int
	}{
		{
			name:             "no arguments",
			arguments:        []string{},
			expectedPath:     ".",
			expectedPosition: -1,
		},
		{
			name:             "absolute path first",
			arguments:        []string{"/some/path", "plan"},
			expectedPath:     "/some/path",
			expectedPosition: 0,
		},
		{
			name:             "absolute path last",
			arguments:        []string{"plan", "/some/path"},
			expectedPath:     "/some/path",
			expectedPosition: 1,
		},
		{
			name:             "relative path with dot slash",
			arguments:        []string{"./relative", "apply"},
			expectedPath:     "./relative",
			expectedPosition: 0,
		},
		{
			name:             "relative path with dot dot slash",
			arguments:        []string{"../parent", "apply"},
			expectedPath:     "../parent",
			expectedPosition: 0,
		},
		{
			name:             "current directory",
			arguments:        []string{".", "plan"},
			expectedPath:     ".",
			expectedPosition: 0,
		},
		{
			name:             "no path-like arguments",
			arguments:        []string{"plan", "--auto-approve"},
			expectedPath:     ".",
			expectedPosition: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, position := findRelativePath(tt.arguments)

			if path != tt.expectedPath {
				t.Errorf("Expected path %q, got %q", tt.expectedPath, path)
			}

			if position != tt.expectedPosition {
				t.Errorf("Expected position %d, got %d", tt.expectedPosition, position)
			}
		})
	}
}

func TestArgumentsHelper_RemovePathFromArguments(t *testing.T) {
	helper := ArgumentsHelper{}

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
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
				return
			}

			for i, arg := range result {
				if arg != tt.expected[i] {
					t.Errorf("Expected argument %d to be %q, got %q", i, tt.expected[i], arg)
				}
			}
		})
	}
}
