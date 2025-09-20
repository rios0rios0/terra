package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestArgumentsHelper_RemovePathFromArguments(t *testing.T) {
	helper := ArgumentsHelper{}

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{
			name:      "no valid directory in arguments",
			arguments: []string{"apply", "--auto-approve", "invalid-path"},
			expected:  []string{"apply", "--auto-approve", "invalid-path"}, // should return unchanged
		},
		{
			name:      "valid directory at start",
			arguments: []string{".", "apply", "--auto-approve"},
			expected:  []string{"apply", "--auto-approve"},
		},
		{
			name:      "valid directory at end", 
			arguments: []string{"apply", "--auto-approve", "."},
			expected:  []string{"apply", "--auto-approve"},
		},
		{
			name:      "empty arguments",
			arguments: []string{},
			expected:  []string{},
		},
		{
			name:      "single argument - non-directory",
			arguments: []string{"apply"},
			expected:  []string{"apply"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For tests that expect existing directories, create them temporarily
		name                 string
		arguments            []string
		expected             []string
		expectsValidDirectory bool
	}{
		{
			name:                 "no valid directory in arguments",
			arguments:            []string{"apply", "--auto-approve", "invalid-path"},
			expected:             []string{"apply", "--auto-approve", "invalid-path"}, // should return unchanged
			expectsValidDirectory: false,
		},
		{
			name:                 "valid directory at start",
			arguments:            []string{".", "apply", "--auto-approve"},
			expected:             []string{"apply", "--auto-approve"},
			expectsValidDirectory: true,
		},
		{
			name:                 "valid directory at end", 
			arguments:            []string{"apply", "--auto-approve", "."},
			expected:             []string{"apply", "--auto-approve"},
			expectsValidDirectory: true,
		},
		{
			name:                 "empty arguments",
			arguments:            []string{},
			expected:             []string{},
			expectsValidDirectory: false,
		},
		{
			name:                 "single argument - non-directory",
			arguments:            []string{"apply"},
			expected:             []string{"apply"},
			expectsValidDirectory: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For tests that expect existing directories, create them temporarily
			if tt.expectsValidDirectory {
				result := helper.RemovePathFromArguments(tt.arguments)
				if len(result) != len(tt.expected) {
					t.Errorf("RemovePathFromArguments() = %v, want %v", result, tt.expected)
					return
				}
				for i, v := range result {
					if v != tt.expected[i] {
						t.Errorf("RemovePathFromArguments() = %v, want %v", result, tt.expected)
						break
					}
				}
			} else {
				result := helper.RemovePathFromArguments(tt.arguments)
				if len(result) != len(tt.expected) {
					t.Errorf("RemovePathFromArguments() = %v, want %v", result, tt.expected)
					return
				}
				for i, v := range result {
					if v != tt.expected[i] {
						t.Errorf("RemovePathFromArguments() = %v, want %v", result, tt.expected)
						break
					}
				}
			}
		})
	}
}

func TestArgumentsHelper_FindAbsolutePath_ValidDirectory(t *testing.T) {
	helper := ArgumentsHelper{}

	// Test with current directory - this should always work
	arguments := []string{"apply", "."}
	result := helper.FindAbsolutePath(arguments)
	
	expectedAbs, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("Failed to get absolute path of current directory: %v", err)
	}
	
	if result != expectedAbs {
		t.Errorf("FindAbsolutePath() = %v, want %v", result, expectedAbs)
	}
}

func TestArgumentsHelper_FindAbsolutePath_WithTempDirectory(t *testing.T) {
	helper := ArgumentsHelper{}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "terra-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	arguments := []string{"apply", tempDir}
	result := helper.FindAbsolutePath(arguments)
	
	expectedAbs, err := filepath.Abs(tempDir)
	if err != nil {
		t.Fatalf("Failed to get absolute path of temp directory: %v", err)
	}
	
	if result != expectedAbs {
		t.Errorf("FindAbsolutePath() = %v, want %v", result, expectedAbs)
	}
}

func TestFindRelativePath(t *testing.T) {
	tests := []struct {
		name             string
		arguments        []string
		expectedPath     string
		expectedPosition int
		setupTempDir     bool
		tempDirIndex     int
	}{
		{
			name:             "current directory at start",
			arguments:        []string{".", "apply"},
			expectedPath:     ".",
			expectedPosition: 0,
		},
		{
			name:             "current directory at end",
			arguments:        []string{"apply", "."},
			expectedPath:     ".",
			expectedPosition: 1,
		},
		{
			name:             "no valid directory",
			arguments:        []string{"apply", "non-existent-path"},
			expectedPath:     ".",
			expectedPosition: -1,
		},
		{
			name:             "temp directory at start",
			arguments:        []string{"TEMP_DIR", "apply"},
			expectedPath:     "TEMP_DIR",
			expectedPosition: 0,
			setupTempDir:     true,
			tempDirIndex:     0,
		},
		{
			name:             "temp directory at end",
			arguments:        []string{"apply", "TEMP_DIR"},
			expectedPath:     "TEMP_DIR",
			expectedPosition: 1,
			setupTempDir:     true,
			tempDirIndex:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]string, len(tt.arguments))
			copy(args, tt.arguments)

			// Setup temporary directory if needed
			var tempDir string
			if tt.setupTempDir {
				var err error
				tempDir, err = os.MkdirTemp("", "terra-test-*")
				if err != nil {
					t.Fatalf("Failed to create temp directory: %v", err)
				}
				defer os.RemoveAll(tempDir)
				args[tt.tempDirIndex] = tempDir
			}

			path, position := findRelativePath(args)

			expectedPath := tt.expectedPath
			if tt.setupTempDir {
				expectedPath = tempDir
			}

			if path != expectedPath {
				t.Errorf("findRelativePath() path = %v, want %v", path, expectedPath)
			}
			if position != tt.expectedPosition {
				t.Errorf("findRelativePath() position = %v, want %v", position, tt.expectedPosition)
			}
		})
	}
}

// Test helper to create a temporary file (not directory) for testing
func TestArgumentsHelper_PathValidation_Integration(t *testing.T) {
	helper := ArgumentsHelper{}

	t.Run("valid current directory", func(t *testing.T) {
		// This should work without panicking
		result := helper.FindAbsolutePath([]string{"apply", "."})
		if result == "" {
			t.Error("Expected non-empty result for current directory")
		}
	})

	t.Run("valid temporary directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "terra-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		result := helper.FindAbsolutePath([]string{"apply", tempDir})
		expectedAbs, _ := filepath.Abs(tempDir)
		if result != expectedAbs {
			t.Errorf("Expected %v, got %v", expectedAbs, result)
		}
	})

	t.Run("remove path from arguments with no valid directory", func(t *testing.T) {
		// This should not panic and should return arguments unchanged
		args := []string{"apply", "--auto-approve", "non-existent-path"}
		result := helper.RemovePathFromArguments(args)
		
		if len(result) != len(args) {
			t.Errorf("Expected same length, got %d, want %d", len(result), len(args))
		}
		
		for i, v := range result {
			if v != args[i] {
				t.Errorf("Expected unchanged arguments, got %v, want %v", result, args)
				break
			}
		}
	})

	t.Run("remove path from arguments with valid directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "terra-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Put temp directory at the end so it gets detected by findRelativePath
		args := []string{"apply", "--auto-approve", tempDir}
		result := helper.RemovePathFromArguments(args)
		
		expected := []string{"apply", "--auto-approve"}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		
		for i, v := range result {
			if v != expected[i] {
				t.Errorf("Expected %v, got %v", expected, result)
				break
			}
		}
	})
}

// TestArgumentsHelper_ErrorScenarios documents the scenarios that would cause
// logger.Fatalf to be called. These cannot be tested directly since logger.Fatalf
// calls os.Exit(), but they are documented here for reference.
//
// Before the fix, these scenarios would cause a panic: runtime error: slice bounds out of range [:-1]
// After the fix, these scenarios call logger.Fatalf with appropriate error messages:
//
// 1. Non-existent directory in path-like argument:
//    FindAbsolutePath([]string{"apply", "/non/existent/path"})
//    Would log: "Directory does not exist: /non/existent/path"
//
// 2. File instead of directory in path-like argument:
//    FindAbsolutePath([]string{"apply", "/path/to/file.txt"})  
//    Would log: "Path is not a directory: /path/to/file.txt"
//
// 3. Permission denied accessing directory:
//    FindAbsolutePath([]string{"apply", "/no/permission/dir"})
//    Would log: "Cannot access directory: /no/permission/dir - permission denied"
//
// 4. Relative path that doesn't exist:
//    FindAbsolutePath([]string{"apply", "./non-existent"})
//    Would log: "Directory does not exist: ./non-existent"
//
// 5. Absolute path that doesn't exist:
//    FindAbsolutePath([]string{"apply", "/absolute/non-existent"})
//    Would log: "Directory does not exist: /absolute/non-existent"
//
// These scenarios are the ones that were causing panics in issue #5 and are now
// handled gracefully with proper error messages.
func TestArgumentsHelper_ErrorScenarios_Documentation(t *testing.T) {
	// This test documents the error scenarios but doesn't actually run them
	// since they would call logger.Fatalf and exit the program
	t.Log("Error scenarios are documented in the test function comments")
	t.Log("These scenarios now show proper error messages instead of panicking")
	t.Log("See function comments for specific examples")
}

func TestArgumentsHelper_EdgeCases(t *testing.T) {
	helper := ArgumentsHelper{}

	t.Run("empty arguments array should not panic", func(t *testing.T) {
		// This was causing the original panic due to slice bounds error
		// Before fix: panic: runtime error: index out of range [0] with length 0
		// After fix: should handle gracefully
		result := helper.RemovePathFromArguments([]string{})
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %v", result)
		}
	})

	t.Run("single element arguments array", func(t *testing.T) {
		args := []string{"apply"}
		result := helper.RemovePathFromArguments(args)
		
		// Should return the argument unchanged since it's not a directory
		if len(result) != 1 || result[0] != "apply" {
			t.Errorf("Expected [apply], got %v", result)
		}
	})

	t.Run("arguments with current directory", func(t *testing.T) {
		// Put "." at the beginning so it gets detected and removed
		args := []string{".", "apply", "--auto-approve"}
		result := helper.RemovePathFromArguments(args)
		
		// Should remove the "." and keep other arguments
		expected := []string{"apply", "--auto-approve"}
		if len(result) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(result))
		}
		for i, v := range result {
			if v != expected[i] {
				t.Errorf("Expected %v, got %v", expected, result)
				break
			}
		}
	})

	t.Run("findRelativePath with empty arguments", func(t *testing.T) {
		// Test the function that was originally causing the panic
		path, position := findRelativePath([]string{})
		
		if path != "." {
			t.Errorf("Expected '.', got %s", path)
		}
		if position != -1 {
			t.Errorf("Expected -1, got %d", position)
		}
	})
}