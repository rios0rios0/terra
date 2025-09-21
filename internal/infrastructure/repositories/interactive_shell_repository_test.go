package repositories_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/infrastructure/repositories"
)

func TestNewInteractiveShellRepository(t *testing.T) {
	repo := repositories.NewInteractiveShellRepository()

	if repo == nil {
		t.Fatal("NewInteractiveShellRepository returned nil")
	}
}

func TestInteractiveShellRepository_ExecuteCommand_ValidCommand(t *testing.T) {
	repo := repositories.NewInteractiveShellRepository()

	// Test with a simple command that should work on any system
	err := repo.ExecuteCommand("echo", []string{"test"}, ".")

	if err != nil {
		t.Errorf("Expected no error for valid command, got: %v", err)
	}
}

func TestInteractiveShellRepository_ExecuteCommand_InvalidCommand(t *testing.T) {
	repo := repositories.NewInteractiveShellRepository()

	// Test with an invalid command
	err := repo.ExecuteCommand("nonexistentcommand12345", []string{}, ".")

	if err == nil {
		t.Error("Expected error for invalid command, got nil")
	}
}

func TestInteractiveShellRepository_ExecuteCommand_InvalidDirectory(t *testing.T) {
	repo := repositories.NewInteractiveShellRepository()

	// Test with an invalid directory
	err := repo.ExecuteCommand("echo", []string{"test"}, "/nonexistent/directory/12345")

	if err == nil {
		t.Error("Expected error for invalid directory, got nil")
	}
}