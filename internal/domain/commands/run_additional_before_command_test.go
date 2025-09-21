package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockShellRepositoryForAdditional is a mock implementation of repositories.ShellRepository.
type MockShellRepositoryForAdditional struct {
	ExecuteCallCount int
	LastCommand      string
	LastArguments    []string
	LastDirectory    string
	ExecuteErrors    []error
	CallHistory      []struct {
		Command   string
		Arguments []string
		Directory string
	}
}

func (m *MockShellRepositoryForAdditional) ExecuteCommand(
	command string,
	arguments []string,
	directory string,
) error {
	m.CallHistory = append(m.CallHistory, struct {
		Command   string
		Arguments []string
		Directory string
	}{
		Command:   command,
		Arguments: arguments,
		Directory: directory,
	})

	m.ExecuteCallCount++
	m.LastCommand = command
	m.LastArguments = arguments
	m.LastDirectory = directory

	if len(m.ExecuteErrors) > 0 {
		err := m.ExecuteErrors[0]
		m.ExecuteErrors = m.ExecuteErrors[1:]
		return err
	}
	return nil
}

// MockCLI is a mock implementation of entities.CLI.
type MockCLI struct {
	Name                  string
	CanChangeAccountValue bool
	CommandChangeAccount  []string
}

func (m *MockCLI) GetName() string {
	return m.Name
}

func (m *MockCLI) CanChangeAccount() bool {
	return m.CanChangeAccountValue
}

func (m *MockCLI) GetCommandChangeAccount() []string {
	return m.CommandChangeAccount
}

// MockAdditionalError implements the error interface.
type MockAdditionalError struct {
	message string
}

func (e *MockAdditionalError) Error() string {
	return e.message
}

func TestNewRunAdditionalBeforeCommand_ShouldCreateInstance_WhenValidDependenciesProvided(t *testing.T) {
	// GIVEN: Valid dependencies for creating the command
	settings := &entities.Settings{
		TerraCloud:              "aws",
		TerraTerraformWorkspace: "dev",
	}
	cli := &MockCLI{
		Name:                  "aws",
		CanChangeAccountValue: true,
		CommandChangeAccount:  []string{"sts", "assume-role"},
	}
	repository := &MockShellRepositoryForAdditional{}

	// WHEN: Creating a new RunAdditionalBeforeCommand
	cmd := commands.NewRunAdditionalBeforeCommand(settings, cli, repository)

	// THEN: Should return a valid command instance
	require.NotNil(t, cmd)
}

func TestRunAdditionalBeforeCommand_ShouldChangeAccount_WhenCLICanChangeAccount(t *testing.T) {
	// GIVEN: A command with CLI that can change account
	settings := &entities.Settings{
		TerraCloud: "aws",
	}
	cli := &MockCLI{
		Name:                  "aws",
		CanChangeAccountValue: true,
		CommandChangeAccount:  []string{"sts", "assume-role", "--role-arn", "test-role"},
	}
	repository := &MockShellRepositoryForAdditional{}
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
}

func TestRunAdditionalBeforeCommand_ShouldNotChangeAccount_WhenCLICannotChangeAccount(t *testing.T) {
	// GIVEN: A command with CLI that cannot change account
	settings := &entities.Settings{
		TerraCloud: "aws",
	}
	cli := &MockCLI{
		Name:                  "aws",
		CanChangeAccountValue: false,
		CommandChangeAccount:  []string{},
	}
	repository := &MockShellRepositoryForAdditional{}
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
}

func TestRunAdditionalBeforeCommand_ShouldNotChangeAccount_WhenCLIIsNil(t *testing.T) {
	// GIVEN: A command with nil CLI
	settings := &entities.Settings{
		TerraCloud: "aws",
	}
	repository := &MockShellRepositoryForAdditional{}
	cmd := commands.NewRunAdditionalBeforeCommand(settings, nil, repository)
	targetPath := "/test/path"
	arguments := []string{"plan"}

	// WHEN: Executing the command
	cmd.Execute(targetPath, arguments)

	// THEN: Should not execute any account change command
	for _, call := range repository.CallHistory {
		assert.NotEqual(t, "aws", call.Command, "Should not execute account change command")
	}
}

func TestRunAdditionalBeforeCommand_ShouldInitEnvironment_WhenArgumentsRequireInit(t *testing.T) {
	// GIVEN: A command with arguments that require environment initialization
	settings := &entities.Settings{
		TerraCloud: "aws",
	}
	repository := &MockShellRepositoryForAdditional{}
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
}

func TestRunAdditionalBeforeCommand_ShouldNotInitEnvironment_WhenArgumentsAreInit(t *testing.T) {
	// GIVEN: A command with 'init' argument
	settings := &entities.Settings{
		TerraCloud: "aws",
	}
	repository := &MockShellRepositoryForAdditional{}
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
}

func TestRunAdditionalBeforeCommand_ShouldNotInitEnvironment_WhenArgumentsAreRunAll(t *testing.T) {
	// GIVEN: A command with 'run-all' argument
	settings := &entities.Settings{
		TerraCloud: "aws",
	}
	repository := &MockShellRepositoryForAdditional{}
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
}

func TestRunAdditionalBeforeCommand_ShouldChangeWorkspace_WhenWorkspaceIsConfigured(t *testing.T) {
	// GIVEN: A command with configured workspace
	settings := &entities.Settings{
		TerraCloud:              "aws",
		TerraTerraformWorkspace: "production",
	}
	repository := &MockShellRepositoryForAdditional{}
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
}

func TestRunAdditionalBeforeCommand_ShouldNotChangeWorkspace_WhenWorkspaceIsEmpty(t *testing.T) {
	// GIVEN: A command with empty workspace
	settings := &entities.Settings{
		TerraCloud:              "aws",
		TerraTerraformWorkspace: "",
	}
	repository := &MockShellRepositoryForAdditional{}
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
}

func TestRunAdditionalBeforeCommand_ShouldExecuteAllSteps_WhenAllConditionsMet(t *testing.T) {
	// GIVEN: A command with all conditions met (account change, init, workspace change)
	settings := &entities.Settings{
		TerraCloud:              "aws",
		TerraTerraformWorkspace: "staging",
	}
	cli := &MockCLI{
		Name:                  "aws",
		CanChangeAccountValue: true,
		CommandChangeAccount:  []string{"sts", "assume-role"},
	}
	repository := &MockShellRepositoryForAdditional{}
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
}