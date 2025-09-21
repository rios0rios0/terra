package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/terra/internal/infrastructure/repositories"
)

func TestNewStdShellRepository(t *testing.T) {
	t.Parallel()
	
	t.Run("should create instance when called", func(t *testing.T) {
		// GIVEN: No preconditions needed

		// WHEN: Creating a new std shell repository
		repo := repositories.NewStdShellRepository()

		// THEN: Should return a valid repository instance
		require.NotNil(t, repo, "NewStdShellRepository should not return nil")
	})
}

func TestStdShellRepository_ExecuteCommand(t *testing.T) {
	t.Parallel()
	
	t.Run("should execute successfully when valid command provided", func(t *testing.T) {
		// GIVEN: A repository instance and valid command parameters
		repo := repositories.NewStdShellRepository()
		command := "echo"
		args := []string{"test"}
		workingDir := "."

		// WHEN: Executing a valid command
		err := repo.ExecuteCommand(command, args, workingDir)

		// THEN: Should execute without error
		assert.NoError(t, err, "Expected no error for valid command execution")
	})
	
	t.Run("should return error when invalid command provided", func(t *testing.T) {
		// GIVEN: A repository instance and invalid command
		repo := repositories.NewStdShellRepository()
		invalidCommand := "nonexistentcommand12345"
		args := []string{}
		workingDir := "."

		// WHEN: Executing an invalid command
		err := repo.ExecuteCommand(invalidCommand, args, workingDir)

		// THEN: Should return an error with expected message
		require.Error(t, err, "Expected error for invalid command")
		assert.Contains(t, err.Error(), "failed to perform command execution",
			"Error message should contain expected text")
	})
	
	t.Run("should return error when invalid directory provided", func(t *testing.T) {
		// GIVEN: A repository instance and invalid working directory
		repo := repositories.NewStdShellRepository()
		command := "echo"
		args := []string{"test"}
		invalidDir := "/nonexistent/directory/12345"

		// WHEN: Executing command in invalid directory
		err := repo.ExecuteCommand(command, args, invalidDir)

		// THEN: Should return an error with expected message
		require.Error(t, err, "Expected error for invalid directory")
		assert.Contains(t, err.Error(), "failed to perform command execution",
			"Error message should contain expected text")
	})
	
	t.Run("should execute successfully when empty arguments provided", func(t *testing.T) {
		// GIVEN: A repository instance and command with empty arguments
		repo := repositories.NewStdShellRepository()
		command := "echo"
		emptyArgs := []string{}
		workingDir := "."

		// WHEN: Executing command with empty arguments
		err := repo.ExecuteCommand(command, emptyArgs, workingDir)

		// THEN: Should execute without error
		assert.NoError(t, err, "Expected no error for command with empty arguments")
	})
	
	t.Run("should execute successfully when multiple arguments provided", func(t *testing.T) {
		// GIVEN: A repository instance and command with multiple arguments
		repo := repositories.NewStdShellRepository()
		command := "echo"
		multipleArgs := []string{"hello", "world", "test"}
		workingDir := "."

		// WHEN: Executing command with multiple arguments
		err := repo.ExecuteCommand(command, multipleArgs, workingDir)

		// THEN: Should execute without error
		assert.NoError(t, err, "Expected no error for command with multiple arguments")
	})
}
