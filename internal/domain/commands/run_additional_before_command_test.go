//go:build unit

package commands_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"

	"github.com/rios0rios0/terra/test/domain/entitybuilders"
	"github.com/rios0rios0/terra/test/domain/entitydoubles"

	"github.com/rios0rios0/terra/test/infrastructure/repositorydoubles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRunAdditionalBeforeCommand(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when valid dependencies provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN: Valid dependencies for creating the command
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			WithTerraTerraformWorkspace("dev").
			BuildSettings()
		cli := &entitydoubles.StubCLI{
			Name:                  "aws",
			CanChangeAccountValue: true,
			CommandChangeAccount:  []string{"sts", "assume-role"},
		}
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}

		// WHEN: Creating a new RunAdditionalBeforeCommand
		cmd := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)

		// THEN: Should return a valid command instance
		require.NotNil(t, cmd)
	})
}

func TestRunAdditionalBeforeCommand_Execute_AccountChange(t *testing.T) {
	t.Parallel()

	t.Run("should change account when CLI can change account", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with CLI that can change account
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		cli := &entitydoubles.StubCLI{
			Name:                  "aws",
			CanChangeAccountValue: true,
			CommandChangeAccount:  []string{"sts", "assume-role", "--role-arn", "test-role"},
		}
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)
		targetPath := t.TempDir()
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should execute account change command
		assert.GreaterOrEqual(t, repository.ExecuteCallCount, 1)
		assert.Equal(t, "aws", repository.CallHistory[0].Command)
		assert.Equal(
			t,
			[]string{"sts", "assume-role", "--role-arn", "test-role"},
			repository.CallHistory[0].Arguments,
		)
		assert.Equal(t, targetPath, repository.CallHistory[0].Directory)
	})

	t.Run("should not change account when CLI cannot change account", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with CLI that cannot change account
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		cli := &entitydoubles.StubCLI{
			Name:                  "aws",
			CanChangeAccountValue: false,
			CommandChangeAccount:  []string{},
		}
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)
		targetPath := t.TempDir()
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute account change command but may execute init/workspace commands
		// Only check that no account change command was executed
		for _, call := range repository.CallHistory {
			assert.NotEqual(t, "aws", call.Command, "Should not execute account change command")
		}
	})

	t.Run("should not change account when CLI is nil", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with nil CLI
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute any account change command
		for _, call := range repository.CallHistory {
			assert.NotEqual(t, "aws", call.Command, "Should not execute account change command")
		}
	})
}

//nolint:gocognit,gocyclo,cyclop // Large test function with comprehensive coverage - complex scenarios
func TestRunAdditionalBeforeCommand_Execute_EnvironmentInit(t *testing.T) {
	t.Parallel()

	t.Run("should init environment when arguments require init", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with arguments that require environment initialization
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		arguments := []string{"plan", "--detailed-exitcode"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should execute terragrunt init command (indirectly tests shouldInitEnvironment)
		initCommandExecuted := false
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 &&
				call.Arguments[0] == "init" {
				initCommandExecuted = true
				assert.Equal(t, targetPath, call.Directory)
				break
			}
		}
		assert.True(t, initCommandExecuted, "Should execute terragrunt init command")
	})

	t.Run("should not init environment when arguments are init", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with 'init' argument
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		arguments := []string{"init"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute terragrunt init command (indirectly tests shouldInitEnvironment)
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 &&
				call.Arguments[0] == "init" {
				assert.Fail(t, "Should not execute terragrunt init when argument is already init")
			}
		}
	})

	t.Run("should not init environment when arguments contain --all flag", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with '--all' flag
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		arguments := []string{"apply", "--all"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute terragrunt init command (indirectly tests shouldInitEnvironment)
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 &&
				call.Arguments[0] == "init" {
				assert.Fail(t, "Should not execute terragrunt init when using --all flag")
			}
		}
	})

	t.Run("should not init environment when --all flag is in different position", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with '--all' flag in different position
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		arguments := []string{"plan", "--detailed-exitcode", "--all", "--out=plan.out"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute terragrunt init command (indirectly tests shouldInitEnvironment)
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 &&
				call.Arguments[0] == "init" {
				assert.Fail(t, "Should not execute terragrunt init when --all flag is present")
			}
		}
	})

	t.Run("should change workspace when workspace is configured", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with configured workspace
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			WithTerraTerraformWorkspace("production").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should execute workspace change command (indirectly tests shouldChangeWorkspace)
		workspaceCommandExecuted := false
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) >= 4 &&
				call.Arguments[0] == "workspace" && call.Arguments[1] == "select" &&
				call.Arguments[2] == "-or-create" && call.Arguments[3] == "production" {
				workspaceCommandExecuted = true
				assert.Equal(t, targetPath, call.Directory)
				break
			}
		}
		assert.True(t, workspaceCommandExecuted, "Should execute workspace change command")
	})

	t.Run("should not change workspace when workspace is empty", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with empty workspace
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			WithTerraTerraformWorkspace("").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute workspace change command (indirectly tests shouldChangeWorkspace)
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) >= 1 &&
				call.Arguments[0] == "workspace" {
				assert.Fail(t, "Should not execute workspace command when workspace is empty")
			}
		}
	})

	t.Run("should not init environment when .terraform directory exists", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A directory that already has a .terraform subdirectory
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(targetPath, ".terraform"), 0o755))
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute terragrunt init because .terraform already exists
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 &&
				call.Arguments[0] == "init" {
				assert.Fail(t, "Should not execute terragrunt init when .terraform directory exists")
			}
		}
	})

	t.Run("should init environment when .terraform directory does not exist", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A directory without .terraform subdirectory
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should execute terragrunt init because .terraform does not exist
		initCommandExecuted := false
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 &&
				call.Arguments[0] == "init" {
				initCommandExecuted = true
				assert.Equal(t, targetPath, call.Directory)
				break
			}
		}
		assert.True(t, initCommandExecuted, "Should execute terragrunt init when .terraform directory is absent")
	})

	t.Run("should not init environment when .terragrunt-cache directory exists", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A directory that already has a .terragrunt-cache subdirectory
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(targetPath, ".terragrunt-cache"), 0o755))
		arguments := []string{"apply"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute terragrunt init because .terragrunt-cache already exists
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 &&
				call.Arguments[0] == "init" {
				assert.Fail(t, "Should not execute terragrunt init when .terragrunt-cache directory exists")
			}
		}
	})

	t.Run("should not init environment when terragrunt-cache directory exists", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A directory that already has a terragrunt-cache subdirectory (no leading dot)
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(targetPath, "terragrunt-cache"), 0o755))
		arguments := []string{"apply"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute terragrunt init because terragrunt-cache already exists
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 &&
				call.Arguments[0] == "init" {
				assert.Fail(t, "Should not execute terragrunt init when terragrunt-cache directory exists")
			}
		}
	})

	t.Run("should not change workspace when TERRA_NO_WORKSPACE is true", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with workspace configured but TERRA_NO_WORKSPACE enabled
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			WithTerraTerraformWorkspace("production").
			WithTerraNoWorkspace(true).
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute workspace change command even though workspace is set
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) >= 1 &&
				call.Arguments[0] == "workspace" {
				assert.Fail(t, "Should not execute workspace command when TERRA_NO_WORKSPACE is true")
			}
		}
	})

	t.Run("should change workspace when TERRA_NO_WORKSPACE is false", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with workspace configured and TERRA_NO_WORKSPACE disabled
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			WithTerraTerraformWorkspace("staging").
			WithTerraNoWorkspace(false).
			BuildSettings()
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := t.TempDir()
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should execute workspace change command
		workspaceCommandExecuted := false
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) >= 4 &&
				call.Arguments[0] == "workspace" && call.Arguments[1] == "select" &&
				call.Arguments[2] == "-or-create" && call.Arguments[3] == "staging" {
				workspaceCommandExecuted = true
				assert.Equal(t, targetPath, call.Directory)
				break
			}
		}
		assert.True(t, workspaceCommandExecuted, "Should execute workspace change command when TERRA_NO_WORKSPACE is false")
	})

	t.Run("should execute all steps when all conditions met", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command with all conditions met (account change, init, workspace change)
		settings := entitybuilders.NewSettingsBuilder().
			WithTerraCloud("aws").
			WithTerraTerraformWorkspace("staging").
			BuildSettings()
		cli := &entitydoubles.StubCLI{
			Name:                  "aws",
			CanChangeAccountValue: true,
			CommandChangeAccount:  []string{"sts", "assume-role"},
		}
		repository := &repositorydoubles.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)
		targetPath := t.TempDir()
		arguments := []string{"apply", "-auto-approve"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should execute all three types of commands
		assert.GreaterOrEqual(
			t,
			repository.ExecuteCallCount,
			3,
			"Should execute at least 3 commands",
		)

		// Verify account change command
		accountChangeFound := false
		initFound := false
		workspaceFound := false

		for _, call := range repository.CallHistory {
			if call.Command == "aws" {
				accountChangeFound = true
			}
			if call.Command == "terragrunt" && len(call.Arguments) > 0 &&
				call.Arguments[0] == "init" {
				initFound = true
			}
			if call.Command == "terragrunt" && len(call.Arguments) >= 4 &&
				call.Arguments[0] == "workspace" && call.Arguments[3] == "staging" {
				workspaceFound = true
			}
		}

		assert.True(t, accountChangeFound, "Should execute account change command")
		assert.True(t, initFound, "Should execute init command")
		assert.True(t, workspaceFound, "Should execute workspace command")
	})
}
