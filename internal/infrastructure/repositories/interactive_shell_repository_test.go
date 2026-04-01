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

func TestInteractiveShellRepository_ExecuteCommandWithAnswer(t *testing.T) {
	t.Parallel()

	t.Run("should execute successfully with auto-answer when command exits quickly", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An interactive shell repository and a command that exits immediately
		repo := repositories.NewInteractiveShellRepository()

		// WHEN: Executing with auto-answer
		err := repo.ExecuteCommandWithAnswer("echo", []string{"hello"}, ".", "y")

		// THEN: Should execute without error
		assert.NoError(t, err)
	})

	t.Run("should return error when invalid command provided with auto-answer", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An interactive shell repository and an invalid command
		repo := repositories.NewInteractiveShellRepository()

		// WHEN: Executing with auto-answer
		err := repo.ExecuteCommandWithAnswer("nonexistentcommand12345", []string{}, ".", "n")

		// THEN: Should return an error
		require.Error(t, err)
	})

	t.Run("should handle multi-line output with ANSI codes when auto-answering", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An interactive shell repository and a command producing ANSI-coded output
		repo := repositories.NewInteractiveShellRepository()

		// WHEN: Executing a command that produces colored output
		err := repo.ExecuteCommandWithAnswer("sh", []string{"-c", "printf '\\033[32mgreen\\033[0m\\n'; printf '\\033[31mred\\033[0m\\n'"}, ".", "y")

		// THEN: Should handle ANSI codes without error
		assert.NoError(t, err)
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

	t.Run("should complete execution without hanging when command produces output", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An interactive shell repository instance and command that produces multiple lines of output
		repo := repositories.NewInteractiveShellRepository()
		command := "sh"
		args := []string{"-c", "echo 'line1'; echo 'line2'; echo 'line3'"}
		workingDir := "."

		// WHEN: Executing command that produces output
		err := repo.ExecuteCommand(command, args, workingDir)

		// THEN: Should complete execution without hanging and without error
		assert.NoError(t, err, "Expected no error for command with output")
	})

	t.Run("should complete execution without hanging when command produces no output", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An interactive shell repository instance and command that produces no output
		repo := repositories.NewInteractiveShellRepository()
		command := "sh"
		args := []string{"-c", "exit 0"}
		workingDir := "."

		// WHEN: Executing command that produces no output
		err := repo.ExecuteCommand(command, args, workingDir)

		// THEN: Should complete execution without hanging and without error
		assert.NoError(t, err, "Expected no error for command with no output")
	})
}
