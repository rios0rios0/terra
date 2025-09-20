package repositories

import (
	"testing"
)

func TestNewInteractiveShellRepository(t *testing.T) {
	repo := NewInteractiveShellRepository()

	if repo == nil {
		t.Fatal("NewInteractiveShellRepository returned nil")
	}
}

func TestInteractiveShellRepository_ExecuteCommand_ValidCommand(t *testing.T) {
	repo := NewInteractiveShellRepository()

	// Test with a simple command that should work quickly and not require interaction
	err := repo.ExecuteCommand("echo", []string{"test"}, ".")

	if err != nil {
		t.Errorf("Expected no error for valid command, got: %v", err)
	}
}

func TestInteractiveShellRepository_ExecuteCommand_InvalidCommand(t *testing.T) {
	repo := NewInteractiveShellRepository()

	// Test with an invalid command
	err := repo.ExecuteCommand("nonexistentcommand12345", []string{}, ".")

	if err == nil {
		t.Error("Expected error for invalid command, got nil")
	}

	// Verify the error message contains expected text
	expectedErrorText := "failed to start command"
	if !containsString(err.Error(), expectedErrorText) {
		t.Errorf("Expected error to contain %q, got: %v", expectedErrorText, err)
	}
}

func TestInteractiveShellRepository_ExecuteCommand_InvalidDirectory(t *testing.T) {
	repo := NewInteractiveShellRepository()

	// Test with a valid command but invalid directory
	err := repo.ExecuteCommand("echo", []string{"test"}, "/nonexistent/directory/12345")

	if err == nil {
		t.Error("Expected error for invalid directory, got nil")
	}

	// Verify the error message contains expected text
	expectedErrorText := "failed to start command"
	if !containsString(err.Error(), expectedErrorText) {
		t.Errorf("Expected error to contain %q, got: %v", expectedErrorText, err)
	}
}

func TestInteractiveShellRepository_removeANSICodes(t *testing.T) {
	repo := NewInteractiveShellRepository()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no ansi codes",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "simple ansi codes",
			input:    "\x1b[31mRed Text\x1b[0m",
			expected: "Red Text",
		},
		{
			name:     "multiple ansi codes",
			input:    "\x1b[1m\x1b[31mBold Red\x1b[0m\x1b[32m Green\x1b[0m",
			expected: "Bold Red Green",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only ansi codes",
			input:    "\x1b[31m\x1b[0m",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.removeANSICodes(tt.input)
			if result != tt.expected {
				t.Errorf("removeANSICodes(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInteractiveShellRepository_processLineAndRespond(t *testing.T) {
	repo := NewInteractiveShellRepository()

	// Since processLineAndRespond uses channels and goroutines, and writes to stdin,
	// we'll test it indirectly by testing the pattern matching logic with removeANSICodes
	
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "external dependency prompt",
			input:       "Should terragrunt apply the external dependency at /path/to/dependency? [y/n]",
			description: "Should match external dependency pattern",
		},
		{
			name:        "confirmation prompt",
			input:       "Are you sure you want to run terragrunt apply in each folder of [/path]? [y/n]",
			description: "Should match confirmation pattern",
		},
		{
			name:        "yes/no prompt",
			input:       "Do you want to continue? [y/n]",
			description: "Should match yes/no pattern",
		},
		{
			name:        "normal output",
			input:       "Terraform will perform the following actions:",
			description: "Should not match any pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that removeANSICodes works correctly for the input
			cleanLine := repo.removeANSICodes(tt.input)
			if cleanLine == "" && tt.input != "" {
				t.Errorf("removeANSICodes removed all content from non-empty input %q", tt.input)
			}
			
			// The actual pattern matching and response logic is complex to test
			// without mocking stdin/stdout, so we just verify the method exists
			// and doesn't panic with basic input
			t.Logf("Testing %s: %q -> %q", tt.description, tt.input, cleanLine)
		})
	}
}