//nolint:testpackage // Testing private functions and fields requires same package
package repositories

import (
	"testing"
)

func TestNewStdShellRepository(t *testing.T) {
	repo := NewStdShellRepository()

	if repo == nil {
		t.Fatal("NewStdShellRepository returned nil")
	}
}

func TestStdShellRepository_ExecuteCommand_ValidCommand(t *testing.T) {
	repo := NewStdShellRepository()

	// Test with a simple command that should work on any system
	err := repo.ExecuteCommand("echo", []string{"test"}, ".")

	if err != nil {
		t.Errorf("Expected no error for valid command, got: %v", err)
	}
}

func TestStdShellRepository_ExecuteCommand_InvalidCommand(t *testing.T) {
	repo := NewStdShellRepository()

	// Test with an invalid command
	err := repo.ExecuteCommand("nonexistentcommand12345", []string{}, ".")

	if err == nil {
		t.Error("Expected error for invalid command, got nil")
	}

	// Verify the error message contains expected text
	expectedErrorText := "failed to perform command execution"
	if !containsString(err.Error(), expectedErrorText) {
		t.Errorf("Expected error to contain %q, got: %v", expectedErrorText, err)
	}
}

func TestStdShellRepository_ExecuteCommand_InvalidDirectory(t *testing.T) {
	repo := NewStdShellRepository()

	// Test with a valid command but invalid directory
	err := repo.ExecuteCommand("echo", []string{"test"}, "/nonexistent/directory/12345")

	if err == nil {
		t.Error("Expected error for invalid directory, got nil")
	}

	// Verify the error message contains expected text
	expectedErrorText := "failed to perform command execution"
	if !containsString(err.Error(), expectedErrorText) {
		t.Errorf("Expected error to contain %q, got: %v", expectedErrorText, err)
	}
}

func TestStdShellRepository_ExecuteCommand_EmptyArguments(t *testing.T) {
	repo := NewStdShellRepository()

	// Test with empty arguments
	err := repo.ExecuteCommand("echo", []string{}, ".")

	if err != nil {
		t.Errorf("Expected no error for command with empty arguments, got: %v", err)
	}
}

func TestStdShellRepository_ExecuteCommand_MultipleArguments(t *testing.T) {
	repo := NewStdShellRepository()

	// Test with multiple arguments
	err := repo.ExecuteCommand("echo", []string{"hello", "world", "test"}, ".")

	if err != nil {
		t.Errorf("Expected no error for command with multiple arguments, got: %v", err)
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsStringHelper(s, substr))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
