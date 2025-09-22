//go:build unit

package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/terra/internal/infrastructure/repositories"
)

func TestNewInteractiveShellRepository(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: No preconditions needed

		// WHEN: Creating a new interactive shell repository
		repo := repositories.NewInteractiveShellRepository()

		// THEN: Should return a valid repository instance
		require.NotNil(t, repo, "NewInteractiveShellRepository should not return nil")
	})
}

func TestInteractiveShellRepository_ExecuteCommand(t *testing.T) {
	t.Parallel()

	t.Run("should execute successfully when valid command provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An interactive shell repository instance and valid command parameters
		repo := repositories.NewInteractiveShellRepository()
		command := "echo"
		args := []string{"test"}
		workingDir := "."

		// WHEN: Executing a valid command
		err := repo.ExecuteCommand(command, args, workingDir)

		// THEN: Should execute without error
		assert.NoError(t, err, "Expected no error for valid command execution")
	})

	t.Run("should return error when invalid command provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An interactive shell repository instance and invalid command
		repo := repositories.NewInteractiveShellRepository()
		invalidCommand := "nonexistentcommand12345"
		args := []string{}
		workingDir := "."

		// WHEN: Executing an invalid command
		err := repo.ExecuteCommand(invalidCommand, args, workingDir)

		// THEN: Should return an error
		require.Error(t, err, "Expected error for invalid command")
	})

	t.Run("should return error when invalid directory provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An interactive shell repository instance and invalid working directory
		repo := repositories.NewInteractiveShellRepository()
		command := "echo"
		args := []string{"test"}
		invalidDir := "/nonexistent/directory/12345"

		// WHEN: Executing command in invalid directory
		err := repo.ExecuteCommand(command, args, invalidDir)

		// THEN: Should return an error
		require.Error(t, err, "Expected error for invalid directory")
	})

	t.Run("should handle empty arguments when no arguments provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An interactive shell repository instance with empty arguments
		repo := repositories.NewInteractiveShellRepository()
		command := "echo"
		args := []string{}
		workingDir := "."

		// WHEN: Executing command with empty arguments
		err := repo.ExecuteCommand(command, args, workingDir)

		// THEN: Should execute without error
		assert.NoError(t, err, "Expected no error for command with empty arguments")
	})
}

// Note: The pattern matching logic in processLineAndRespond and removeANSICodes methods
// are private and tested through integration tests. The core functionality includes:
// 1. External dependency prompt pattern: "should terragrunt apply the external dependency.*?"
// 2. Confirmation prompt pattern: "are you sure you want to run.*" (switches to manual mode)
// 3. Yes/No prompt pattern: ".*?.*[y/n]" (responds with 'n')
// 4. ANSI code removal for accurate pattern matching
//
// These patterns ensure auto-answering works correctly with various Terragrunt prompts
// while preserving user control for critical confirmation dialogs.
