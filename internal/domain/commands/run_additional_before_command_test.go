package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRunAdditionalBeforeCommand(t *testing.T) {
	t.Parallel()
	
	t.Run("should create instance when valid dependencies provided", func(t *testing.T) {
		// GIVEN: Valid dependencies for creating the command
		settings := &entities.Settings{
			TerraCloud:              "aws",
			TerraTerraformWorkspace: "dev",
		}
		cli := &test.StubCLI{
			Name:                  "aws",
			CanChangeAccountValue: true,
			CommandChangeAccount:  []string{"sts", "assume-role"},
		}
		repository := &test.StubShellRepositoryForAdditional{}

		// WHEN: Creating a new RunAdditionalBeforeCommand
		cmd := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)

		// THEN: Should return a valid command instance
		require.NotNil(t, cmd)
	})
}

func TestRunAdditionalBeforeCommand_Execute(t *testing.T) {
	t.Run("should change account when CLI can change account", func(t *testing.T) {
		// GIVEN: A command with CLI that can change account
		settings := &entities.Settings{
			TerraCloud: "aws",
		}
		cli := &test.StubCLI{
			Name:                  "aws",
			CanChangeAccountValue: true,
			CommandChangeAccount:  []string{"sts", "assume-role", "--role-arn", "test-role"},
		}
		repository := &test.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)
		targetPath := "/test/path"
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should execute account change command
		assert.GreaterOrEqual(t, repository.ExecuteCallCount, 1)
		assert.Equal(t, "aws", repository.CallHistory[0].Command)
		assert.Equal(t, []string{"sts", "assume-role", "--role-arn", "test-role"}, repository.CallHistory[0].Arguments)
		assert.Equal(t, targetPath, repository.CallHistory[0].Directory)
	})
	
	t.Run("should not change account when CLI cannot change account", func(t *testing.T) {
		// GIVEN: A command with CLI that cannot change account
		settings := &entities.Settings{
			TerraCloud: "aws",
		}
		cli := &test.StubCLI{
			Name:                  "aws",
			CanChangeAccountValue: false,
			CommandChangeAccount:  []string{},
		}
		repository := &test.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)
		targetPath := "/test/path"
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
		// GIVEN: A command with nil CLI
		settings := &entities.Settings{
			TerraCloud: "aws",
		}
		repository := &test.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := "/test/path"
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute any account change command
		for _, call := range repository.CallHistory {
			assert.NotEqual(t, "aws", call.Command, "Should not execute account change command")
		}
	})
	
	t.Run("should init environment when arguments require init", func(t *testing.T) {
		// GIVEN: A command with arguments that require environment initialization
		settings := &entities.Settings{
			TerraCloud: "aws",
		}
		repository := &test.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := "/test/path"
		arguments := []string{"plan", "--detailed-exitcode"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should execute terragrunt init command (indirectly tests shouldInitEnvironment)
		initCommandExecuted := false
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 && call.Arguments[0] == "init" {
				initCommandExecuted = true
				assert.Equal(t, targetPath, call.Directory)
				break
			}
		}
		assert.True(t, initCommandExecuted, "Should execute terragrunt init command")
	})
	
	t.Run("should not init environment when arguments are init", func(t *testing.T) {
		// GIVEN: A command with 'init' argument
		settings := &entities.Settings{
			TerraCloud: "aws",
		}
		repository := &test.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := "/test/path"
		arguments := []string{"init"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute terragrunt init command (indirectly tests shouldInitEnvironment)
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 && call.Arguments[0] == "init" {
				assert.Fail(t, "Should not execute terragrunt init when argument is already init")
			}
		}
	})
	
	t.Run("should not init environment when arguments are run-all", func(t *testing.T) {
		// GIVEN: A command with 'run-all' argument
		settings := &entities.Settings{
			TerraCloud: "aws",
		}
		repository := &test.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := "/test/path"
		arguments := []string{"run-all", "plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute terragrunt init command (indirectly tests shouldInitEnvironment)
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) > 0 && call.Arguments[0] == "init" {
				assert.Fail(t, "Should not execute terragrunt init when argument is run-all")
			}
		}
	})
	
	t.Run("should change workspace when workspace is configured", func(t *testing.T) {
		// GIVEN: A command with configured workspace
		settings := &entities.Settings{
			TerraCloud:              "aws",
			TerraTerraformWorkspace: "production",
		}
		repository := &test.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := "/test/path"
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
		// GIVEN: A command with empty workspace
		settings := &entities.Settings{
			TerraCloud:              "aws",
			TerraTerraformWorkspace: "",
		}
		repository := &test.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
		targetPath := "/test/path"
		arguments := []string{"plan"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should not execute workspace change command (indirectly tests shouldChangeWorkspace)
		for _, call := range repository.CallHistory {
			if call.Command == "terragrunt" && len(call.Arguments) >= 1 && call.Arguments[0] == "workspace" {
				assert.Fail(t, "Should not execute workspace command when workspace is empty")
			}
		}
	})
	
	t.Run("should execute all steps when all conditions met", func(t *testing.T) {
		// GIVEN: A command with all conditions met (account change, init, workspace change)
		settings := &entities.Settings{
			TerraCloud:              "aws",
			TerraTerraformWorkspace: "staging",
		}
		cli := &test.StubCLI{
			Name:                  "aws",
			CanChangeAccountValue: true,
			CommandChangeAccount:  []string{"sts", "assume-role"},
		}
		repository := &test.StubShellRepositoryForAdditional{}
		cmd := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)
		targetPath := "/test/path"
		arguments := []string{"apply", "-auto-approve"}

		// WHEN: Executing the command
		cmd.Execute(targetPath, arguments)

		// THEN: Should execute all three types of commands
		assert.GreaterOrEqual(t, repository.ExecuteCallCount, 3, "Should execute at least 3 commands")

		// Verify account change command
		accountChangeFound := false
		initFound := false
		workspaceFound := false

		for _, call := range repository.CallHistory {
			if call.Command == "aws" {
				accountChangeFound = true
			}
			if call.Command == "terragrunt" && len(call.Arguments) > 0 && call.Arguments[0] == "init" {
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